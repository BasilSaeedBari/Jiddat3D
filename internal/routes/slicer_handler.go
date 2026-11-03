package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"jiddat3d/internal/compiler"
	"jiddat3d/internal/slicer"
)

// RegisterSlicerRoutes wires up the server-side instant slicer and HTMX endpoints
func RegisterSlicerRoutes(app core.App, e *core.ServeEvent) {

	// 1. Full Page Route: /instant-quote
	e.Router.GET("/instant-quote", func(c *core.RequestEvent) error {
		cfg := getActiveSlicerConfig(app)
		filaments := getActiveFilaments(app)

		siteMap := map[string]any{
			"whatsapp_business_number": "03001318112",
		}
		if siteRec, err := app.FindFirstRecordByFilter("site_settings", "1=1"); err == nil && siteRec != nil {
			siteMap["whatsapp_business_number"] = siteRec.GetString("whatsapp_business_number")
		}

		data := compiler.TemplateData{
			Title:       "Instant 3D Model Slicer & Quote | Jiddat3D Karachi",
			Description: "Upload your 3MF, OBJ, or STL 3D model for instant server-side slicing, print time calculation, and authentic rupee quotation from Jiddat3D Karachi.",
			Site:        siteMap,
			Record: map[string]any{
				"Config":    cfg,
				"Filaments": filaments,
			},
		}

		html, err := compiler.RenderTemplate("instant_quote.html", data)
		if err != nil {
			app.Logger().Error("Template render error for instant_quote.html: " + err.Error())
			return c.String(http.StatusInternalServerError, "Failed to render instant quote: "+err.Error())
		}

		return c.HTML(http.StatusOK, html)
	})

	// 2. Preemptive / Speculative Upload Endpoint: POST /api/instant-slice
	e.Router.POST("/api/instant-slice", func(c *core.RequestEvent) error {
		// Limit upload to 40 MB
		if err := c.Request.ParseMultipartForm(40 << 20); err != nil {
			return c.HTML(http.StatusBadRequest, renderErrorCard("File too large. Maximum supported upload is 40MB."))
		}

		file, header, err := c.Request.FormFile("model")
		if err != nil {
			return c.HTML(http.StatusBadRequest, renderErrorCard("No 3D model file received. Please select a .3mf, .obj, or .stl file."))
		}
		defer file.Close()

		filename := header.Filename
		ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])
		if ext != ".3mf" && ext != ".obj" && ext != ".stl" {
			return c.HTML(http.StatusBadRequest, renderErrorCard("Unsupported file extension ("+ext+"). Please upload a .3mf, .obj, or .stl file."))
		}

		// Read into memory buffer for parsing
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return c.HTML(http.StatusInternalServerError, renderErrorCard("Failed to read uploaded file buffer."))
		}
		fileBytes := buf.Bytes()

		// Parse 3D Geometry
		mesh, err := slicer.ParseModel(filename, bytes.NewReader(fileBytes), int64(len(fileBytes)))
		if err != nil {
			return c.HTML(http.StatusBadRequest, renderErrorCard("Could not parse 3D mesh: "+err.Error()))
		}

		// Load Slicer Configuration from Admin Settings
		cfg := getActiveSlicerConfig(app)
		filaments := getActiveFilaments(app)

		// Determine filament (default is White Basic PLA)
		filamentID := c.Request.FormValue("filament")
		chosenFilament := filaments[0] // fallback
		for _, f := range filaments {
			if f.ID == filamentID || (filamentID == "" && f.IsDefault) {
				chosenFilament = f
				break
			}
		}

		// Allow optional infill override from form
		if infillStr := c.Request.FormValue("infill"); infillStr != "" {
			if infillVal, err := strconv.ParseFloat(infillStr, 64); err == nil && infillVal >= 5 && infillVal <= 100 {
				cfg.InfillPercentage = infillVal
			}
		}

		// Run Slicing & LayerInvoice Calculation
		result := slicer.SliceAndCalculate(*mesh, chosenFilament, cfg)

		// Render the reactive HTMX Result Partial
		return c.HTML(http.StatusOK, renderQuoteResultHTML(result, filaments))
	})

	// 3. Reactive Recalculation Endpoint: POST /api/recalculate-quote
	e.Router.POST("/api/recalculate-quote", func(c *core.RequestEvent) error {
		if err := c.Request.ParseForm(); err != nil {
			return c.HTML(http.StatusBadRequest, renderErrorCard("Invalid parameters."))
		}

		// Rebuild Mesh from hidden inputs
		sizeX, _ := strconv.ParseFloat(c.Request.FormValue("mesh_size_x"), 64)
		sizeY, _ := strconv.ParseFloat(c.Request.FormValue("mesh_size_y"), 64)
		sizeZ, _ := strconv.ParseFloat(c.Request.FormValue("mesh_size_z"), 64)
		volCm3, _ := strconv.ParseFloat(c.Request.FormValue("mesh_vol_cm3"), 64)
		areaCm2, _ := strconv.ParseFloat(c.Request.FormValue("mesh_area_cm2"), 64)
		triCount, _ := strconv.Atoi(c.Request.FormValue("mesh_tri_count"))
		filename := c.Request.FormValue("mesh_filename")
		format := c.Request.FormValue("mesh_format")

		mesh := slicer.MeshStats{
			Filename:       filename,
			Format:         format,
			SizeX:          sizeX,
			SizeY:          sizeY,
			SizeZ:          sizeZ,
			VolumeCm3:      volCm3,
			SurfaceAreaCm2: areaCm2,
			TriangleCount:  triCount,
		}

		cfg := getActiveSlicerConfig(app)
		filaments := getActiveFilaments(app)

		if infillStr := c.Request.FormValue("infill"); infillStr != "" {
			if infillVal, err := strconv.ParseFloat(infillStr, 64); err == nil && infillVal >= 5 && infillVal <= 100 {
				cfg.InfillPercentage = infillVal
			}
		}

		filamentID := c.Request.FormValue("filament")
		chosenFilament := filaments[0]
		for _, f := range filaments {
			if f.ID == filamentID {
				chosenFilament = f
				break
			}
		}

		result := slicer.SliceAndCalculate(mesh, chosenFilament, cfg)
		return c.HTML(http.StatusOK, renderQuoteResultHTML(result, filaments))
	})
}

func getActiveSlicerConfig(app core.App) slicer.SlicerConfig {
	cfg := slicer.DefaultSlicerConfig()

	rec, err := app.FindFirstRecordByFilter("slicer_settings", "active = true")
	if err == nil && rec != nil {
		if tMult := rec.GetFloat("time_multiplier"); tMult > 0 {
			cfg.TimeMultiplier = tMult
		}
		if wMult := rec.GetFloat("weight_multiplier"); wMult > 0 {
			cfg.WeightMultiplier = wMult
		}
		if rate := rec.GetInt("machine_hourly_rate"); rate > 0 {
			cfg.MachineHourlyRate = int64(rate)
		}
		if labor := rec.GetInt("labor_overhead"); labor >= 0 {
			cfg.LaborOverhead = int64(labor)
		}
		if elec := rec.GetInt("electricity_hourly"); elec >= 0 {
			cfg.ElectricityHourly = int64(elec)
		}
		if fail := rec.GetInt("failure_rate_pct"); fail >= 0 {
			cfg.FailureRatePct = int64(fail)
		}
		if profit := rec.GetInt("profit_multiplier_pct"); profit >= 0 {
			cfg.ProfitMultiplierPct = int64(profit)
		}
		if defPrice := rec.GetInt("default_filament_price"); defPrice > 0 {
			cfg.DefaultFilamentPrice = int64(defPrice)
		}
		if name := rec.GetString("name"); name != "" {
			cfg.Name = name
		}
	}

	return cfg
}

func getActiveFilaments(app core.App) []slicer.FilamentProfile {
	defaults := slicer.DefaultFilaments()

	records, err := app.FindRecordsByFilter("filament_types", "active = true", "price_per_gram", 20, 0)
	if err != nil || len(records) == 0 {
		return defaults
	}

	var list []slicer.FilamentProfile
	for _, rec := range records {
		list = append(list, slicer.FilamentProfile{
			ID:           rec.Id,
			Name:         rec.GetString("name"),
			Material:     rec.GetString("material"),
			ColorName:    rec.GetString("color_name"),
			ColorHex:     rec.GetString("color_hex"),
			DensityGCm3:  rec.GetFloat("density_g_cm3"),
			PricePerGram: int64(rec.GetInt("price_per_gram")),
			IsDefault:    rec.GetBool("is_default"),
		})
	}

	if len(list) == 0 {
		return defaults
	}
	return list
}

func renderErrorCard(message string) string {
	return fmt.Sprintf(`
		<div class="nouveau-card p-6 bg-red-500/10 border border-red-500/40 text-red-500 dark:text-red-400 space-y-2 animate-fade-in">
			<div class="flex items-center gap-2 font-display font-bold text-sm">
				<span>⚠️ Slicing Error</span>
			</div>
			<p class="text-xs leading-relaxed">%s</p>
			<p class="text-[11px] text-[#635D54] dark:text-[#9DA3AF]">Make sure the model is not corrupted, is a valid .3mf, .obj, or .stl file, and has watertight faces.</p>
		</div>
	`, template.HTMLEscapeString(message))
}

func renderQuoteResultHTML(res slicer.SlicedQuoteResult, filaments []slicer.FilamentProfile) string {
	whatsappMsg := fmt.Sprintf(
		"Assalam o Alaikum Jiddat3D! I just used your instant server slicer for my model: %s\n"+
			"• Dimensions: %.1f x %.1f x %.1f mm\n"+
			"• Filament: %s [%.1fg]\n"+
			"• Estimated Print Time: %s (with 1.5x buffer)\n"+
			"• Instant Quote: Rs %d (Prepaid)\n"+
			"I want to proceed with this print!",
		res.Mesh.Filename,
		res.Mesh.SizeX, res.Mesh.SizeY, res.Mesh.SizeZ,
		res.Filament.Name, res.BufferedGrams,
		res.FormattedTime,
		res.FinalPrice,
	)

	waURL := "https://wa.me/923001318112?text=" + url.QueryEscape(whatsappMsg)

	// Build Filament Radio Options HTML
	filamentRadios := ""
	for _, f := range filaments {
		checked := ""
		activeClasses := "border-[#E5DEC9] dark:border-[#28303F] bg-white/70 dark:bg-[#14181F]/70"
		if f.ID == res.Filament.ID {
			checked = "checked"
			activeClasses = "border-filament bg-filament/10 dark:bg-filament/15 shadow-sm ring-1 ring-filament"
		}

		filamentRadios += fmt.Sprintf(`
			<label class="p-3 rounded-[12px_2px_12px_2px] border cursor-pointer transition-all flex items-center justify-between text-xs %s">
				<div class="flex items-center gap-2.5">
					<input type="radio" name="filament" value="%s" %s 
						hx-post="/api/recalculate-quote" 
						hx-include="#recalc-params" 
						hx-target="#quote-result" 
						class="accent-filament">
					<div>
						<div class="font-semibold text-parchment-ink dark:text-white flex items-center gap-1.5">
							<span class="w-2.5 h-2.5 rounded-full inline-block border border-black/20" style="background-color: %s;"></span>
							<span>%s</span>
						</div>
						<div class="text-[10px] text-[#635D54] dark:text-[#9DA3AF]">%s &bull; Rs %d/g</div>
					</div>
				</div>
				<span class="text-[11px] font-mono text-filament font-bold">Rs %d/g</span>
			</label>
		`, activeClasses, f.ID, checked, f.ColorHex, f.Name, f.Material, f.PricePerGram, f.PricePerGram)
	}

	return fmt.Sprintf(`
	<div class="space-y-6 animate-fade-in" id="quote-card">
		
		<!-- Hidden params for reactive HTMX recalculations -->
		<form id="recalc-params">
			<input type="hidden" name="mesh_filename" value="%s">
			<input type="hidden" name="mesh_format" value="%s">
			<input type="hidden" name="mesh_size_x" value="%.4f">
			<input type="hidden" name="mesh_size_y" value="%.4f">
			<input type="hidden" name="mesh_size_z" value="%.4f">
			<input type="hidden" name="mesh_vol_cm3" value="%.4f">
			<input type="hidden" name="mesh_area_cm2" value="%.4f">
			<input type="hidden" name="mesh_tri_count" value="%d">
		</form>

		<!-- Model Diagnostics & Geometry Card -->
		<div class="nouveau-card p-6 bg-parchment-surface dark:bg-[#14181F] border border-[#E5DEC9] dark:border-[#28303F] shadow-lg">
			<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-[#E2D7C7] dark:border-[#28303F] pb-4 mb-4">
				<div class="flex items-center gap-3">
					<div class="w-10 h-10 rounded-[10px_2px_10px_2px] bg-filament/10 text-filament flex items-center justify-center text-xl shrink-0">
						📦
					</div>
					<div>
						<div class="font-display font-bold text-base text-parchment-ink dark:text-white line-clamp-1">%s</div>
						<div class="text-[11px] font-mono text-[#635D54] dark:text-[#9DA3AF]">%s &bull; <span class="font-numeral font-bold">%d</span> triangles</div>
					</div>
				</div>
				<span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-[8px_2px_8px_2px] bg-[#3DDC84]/10 text-[#3DDC84] border border-[#3DDC84]/30 text-[10px] font-mono font-bold uppercase self-start sm:self-auto">
					✓ Geometry Verified
				</span>
			</div>

			<!-- Bounding Dimensions Matrix -->
			<div class="grid grid-cols-3 gap-3 text-center">
				<div class="p-2.5 rounded-[10px_2px_10px_2px] bg-white/60 dark:bg-[#0C0E11]/60 border border-[#E5DEC9] dark:border-[#28303F]">
					<span class="text-[10px] font-mono text-[#635D54] dark:text-[#9DA3AF] block uppercase">Width (X)</span>
					<span class="text-sm font-bold font-numeral text-parchment-ink dark:text-white">%.1f mm</span>
				</div>
				<div class="p-2.5 rounded-[10px_2px_10px_2px] bg-white/60 dark:bg-[#0C0E11]/60 border border-[#E5DEC9] dark:border-[#28303F]">
					<span class="text-[10px] font-mono text-[#635D54] dark:text-[#9DA3AF] block uppercase">Depth (Y)</span>
					<span class="text-sm font-bold font-numeral text-parchment-ink dark:text-white">%.1f mm</span>
				</div>
				<div class="p-2.5 rounded-[10px_2px_10px_2px] bg-white/60 dark:bg-[#0C0E11]/60 border border-[#E5DEC9] dark:border-[#28303F]">
					<span class="text-[10px] font-mono text-[#635D54] dark:text-[#9DA3AF] block uppercase">Height (Z)</span>
					<span class="text-sm font-bold font-numeral text-parchment-ink dark:text-white">%.1f mm</span>
				</div>
			</div>
		</div>

		<!-- Sliced Toolpath & Material Estimation Box -->
		<div class="nouveau-card p-6 bg-gradient-to-br from-parchment-surface via-[#FAF5EC] to-parchment-surface dark:from-[#14181F] dark:via-[#191524] dark:to-[#14181F] border-2 border-filament/40 shadow-xl space-y-6">
			
			<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
				
				<!-- Estimated Print Time with 1.5x Buffer -->
				<div class="p-5 rounded-[14px_3px_14px_3px] bg-white/80 dark:bg-[#0C0E11]/80 border border-[#E5DEC9] dark:border-[#28303F] space-y-1">
					<div class="flex items-center justify-between">
						<span class="text-[11px] font-mono uppercase tracking-wider text-[#635D54] dark:text-[#9DA3AF]">Estimated Print Time</span>
						<span class="text-[9px] font-mono bg-filament/10 text-filament px-1.5 py-0.5 rounded font-bold">1.5x Buffer Applied</span>
					</div>
					<div class="text-3xl font-display font-bold text-filament font-numeral">%s</div>
					<div class="text-[11px] text-[#635D54] dark:text-[#9DA3AF] leading-tight pt-1">
						Calculated for <span class="font-numeral font-bold">%d</span> layers @ 0.20mm standard resolution
					</div>
				</div>

				<!-- Estimated Filament Weight with 1.5x Buffer -->
				<div class="p-5 rounded-[14px_3px_14px_3px] bg-white/80 dark:bg-[#0C0E11]/80 border border-[#E5DEC9] dark:border-[#28303F] space-y-1">
					<div class="flex items-center justify-between">
						<span class="text-[11px] font-mono uppercase tracking-wider text-[#635D54] dark:text-[#9DA3AF]">Plastic Required</span>
						<span class="text-[9px] font-mono bg-[#17C3B2]/10 text-[#17C3B2] px-1.5 py-0.5 rounded font-bold">1.5x Buffer Applied</span>
					</div>
					<div class="text-3xl font-display font-bold text-[#17C3B2] font-numeral">%.1fg</div>
					<div class="text-[11px] text-[#635D54] dark:text-[#9DA3AF] leading-tight pt-1">
						Extruded shell + 20%% infill (%s)
					</div>
				</div>

			</div>

			<!-- Interactive Infill Density Slider -->
			<div class="p-4 rounded-[12px_2px_12px_2px] bg-white/60 dark:bg-[#0C0E11]/60 border border-[#E5DEC9] dark:border-[#28303F]">
				<div class="flex items-center justify-between mb-2">
					<label class="text-xs font-semibold text-parchment-ink dark:text-white">Internal Infill Density:</label>
					<span class="text-xs font-mono font-bold text-filament"><span class="font-numeral">%d</span>%% Infill</span>
				</div>
				<div class="grid grid-cols-4 gap-2 text-xs font-mono">
					<button type="button" 
						hx-post="/api/recalculate-quote?infill=15&filament=%s" 
						hx-include="#recalc-params" 
						hx-target="#quote-result" 
						class="py-1.5 px-2 rounded border text-center transition-all %s">
						15%% (Light)
					</button>
					<button type="button" 
						hx-post="/api/recalculate-quote?infill=20&filament=%s" 
						hx-include="#recalc-params" 
						hx-target="#quote-result" 
						class="py-1.5 px-2 rounded border text-center transition-all %s">
						20%% (Standard)
					</button>
					<button type="button" 
						hx-post="/api/recalculate-quote?infill=40&filament=%s" 
						hx-include="#recalc-params" 
						hx-target="#quote-result" 
						class="py-1.5 px-2 rounded border text-center transition-all %s">
						40%% (Tough)
					</button>
					<button type="button" 
						hx-post="/api/recalculate-quote?infill=100&filament=%s" 
						hx-include="#recalc-params" 
						hx-target="#quote-result" 
						class="py-1.5 px-2 rounded border text-center transition-all %s">
						100%% (Solid)
					</button>
				</div>
			</div>

			<!-- Selectable Filament Material Choices -->
			<div>
				<div class="flex items-center justify-between mb-2">
					<span class="text-xs font-semibold text-parchment-ink dark:text-white">Choose Production Filament:</span>
					<span class="text-[10px] font-mono text-[#635D54] dark:text-[#9DA3AF]">Default: White Basic PLA</span>
				</div>
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
					%s
				</div>
			</div>

			<!-- Comprehensive LayerInvoice Price Breakdown -->
			<div class="border-t border-[#E2D7C7] dark:border-[#28303F] pt-5 space-y-3">
				
				<div class="flex items-center justify-between text-xs text-[#635D54] dark:text-[#9DA3AF]">
					<span>Raw Filament Cost (%.1fg @ Rs %d/g):</span>
					<span class="font-mono font-numeral">Rs %d</span>
				</div>
				<div class="flex items-center justify-between text-xs text-[#635D54] dark:text-[#9DA3AF]">
					<span>Machine Time Cost (%.1f hrs @ Rs %d/hr):</span>
					<span class="font-mono font-numeral">Rs %d</span>
				</div>
				<div class="flex items-center justify-between text-xs text-[#635D54] dark:text-[#9DA3AF]">
					<span>Electricity & Bench Preparation:</span>
					<span class="font-mono font-numeral">Rs %d</span>
				</div>
				<div class="flex items-center justify-between text-xs text-[#635D54] dark:text-[#9DA3AF]">
					<span>Quality Assurance & Risk Buffer (%d%%):</span>
					<span class="font-mono font-numeral">Rs %d</span>
				</div>

				<div class="p-4 rounded-[12px_2px_12px_2px] bg-white dark:bg-[#0C0E11] border border-filament/40 flex items-center justify-between mt-3">
					<div>
						<span class="text-[10px] font-mono uppercase tracking-wider text-[#635D54] dark:text-[#9DA3AF] block">Total Prepaid Estimate</span>
						<span class="text-3xl font-display font-bold text-parchment-ink dark:text-white font-numeral-italic">Rs %d</span>
					</div>
					<div class="text-right">
						<span class="inline-block px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-[#3DDC84]/15 text-[#3DDC84] border border-[#3DDC84]/30 uppercase">
							100%% Prepaid Order
						</span>
						<span class="block text-[10px] text-[#635D54] dark:text-[#9DA3AF] mt-0.5">Karachi Dispatch</span>
					</div>
				</div>

				<!-- WhatsApp Action Button -->
				<div class="pt-2">
					<a href="%s" target="_blank" rel="noopener noreferrer" class="btn-whatsapp w-full py-4 text-sm font-bold justify-center shadow-xl gap-2">
						<svg class="w-5 h-5 fill-current" viewBox="0 0 24 24"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51a12.8 12.8 0 0 0-.57-.01c-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 0 1-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 0 1-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 0 1 2.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0 0 12.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 0 0 5.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 0 0-3.48-8.413Z"/></svg>
						<span>Lock In Order on WhatsApp (<span class="font-numeral">0300-1318112</span>)</span>
					</a>
					<p class="text-[11px] text-center text-[#635D54] dark:text-[#9DA3AF] mt-2">
						Your quote specs are copied to WhatsApp. Send your file for immediate engineer confirmation.
					</p>
				</div>

			</div>

		</div>

	</div>
	`,
		template.HTMLEscapeString(res.Mesh.Filename),
		template.HTMLEscapeString(res.Mesh.Format),
		res.Mesh.SizeX, res.Mesh.SizeY, res.Mesh.SizeZ,
		res.Mesh.VolumeCm3, res.Mesh.SurfaceAreaCm2,
		res.Mesh.TriangleCount,

		template.HTMLEscapeString(res.Mesh.Filename),
		template.HTMLEscapeString(res.Mesh.Format),
		res.Mesh.TriangleCount,

		res.Mesh.SizeX, res.Mesh.SizeY, res.Mesh.SizeZ,

		res.FormattedTime,
		res.LayerCount,

		res.BufferedGrams,
		res.Filament.Name,

		int(math.Round(res.Config.InfillPercentage)),

		res.Filament.ID, getInfillBtnClass(res.Config.InfillPercentage, 15),
		res.Filament.ID, getInfillBtnClass(res.Config.InfillPercentage, 20),
		res.Filament.ID, getInfillBtnClass(res.Config.InfillPercentage, 40),
		res.Filament.ID, getInfillBtnClass(res.Config.InfillPercentage, 100),

		filamentRadios,

		res.BufferedGrams, res.Filament.PricePerGram, res.RawFilamentCost,
		res.BufferedHours, res.Config.MachineHourlyRate, res.MachineTimeCost,
		res.LaborCost+res.ElectricityCost,
		res.Config.FailureRatePct, res.FailureAllowance,

		res.FinalPrice,
		waURL,
	)
}

func getInfillBtnClass(currentInfill, targetInfill float64) string {
	if math.Abs(currentInfill-targetInfill) < 1.0 {
		return "bg-filament text-white border-filament font-bold"
	}
	return "bg-white/80 dark:bg-[#14181F] text-[#635D54] dark:text-[#9DA3AF] border-[#E5DEC9] dark:border-[#28303F] hover:border-filament"
}

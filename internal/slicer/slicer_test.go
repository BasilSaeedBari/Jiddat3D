package slicer

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"testing"
)

// 10x10x10 mm cube: 12 triangles, volume = 1.0 cm^3, surface area = 6.0 cm^2.
var unitCube10mm = [][3][3]float32{
	// bottom (z=0)
	{{0, 0, 0}, {10, 10, 0}, {10, 0, 0}},
	{{0, 0, 0}, {0, 10, 0}, {10, 10, 0}},
	// top (z=10)
	{{0, 0, 10}, {10, 0, 10}, {10, 10, 10}},
	{{0, 0, 10}, {10, 10, 10}, {0, 10, 10}},
	// front (y=0)
	{{0, 0, 0}, {10, 0, 0}, {10, 0, 10}},
	{{0, 0, 0}, {10, 0, 10}, {0, 0, 10}},
	// back (y=10)
	{{0, 10, 0}, {10, 10, 10}, {10, 10, 0}},
	{{0, 10, 0}, {0, 10, 10}, {10, 10, 10}},
	// left (x=0)
	{{0, 0, 0}, {0, 0, 10}, {0, 10, 10}},
	{{0, 0, 0}, {0, 10, 10}, {0, 10, 0}},
	// right (x=10)
	{{10, 0, 0}, {10, 10, 10}, {10, 0, 10}},
	{{10, 0, 0}, {10, 10, 0}, {10, 10, 10}},
}

func assertCube10mm(t *testing.T, stats *MeshStats) {
	t.Helper()
	if stats.TriangleCount != 12 {
		t.Fatalf("expected 12 triangles, got %d", stats.TriangleCount)
	}
	for name, got := range map[string]float64{"SizeX": stats.SizeX, "SizeY": stats.SizeY, "SizeZ": stats.SizeZ} {
		if math.Abs(got-10.0) > 1e-4 {
			t.Fatalf("expected %s=10.0, got %f", name, got)
		}
	}
	if math.Abs(stats.VolumeCm3-1.0) > 1e-4 {
		t.Fatalf("expected VolumeCm3=1.0, got %f", stats.VolumeCm3)
	}
	if math.Abs(stats.SurfaceAreaCm2-6.0) > 1e-4 {
		t.Fatalf("expected SurfaceAreaCm2=6.0, got %f", stats.SurfaceAreaCm2)
	}
}

func buildBinarySTL(tris [][3][3]float32, headerPrefix string) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 84+len(tris)*50))
	header := make([]byte, 80)
	copy(header, headerPrefix)
	buf.Write(header)
	_ = binary.Write(buf, binary.LittleEndian, uint32(len(tris)))
	for _, tri := range tris {
		_ = binary.Write(buf, binary.LittleEndian, [3]float32{0, 0, 0}) // normal
		for _, v := range tri {
			_ = binary.Write(buf, binary.LittleEndian, v)
		}
		_ = binary.Write(buf, binary.LittleEndian, uint16(0))
	}
	return buf.Bytes()
}

func TestParseModel_BinarySTL(t *testing.T) {
	// Include "solid " in binary header to ensure binary detection doesn't misclassify it.
	raw := buildBinarySTL(unitCube10mm, "solid binary_cube")
	stats, err := ParseModel("cube.stl", bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("ParseModel binary STL failed: %v", err)
	}
	assertCube10mm(t, stats)
}

func TestParseModel_AsciiSTL(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("  solid cube\n")
	for _, tri := range unitCube10mm {
		sb.WriteString("  facet normal 0 0 0\n    outer loop\n")
		for _, v := range tri {
			fmt.Fprintf(&sb, "      VERTEX %.2f %.2f %.2f\n", v[0], v[1], v[2])
		}
		sb.WriteString("    endloop\n  endfacet\n")
	}
	sb.WriteString("endsolid cube\n")
	raw := []byte(sb.String())

	stats, err := ParseModel("cube.stl", bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("ParseModel ASCII STL failed: %v", err)
	}
	assertCube10mm(t, stats)
}

func TestParseModel_OBJ(t *testing.T) {
	obj := `
v 0 0 0
v 10 0 0
v 10 10 0
v 0 10 0
v 0 0 10
v 10 0 10
v 10 10 10
v 0 10 10
f 1 4 3 2
f 5 6 7 8
f 1 2 6 5
f 4 8 7 3
f 1 5 8 4
f -7 -5 -2 -6
`
	// Fix last face winding with standard indices (right face x=10: v2,v3,v7,v6)
	obj = strings.Replace(obj, "f -7 -5 -2 -6", "f -7 -6 -2 -3", 1)
	raw := []byte(obj)

	stats, err := ParseModel("cube.obj", bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("ParseModel OBJ failed: %v", err)
	}
	assertCube10mm(t, stats)
}

func TestParseModel_3MF(t *testing.T) {
	modelXML := `<?xml version="1.0" encoding="UTF-8"?>
<model unit="millimeter">
  <resources>
    <object id="1" type="model">
      <mesh>
        <vertices>
          <vertex x="0" y="0" z="0"/>
          <vertex x="10" y="0" z="0"/>
          <vertex x="10" y="10" z="0"/>
          <vertex x="0" y="10" z="0"/>
          <vertex x="0" y="0" z="10"/>
          <vertex x="10" y="0" z="10"/>
          <vertex x="10" y="10" z="10"/>
          <vertex x="0" y="10" z="10"/>
        </vertices>
        <triangles>
          <triangle v1="0" v2="2" v3="1"/>
          <triangle v1="0" v2="3" v3="2"/>
          <triangle v1="4" v2="5" v3="6"/>
          <triangle v1="4" v2="6" v3="7"/>
          <triangle v1="0" v2="1" v3="5"/>
          <triangle v1="0" v2="5" v3="4"/>
          <triangle v1="3" v2="6" v3="2"/>
          <triangle v1="3" v2="7" v3="6"/>
          <triangle v1="0" v2="4" v3="7"/>
          <triangle v1="0" v2="7" v3="3"/>
          <triangle v1="1" v2="6" v3="5"/>
          <triangle v1="1" v2="2" v3="6"/>
        </triangles>
      </mesh>
    </object>
  </resources>
</model>`
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	fw, err := zw.Create("3D/3dmodel.model")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = fw.Write([]byte(modelXML))
	_ = zw.Close()

	raw := buf.Bytes()
	stats, err := ParseModel("cube.3mf", bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("ParseModel 3MF failed: %v", err)
	}
	assertCube10mm(t, stats)
}

func TestParseModel_RejectsMalformedInputsWithoutPanicking(t *testing.T) {
	// 1. Truncated binary STL
	trunc := buildBinarySTL(unitCube10mm, "trunc")[:100]
	if _, err := ParseModel("trunc.stl", bytes.NewReader(trunc), int64(len(trunc))); err == nil {
		t.Error("expected error for truncated binary STL")
	}

	// 2. OBJ with out-of-bounds positive & negative indices
	badOBJ := []byte("v 0 0 0\nv 1 0 0\nv 0 1 0\nf 1 2 99\n")
	if _, err := ParseModel("bad.obj", bytes.NewReader(badOBJ), int64(len(badOBJ))); err == nil {
		t.Error("expected error for out-of-bounds OBJ vertex index")
	}
	badOBJNeg := []byte("v 0 0 0\nf -5 -4 -3\n")
	if _, err := ParseModel("badneg.obj", bytes.NewReader(badOBJNeg), int64(len(badOBJNeg))); err == nil {
		t.Error("expected error for negative out-of-bounds OBJ vertex index")
	}

	// 3. 3MF with negative vertex index
	badXML := `<model><resources><object id="1"><mesh><vertices><vertex x="0" y="0" z="0"/></vertices><triangles><triangle v1="-1" v2="0" v3="0"/></triangles></mesh></object></resources></model>`
	var zbuf bytes.Buffer
	zw := zip.NewWriter(&zbuf)
	fw, _ := zw.Create("3D/3dmodel.model")
	_, _ = fw.Write([]byte(badXML))
	_ = zw.Close()
	if _, err := ParseModel("bad.3mf", bytes.NewReader(zbuf.Bytes()), int64(zbuf.Len())); err == nil {
		t.Error("expected error for negative 3MF vertex index")
	}

	// 4. Unsupported extension
	if _, err := ParseModel("noext", bytes.NewReader(nil), 0); err == nil {
		t.Error("expected error for extensionless filename")
	}
}

func TestSliceAndCalculate_PricingAndInfill(t *testing.T) {
	// 50x50x50 mm cube = 125 cm^3 volume, 150 cm^2 surface area
	mesh := MeshStats{
		Filename:       "part.stl",
		Format:         "Binary STL",
		SizeX:          50,
		SizeY:          50,
		SizeZ:          50,
		VolumeCm3:      125,
		SurfaceAreaCm2: 150,
		TriangleCount:  12,
	}
	cfg20 := DefaultSlicerConfig()
	fil := DefaultFilaments()[0] // White Basic PLA, 1.24 g/cm3, Rs 8/g

	q20 := SliceAndCalculate(mesh, fil, cfg20)
	if q20.LayerCount != 250 {
		t.Fatalf("expected 250 layers for 50mm @ 0.20mm, got %d", q20.LayerCount)
	}
	if q20.FinalPrice%10 != 0 || q20.FinalPrice < 150 {
		t.Fatalf("unexpected FinalPrice %d", q20.FinalPrice)
	}

	// 100% infill must use strictly more plastic and cost more than 20% infill
	cfg100 := DefaultSlicerConfig()
	cfg100.InfillPercentage = 100
	q100 := SliceAndCalculate(mesh, fil, cfg100)
	if q100.BufferedGrams <= q20.BufferedGrams {
		t.Fatalf("expected 100%% infill grams (%.1f) > 20%% infill grams (%.1f)", q100.BufferedGrams, q20.BufferedGrams)
	}
	if q100.FinalPrice <= q20.FinalPrice {
		t.Fatalf("expected 100%% infill price (%d) > 20%% infill price (%d)", q100.FinalPrice, q20.FinalPrice)
	}

	// Tiny model hits minimum 5g and 15min floors -> Rs 260 under default LayerInvoice rates
	tiny := MeshStats{SizeX: 2, SizeY: 2, SizeZ: 2, VolumeCm3: 0.008, SurfaceAreaCm2: 0.24, TriangleCount: 12}
	qTiny := SliceAndCalculate(tiny, fil, cfg20)
	if qTiny.BufferedGrams < 5.0 {
		t.Fatalf("expected minimum 5g floor, got %.2f", qTiny.BufferedGrams)
	}
	if qTiny.FinalPrice != 260 {
		t.Fatalf("expected Rs 260 for tiny model under default rates, got %d", qTiny.FinalPrice)
	}

	// Low-rate config hits minimum Rs 150 shop order floor
	lowCfg := cfg20
	lowCfg.MachineHourlyRate = 50
	lowCfg.LaborOverhead = 0
	lowCfg.ElectricityHourly = 0
	lowCfg.ProfitMultiplierPct = 0
	if qLow := SliceAndCalculate(tiny, fil, lowCfg); qLow.FinalPrice != 150 {
		t.Fatalf("expected minimum Rs 150 floor, got %d", qLow.FinalPrice)
	}
}

func TestCbotCN1Preset_RulesSpeedsAccelerationsAndLimits(t *testing.T) {
	cfg := DefaultSlicerConfig()
	p := cfg.Preset

	// Verify exact Cbot CN-1 0.4 nozzle + V11 preset specs
	if p.NozzleDiameterMm != 0.4 || p.BuildXMm != 200 || p.BuildYMm != 200 || p.BuildZMm != 200 {
		t.Fatalf("unexpected nozzle/build volume: %.2fmm, %.0fx%.0fx%.0f", p.NozzleDiameterMm, p.BuildXMm, p.BuildYMm, p.BuildZMm)
	}
	if p.OuterWallSpeedMmS != 50 || p.InnerWallSpeedMmS != 80 || p.SparseInfillSpeedMmS != 150 || p.TopSurfaceSpeedMmS != 40 || p.TravelSpeedMmS != 150 {
		t.Fatalf("unexpected V11 process speeds: %+v", p)
	}
	if p.OuterWallAccelMmS2 != 800 || p.InnerWallAccelMmS2 != 1200 || p.SparseInfillAccelMmS2 != 1500 || p.MaxAccelXMmS2 != 2000 || p.MaxAccelExtrudingMmS2 != 1500 {
		t.Fatalf("unexpected V11 accelerations: %+v", p)
	}
	if p.RetractionLengthMm != 6.0 || p.RetractionSpeedMmS != 35 || p.ZHopMm != 0.4 || p.PrepareCompensationSec != 260 {
		t.Fatalf("unexpected retraction/zhop/prepare values: %+v", p)
	}
	if cfg.InfillPercentage != 10.0 || cfg.WallPerimeters != 2 {
		t.Fatalf("expected V11 default 10%% infill & 2 walls, got %.1f%% & %d walls", cfg.InfillPercentage, cfg.WallPerimeters)
	}

	// Verify volumetric flow clamping: 150 mm/s sparse infill @ 0.45 x 0.20mm = 13.5 mm^3/s
	// must clamp to <= 12.0 / (0.45 * 0.20) = 133.33 mm/s under White PLA 12 mm^3/s cap.
	clampedSpeed := effectiveSegmentSpeedMmS(150, 200, 1500, 1500, 20, 35, 0.45, 0.20, 12.0, 500)
	if clampedSpeed > 133.34 {
		t.Fatalf("expected volumetric-clamped speed <= 133.33 mm/s, got %.2f mm/s", clampedSpeed)
	}

	// Verify kinematics sensitivity: halving accelerations must strictly increase RawHours.
	mesh := MeshStats{
		Filename:       "bracket.stl",
		Format:         "Binary STL",
		SizeX:          50,
		SizeY:          50,
		SizeZ:          50,
		VolumeCm3:      125,
		SurfaceAreaCm2: 150,
		TriangleCount:  12,
	}
	fil := DefaultFilaments()[0]
	baseQuote := SliceAndCalculate(mesh, fil, cfg)

	slowAccelCfg := DefaultSlicerConfig()
	slowAccelCfg.Preset.OuterWallAccelMmS2 = 400
	slowAccelCfg.Preset.InnerWallAccelMmS2 = 600
	slowAccelCfg.Preset.SparseInfillAccelMmS2 = 750
	slowAccelCfg.Preset.MaxAccelExtrudingMmS2 = 750
	slowQuote := SliceAndCalculate(mesh, fil, slowAccelCfg)
	if slowQuote.RawHours <= baseQuote.RawHours {
		t.Fatalf("expected lower acceleration to increase print time (slow=%.3fh <= base=%.3fh)", slowQuote.RawHours, baseQuote.RawHours)
	}

	// Verify build volume boundary check (200x200x200mm)
	if baseQuote.ExceedsBuildVolume {
		t.Fatalf("50mm cube should fit within 200x200x200mm build volume")
	}
	oversize := mesh
	oversize.SizeX = 240
	overQuote := SliceAndCalculate(oversize, fil, cfg)
	if !overQuote.ExceedsBuildVolume {
		t.Fatalf("240mm wide part should flag ExceedsBuildVolume for 200mm Cbot CN-1 bed")
	}
}

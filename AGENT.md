# AGENT.md — Autonomous Agent Brief & Reference Guide

> **Welcome, Agent.** This file is the primary briefing document for any AI coding agent working in this repository.
> 
> **For the full, detailed architecture guide, directory map, and math specifications, see:**
> 📄 **[`CONTEXT.md`](./CONTEXT.md)**

---

## 1. Golden Rules (Non-Negotiables)

1. **Zero Node.js / NPM:** Do NOT create `package.json` or install Node packages. Tailwind CSS is compiled strictly via the bundled standalone binary (`./tailwindcss`).
2. **Karachi Exclusivity:** Jiddat3D is strictly **Karachi-based**. Never introduce references to Lahore or other cities as our headquarters.
3. **In-House Specialized Hardware Only:** We do **not** sell, dropship, or resell ready-made international consumer 3D printers (no Creality, Bambu Lab, Anycubic, etc.). We exclusively build our own specialized CoreXY and Cartesian FDM/FFM printers.
4. **Prepaid Orders Only (No COD):** We do not offer Cash on Delivery. All purchases are prepaid via direct bank transfer (IBFT) or Raast with official NTN invoices.
5. **Warranties:** 6-Month Working Warranty (with proper bench claim clauses) and STEAM Education Warranty for academic institutions. No generic 12-month claims.
6. **100% Open Source Sovereignty:** All hardware CAD, Marlin/Klipper profiles, and schematics remain open source.
7. **Filament & Resin Roadmap:** In-house filament extrusion and liquid UV resin (SLA/MSLA) machines are in active development (waitlists and roadmap transparency).
8. **Art Nouveau & Memphis Design Language:** Maintain the organic curves (`rounded-nouveau`), Memphis squiggle dividers, calligraphic numerals (`font-numeral`), and the mounted logo plinth (`bg-white/[0.97]`).
9. **WhatsApp Sales Pipeline:** The primary sales and inquiry conversion point is WhatsApp Business at `0300-1318112` (`+92 300 1318112`).

---

## 2. Core File Locations & What They Do

- **`cmd/jiddat/main.go`:** Main Go entrypoint. Mounts PocketBase, registers routes, hooks, and static file serving.
- **`internal/slicer/geometry.go`:** Fast geometric mesh parser for `.3mf`, `.obj`, and `.stl` files. Computes volume and surface area.
- **`internal/slicer/calculator.go`:** LayerInvoice math engine calculating filament weight, print hours, internal buffers, and rupee quotes.
- **`internal/routes/slicer_handler.go`:** Speculative HTMX upload and slicing endpoint (`POST /api/instant-slice`).
- **`internal/routes/public.go`:** Public route dispatcher, dynamic page rendering, and static file server.
- **`ui/templates/pages/`:** Go `html/template` page templates (`home.html`, `instant_quote.html`, `products.html`, `print_service.html`, `filament.html`, `schools.html`, `workshops.html`, `learn.html`, `about.html`, `gallery.html`, `support.html`, `quiz.html`).
- **`ui/templates/partials/nav.html`:** Responsive header with dropdown categories (`Hardware & Services`, `Labs & Learning`), theme toggle, and logo plinth.
- **`ui/templates/partials/instant_quote_result.html`:** HTMX response template with model dimensions, print duration, weight, LayerInvoice price, and prefilled WhatsApp lock-in button.
- **`ui/static/css/input.css` & `output.css`:** Tailwind source and minified output.
- **`compose.yaml` & `Dockerfile`:** Container definitions for production deployment on port `8090`.

---

## 3. How to Make and Verify Changes

```bash
# Rebuild CSS if templates or styles change:
./tailwindcss -i ui/static/css/input.css -o ui/static/css/output.css --minify

# Start the server locally:
./jiddat3d serve --http="127.0.0.1:8090"
# (Front-end at http://127.0.0.1:8090, PocketBase Admin at http://127.0.0.1:8090/_/)

# Check git status:
git status
```

For complete details on database migrations, accounting formulas, and deployment, please consult **[`CONTEXT.md`](./CONTEXT.md)**.

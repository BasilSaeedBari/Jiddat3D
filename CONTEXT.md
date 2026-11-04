# Jiddat3D — Complete Project Context & Architecture Guide (CONTEXT.md)

> **Purpose of this Document:** This document serves as the master context and architecture guide for any AI coding agent or human developer working on the **Jiddat3D Website & Digital Engineering Hub**. It records every architectural decision, file location, business rule, and technical pattern implemented across the project.

---

## 1. Executive Summary & Brand Identity

**Jiddat3D** (`Jiddat` = Innovation / Novelty in Urdu) is an engineering innovation lab based exclusively in **Karachi, Sindh, Pakistan**.

The platform is designed around **sovereign creation** rather than passive consumerism:
1. **Core Purpose:** Front-door sales engine, instant 3D model quoting tool, educational hub, and direct WhatsApp sales pipeline (`0300-1318112`).
2. **Specialized In-House Hardware:** Jiddat3D designs, bench-assembles, and print-tests its own specialized CoreXY and Cartesian FDM/FFM 3D printers. **We do not resell, dropship, or distribute ready-made international consumer printers (e.g., Creality, Bambu Lab, Anycubic).**
3. **Engineering Innovation Lab Vision:** The company is expanding beyond 3D printing into **Rapid PCB Prototyping** and **Desktop CNC Milling Machines**.
4. **100% Open Source Guarantee:** All CAD models, Marlin/Klipper profiles, and wiring schematics are open-source. If Jiddat3D ever ceases to exist, customer machines remain fully repairable with standard off-the-shelf components.
5. **Prepaid Orders Only (No COD):** Due to the high degree of custom assembly and bench calibration, Cash on Delivery is strictly not offered. All transactions are prepaid via Raast or direct bank transfer (IBFT) with official NTN invoices.
6. **Warranties:** 
   - **6-Month Working Warranty:** Covers manufacturing defects under clear bench claim clauses.
   - **STEAM Education Warranty:** Specialized packages and SLA support for schools, makerspaces, and engineering universities.
   - *(Note: We do not offer generic 12-month warranties).*
7. **Filaments & Liquid UV Resin:**
   - **Filaments:** In-house filament extrusion is in active development on our roadmap. Machines are calibrated for standard 1kg precision spools.
   - **Resin Printers:** High-resolution stereolithography (SLA/MSLA) systems are in active development; visitors can join a waitlist via WhatsApp.

---

## 2. Technology Stack & Design Principles

The stack is intentionally lean, self-contained, and devoid of Node.js / JavaScript bundle bloat:

| Layer | Technology | Details |
| :--- | :--- | :--- |
| **Backend & DB** | [Go 1.22+](https://go.dev/) + [PocketBase](https://pocketbase.io/) | PocketBase is embedded as a Go library; SQLite database stored in `pb_data/`. |
| **HTML Templating** | Go `html/template` + [Goldmark](https://github.com/yuin/goldmark) | Dynamic SSR with cached pre-rendering capabilities and custom markdown extensions. |
| **Styling** | [Tailwind CSS v3](https://tailwindcss.com/) | Built using the standalone Linux CLI (`./tailwindcss`), zero `node_modules` or npm. |
| **Frontend State** | [HTMX](https://htmx.org/) + [Alpine.js](https://alpinejs.dev/) | Vendored vanilla JS for hypermedia swapping and declarative component state. |
| **Deployment** | Docker & Docker Compose | Multi-stage build producing an ultra-light ~40MB Alpine container running on port `8090`. |

### Design System & Visual Language
- **Themes:** Dynamic Light (*Artisan Parchment*) and Dark (*Modern Midnight Workshop*) with zero-flicker client-side initialization via `localStorage`.
- **Curves & Borders:** Art Nouveau whiplash curves (`rounded-nouveau`, `rounded-[18px_4px_18px_4px]`) instead of generic bubbly borders.
- **Dividers:** Decorative Memphis squiggles and subtle grid lattices.
- **Typography:** Calligraphic figures (`font-numeral`) for rupee prices and metrics; serif headlines (`font-display`) and clean sans body (`font-sans`).
- **Logo Mounting:** The official dark violet and magenta logo is mounted on a crisp plinth (`bg-white/[0.97]` with delicate borders and violet ambient glow) to guarantee visual prominence across both light and dark themes.

---

## 3. Directory & File Map

```
Jiddat3D/
├── cmd/
│   └── jiddat/
│       └── main.go                 # Application entrypoint: registers PocketBase, routes, hooks, and slicer
├── internal/
│   ├── compiler/
│   │   ├── images.go               # Responsive WebP image processing via cwebp
│   │   ├── markdown.go             # Goldmark markdown compiler with syntax highlighting & typographer
│   │   └── render.go               # Go html/template renderer & static caching pipeline
│   ├── hooks/
│   │   └── content_hooks.go        # PocketBase event hooks (auto-render/cache HTML on CMS updates)
│   ├── mailer/
│   │   └── mail.go                 # Contact form email dispatcher (SMTP/Mailtrap)
│   ├── routes/
│   │   ├── public.go               # Public HTTP endpoints, static asset router, Service Worker handler
│   │   └── slicer_handler.go       # Speculative instant 3D model slicing endpoint (/api/instant-slice)
│   └── slicer/
│       ├── geometry.go             # 3D mesh volume & surface area parser (.3MF, .OBJ, .STL)
│       └── calculator.go           # LayerInvoice accounting & cost estimation formula engine
├── pb_data/                        # SQLite databases (data.db, auxiliary.db) & uploaded media
├── pb_migrations/                  # PocketBase schema & data migrations in Go
│   ├── 1700000000_collections.go   # Base collections (content, products, subscribers, contact, settings)
│   ├── 1700000001_seed_settings.go # Default settings seed
│   ├── 1700000002_seed_blog.go     # Initial knowledgebase content
│   ├── 1700000003_update_settings_schema.go
│   ├── 1700000004_business_seed.go # Karachi lab info, open-source guarantee, warranty terms
│   └── 1700000005_slicer_settings.go # Slicer configuration and filament pricing collections
├── ui/
│   ├── static/
│   │   ├── css/
│   │   │   ├── input.css           # Tailwind source CSS with custom Art Nouveau/Memphis components
│   │   │   └── output.css          # Minified production CSS generated by standalone Tailwind CLI
│   │   ├── js/                     # Vendored client scripts (alpine.min.js, htmx.min.js, sw-register.js, etc.)
│   │   ├── img/                    # Logos, plinths, and workshop photography
│   │   └── sw.js                   # Service Worker providing offline caching
│   └── templates/
│       ├── layouts/
│       │   └── base.html           # Master HTML shell, theme scripts, header/footer injection
│       ├── partials/
│       │   ├── nav.html            # Sticky header with grouped dropdowns and brand plinth
│       │   ├── footer.html         # Global footer with lab credentials and NTN
│       │   └── instant_quote_result.html # HTMX fragment with calculation breakdown & WhatsApp lock-in button
│       └── pages/
│           ├── home.html           # Main landing page (hero, printers, comparison matrix, open source)
│           ├── instant_quote.html  # Instant 3D Model Slicer interface (.3mf, .obj, .stl)
│           ├── products.html       # Specialized CoreXY & Cartesian 3D printer catalog (/printers)
│           ├── product_detail.html # Deep dive specification sheet for individual machines
│           ├── print_service.html  # On-demand engineering print service & FYP clinic
│           ├── filament.html       # Climate-calibrated spools (in plan) & stock replacement parts
│           ├── schools.html        # Campus makerspaces, university tenders, and STEAM warranty
│           ├── workshops.html      # Hands-on hardware assembly and Klipper firmware bootcamps
│           ├── learn.html          # Pakistan 3D Printing Handbook & UPS load-shedding sizing
│           ├── learn_detail.html   # Individual handbook article view
│           ├── about.html          # Maker rhetoric, sovereign engineering, open source guarantee
│           ├── gallery.html        # Real workshop print portfolio with weight, time, and PKR cost
│           ├── support.html        # 6-Month Working Warranty & STEAM Warranty claim policies
│           ├── quiz.html           # 1-Minute Machine Finder Quiz (/find-my-printer)
│           ├── contact.html        # Contact inquiry form
│           ├── community.html      # WhatsApp and Discord community links
│           └── 404.html            # Custom error template
├── scripts/                        # Database verification and maintenance Go scripts
├── compose.yaml                    # Production Docker Compose orchestration
├── Dockerfile                      # Multi-stage Docker build file
├── README.md                       # Public-facing repository documentation
├── tailwind.config.js              # Custom Tailwind configuration (colors, fonts, shapes)
└── tailwindcss                     # Standalone Linux x64 Tailwind compiler binary
```

---

## 4. Instant 3D Slicer & Quotation Engine

One of the flagship technical features is the **server-side instant slicer** located at `/instant-quote`:

### Preemptive / Speculative Execution
Using HTMX (`hx-post="/api/instant-slice"` and `hx-trigger="change"`), as soon as the user selects or drops a `.3mf`, `.obj`, or `.stl` file, the browser immediately transmits the payload to the server. Slicing begins speculatively before the user even clicks a button.

### Geometry Parsing (`internal/slicer/geometry.go`)
- **.3MF:** Reads the ZIP container, parses `3D/3dmodel.model` XML, extracts vertices and `<triangle>` definitions, and computes discrete signed volume:
  $$\text{Volume} = \frac{1}{6} \sum (v_1 \times v_2) \cdot v_3$$
- **.OBJ:** Fast stream scanner parsing `v` (vertices) and `f` (polygon faces). Supports arbitrary polygons (fan triangulation) and negative relative indexing.
- **.STL:** Auto-detects binary (80-byte header + 32-bit triangle count) vs ASCII formats and accumulates tetrahedral volume.

### Accounting & Multiplier Math (`internal/slicer/calculator.go`)
Based on the formulas in [`BasilSaeedBari/LayerInvoice`](https://github.com/BasilSaeedBari/LayerInvoice):
1. **Raw Weight:** Calculated using material density (e.g. PLA $\approx 1.24\text{ g/cm}^3$) multiplied by estimated infill and wall structures.
2. **Machine Hours:** Derived from path length, travel speeds, acceleration, and layer heights.
3. **Internal Calibration Multipliers:**
   - `time_multiplier` (default `1.5`): Accounts for Karachi grid load-shedding, warmup, and UPS transition buffers.
   - `weight_multiplier` (default `1.5`): Accounts for supports, brims, rafts, and purge line waste.
4. **Final Quotation Formula:**
   $$\text{Filament Cost} = \text{Adjusted Grams} \times \text{Filament Rate per Gram}$$
   $$\text{Machine Cost} = \text{Adjusted Hours} \times \text{Hourly Rate (Rs 250/hr)}$$
   $$\text{Production Cost} = (\text{Filament Cost} + \text{Machine Cost} + \text{Labor} + \text{Electricity}) \times 1.05 \text{ (Failure buffer)}$$
   $$\text{Final Client Price} = \text{Production Cost} + \text{Markup}$$

> [!NOTE]
> The client-facing quote renders clean, professional values (dimensions, print duration, filament weight, price). Internal formulas and buffer multipliers are kept strictly in the admin panel and are not exposed to the user.

### PocketBase Slicer Collections
- **`slicer_settings`:** Contains active configuration records (`time_multiplier`, `weight_multiplier`, `machine_hourly_rate`, `default_filament_price`, `profile_file`).
- **`filament_types`:** Library of available printing materials with their density, default temperature, and PKR price per gram (e.g., *White Basic PLA*, *Tough PLA+*, *High-Temp PETG*, *TPU 95A*).

---

## 5. Navigation & Routes

The site's navigation is divided into two primary dropdowns and direct action links:

1. **Hardware & Services (Dropdown):**
   - `/printers` (or `/products`) — Our Specialized Printers
   - `/instant-quote` — Instant Slicer & Rupee Quote
   - `/print-service` — On-Demand Engineering Print Service
   - `/filament` — Filament Roadmap & Stock Spares
2. **Labs & Learning (Dropdown):**
   - `/schools` — Schools & University Labs
   - `/workshops` — Hands-On Workshops
   - `/learn` — Pakistan 3D Printing Handbook
3. **Direct Navigation Links:**
   - `/gallery` — Authentic Workshop Portfolio
   - `/support` — 6-Month Working Warranty & STEAM Terms
   - `/about` — Sovereign Engineering Mission & Open Source Guarantee
   - `/find-my-printer` — 1-Minute Machine Finder Quiz
4. **Contact & Community:**
   - WhatsApp Hotline: `+92 300 1318112` (`0300-1318112`)
   - PocketBase Admin Panel: `http://127.0.0.1:8090/_/`

---

## 6. Developer & Agent Cheat Sheet

### Common Commands

```bash
# 1. Compile Tailwind CSS
./tailwindcss -i ui/static/css/input.css -o ui/static/css/output.css --minify

# 2. Run the application locally using the precompiled binary
./jiddat3d serve --http="127.0.0.1:8090"

# 3. Run from Go source (if Go 1.22+ is installed in PATH)
go run ./cmd/jiddat serve --http="127.0.0.1:8090"

# 4. Build Docker container image locally
docker build -t jiddat3d:latest .

# 5. Run via Docker Compose
docker compose up -d --build
```

### Critical Rules for Future Changes
1. **Never introduce Node.js, npm, or JS bundlers.** Use the standalone `./tailwindcss` binary.
2. **No mentions of Lahore or other cities as headquarters.** Jiddat3D is strictly **Karachi-based**.
3. **No dropshipped/imported printers.** We only build our own specialized in-house designs.
4. **No Cash on Delivery (COD).** Prepaid via Raast/IBFT only.
5. **Preserve Open Source Guarantee.** Emphasize that all CAD, BOMs, and configs are open.
6. **Preserve Logo Plinth.** The logo must remain on its dedicated plinth (`bg-white/[0.97]` with subtle borders) so the dark purple and magenta brand colors do not get lost against dark backgrounds.
7. **Keep Template Syntax Clean.** All templates use Go standard `html/template`. Do not break Go template actions (`{{ ... }}`).

# Jiddat3D - The Website & Digital Hub

> **A Note on Scope:** This repository contains the source code, architecture, and deployment mechanisms exclusively for the **Jiddat3D Website & Digital Engineering Hub**. It serves as the digital front-door, content management system, instant 3D model slicer, and quotation engine for the initiative. The hardware CAD models, firmware configurations, and BOMs for our physical machines are hosted in separate, dedicated open-source repositories.

---

## 1. The Rhetoric: Ideology & Origin

### The Maker Movement and RepRap
In 2005, the RepRap project—an acronym for *Replicating Rapid-prototyper*—ignited a revolution in digital fabrication. By open-sourcing hardware designs and creating a machine capable of printing its own components, it shattered the corporate monopoly on 3D printing. It proved that a collective of passionate people freely sharing knowledge can out-innovate massive corporate entities. It democratized manufacturing, bringing it from corporate labs straight into the hands of hackers, tinkerers, and students around the world.

### Where Does Pakistan Fit In?
Historically, the region has been a cradle of invention—from the urban planning and metallurgy of the Indus Valley to the pioneering of optics and algebra during the Islamic Golden Age. Innovation is in our blood. Yet, in the modern era, Pakistan largely finds itself positioned as a consumer rather than a creator, often lagging 10-20 years behind the cutting edge of manufacturing technology. To much of the world, we are in the era of AI and advanced robotics; to many in Pakistan, basic digital fabrication is still a novelty.

### The Genesis of Jiddat3D
Jiddat3D was born from a singular vision: **to restore the culture of sovereign creation in Pakistan.** Buying closed-source, imported turn-key appliances reinforces a passive consumer mindset. When those machines break, users must wait weeks for expensive proprietary parts and lack local support. 

By engineering our own specialized hardware right here in **Karachi**, we understand every nut, bolt, stepper motor, and firmware parameter:
- **100% In-House Specialized Machines:** We do not dropship or sell ready-made international consumer brands (Creality, Bambu Lab, etc.). We design, bench-assemble, and print-test our own specialized CoreXY and Cartesian FDM/FFM 3D printers.
- **100% Open Source Architecture Guarantee:** If tomorrow Jiddat3D ceases to exist, your machine will never become an obsolete brick. All CAD models, Marlin/Klipper profiles, and wiring schematics remain public and repairable with standard off-the-shelf components.
- **Prepaid Custom Builds (No COD):** Every unit is custom-calibrated with multi-hour bench testing before dispatch. We operate strictly on prepaid direct bank wire (IBFT) and Raast with formal NTN invoices.
- **6-Month Working Warranty & STEAM Education Warranty:** Clear bench claim clauses for manufacturing defects, plus dedicated support packages for partner schools and engineering universities.
- **Filament & Resin Roadmap:** In-house filament extrusion and ultra-high resolution liquid UV resin (SLA/MSLA) systems are in active engineering development. All machines are calibrated for standard, full-size 1kg spools.
- **Engineering Innovation Lab Vision:** Active roadmap development expanding beyond 3D printing into Rapid PCB Prototyping and Desktop CNC Milling Machines.

---

## 2. Technology Stack & Rationale

To build a platform robust enough to scale, yet simple enough to maintain without massive operational overhead, we intentionally avoided the complexity of traditional JavaScript-heavy Single Page Applications (SPAs). 

### The Stack:
- **Backend & Database:** [Go (Golang)](https://go.dev/) + [PocketBase](https://pocketbase.io/)
- **Templating:** Go's native `html/template` integrated with a custom markdown compiler ([Goldmark](https://github.com/yuin/goldmark)).
- **Styling:** [Tailwind CSS v3](https://tailwindcss.com/) (using the standalone CLI, avoiding Node.js runtime bloat)
- **Frontend Interactivity:** [HTMX](https://htmx.org/) (hypermedia-driven reactive updates) + [Alpine.js](https://alpinejs.dev/) (declarative client state for dropdowns, modals, and themes)
- **Design Language:** Dynamic Light (Artisan Parchment) & Dark (Modern Midnight Workshop) themes, Art Nouveau whiplash curves, Memphis squiggles, and calligraphic figure styling.

---

## 3. Instant 3D Slicer & Quotation Engine

Jiddat3D features an integrated, server-side instant slicing engine that gives clients immediate print time and rupee quotations without manual CAD slicing:

### Supported Formats
- **.3MF** (3D Manufacturing Format - ZIP + XML mesh parsing)
- **.OBJ** (Wavefront OBJ geometric vertices & polygonal faces)
- **.STL** (Binary & ASCII stereolithography facets)

### Preemptive / Speculative Execution
Using HTMX (`hx-post="/api/instant-slice"`), the moment a user selects or drops a 3D model, the file immediately uploads and slices on the server before they even click start. The backend geometry engine calculates discrete signed tetrahedral volume and surface area in microseconds.

### Slicing Calculations & LayerInvoice Formulation
Slicing logic directly integrates the accounting formulas from [BasilSaeedBari/LayerInvoice](https://github.com/BasilSaeedBari/LayerInvoice):
$$\text{Raw Filament Cost} = \text{Grams} \times \text{Cost per Gram}$$
$$\text{Machine Time Cost} = \text{Hours} \times \text{Hourly Rate (Rs 250/hr)}$$
$$\text{Base Cost} = \text{Filament} + \text{Machine Time} + \text{Labor} + \text{Electricity}$$
$$\text{Total Production Cost} = \text{Base Cost} + \text{Failure Allowance (5\%)}$$
$$\text{Final Quote} = \text{Total Production Cost} + \text{Profit Margin}$$

Once quoted, customers can lock in their order via WhatsApp (`0300-1318112`) with model dimensions, print duration, and price prefilled for direct invoice confirmation.

### Admin Panel Slicer Controls (`/_/`)
Administrators can log into the PocketBase Admin UI (`http://127.0.0.1:8090/_/`) and configure the `slicer_settings` and `filament_types` collections:
- **`time_multiplier`:** Default `1.5` (internal safety buffer to account for UPS transitions, warmup, and guaranteed print turnaround).
- **`weight_multiplier`:** Default `1.5` (internal production buffer for supports, brims, wall perimeters, and purge calibration).
- **`machine_hourly_rate`:** Default `Rs 250/hr`.
- **`default_filament_price`:** Default `Rs 8/g` for White Basic PLA.
- **`profile_file`:** Upload custom `.ini`, `.json`, or `.curaprofile` slicing configurations.
- **Filament Library:** Add or update material pricing (White Basic PLA, Tough PLA+, High-Temp PETG, TPU 95A).

---

## 4. Header & Navigation Structure

The navigation header is organized into clean, responsive dropdowns with Art Nouveau geometry and mobile slide-out navigation:
1. **Hardware & Services (Dropdown):**
   - **Our Specialized Printers:** In-house CoreXY & Cartesian FDM/FFM lineup (`/printers`)
   - **Instant Slicer & Quote:** Preemptive 3D model upload & live pricing (`/instant-quote`)
   - **On-Demand Print Service:** Custom mechanical prototyping & FYP clinic (`/print-service`)
   - **Filament & Spares:** Climate-tested spools (in plan) & stock spares (`/filament`)
2. **Labs & Learning (Dropdown):**
   - **Schools & University Labs:** Turn-key campus makerspaces & PPRA tenders (`/schools`)
   - **Hands-On Workshops:** Hardware assembly & Klipper bootcamps (`/workshops`)
   - **Pakistan 3D Printing Handbook:** Knowledgebase & UPS sizing (`/learn`)
3. **Direct Navigation:**
   - **Gallery:** Authentic lab portfolio with weights, print hours & rupee costs (`/gallery`)
   - **Warranty & Claims:** 6-Month Working Warranty & STEAM Education terms (`/support`)
   - **About Lab:** Engineering mission & open-source sovereignty (`/about`)
4. **Header Utilities:**
   - Dark / Light Theme Toggle (Sun / Moon)
   - 1-Minute Machine Finder Quiz (`/find-my-printer`)
   - Direct WhatsApp Business Hotline (`0300-1318112`)

---

## 5. Local Development Guide

### Prerequisites
1. **Go:** Install [Go 1.22+](https://go.dev/dl/). Verify via `go version`.
2. **Tailwind CLI:** Standalone Linux x64 binary is bundled in the repo (`./tailwindcss`).

### Running Locally
```bash
# Clone the repository
git clone https://github.com/BasilSaeedBari/Jiddat3D.git
cd Jiddat3D

# Compile Tailwind CSS bundle
./tailwindcss -i ui/static/css/input.css -o ui/static/css/output.css --minify

# Start using the Go source
go run ./cmd/jiddat serve --http="127.0.0.1:8090"

# Alternatively, run the pre-compiled binary directly:
./jiddat3d serve --http="127.0.0.1:8090"
```

- **Website Front-end:** `http://127.0.0.1:8090`
- **PocketBase Admin UI:** `http://127.0.0.1:8090/_/`

---

## 6. Deployment via Docker & Compose

The repository includes a multi-stage `Dockerfile` and `compose.yaml` for production deployments.

### Using Docker Compose
```bash
# Build and start the container in the background
docker compose up -d --build
```

### `compose.yaml` Configuration
```yaml
version: '3.8'

services:
  jiddat3d:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: jiddat3d_app
    restart: unless-stopped
    ports:
      - "8090:8090"
    volumes:
      # Persists SQLite database, slicer settings, and uploaded files
      - ./pb_data:/app/pb_data
    environment:
      - ENV=production
      - TZ=Asia/Karachi
```

### Manual Docker Build & Run
```bash
# Build container image
docker build -t jiddat3d:latest .

# Run container with persistent data volume
docker run -d \
  --name jiddat3d_web \
  -p 8090:8090 \
  -v $(pwd)/pb_data:/app/pb_data \
  --restart unless-stopped \
  jiddat3d:latest
```

---

## 7. Community & Inquiries

- **Headquarters:** DHA Phase 1, Karachi, Sindh, Pakistan
- **WhatsApp Business:** [+92 300 1318112](https://wa.me/923001318112) (`0300-1318112`)
- **Email:** [info@jiddat3d.com](mailto:info@jiddat3d.com)
- **WhatsApp Community:** [https://chat.whatsapp.com/Fgvte09elKICDNzlPVuJTe](https://chat.whatsapp.com/Fgvte09elKICDNzlPVuJTe)
- **Discord:** [https://discord.gg/rvsaqF2Q7p](https://discord.gg/rvsaqF2Q7p)

---

## 8. License

All source code, templates, and configurations within this repository are released under the **[GNU Affero General Public License v3.0 (AGPL-3.0)](https://www.gnu.org/licenses/agpl-3.0.en.html)**. We believe in sovereign, transparent, open-source engineering for Pakistan.

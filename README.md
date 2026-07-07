# Jiddat3D - The Website & Digital Hub

> **A Note on Scope:** This repository contains the source code, architecture, and deployment mechanisms exclusively for the **Jiddat3D Website**. It serves as the digital front-door and content management system for the initiative. The hardware CAD models, firmware configurations, and BOMs for our physical machines are hosted in separate, dedicated repositories.

---

## 1. The Rhetoric: Ideology & Origin

### The Maker Movement and RepRap
In 2005, the RepRap project—an acronym for *Replicating Rapid-prototyper*—ignited a revolution in digital fabrication. By open-sourcing hardware designs and creating a machine capable of printing its own components, it shattered the corporate monopoly on 3D printing. It proved that a collective of passionate people freely sharing knowledge can out-innovate massive corporate entities. It democratized manufacturing, bringing it from corporate labs straight into the hands of hackers, tinkerers, and students around the world.

### Where Does Pakistan Fit In?
Historically, the region has been a cradle of invention—from the urban planning and metallurgy of the Indus Valley to the pioneering of optics and algebra during the Islamic Golden Age. Innovation is in our blood. Yet, in the modern era, Pakistan largely finds itself positioned as a consumer rather than a creator, often lagging 10-20 years behind the cutting edge of manufacturing technology. To much of the world, we are in the era of AI and advanced robotics; to many in Pakistan, basic digital fabrication is still a novelty.

### The Genesis of Jiddat3D
Jiddat3D was born from a singular vision: **to restore the culture of creation in Pakistan.** Buying closed-source, imported appliances reinforces a consumer mindset. When those machines break, users must wait for expensive, imported spare parts. 

By building our own open-source hardware, we understand every nut, bolt, and stepper motor. We become creators. Jiddat3D is not merely about selling 3D printers; it is about providing the gateway to engineering literacy. By bringing affordable, repairable, and locally maintainable manufacturing tools to the Pakistani ecosystem, we aim to help the next generation of students bypass decades of missed industrial development and reclaim their heritage as inventors.

---

## 2. Technology Stack & Rationale

To build a platform robust enough to scale, yet simple enough to maintain without massive operational overhead, we intentionally avoided the complexity of traditional Javascript-heavy Single Page Applications (SPAs). 

### The Stack:
- **Backend & Database:** [Go (Golang)](https://go.dev/) + [PocketBase](https://pocketbase.io/) (v0.23+)
- **Templating:** Go's native `html/template` integrated with a custom markdown compiler ([Goldmark](https://github.com/yuin/goldmark)).
- **Styling:** [Tailwind CSS](https://tailwindcss.com/) (using the Standalone CLI)
- **Frontend Interactivity:** [HTMX](https://htmx.org/) (for hypermedia-driven interactions) + [Alpine.js](https://alpinejs.dev/) (for lightweight client-side state)

### Why This Stack?
1. **Zero Node.js Dependency:** By utilizing the standalone Tailwind CLI executable, the entire build process is decoupled from NPM/Node.js, drastically reducing the attack surface and `node_modules` bloat.
2. **Single Binary Deployment:** PocketBase embeds a high-performance SQLite database directly into a Go application. This means the entire backend, API, Auth system, database, and HTML rendering engine compile down to a **single static executable binary**. No separate database server (like Postgres or MySQL) is required.
3. **Performance & Simplicity:** Go provides incredible concurrency and execution speed. Rendering HTML server-side ensures immediate first-paint times and perfect SEO out of the box, without requiring complex hydration techniques.
4. **Declarative UI:** HTMX handles complex form submissions (like the newsletter and contact forms) directly in HTML attributes, while Alpine.js manages simple state (like mobile menu toggling) without requiring a bundler.

---

## 3. Architecture & Request Lifecycle

The application operates as an embedded PocketBase instance, intercepting specific web requests to serve compiled HTML, while allowing PocketBase to handle API and Admin UI requests normally.

### High-Level Flow
1. **Bootstrapping (`cmd/jiddat/main.go`):** The application initializes a PocketBase instance.
2. **Migrations (`pb_migrations/`):** Upon startup, PocketBase checks the `_migrations` table. If new Go migrations exist, they are executed. This sets up the exact database schema and seeds initial content automatically.
3. **Routing (`internal/routes/public.go`):** 
   - Intercepts `/api/newsletter/subscribe` and `/api/contact/submit` to handle custom form logic.
   - Provides a catch-all router `/{path...}`. 
   - It ignores `/api/` and `/_/` (admin) paths.
   - It checks if a requested path maps to a static file in `pb_public`.
   - If not, it falls back to dynamic template rendering.
4. **Rendering (`internal/compiler/render.go`):** 
   - Dynamically reads the `.html` files in `ui/templates/pages/` and wraps them in `ui/templates/layouts/base.html`.
   - Parses Markdown content retrieved from the database into sanitized HTML.
   - Injects global `site_settings` and specific collection records into a `TemplateData` struct for the views to consume.

---

## 4. Database Schema

The SQLite database is managed via PocketBase and contains the following core collections. (These are automatically generated via the `pb_migrations` upon first boot).

### 1. `site_settings` (Single Record Collection)
Stores global configuration variables for the website to allow non-technical administrators to update links without touching code.
- `contact_email` (Email)
- `whatsapp_business_number` (Text)
- `discord_url` (URL)
- `facebook_url` (URL)
- `instagram_url` (URL)
- `youtube_url` (URL)
- `hero_headline` (Text)
- `hero_subheadline` (Text)
- Various `impact_stat_*` fields

### 2. `content` (CMS Collection)
A unified table that stores dynamic pages.
- `type` (Select: `blog`, `project`, `learn`)
- `title` (Text)
- `slug` (Text - Unique)
- `excerpt` (Text)
- `body` (Text - Markdown formatted)
- `hero_image` (File)
- `published` (Bool)
- `published_at` (Date)
- `author_name` (Text)

### 3. `subscribers` (Mailing List)
- `email` (Email - Unique)
- `active` (Bool)
- `source` (Text)

### 4. `contact_submissions` (Inquiries)
- `name` (Text)
- `contact_method` (Text)
- `subject` (Text)
- `message` (Text)
- `status` (Select: `new`, `read`, `replied`)

---

## 5. Local Development & Replication Guide

To replicate and run this project on your local machine, follow these steps meticulously.

### Prerequisites
1. **Install Go:** Download and install [Go 1.22+](https://go.dev/dl/). Verify via `go version`.
2. **Install Tailwind CLI:** Download the standalone Tailwind CSS CLI executable for your OS. Rename the executable to `tailwindcss` (or `tailwindcss.exe` on Windows) and place it in the root directory of the project.

### Step 1: Clone the Repository
```bash
git clone https://github.com/BasilSaeedBari/Jiddat3D.git
cd Jiddat3D
```

### Step 2: Run the Development Server
We have provided a batch script for Windows users that handles downloading Tailwind (if missing), compiling the CSS, and booting the Go server.
```cmd
.\dev.bat
```
*Note: For Unix/Mac users, you can manually run `tailwindcss -i ./ui/static/css/input.css -o ./ui/static/css/index.css --watch` in one terminal, and `go run ./cmd/jiddat serve` in another.*

### Step 3: Access the Application
- **Website Front-end:** Navigate to `http://127.0.0.1:8090`
- **PocketBase Admin UI:** Navigate to `http://127.0.0.1:8090/_/`
  - *Default Admin Login is configured upon first initialization. If you run it for the first time, it will prompt you to create an admin account.*

### Step 4: Modifying Content & Templates
- **HTML Templates:** Edit files in `ui/templates/`. The server reads these dynamically on every request, so you simply need to refresh your browser.
- **CSS:** Edit `ui/static/css/input.css` or add Tailwind classes in the HTML. The `dev.bat` script runs Tailwind in `--watch` mode, automatically recompiling `index.css`.
- **Backend Logic:** Edit `internal/routes/public.go`. You will need to stop (`Ctrl+C`) and restart `dev.bat` for Go code changes to take effect.

---

## 6. Build for Production

To deploy the application to a production server, you must compile it into a standalone binary.

```cmd
.\build.bat
```
This script will:
1. Minify and build the production CSS via Tailwind.
2. Run `go build` to generate a compiled `jiddat.exe` (or `jiddat` binary).

You can then upload this single binary to your server (alongside the `ui/templates`, `ui/static`, and `pb_public` folders) and execute it:
```bash
./jiddat serve --http="0.0.0.0:80"
```

---

## 7. Community & Support

Jiddat3D thrives on its community. If you need help with a build, firmware modding, or just want to share your creations, join us:

- **Discord Server:** [https://discord.gg/rvsaqF2Q7p](https://discord.gg/rvsaqF2Q7p) (Deep dives, firmware help, and modding)
- **WhatsApp Community:** [https://chat.whatsapp.com/Fgvte09elKICDNzlPVuJTe](https://chat.whatsapp.com/Fgvte09elKICDNzlPVuJTe) (Join the daily conversation)
- **Direct WhatsApp Contact:** [+923001318112](https://wa.me/923001318112)

> *"By embracing open-source hardware, we can bypass the decades of industrial development we missed. We can manufacture our own tools, repair our own machines, and build our own solutions. Jiddat3D isn't just a 3D printer; it is a gateway to engineering literacy."*

# TECHSTACK.md — Jiddat3D Technology & Architecture Specification

> **Audience:** An autonomous coding agent with zero prior context on this project.
> **Purpose:** Define, with no ambiguity, every technology, library, folder, build step, and deployment detail required to build the Jiddat3D website.
> **Golden Rule:** No Node.js, no npm, no npx, no JavaScript runtime anywhere in the backend or build pipeline. Go is the only compiled language. JavaScript is permitted ONLY as plain `<script>` files shipped straight to the browser — never compiled, bundled, transpiled, or installed via a package manager.

---

## 1. Non-Negotiable Constraints

These rules override any conflicting suggestion found anywhere else, including in reference documents.

1. **No frontend build step.** There is no `npm install`, no `npx`, no Webpack, no Vite, no esbuild, no PostCSS pipeline that requires Node. The only "build" step permitted is `go build` and the standalone Tailwind CLI binary (which is a self-contained executable, not an npm package).
2. **Go is the backend language.** All server logic, routing, templating, markdown parsing, caching, email, and admin logic is written in Go.
3. **JavaScript is frontend-only, and minimal.** Alpine.js and HTMX are shipped as static `<script src>` files (downloaded once, vendored locally, never fetched from a CDN at runtime by default — see §7). Any custom JS is hand-written vanilla JS, no frameworks, no transpilation.
4. **Caching is mandatory, not optional.** Every content page (Learn, Blog, Projects, Products) must be pre-rendered to static HTML on save and served directly from disk. The database and Go templating engine must never be touched on a normal visitor's read request.
5. **Total page weight budget:** First meaningful paint under 300KB total transferred (HTML + CSS + JS + above-the-fold image), fully interactive under 1MB, on a simulated Fast 3G connection. This budget is a hard design constraint that every later decision must respect.
6. **Everything must be backup-able as a single file or a single folder copy.** No exotic managed services required to run this. It must run entirely on one small VPS.

---

## 2. The Stack At a Glance

| Layer | Technology | Why |
|---|---|---|
| Language / Runtime | **Go** (latest stable, 1.22+) | Single static binary, no runtime dependency, extremely low memory footprint, fast compilation, excellent stdlib for HTTP/templating. |
| Backend Framework / CMS | **PocketBase** (used as a Go framework, not the prebuilt binary) | Gives an embedded SQLite database, a full Admin UI, authentication, file storage, REST API, and a Go hooks/event system — all in a single importable Go package (`github.com/pocketbase/pocketbase`). |
| Database | **SQLite** (embedded inside PocketBase) | Zero external DB server, zero configuration, backs up by copying one `.db` file, more than sufficient for the read/write volume of this site. |
| Markdown Parsing | **goldmark** (Go library, `github.com/yuin/goldmark`) | Native Go markdown-to-HTML converter, extensible (tables, footnotes, syntax highlighting via `goldmark-highlighting`), no external process required. |
| LaTeX Support (optional content type) | **KaTeX** rendered client-side only for pages that opt in, OR pre-render via a Go-callable helper — see §8 | LaTeX is rare-use; do not add server-side LaTeX rendering complexity unless content demands it. Default to Markdown + KaTeX for math notation. |
| HTML Templating | **Go `html/template`** | Native, safe (auto-escaping), zero dependencies, composes with layouts/partials natively. |
| Frontend Interactivity (structural) | **HTMX** (single vendored `.js` file) | Enables server-driven partial page updates (forms, filters, pagination) without a JS framework or virtual DOM. |
| Frontend Interactivity (local state) | **Alpine.js** (single vendored `.js` file) | Handles small UI states: mobile nav toggle, image lightbox, accordion, tabs, dropdown menus — declared directly in HTML attributes. |
| Styling | **Tailwind CSS via the official Standalone CLI binary** | Produces a single compiled, purged, minified CSS file. The standalone CLI is a compiled Go/Rust-adjacent binary Tailwind ships specifically to avoid requiring Node — this satisfies the "no npm" rule while keeping the utility-class workflow. |
| Animations | Hand-authored **CSS transitions/keyframes** + a small **vanilla JS IntersectionObserver utility** for scroll-reveal + Alpine `x-transition` for component-level motion | No animation library needed; native CSS + IntersectionObserver covers 95% of needs at near-zero KB cost. |
| Offline / Asset Caching | **Service Worker** (vanilla JS, Cache API) | Caches static assets and previously-viewed pages so repeat visits and flaky-connection revisits are instant; enables background video buffering. |
| Web Server / Reverse Proxy / TLS | **Caddy** | Automatic HTTPS via Let's Encrypt, serves static cached HTML directly from disk with zero backend involvement, simple `Caddyfile` config, reverse-proxies dynamic API routes to the Go binary. |
| Image Handling | **WebP** as primary format, **AVIF** as progressive enhancement where supported, generated at multiple responsive widths using Go's `image` stdlib + `golang.org/x/image/webp` or a build-time Go CLI tool (`cwebp` invoked via Go's `os/exec` at upload time, not at request time) | Aggressive compression is mandatory for the target audience's bandwidth constraints. |
| Video Handling | **HTML5 `<video>`** with `preload="none"`, poster frames, adaptive quality via Range-request byte-serving (native HTTP support in Caddy/Go) + Service Worker caching | No paid video CDN required; range requests + SW caching are sufficient at this scale. |
| Email / SMTP | **PocketBase's built-in `mailer` package**, configured for **Mailtrap** SMTP via the Admin UI Settings panel | Zero custom code needed for transactional email; custom Go code only for the bulk newsletter sender. |
| Containerization | **Docker** + **Docker Compose** | Two services: `caddy` and `jiddat` (the Go+PocketBase binary). Reproducible deployment to any VPS. |
| Deployment Target | Any Linux VPS (e.g. Hostinger, DigitalOcean, Contabo) — 1 vCPU / 1GB RAM is sufficient | Confirms the "nothing bloated" requirement end-to-end, including the hosting bill. |
| Version Control | **Git** | Standard. `.gitignore` must exclude `pb_data/`, compiled binaries, and any `.env` secrets. |

---

## 3. Why This Combination (Rationale for the Agent)

If you (the building agent) are tempted to substitute React, Next.js, Vue, Node, or any bundler-based system: **do not**. The entire point of this stack is:

- **A single Go binary** that embeds the database, the admin panel, the API, and the markdown compiler. This binary is the entire backend.
- **Zero JavaScript runtime dependency** anywhere in the toolchain that produces the final artifact. A visitor's browser is the only JS runtime involved, and only for small, hand-authored, unbundled scripts.
- **Pre-rendering over server-side rendering per request.** Every content page is generated once (on save) and served as static HTML forever after — this is the single biggest lever for speed on Pakistan's average mobile network conditions (3G/4G with variable latency, frequent packet loss, higher costs per MB).

This is a **static-first, event-driven regeneration** architecture, not a traditional dynamic web app.

---

## 4. Directory Structure

```text
Jiddat3D/
│
├── README.md                        # High-level project overview, links to /docs
├── go.mod
├── go.sum
├── Dockerfile
├── compose.yaml
├── Caddyfile
├── .gitignore
├── .env.example                     # Template for secrets (SMTP creds, admin seed email, etc.)
│
├── docs/                            # This documentation set lives here in the repo too
│   ├── agent.md
│   ├── design.md
│   ├── techstack.md
│   └── content.md
│
├── cmd/
│   └── jiddat/
│       └── main.go                  # Entry point: boots PocketBase, registers hooks, registers routes
│
├── internal/
│   ├── hooks/
│   │   ├── content_hooks.go         # OnRecordAfterCreate/Update/Delete for "content" collection
│   │   ├── product_hooks.go         # Same for "products" collection
│   │   └── subscriber_hooks.go      # Newsletter signup side-effects
│   │
│   ├── compiler/
│   │   ├── markdown.go              # goldmark setup, extensions, sanitization
│   │   ├── render.go                # Loads Go templates, injects compiled HTML, writes to disk
│   │   └── images.go                # Image ingestion: resize, convert to WebP/AVIF, generate srcset variants
│   │
│   ├── mailer/
│   │   └── newsletter.go            # Custom bulk-send endpoint logic using PocketBase's mailer
│   │
│   ├── routes/
│   │   ├── public.go                # Non-cached dynamic routes (forms, search/filter, newsletter signup)
│   │   ├── admin_extra.go           # Any custom admin-only API routes beyond PocketBase defaults
│   │   └── sitemap.go               # Auto-generates sitemap.xml and robots.txt from live collections
│   │
│   └── cache/
│       └── invalidate.go            # Deletes/rebuilds stale cached HTML files
│
├── ui/
│   ├── templates/
│   │   ├── layouts/
│   │   │   ├── base.html            # <html>, <head>, preloader, header, footer, closing tags
│   │   │   └── admin_wrapper.html   # (Only if PocketBase admin UI is themed/wrapped — usually skip)
│   │   ├── partials/
│   │   │   ├── nav.html
│   │   │   ├── footer.html
│   │   │   ├── card_product.html
│   │   │   ├── card_project.html
│   │   │   ├── card_blog.html
│   │   │   ├── newsletter_form.html
│   │   │   └── contact_form.html
│   │   └── pages/
│   │       ├── home.html
│   │       ├── products_index.html
│   │       ├── product_detail.html
│   │       ├── learn_index.html
│   │       ├── learn_detail.html
│   │       ├── projects_index.html
│   │       ├── project_detail.html
│   │       ├── blog_index.html
│   │       ├── blog_detail.html
│   │       ├── about.html
│   │       ├── community.html
│   │       ├── contact.html
│   │       └── 404.html
│   │
│   ├── static/
│   │   ├── css/
│   │   │   ├── input.css            # Tailwind directives + custom @layer components
│   │   │   └── output.css           # COMPILED, MINIFIED — generated by Tailwind standalone CLI, gitignored
│   │   ├── js/
│   │   │   ├── htmx.min.js          # Vendored, not CDN-fetched, exact version pinned
│   │   │   ├── alpine.min.js        # Vendored, exact version pinned
│   │   │   ├── preloader.js         # Hand-written vanilla JS
│   │   │   ├── scroll-reveal.js     # Hand-written IntersectionObserver-based animation trigger
│   │   │   ├── lightbox.js          # Hand-written, used with Alpine for gallery/project images
│   │   │   └── sw-register.js       # Registers the service worker
│   │   ├── sw.js                    # The Service Worker itself (must live at domain root, see §9)
│   │   ├── fonts/                   # Self-hosted font files (WOFF2 only, subset)
│   │   └── img/
│   │       ├── logo/                # Logo assets in multiple sizes/formats
│   │       └── icons/               # Favicon, PWA icons, social share images
│   │
│   └── tailwind.config.js           # Config file for the standalone CLI (this is NOT an npm project file)
│
├── pb_public/
│   ├── cached/                      # GENERATED. Caddy serves directly from here. Mirrors URL structure.
│   │   ├── learn/{slug}.html
│   │   ├── blog/{slug}.html
│   │   ├── projects/{slug}.html
│   │   └── products/{slug}.html
│   ├── static/                      # Symlink or copy target of ui/static at build time
│   └── uploads/                     # PocketBase-managed file storage (images, PDFs, gallery photos)
│
├── pb_data/                         # SQLite database + PocketBase internal files. GIT-IGNORED. Backed up separately.
│
├── scripts/
│   ├── build.sh                     # Runs Tailwind CLI, copies static assets, go build
│   ├── dev.sh                       # Runs Tailwind CLI in --watch mode + `go run` with live reload (see §6)
│   └── backup.sh                    # Cron-friendly script: copies pb_data/data.db + pb_public/uploads to backup location
│
└── tools/
    └── tailwindcss                  # The downloaded standalone Tailwind CLI binary (gitignored, fetched by build.sh)
```

---

## 5. The Go Backend In Detail

### 5.1 PocketBase as a Library, Not a Binary

Do **not** download the prebuilt PocketBase executable from GitHub releases. Instead:

```bash
go get github.com/pocketbase/pocketbase
```

Then in `cmd/jiddat/main.go`, construct a `pocketbase.New()` app instance, register custom hooks and routes on it, and call `app.Start()`. This is what allows custom Go logic (markdown compilation, caching, custom mail routes) to run inside the same binary and process as the CMS itself, sharing the same SQLite connection.

### 5.2 Collections (Database Schema)

Defined in full in `content.md`, but summarized here for the agent's architectural awareness. Collections are created either by:
- Writing a Go migration file under `pb_migrations/` (preferred — version-controlled, reproducible), or
- Manually via the Admin UI on first boot (acceptable for a solo/small team, but must then be documented and exported as a migration snapshot for reproducibility).

**Use Go migrations.** This keeps the schema in version control and lets the agent building this site produce a fully reproducible `go run` → working site with zero manual admin clicking required for schema setup. Only *content* (actual posts, products) is entered manually via the Admin UI afterward.

### 5.3 Event Hooks (The Core Mechanism)

This is the single most important piece of backend logic. Attach to PocketBase's record lifecycle hooks:

```
OnRecordAfterCreateSuccess("content")
OnRecordAfterUpdateSuccess("content")
OnRecordAfterDeleteSuccess("content")
```//(mirror the same three for "products")

**On Create/Update:**
1. Fetch the record's fields (title, slug, body/markdown, type, hero image, metadata, published flag).
2. If `published == false`, skip generation (or generate to a `drafts/` cache path not exposed publicly) and stop.
3. Run the raw markdown body through the `compiler/markdown.go` goldmark pipeline → sanitized HTML fragment.
4. Select the correct Go template based on `type` (Learn / Blog / Project / Product).
5. Execute the template with the fragment + metadata (title, date, hero image URL, tags, SEO description) injected.
6. Write the resulting full HTML page to `pb_public/cached/{collection}/{slug}.html`, creating directories as needed.
7. If the slug changed on update (old slug ≠ new slug), delete the old cached file.
8. Regenerate any index pages that list this content (e.g., `/blog`, `/learn`) if they show excerpts — OR, better, make index pages themselves lightweight enough to query the DB directly with pagination rather than being fully static (see §5.4 for this nuance).

**On Delete:**
1. Delete the corresponding cached `.html` file.
2. Regenerate affected index pages if applicable.

### 5.4 What Gets Statically Cached vs. What Stays Dynamic

Be precise about this distinction — it is easy to over-cache or under-cache:

| Page Type | Strategy |
|---|---|
| Individual Blog / Learn / Project / Product detail pages | **Fully static-cached.** Regenerated only on save. Served by Caddy directly, Go never touched. |
| Home page | **Static-cached**, regenerated whenever featured content changes (a hook on Products/Content marked "featured", or manually re-triggerable from an admin button). |
| Index/listing pages (`/blog`, `/learn`, `/projects`, `/products`) | **Hybrid.** Render a static cached "page 1" for the fast initial load (this is what most visitors see and what needs to be instant). Deeper pagination or filtering (e.g., filter Learn articles by tag) is handled via an HTMX GET request to a lightweight Go route that queries SQLite directly and returns just the HTML fragment for the results grid. This keeps the *first paint* static and instant while allowing dynamic filtering without a full framework. |
| Forms (newsletter signup, contact form) | **Always dynamic.** POST to a Go/PocketBase API route. Never cached. |
| Admin panel | **Always dynamic**, served by PocketBase itself at `/_/`. Not part of the public caching strategy at all. |

### 5.5 Markdown Compiler Configuration

`goldmark` must be configured with these extensions enabled:
- **GFM (GitHub Flavored Markdown)** table support, strikethrough, autolinks.
- **Footnotes** for citation-style content (useful for the Learn section's technical articles).
- **Syntax highlighting** via `github.com/yuin/goldmark-highlighting/v2` for any code blocks (e.g., firmware config snippets).
- **Typographer** extension for smart quotes/dashes (better reading experience).
- **Unsafe HTML disabled** — sanitize output. Since content is admin-authored (trusted), raw HTML passthrough may be enabled, but images/embeds should still be checked for valid, local, or explicitly whitelisted external sources (e.g., YouTube embed for project videos).

### 5.6 LaTeX Support

Do not build a server-side LaTeX-to-image pipeline (heavy, adds bloat, contradicts the lightweight goal). Instead:
- Allow authors to write inline/block math using standard `$...$` / `$$...$$` delimiters inside their Markdown.
- Ship **KaTeX** (self-hosted, not CDN) as an optional per-page CSS/JS include, only loaded on pages whose frontmatter/metadata flags `has_math: true`. This keeps math rendering entirely opt-in and out of the default page weight budget for the 99% of pages that don't need it.

---

## 6. Local Development Workflow (No Node, Ever)

1. Install Go (1.22+) locally.
2. Download the **Tailwind CSS standalone CLI binary** matching the OS/architecture once, place it in `tools/tailwindcss`, mark executable. This is a direct binary download from Tailwind's GitHub releases — it is not `npm install tailwindcss`, and it requires no `package.json`.
3. `scripts/dev.sh` runs two processes concurrently (via a Go-native process, `air` for Go live-reload is optional but itself Go-based and acceptable, OR simply two terminal tabs):
   - `./tools/tailwindcss -i ui/static/css/input.css -o ui/static/css/output.css --watch`
   - `go run ./cmd/jiddat`
4. Visit `http://localhost:8090` for the public site and `http://localhost:8090/_/` for the PocketBase Admin UI.
5. First boot: PocketBase prompts for creation of a superuser admin account (email/password) directly in the terminal or via the Admin UI's first-run screen.

**No `package.json`, `node_modules`, or lockfile of any kind should ever exist in this repository.**

---

## 7. Vendoring Third-Party JS/CSS (No CDN Dependency)

Because the target audience has unreliable connections, and because relying on third-party CDNs introduces an external point of failure and an unnecessary DNS/TLS handshake:

- Download **HTMX** (single minified `.js` file, pin a specific version) and **Alpine.js** (single minified `.js` file, pin a specific version) once, place them in `ui/static/js/`, and serve them from the same origin as everything else.
- Same for any web font files — self-host WOFF2, do not link to Google Fonts' hosted CSS (that would be an external round-trip).
- This also means the site works correctly even if a third-party CDN is blocked or slow in a given region — a real consideration for Pakistan.

---

## 8. Image & Media Pipeline

Handled at **upload time**, not request time, to keep runtime overhead at zero:

1. When an image is uploaded via the Admin UI (into a Content or Product record's file field), a Go hook (`internal/compiler/images.go`) intercepts it.
2. Generate 3 responsive widths (e.g., 480px, 960px, 1600px) plus the original.
3. Convert each to **WebP** (broad support, big size win over JPEG/PNG) via `os/exec` calling `cwebp` (install this system dependency in the Dockerfile), or a pure-Go WebP encoder if avoiding the external binary is preferred for portability.
4. Store all variants in PocketBase's file storage; reference them via a `srcset` in the rendered `<img>` tag so the browser picks the appropriate size.
5. Always set explicit `width`/`height` attributes on `<img>` tags to prevent layout shift (a Core Web Vitals requirement and a real UX issue on slow connections where images pop in late).
6. Lazy-load every image below the fold using the native `loading="lazy"` HTML attribute — zero JS cost.

**Video** follows a similar philosophy:
- Always specify a lightweight poster image (`poster="..."`) and `preload="none"` so a video never auto-downloads.
- Serve via native HTTP Range requests (Go's `http.ServeContent` supports this out of the box) so users can seek without downloading the whole file.
- Cache played segments via the Service Worker for instant replay/revisits.

---

## 9. Service Worker & Preloader Behavior

### 9.1 Preloader (`preloader.js`)
- A tiny inline `<script>` in `<head>`, or a very small external file loaded blocking, that:
  1. Immediately shows a full-viewport branded loading screen (see `design.md` for exact visual spec) the instant the HTML starts parsing.
  2. Listens for `window.addEventListener('load', ...)` (not just `DOMContentLoaded`, to ensure images are ready too, but with a maximum timeout fallback of ~1.5s so the loader never blocks perceived usability on a slow image).
  3. Fades the loader out via a CSS transition class toggle.
- This guarantees the visitor sees *something* branded and alive within milliseconds, even before CSS/fonts fully load — critical for perceived performance on high-latency connections.

### 9.2 Service Worker (`sw.js`)
- Registered from `sw-register.js` on every page.
- **Cache strategy:** Cache-first for static assets (CSS, JS, fonts, logo) with a versioned cache name (bump version string on deploy to bust stale caches). Network-first with cache fallback for HTML pages (so content updates are seen quickly, but the site still works offline/on flaky connections with last-known content).
- **Video/image caching:** Cache visited media assets opportunistically so a user revisiting a Learn article or Project page on a poor connection gets instant repeat loads.
- Do not cache POST requests, form submissions, or the Admin UI (`/_/*`) or API (`/api/*`) routes — those must always hit the network.

---

## 10. Caddy Configuration Outline

`Caddyfile` responsibilities:
1. Automatic HTTPS for the production domain (Let's Encrypt).
2. Serve `pb_public/cached/*` and `pb_public/static/*` directly as static files with long `Cache-Control` headers (e.g., `max-age=31536000, immutable` for hashed/static assets; shorter for HTML pages to allow updates to propagate).
3. Reverse-proxy anything under `/api/*` and `/_/*` (PocketBase's REST API and Admin UI) to the Go binary's internal port (e.g., `127.0.0.1:8090`).
4. Reverse-proxy any custom dynamic routes (newsletter signup, contact form, filter/search fragments) to the same Go binary.
5. Fallback 404 handling pointing to the custom `404.html` cached page.
6. GZIP/Brotli compression enabled globally (Caddy supports this natively with a single directive) — this alone can cut transferred HTML/CSS/JS size by 60-80%, a major win for the target audience.

---

## 11. Deployment & Backups

### 11.1 Docker Compose Services
- **`caddy`**: official Caddy image, mounts the `Caddyfile`, exposes 80/443, mounts a volume for `pb_public/` (read access to serve cached files and static assets) and a volume for Caddy's own data/certs.
- **`jiddat`**: built from the project `Dockerfile` (multi-stage: build the Go binary in a `golang:1.22` builder stage, copy into a minimal `distroless` or `alpine` final image), exposes an internal port only (not published to the host directly — Caddy is the only public entry point), mounts volumes for `pb_data/` and `pb_public/`.

### 11.2 Backup Strategy
- A cron job (`scripts/backup.sh`) runs daily on the host:
  1. Copies `pb_data/data.db` (and its WAL/SHM files if present) to a timestamped backup folder.
  2. Copies `pb_public/uploads/` (all uploaded images/media) alongside it.
  3. Optionally syncs the backup folder to off-site storage (rsync to another server, or an S3-compatible bucket) — keep this simple and document the exact command used, do not introduce a complex backup service.
- Because SQLite is a single file, disaster recovery is: stop the container, replace `data.db`, restart. This must be explicitly documented in the deployment runbook.

---

## 12. Explicit Anti-Patterns (Things the Agent Must Not Do)

- Do **not** introduce React, Vue, Svelte, Next.js, or any SPA framework.
- Do **not** run `npm`, `npx`, `yarn`, or `pnpm` at any point, for any reason, including "just for Tailwind."
- Do **not** fetch fonts, icons, HTMX, or Alpine from a public CDN at runtime — vendor everything locally.
- Do **not** render pages from the database on every request for content that rarely changes — that defeats the entire caching architecture.
- Do **not** use a heavyweight ORM; PocketBase's built-in query builder / raw SQL via its Go DAO is sufficient.
- Do **not** add authentication/user accounts for the public site (no login, no cart, no checkout) — this is explicitly a "contact us / WhatsApp to order" Version 1 site, per the project brief. Only the Admin UI requires authentication (and PocketBase provides this natively).
- Do **not** over-engineer the LaTeX pipeline; Markdown + optional client-side KaTeX is sufficient.

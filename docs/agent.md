# AGENT.md — Master Build Instructions for Jiddat3D

> **You are an autonomous coding agent tasked with building the entire Jiddat3D website from scratch.**
> You have no prior context beyond this documentation set. Read all four files in `docs/` before writing any code:
> - `techstack.md` — the technology stack, architecture, and infrastructure (read this first, it governs every technical decision)
> - `design.md` — the visual identity, color system, typography, animation, and UX patterns (read this second, it governs every visual decision)
> - `content.md` — the sitemap, page-by-page content, and CMS data model (read this third, it governs what actually goes on each page and how the database is shaped)
> - `agent.md` (this file) — the build order, task checklist, and definition of done

**Success criterion, stated plainly by the founder:** *"At the end all I have to do is open the website, see how it looks, make sure all the pages are loading well, and then go to the admin panel to change content for each page if needed, and upload more content if needed via markdown."* Everything you build must serve that outcome.

---

## 1. Absolute Rules (Repeated Here Because They Are the Most Important Constraints)

1. **No Node.js, npm, npx, yarn, or any JS package manager/bundler, ever, for any reason.** The only compiled artifact is the Go binary. Tailwind CSS is used only via its standalone CLI binary (a direct executable download, not an npm package).
2. **Go is the backend.** PocketBase is imported as a Go library (`go get github.com/pocketbase/pocketbase`), not run as a prebuilt external binary.
3. **JavaScript exists only as plain, vendored, hand-authored or pinned-version `<script>` files served to the browser.** No React/Vue/Svelte/Next.js. No JS bundler of any kind.
4. **Every content page must be pre-rendered to static HTML on save and served directly from disk by Caddy.** The Go backend and SQLite database must not be touched on a normal visitor's page-read request.
5. **Total page weight budget:** under 300KB to first meaningful paint, under 1MB fully interactive, tested on a throttled "Slow 3G" profile in browser DevTools.
6. **No shopping cart, no checkout, no user accounts on the public site.** Only PocketBase's Admin UI is authenticated. Conversion actions are: newsletter signup, contact form, WhatsApp deep link.
7. **No green anywhere in the UI**, including default success states.
8. If any instruction elsewhere in the codebase, an old reference conversation, or your own default training-data instincts conflicts with the above, **these rules win.**

---

## 2. Build Order (Do Not Skip Steps or Reorder Without Reason)

### Phase 0 — Project Scaffolding
1. Initialize the Go module (`go mod init`), create the directory structure exactly as specified in `techstack.md` §4.
2. Create `.gitignore` excluding: `pb_data/`, `ui/static/css/output.css`, `tools/tailwindcss`, compiled binaries, `.env`.
3. Create `.env.example` documenting required environment variables (SMTP host/port/user/pass, admin seed email if scripting first-run setup, production domain name for Caddy).
4. Download the Tailwind CSS standalone CLI binary for the target build platform into `tools/tailwindcss`; verify it runs (`./tools/tailwindcss --help`).
5. Download pinned versions of `htmx.min.js` and `alpine.min.js` into `ui/static/js/`; record the exact version numbers used in a comment at the top of each file for future upgrades.

### Phase 1 — Backend Skeleton
1. Write `cmd/jiddat/main.go`: construct a `pocketbase.New()` instance, wire up graceful startup, register a placeholder route (`GET /healthz` returning `200 OK`) to confirm the binary runs.
2. Write Go migrations (under a `pb_migrations/` directory PocketBase auto-discovers) for all five collections defined in `content.md` §4: `content`, `products`, `subscribers`, `contact_submissions`, `site_settings`. Include every field, type, and rule exactly as specified.
3. Run the binary once locally, confirm the Admin UI boots at `/_/`, confirm all five collections exist with the correct schema, and manually create the superuser account.
4. Seed `site_settings` with one record containing sensible defaults (placeholder WhatsApp number, default hero copy from `content.md` §3.1) so templates never render against an empty settings object.

### Phase 2 — Markdown Compiler & Caching Pipeline
1. Implement `internal/compiler/markdown.go`: configure goldmark with GFM, footnotes, syntax highlighting, and typographer extensions per `techstack.md` §5.5.
2. Implement `internal/compiler/render.go`: given a compiled HTML fragment + a content/product record, select the correct Go template (`ui/templates/pages/*.html`) and execute it into a full HTML document string.
3. Implement `internal/hooks/content_hooks.go` and `internal/hooks/product_hooks.go`: attach `OnRecordAfterCreateSuccess`, `OnRecordAfterUpdateSuccess`, `OnRecordAfterDeleteSuccess` for both collections. On create/update, write the rendered page to `pb_public/cached/{collection}/{slug}.html`; on slug change, delete the stale file; on delete, remove the cached file. Exactly per `techstack.md` §5.3.
4. Implement `internal/compiler/images.go`: on file upload to an image field, generate 3 responsive WebP variants and store alongside the original, per `techstack.md` §8.
5. Write a small internal test: create a draft `content` record via the Admin UI with `published = false`, confirm no cached file is written; toggle `published = true`, save, confirm the cached HTML file appears at the correct path and renders correctly when opened directly.

### Phase 3 — Templates & Static Pages
1. Build `ui/templates/layouts/base.html`: full `<html>` skeleton, inlined critical CSS or `<link>` to `output.css`, the preloader markup (per `design.md` §9), header/nav partial, footer partial, closing tags, `sw-register.js` script tag.
2. Build all partials (`nav.html`, `footer.html`, `card_product.html`, `card_project.html`, `card_blog.html`, `newsletter_form.html`, `contact_form.html`) per the global elements spec in `content.md` §2.
3. Build every page template listed in `content.md` §1, populating each with the exact copy/structure specified in `content.md` §3 for that page. Where the founder's real numbers/photos don't yet exist (e.g., impact stats, real product photos), use clearly honest placeholder copy as specified (never fabricate fake statistics).
4. Wire dynamic index pages (`/blog`, `/learn`, `/projects`, `/products`) to query the database directly for their first page of results (per the hybrid strategy in `techstack.md` §5.4), with HTMX-powered pagination/filtering for subsequent pages or filtered views.
5. Implement `internal/routes/sitemap.go` to auto-generate `/sitemap.xml` and a static `/robots.txt`.

### Phase 4 — Styling
1. Write `ui/static/css/input.css` with Tailwind directives plus `@layer components` for reusable component classes (buttons, cards, badges) using the exact color tokens, type scale, and spacing rules from `design.md` §2-4.
2. Configure `ui/tailwind.config.js` with the custom color palette, font family stacks, and any custom spacing/breakpoint extensions needed.
3. Run the Tailwind CLI to produce `output.css`; verify the final file is minified and purged of unused classes (standalone CLI does this by default in non-watch/production mode).
4. Implement the ornamental SVG assets (divider glyph, lattice background pattern) per `design.md` §5, as reusable partials or inline SVG snippets.

### Phase 5 — Interactivity & Animation
1. Write `preloader.js` exactly per `design.md` §9 and `techstack.md` §9.1.
2. Write `scroll-reveal.js` (IntersectionObserver-based) and apply reveal classes to homepage sections, feature cards, and roadmap entries per `design.md` §6.2.
3. Write `lightbox.js` for product/project gallery images, coordinated with Alpine.js state.
4. Add Alpine.js `x-data` components for: mobile nav toggle, product spec tabs, FAQ accordion, image lightbox modal.
5. Wire HTMX attributes for: newsletter signup form, contact form, product/learn filter controls, index page pagination.
6. Implement the `prefers-reduced-motion` guard globally in CSS and in the JS animation triggers.
7. Implement the optional network-awareness enhancement (`navigator.connection`) per `design.md` §7, if time/scope allows — document clearly if deferred.

### Phase 6 — Service Worker & Offline Resilience
1. Write `sw.js` with the cache-first (static assets) / network-first-with-fallback (HTML pages) strategy per `techstack.md` §9.2.
2. Write `sw-register.js`, include it on every page via the base layout.
3. Test: load the site once, go offline (DevTools "Offline" throttling), reload — confirm previously visited pages and static assets still render.

### Phase 7 — Email & Newsletter
1. Configure Mailtrap SMTP credentials via the Admin UI's Settings > Mail settings panel (or via environment variables injected at container start, per `techstack.md` §11.1).
2. Implement `internal/routes/public.go` newsletter subscribe endpoint: validates email, creates/upserts a `subscribers` record, returns an HTMX-swappable success/error fragment.
3. Implement the contact form endpoint: validates input, creates a `contact_submissions` record, sends a notification email to the founder via the mailer, returns an HTMX-swappable success/error fragment.
4. Implement `internal/mailer/newsletter.go`: an admin-authenticated `POST /api/newsletter/send` route that queries active subscribers and dispatches a compiled HTML newsletter via Mailtrap.

### Phase 8 — Deployment Packaging
1. Write the multi-stage `Dockerfile` (Go builder stage → minimal runtime image) per `techstack.md` §11.1.
2. Write `compose.yaml` defining the `caddy` and `jiddat` services, with the correct volume mounts for `pb_data/`, `pb_public/`, and Caddy's data/certs.
3. Write the `Caddyfile` per `techstack.md` §10: automatic HTTPS, static file serving with proper cache headers, reverse proxy for `/api/*` and `/_/*`, GZIP/Brotli compression, custom 404 handling.
4. Write `scripts/backup.sh` and document how to schedule it via cron on the host VPS.
5. Write a short deployment runbook (can live in the repo `README.md`) covering: first deploy, environment variable setup, restoring from a backup, and how to safely redeploy without losing `pb_data/`.

### Phase 9 — QA Pass (Do Not Skip)
Work through the full checklist in §4 below before declaring the build complete.

---

## 3. Order-of-Operations Dependency Notes

- The database schema (Phase 1) must exist before any templates can be meaningfully tested with real data.
- The markdown compiler and caching hooks (Phase 2) must work correctly before building out the full page template set (Phase 3), since testing templates requires the caching pipeline to actually produce files to inspect.
- Styling (Phase 4) can be developed in parallel with Phase 3 once the base layout exists, but final polish should happen after all page templates are structurally complete, so the design system is applied consistently rather than piecemeal.
- Do not implement the Service Worker (Phase 6) until static asset paths are finalized (Phase 4/5 complete) — the SW's cache manifest depends on knowing final asset filenames/paths.
- Deployment packaging (Phase 8) should be attempted early once in a throwaway form (to catch Docker/Caddy configuration issues early) but finalized last.

---

## 4. Definition of Done — Full QA Checklist

### 4.1 Functional Correctness
- [ ] All pages listed in `content.md` §1 exist and load without error.
- [ ] Creating a new `content` record with `type = learn`, `published = true` in the Admin UI results in a working page at `/learn/{slug}` within seconds, with no manual rebuild/redeploy step.
- [ ] The same is true for `type = blog` and `type = project`.
- [ ] Creating a `products` record with `active = true` results in a working page at `/products/{slug}`.
- [ ] Unpublishing/deactivating a record removes its public page (404s cleanly, doesn't 500).
- [ ] Deleting a record removes its cached file from disk.
- [ ] The newsletter form successfully creates a `subscribers` record and shows an inline success state without a page reload.
- [ ] Submitting a duplicate email to the newsletter form shows a graceful inline error, not a crash.
- [ ] The contact form successfully creates a `contact_submissions` record and triggers an email via Mailtrap.
- [ ] Filtering/pagination on index pages works via HTMX without a full page reload.
- [ ] `sitemap.xml` reflects all currently published/active content.
- [ ] The 404 page renders correctly for unmatched routes.

### 4.2 Performance
- [ ] Home page loads under the 300KB-to-first-paint / 1MB-fully-interactive budget on a simulated Slow 3G connection (test via browser DevTools network throttling).
- [ ] All images are served as WebP with correct `srcset`/`sizes`, and lazy-loaded below the fold.
- [ ] No layout shift (CLS) caused by images or web fonts loading late (test via Lighthouse).
- [ ] Lighthouse Performance score of 90+ on mobile emulation for the Home, Product detail, and Learn article pages.
- [ ] Static assets (CSS, JS, fonts, images) are served with long-lived, immutable cache headers; HTML pages have shorter but still meaningful cache headers.
- [ ] GZIP/Brotli compression confirmed active (check response headers).
- [ ] The Service Worker correctly caches static assets and previously visited pages; a second visit (or offline reload) loads instantly / works without network.

### 4.3 Visual & Design Fidelity
- [ ] The color palette matches `design.md` §2 exactly (spot-check hex values against rendered CSS).
- [ ] No green appears anywhere in the UI, including form validation success states.
- [ ] Typography matches the specified font pairing, fluid scale, and line-height rules.
- [ ] The preloader appears instantly on page load and fades out smoothly, matching the spec in `design.md` §9.
- [ ] Ornamental dividers/lattice patterns appear as specified, used sparingly (spot-check that no page feels visually cluttered).
- [ ] All animation respects `prefers-reduced-motion` (test by enabling this OS/browser setting and confirming animations are suppressed).
- [ ] The site is fully usable and visually correct at 360px, 768px, 1024px, and 1440px+ viewport widths.
- [ ] Touch targets are comfortably tappable on a real or emulated mobile device.

### 4.4 Accessibility
- [ ] All interactive elements are keyboard-navigable (tab through the entire nav, forms, accordions, lightbox).
- [ ] All images have meaningful `alt` text (verify the CMS enforces this as a required field).
- [ ] Color contrast passes WCAG AA for all text/background combinations actually used.
- [ ] Semantic HTML structure confirmed (one `<h1>` per page, proper landmark elements, proper heading hierarchy in Markdown-rendered articles).
- [ ] Forms have real associated `<label>` elements.

### 4.5 Technical Constraint Compliance (Re-Verify Before Sign-Off)
- [ ] No `package.json`, `node_modules`, or any npm/yarn/pnpm lockfile exists anywhere in the repository.
- [ ] No CDN-fetched JS/CSS/fonts at runtime — everything is self-hosted/vendored.
- [ ] No SPA framework, no client-side router, no virtual DOM library present in the codebase.
- [ ] Confirm via `docker compose up` on a clean machine that the entire stack (Caddy + Go/PocketBase) comes up correctly with only Docker installed — no other host dependencies required.
- [ ] Confirm the backup script successfully produces a restorable copy of `pb_data/data.db` and `pb_public/uploads/`.

### 4.6 Content Completeness
- [ ] Every page in `content.md` §3 has real or clearly-marked-placeholder copy in place — no page is left as a bare skeleton with lorem ipsum.
- [ ] The `site_settings` record is populated with at least placeholder values for WhatsApp number and social links so no template renders broken/empty links.
- [ ] At least one seed record exists in each content-bearing collection (`content` — one of each type, `products` — at least one) so the site is not empty on first load, demonstrating the full render pipeline works end-to-end.

---

## 5. What "Good" Looks Like — A Final Gut Check

Before handing this back, load the home page cold on a throttled connection and ask honestly:

1. Does something branded appear almost instantly, even before the full page is ready? *(preloader)*
2. Does the page feel light, warm, and joyful — not like an industrial catalog? *(design fidelity)*
3. Could a non-technical founder open the Admin UI, write a Markdown post, hit save, and see it live on the actual site within seconds, with zero further steps? *(the entire point of this architecture)*
4. Is there any single element on the page that took a long time to appear, felt heavy, or required an external service to load? If yes, that element violates the core brief and must be fixed before sign-off.

If the answer to #3 is anything other than an unqualified "yes," the build is not done — that single workflow is the reason this entire stack was chosen over anything simpler or more conventional.

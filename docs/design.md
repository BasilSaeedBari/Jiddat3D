# DESIGN.md — Jiddat3D Visual & Interaction Design System

> **Audience:** An autonomous coding agent with zero prior context on this project.
> **Purpose:** Define, with no ambiguity, the visual identity, color system, typography, spacing, animation behavior, responsive rules, and page-level UX patterns for the Jiddat3D website.
> **Brand Truth:** Jiddat3D does not sell 3D printers. It sells the ability to create. Every visual decision should feel joyful, warm, and human — never industrial, corporate, or cold. The visitor should feel "I could build something here," not "I am browsing a machine catalog."

---

## 1. Brand Foundation

- **Name:** Jiddat3D. "Jiddat (جدت)" means innovation, novelty, a deliberate break from the old way of doing things.
- **Logo:** An elegant serif/calligraphic wordmark "Jiddat" rendered in a deep violet-purple with a soft glowing outline against a white background, topped with a small ornamental motif (a stylized geometric flower/dome shape reminiscent of Mughal architectural tilework and jaali (lattice) patterns) and a small decorative divider of diamond/star shapes beneath the wordmark.
- **Tone of the ornament:** This ornamental motif is the single most important design cue in this project. It signals "Pakistani/South Asian craftsmanship and heritage" without relying on clichés (no generic "eastern" fonts, no flag colors, no literal minarets). The agent must echo this motif language — small diamond/star dividers, gentle geometric lattice patterns as subtle background textures, arch-like shapes for card corners or section dividers — throughout the site, used sparingly as accents, never as heavy decoration that adds weight or clutter.
- **Voice:** Confident, warm, a little bit rebellious (against the "we only consume, we don't build" narrative), never salesy or corporate-jargon-heavy.

---

## 2. Color System

### 2.1 Source of Truth
The palette is derived directly from the supplied logo file (`JiddatBrightTrasnparent.png`), which uses a rich violet-purple glow on a white field. This becomes the anchor hue. **Green is explicitly forbidden anywhere in the UI** (including "success" states — use a warm gold/amber or the primary violet instead of the conventional green checkmark).

### 2.2 Primary Palette

| Token | Hex (approx.) | Usage |
|---|---|---|
| `--color-primary` | `#5B2A86` (deep royal violet, matched to the logo's core stroke color) | Logo, primary headings, primary buttons' text/icon on light backgrounds, nav active states. |
| `--color-primary-light` | `#9B5DE5` (bright orchid/violet glow) | Gradients, hover glows, decorative accents, the "glow" effect behind hero elements — mirrors the logo's outer glow. |
| `--color-primary-dark` | `#3A1A57` | Text on light backgrounds needing extra contrast, footer background option. |

### 2.3 Secondary / Accent Palette — "Truck Art & Mughal Tile" Inspired

Rather than reach for generic "flag green," pull from Pakistan's actual vernacular visual culture: truck art, Multani/Mughal tilework, and traditional textile dyes. These read as authentically Pakistani without being a literal flag reference.

| Token | Hex (approx.) | Usage |
|---|---|---|
| `--color-accent-gold` | `#E8A93A` (warm marigold/saffron gold) | Secondary CTA buttons, "featured" badges, star ratings, highlight underlines, success/confirmation states (replaces green). |
| `--color-accent-teal` | `#1E8A8A` (Multani tile teal/turquoise) | Secondary accents, icon backgrounds, tag chips, alternating section backgrounds, links within body copy. |
| `--color-accent-terracotta` | `#D9603B` (warm clay/terracotta orange-red) | Tertiary accent — sparingly, for warning states, "hot"/limited badges, occasional illustrative details. |
| `--color-accent-magenta` | `#C43F6E` (deep rose/magenta, echoes truck-art pinks) | Rare decorative accent for gradients (e.g., violet-to-magenta hero gradients), hover states on secondary elements. |

### 2.4 Neutrals

| Token | Hex | Usage |
|---|---|---|
| `--color-bg` | `#FFFFFF` | Primary background — lots of white space per the brand brief. |
| `--color-bg-subtle` | `#FAF8FC` (barely-tinted violet-white) | Alternating section backgrounds instead of stark white-on-white, gives warmth without weight. |
| `--color-ink` | `#211A2B` (near-black with a violet undertone) | Body text — never pure `#000000`, which feels cold/harsh. |
| `--color-ink-muted` | `#6B6275` | Secondary text, captions, metadata (dates, tags). |
| `--color-border` | `#E8E2EF` | Card borders, dividers — kept extremely light; prefer shadow/spacing over hard borders where possible. |

### 2.5 Usage Rules
- Every page must be predominantly white/near-white background with generous whitespace. Color is used for accents, CTAs, and illustrative moments — not as large flat color blocks that add visual weight without purpose.
- Gradients (violet → orchid → occasional magenta whisper) are reserved for hero sections and major CTA banners only — do not overuse, or the "joyful energy" becomes noise.
- Gold is the primary "this matters" color (badges, prices, featured tags) since it reads as premium/celebratory without being alarmist like red or generic like green.
- Never place body text in pure white on a light accent color (contrast failure) — always check WCAG AA contrast ratios (4.5:1 minimum for body text, 3:1 for large text/UI components).

---

## 3. Typography

### 3.1 Principle
Zero font-loading layout shift is a hard requirement (the previous conversations flagged this explicitly). This means:
- Self-host font files as **WOFF2**, subset to Latin + Urdu-adjacent extended Latin characters (for names/terms like "Jiddat") to keep file size minimal.
- Use `font-display: swap` but pair it with matching fallback font metrics (`size-adjust`, `ascent-override` if needed) to minimize the visible reflow when the webfont swaps in.
- Limit to **two font families maximum**: one display/heading face, one body face. A third "mono" face only if code snippets appear in Learn articles (can reuse a system mono stack — no need to self-host a third font just for this).

### 3.2 Recommended Pairing
- **Display/Headings:** A confident, slightly warm serif or semi-serif with character — something that can visually rhyme with the logo's calligraphic wordmark without literally copying it (e.g., a humanist serif like **Fraunces**, **Lora**, or **Source Serif 4** at heavier weights for H1/H2). This gives the "storytelling, heritage, craft" feeling the brand needs.
- **Body/UI:** A clean, highly legible geometric or humanist sans-serif optimized for screen reading at small sizes on cheap phone displays (e.g., **Inter**, **Manrope**, or **IBM Plex Sans**). Self-host 2-3 weights only (Regular, Medium, Semibold/Bold) — do not import 9 weights "just in case."

### 3.3 Scale (Fluid, using `clamp()`)
Use CSS `clamp()` for fluid type scaling instead of fixed breakpoint overrides — this reduces CSS size and gives naturally adaptive sizing across all screen widths without JS:

```css
--text-hero:   clamp(2.25rem, 5vw + 1rem, 4.5rem);
--text-h1:     clamp(1.875rem, 3vw + 1rem, 3rem);
--text-h2:     clamp(1.5rem, 2vw + 1rem, 2.25rem);
--text-h3:     clamp(1.25rem, 1vw + 1rem, 1.5rem);
--text-body:   1rem;      /* 16px baseline, never smaller for body copy */
--text-small:  0.875rem;
```

### 3.4 Rhythm
- Line-height 1.6-1.7 for body copy (important for Urdu-influenced bilingual readers and for readability on small screens).
- Line-height 1.1-1.2 for large display headings.
- Max content measure (line length) for body text: `65ch` — long-form Learn/Blog articles must never stretch full-width on desktop.

---

## 4. Spacing, Layout & Grid

- Base spacing unit: **4px**, scaling via a standard Tailwind-style spacing scale (4, 8, 12, 16, 24, 32, 48, 64, 96, 128px).
- Section vertical padding: generous — minimum `64px` top/bottom on mobile, `96-128px` on desktop, per the "lots of whitespace, lots of light" brand instruction.
- Content container max-width: `1280px`, centered, with responsive side padding (`16px` mobile → `32px` tablet → `64px`+ desktop gutters).
- Grid system: CSS Grid / Flexbox via Tailwind utility classes exclusively — no custom grid framework needed.
- Card grids (Products, Projects, Blog, Learn index pages): responsive auto-fit grid, e.g. `grid-template-columns: repeat(auto-fit, minmax(280px, 1fr))`, so it naturally reflows from 1 column (mobile) → 2 (tablet) → 3-4 (desktop) without manual breakpoint classes for every case.

---

## 5. The Ornamental Motif System

Derived from the logo's small diamond/star divider and lattice-flower icon:

- **Section dividers:** A thin horizontal line with a small centered diamond/star glyph (SVG, inline, tiny file size) — used between major homepage sections instead of a plain `<hr>`. Mirrors the exact divider style beneath the logo wordmark.
- **Card corner accents:** Optionally, a very subtle geometric corner flourish (a tiny quarter-lattice pattern, low opacity) on hero cards or featured product cards — must be implemented as inline SVG or CSS `background-image` with a tiny optimized SVG, never a large raster image.
- **Background texture:** For hero/section backgrounds, an extremely subtle (5-8% opacity) repeating geometric lattice/jaali pattern as an SVG background — adds cultural texture without adding meaningful file weight (SVG patterns compress to almost nothing) and without competing with foreground content.
- **Iconography:** Use a single consistent icon library style (line icons, 1.5-2px stroke weight) for feature callouts, nav icons, and social links. Prefer inlining used icons as SVG sprites rather than an icon font (icon fonts cause FOUC/layout shift and are heavier than needed for a small icon set).

**Rule of restraint:** These motifs are seasoning, not the meal. If a page feels "busy" or "cluttered," remove ornament before removing whitespace.

---

## 6. Animation & Motion Design

### 6.1 Philosophy
Motion should feel like **craftsmanship revealing itself** — content settling into place, gentle glows pulsing on hover — never like flashy marketing gimmicks that distract from load performance. Every animation must degrade gracefully to "instant, no animation" for users with `prefers-reduced-motion: reduce` set, and must never block content from being visible/readable if JS fails to load.

### 6.2 Core Animation Patterns

| Pattern | Implementation | Where Used |
|---|---|---|
| **Scroll-reveal fade/rise** | Vanilla JS `IntersectionObserver` toggles a `.is-visible` class; CSS handles the actual `opacity`/`transform` transition (`translateY(16px)` → `translateY(0)`, ~500ms ease-out). | Section headings, feature cards, project gallery items, timeline/roadmap entries as they scroll into view. |
| **Hover glow** | CSS `box-shadow`/`filter: drop-shadow()` transition on hover, using `--color-primary-light` — literally echoing the logo's own glow effect. | Buttons, product cards, nav links. |
| **Preloader fade-out** | CSS opacity/visibility transition triggered by `preloader.js` on load event (see `techstack.md` §9.1). | Every page, on initial load only. |
| **Micro-interactions** | Alpine.js `x-transition` directives for dropdown/accordion/modal open-close (built-in Alpine transition classes, zero extra JS). | Mobile nav menu, FAQ accordions, image lightbox modal, tab switches on product spec pages. |
| **Number/stat count-up** (for impact metrics, e.g. "units sold," "makers reached") | Small vanilla JS utility triggered by IntersectionObserver, using `requestAnimationFrame` for a smooth count from 0 to target value. | About page or homepage "impact" stat band, if/when real metrics exist. |
| **Roadmap/timeline reveal** | CSS `scroll-timeline`/`animation-timeline` where supported, with IntersectionObserver fallback — staggered reveal of roadmap phase cards. | "Future Horizons"/roadmap section. |
| **Hero background motion** | Extremely subtle CSS `@keyframes` gradient shift or slow-drifting SVG lattice pattern (looping, GPU-cheap `transform`/`opacity` only — never animate `width`/`height`/`top`/`left` which triggers layout thrashing). | Homepage hero background only. |

### 6.3 Performance Rules for Animation
- Animate only `transform` and `opacity` — never properties that trigger layout recalculation (`width`, `height`, `top`, `left`, `margin`).
- No animation library dependency (no GSAP, no Framer Motion, no AOS.js) — all of the above patterns are achievable in well under 100 lines of hand-written vanilla JS total, and CSS handles the actual motion.
- Respect `@media (prefers-reduced-motion: reduce)` globally — wrap all custom animation triggers in a JS check, and wrap all CSS keyframe/transition declarations in a media query fallback that sets `transition: none`/`animation: none`.

---

## 7. Adaptive & Responsive Behavior

- **Mobile-first.** The majority of the target audience will access this site on mid-range Android phones over mobile data. Design and build for a 360-390px viewport first, then progressively enhance for tablet (768px+) and desktop (1024px+, 1280px+).
- **Fluid typography and spacing** (via `clamp()`, see §3.3) reduce the number of hard breakpoints needed, which reduces CSS size and complexity.
- **Navigation:** Collapses to a hamburger-triggered full-screen or slide-in menu below `768px` (Alpine.js-managed `x-show`/`x-transition`, no JS framework router needed since this is a multi-page site with real URLs, not an SPA).
- **Images:** Use `srcset`/`sizes` attributes so mobile devices download the 480px variant, not the 1600px desktop hero image (see `techstack.md` §8).
- **Touch targets:** Minimum 44x44px tappable area for all buttons/links on mobile, per standard accessibility guidance — important given many users are on touch-only devices.
- **Adaptive network awareness (progressive enhancement, optional but recommended):** Use the `navigator.connection.effectiveType` / `saveData` API where available to conditionally skip loading hero video backgrounds or high-res galleries on detected slow connections (`2g`/`slow-2g` or `saveData: true`), substituting a static poster image instead. This directly serves the "Pakistan has poor internet" requirement at the UX layer, not just the infrastructure layer.

---

## 8. Page-Level UX Patterns (Inspired by the Four Referenced Landing Pages)

The user referenced a **Real Estate**, a **Nonprofit**, and a **Tourism** landing page pattern, plus an industrial ThemeForest factory theme, as design inspiration. Translate each pattern's *underlying UX mechanic* — not its literal industry content — onto Jiddat3D's actual pages as follows:

### 8.1 From the Real Estate pattern → Products Page
- **Listings grid → Product grid.** Card-based grid showing each machine (currently likely just "Jiddat One") with hero photo, name, short spec teaser, and price.
- **Search/filter by location/price/type → Filter by category/price/availability.** Even with one product today, build the filter UI to scale (by machine type: 3D Printer / future CNC / future Laser; by price range; by "in stock" status) so it doesn't need rework at Version 2.
- **Contact forms routed to agent profiles → Contact form routed to a WhatsApp/email CTA.** No shopping cart (per explicit brief); instead, each product detail page ends with a clear "Order via WhatsApp" / "Contact Us to Buy" CTA block, optionally alongside a short inquiry form that emails the team.
- **Calendar view for showings → (Future) Booking calendar for Maker Space visits/workshops.** Not required for V1, but the page template should leave a clean slot for a future "Book a visit" calendar component once the Maker Space pillar goes live.

### 8.2 From the Nonprofit pattern → About / Community / Impact
- **Donation CTAs → "Support the Movement" / "Join the Community" CTAs.** Instead of asking for money, the equivalent conversion action is: join the WhatsApp/Discord community, follow the newsletter, or share a maker project.
- **Volunteer signup → Community signup / newsletter form.** Reuse the same lightweight HTMX-powered form pattern.
- **Event calendar → Community meetups/maker challenges listing** (future-facing; can launch as a simple static list of past/upcoming events before a full calendar system is warranted).
- **Impact showcase with metrics and stories → The "Impact" band on the homepage or About page**, showing real (or honestly framed, growing) numbers — units sold, projects completed, workshops run — paired with short human stories, directly echoing the "Projects" page concept from the original blueprint (student built a drone, someone repaired an appliance, etc.).

### 8.3 From the Tourism pattern → Learn / Projects (Discovery & Storytelling)
- **Destination gallery with rich imagery → Projects gallery.** Every community-built project gets a card with a strong photo, title, and one-line hook, linking to a full story page.
- **Tour packages with pricing tiers → (Not directly applicable to V1)** — skip literal pricing tiers; instead apply the same visual "package card" pattern to structure **Learn content into tracks/series** (e.g., "Beginner: What is FDM Printing," "Intermediate: Designing Your First Part," "Advanced: Building From Scratch") so content feels like a guided journey, not a flat blog list.
- **Booking forms → Not applicable to V1** (no bookings yet); reserve the pattern for the future Maker Space booking feature.
- **Travel guides and tips sections → Learn category structure** directly maps here: practical, evergreen, SEO-valuable guide content.

### 8.4 From the Industrial/Factory Theme reference → Visual Restraint Warning
The referenced ThemeForest "Steelnova Modern Factory Industrial" theme uses dark backgrounds, heavy steel/metal textures, and an aggressive industrial tone. **Do not adopt its color palette or mood.** The only thing to borrow from an industrial theme reference is **structural confidence**: bold section headers, strong grid alignment, and clear visual hierarchy for technical specification tables (useful for the Product detail page's spec sheet). Everything else (dark tones, metal textures, heavy typography) contradicts the explicit "colourful, full of joy" brand brief and must be avoided.

### 8.5 AI Agent / Chat Concierge Pattern (Optional, Not V1)
The reference examples mention AI concierge agents (property inquiry bot, donor engagement bot, travel concierge bot) integrated into those landing pages. This is **explicitly out of scope for Jiddat3D Version 1** per the project brief (no accounts, no complex dynamic systems beyond forms). If revisited in a future version, the correct place for this is a lightweight HTMX-powered chat widget calling a Go backend route — not a third-party embedded script, to preserve the no-bloat principle. Document this as a "Version 3+" idea only; do not build it now.

---

## 9. The Preloader — Exact Visual Spec

Since the loading screen is a headline requirement:
1. Full-viewport, `--color-bg` (white) or a very soft violet-white gradient background.
2. Centered: the Jiddat3D logo mark (SVG, inlined directly in the HTML `<head>` or very first bytes of `<body>` so it renders before any external CSS/JS loads) with a subtle CSS-only pulsing glow animation (`opacity`/`filter: drop-shadow()` keyframe loop, echoing the logo's own glow).
3. Optionally, the small diamond/star divider glyph beneath it, animating with a gentle shimmer.
4. No spinner clichés (no generic circular spinner) — the pulsing logo glow itself communicates "loading" while staying on-brand.
5. Disappears via opacity fade (300-400ms) the moment the page is ready, per the behavior defined in `techstack.md` §9.1.

---

## 10. Accessibility Checklist (Non-Negotiable)

- Color contrast: minimum WCAG AA across all text/background combinations (test the violet-on-white and gold-on-white pairings specifically).
- All interactive elements reachable and operable via keyboard (forms, nav, accordions, lightbox — Alpine.js and HTMX both support this natively if built with real `<button>`/`<a>` elements, not `<div onclick>`).
- All images require descriptive `alt` text (enforced as a required field in the CMS content model for image uploads — see `content.md`).
- Respect `prefers-reduced-motion` (see §6.3).
- Semantic HTML throughout (`<nav>`, `<main>`, `<article>`, `<section>`, `<footer>`, proper heading hierarchy) — this also directly benefits SEO, which matters for the Learn section's organic search goals.
- Form fields have associated `<label>` elements, not placeholder-only labeling.

---

## 11. Design Deliverable Checklist for the Building Agent

Before considering the design system "done," confirm:
- [ ] Color tokens implemented as CSS custom properties in `ui/static/css/input.css`, consumed via Tailwind config `theme.extend.colors`.
- [ ] Two font families self-hosted, subset, loaded with `font-display: swap`.
- [ ] Fluid type scale implemented via `clamp()`.
- [ ] Ornamental divider/lattice SVG assets created and reused as partials, not duplicated inline everywhere.
- [ ] Preloader implemented exactly per §9, tested on throttled "Slow 3G" DevTools setting.
- [ ] All animations respect `prefers-reduced-motion`.
- [ ] Mobile nav, lightbox, accordions all keyboard-accessible.
- [ ] No green used anywhere, including default form validation states (use gold for success, terracotta for error/warning instead of the conventional red/green pairing — or muted red only for true errors, gold for success).

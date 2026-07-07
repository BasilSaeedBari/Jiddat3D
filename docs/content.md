# CONTENT.md — Jiddat3D Sitemap, Page Content & CMS Data Model

> **Audience:** An autonomous coding agent with zero prior context on this project.
> **Purpose:** Define every page on the site, its purpose, its exact content/copy (real or clearly-marked placeholder), and the CMS collection structure that stores editable content behind the admin panel.

---

## 1. Sitemap (Version 1)

```
/                       Home
/products               Products index
/products/{slug}        Product detail (e.g., /products/jiddat-one)
/learn                  Learn index (guides, tracked by difficulty)
/learn/{slug}           Learn article detail
/projects               Community Projects gallery
/projects/{slug}        Project story detail
/blog                   Blog index (the build journey)
/blog/{slug}            Blog post detail
/about                  About / Manifesto
/community              Community (WhatsApp/Discord links, future events)
/contact                Contact page (form + WhatsApp CTA)
/newsletter/thank-you   Simple confirmation page after signup (or inline HTMX success state, no separate page needed — prefer inline)
/404                    Not found page
/sitemap.xml            Auto-generated
/robots.txt             Static/auto-generated
```

No accounts, no cart, no checkout, no login for the public site — per the explicit "Version 1" brief. The Admin UI (`/_/`) is the only authenticated area, provided natively by PocketBase.

---

## 2. Global Elements (Appear on Every Page)

### 2.1 Header / Navigation
- Left: Jiddat3D logo (links to `/`).
- Center/Right (desktop) or slide-in menu (mobile): Home · Products · Learn · Projects · Blog · About · Community · Contact.
- A subtle CTA button in the nav: "Join the Community" or "Get Updates" (opens newsletter signup, either inline or scrolls to footer form).

### 2.2 Footer
- Short brand line: *"Jiddat3D exists to make creation accessible."*
- Quick links (mirrors nav).
- Newsletter signup form (HTMX-powered, POSTs to `/api/newsletter/subscribe`).
- Social links (WhatsApp Business, Instagram, Facebook, YouTube — whichever the founder actually uses; leave as placeholder icons/links to be filled in).
- Copyright line, "Made in Pakistan" mark, small ornamental divider glyph.

### 2.3 Newsletter Signup Component
- Single email input + submit button.
- On submit: HTMX POST to a Go route that creates a record in the `subscribers` collection (deduplicates on unique email).
- Success: inline message replaces the form, e.g., *"You're in. We'll keep you posted."* (no separate thank-you page needed — HTMX swap is enough, keeps friction low).
- Failure (duplicate email, invalid email): inline error message, no page reload.

---

## 3. Page-by-Page Content Specification

### 3.1 Home (`/`)

**Purpose:** Convert a first-time visitor into either (a) a product-curious buyer, (b) a content reader/subscriber, or (c) a community joiner. This is the highest-priority page for both design polish and load speed.

**Sections, top to bottom:**

1. **Hero**
   - Headline: **"Build Something Amazing."**
   - Subheadline: **"Affordable, repairable manufacturing tools — designed, assembled, and supported from right here in Pakistan."**
   - Two CTA buttons: `Explore Printers` (→ `/products`) and `Start Learning` (→ `/learn`).
   - Background: subtle animated gradient + lattice pattern per `design.md` §6.2/§5.

2. **Three-Pillar Cards** ("Learn / Build / Create")
   - **Learn** — "Free guides, videos, and honest lessons on 3D printing — from your first layer to your hundredth failed print." → links to `/learn`.
   - **Build** — "Affordable, repairable printers and parts, priced for the Pakistani maker, not the import market." → links to `/products`.
   - **Create** — "Real projects from real people — students, repair techs, artists, and tinkerers building things nobody imported for them." → links to `/projects`.

3. **Featured Product Spotlight**
   - Pulls the record(s) marked `featured: true` in the `products` collection.
   - Large product photo, name, one-line pitch, price (or "Contact for pricing" if not yet public), CTA → product detail page.

4. **"The Future Is Made Locally" Statement Band**
   - Large, confident typographic statement (full-bleed section, colored/gradient background per `design.md`), paraphrased from the founding narrative:
   - *"For too long, building things in Pakistan meant importing them instead. Jiddat3D exists to change that — one machine, one maker, one workshop at a time."*

5. **Impact / Social Proof Band** (nonprofit-pattern-inspired, §8.2 of `design.md`)
   - 3-4 stat callouts with count-up animation: e.g., "Units in the field," "Active makers," "Projects shared," "Guides published." Use real numbers once available; use honest, non-inflated placeholder copy until then (e.g., "Just getting started — be one of our first," rather than a fabricated number).

6. **Latest From the Blog / Learn** (auto-populated, 3 most recent items each, or interleaved "Latest Stories")
   - Pulls dynamically from the `content` collection, sorted by date, filtered by `type` and `published = true`.

7. **Roadmap Teaser**
   - Compressed visual timeline (see §3.9 About page for the full version) — 3-4 phase labels only, "See the full journey →" links to `/about`.

8. **Final CTA Band**
   - "Ready to build something?" + two buttons: `Browse Products` and `Join the Community`.

---

### 3.2 Products Index (`/products`)

**Purpose:** Real-estate-listing-pattern grid of available (and "coming soon") machines.

- Filter bar (HTMX-powered, non-blocking): filter by category (3D Printer / CNC — future / Laser — future), by price range, by availability (In Stock / Pre-Order / Coming Soon).
- Product grid: card per product — hero image, name, one-line pitch, price or "Contact for pricing," status badge (gold badge for "Featured" or "New").
- Empty/coming-soon state handled gracefully: if only one product exists, do not show an empty-feeling grid — show that one product prominently plus a "More machines coming — join the waitlist" card in the grid to fill visual space intentionally.

### 3.3 Product Detail (`/products/{slug}`)

- Hero: large photo/gallery (lightbox-enabled via Alpine.js), product name, tagline, price.
- Tabs or accordion sections (Alpine.js `x-data` tabs component): **Specifications** (structured spec table — build dimensions, layer resolution, materials supported, frame material, electronics), **Features** (bullet highlights), **Gallery** (additional photos/build process shots), **FAQ** (accordion).
- CTA block: "Order via WhatsApp" (deep link to WhatsApp Business number with a pre-filled message referencing the product name) + a lightweight inquiry form as an alternative (name, phone/email, message → emails the team via the contact form pipeline).
- No cart, no online payment in V1.

### 3.4 Learn Index (`/learn`)

**Purpose:** SEO-valuable, evergreen educational hub. Tourism-pattern-inspired "guided journey" structure (see `design.md` §8.3).

- Organized into tracks/difficulty levels, not a flat chronological blog feed:
  - **Getting Started** — What is FDM printing? Why PLA? How does a stepper motor work?
  - **Designing & Making** — Designing your first part, how to fix stringing, why beds need leveling.
  - **Going Deeper** — How CoreXY works, history of RepRap, how firmware works, building a printer from scratch.
- Filter/tag chips (HTMX fragment swap) to browse by track.
- Each entry: title, one-line description, estimated read time, track badge.

### 3.5 Learn Article Detail (`/learn/{slug}`)

- Rendered from the `content` collection's Markdown body (via goldmark, statically cached per `techstack.md` §5).
- Standard long-form article layout: `65ch` max width, generous line-height, optional embedded images/diagrams/code blocks, optional KaTeX math (opt-in per `techstack.md` §5.6).
- End-of-article: "Related Guides" (same track), and the newsletter signup component.

### 3.6 Projects Index (`/projects`)

**Purpose:** The community showcase — described in the original blueprint as "underrated" and central to long-term engagement.

- Gallery-style grid (tourism-destination-pattern, `design.md` §8.3): strong photo-first cards.
- Examples of the *kind* of content this section holds (write as illustrative placeholder entries until real customer stories exist):
  - "A student built a working drone frame."
  - "Someone repaired a washing machine door hinge instead of buying a new machine."
  - "A cosplayer printed an entire armor set."
  - "An engineer prototyped a prosthetic hand mechanism."
- Each project card links to a full detail page.

### 3.7 Project Detail (`/projects/{slug}`)

- Structure per the original blueprint: **Photos**, **Story** (markdown body), **Files** (optional downloadable STL/design file links, stored via PocketBase file fields), **Lessons Learned** (short callout block).
- Attribution: maker's name/handle (with consent), location (city-level, optional).

### 3.8 Blog Index (`/blog`) & Blog Detail (`/blog/{slug}`)

**Purpose:** The build-in-public journey — "this is where trust comes from," per the founding brief.

- Chronological feed, most recent first.
- Example topics (placeholder/illustrative): "Today we tested new rails," "Reducing printer cost by 12%," "Testing Pakistani-sourced bearings," "Behind the scenes: a failed print and what we learned."
- Detail page: same long-form article layout as Learn, but tagged/styled slightly differently (e.g., a "Blog" badge vs. a "Learn" track badge) to visually distinguish journal-style posts from structured guides, even though both live in the same `content` collection with a `type` field.

### 3.9 About (`/about`)

**Purpose:** Houses the brand manifesto and roadmap. This is the emotional core of the site.

**Content (adapted from the founding vision documents, tightened for a web page — full long-form manifesto can live as a downloadable PDF or a dedicated `/about/manifesto` long-read if the founder wants the full 6-page version available):**

1. **Origin story block** (first-person, not corporate):
   > "I started building printers because I wanted affordable manufacturing in Pakistan. I wanted engineers, students, and independent makers to own real machines — not just admire them in imported catalogs. Jiddat means innovation — a deliberate departure from the old way of doing things. That's what we're building toward."

2. **The Sovereignty Statement** (condensed from the manifesto):
   > "For decades, building physical things in Pakistan has meant importing them instead. When a machine breaks, we hunt for imported parts. When an engineer dreams up a new hardware product, they hit a wall of customs friction and the quiet belief that complex things simply aren't built here. Jiddat3D exists to dismantle that boundary. We're not just a hardware company — we're building the tools that let Pakistan manufacture its own future."

3. **Pillars Grid** (visual, icon + short label + one line, no fixed dates):
   - **Community** — a growing peer network for troubleshooting, sharing designs, and building in the open.
   - **Products** — durable, repairable machines built for local realities.
   - **Documentation** — transparent manuals, guides, and open resources.
   - **Manufacturing** — building toward locally sourced supply chains.
   - **Future Horizons** — from 3D printing today toward CNC, laser cutting, and beyond, as the ecosystem grows.

4. **Engineering Principles Block:**
   - **Right to Repair** — every machine can be disassembled and fixed with basic hand tools, no proprietary lock-in.
   - **Built for Local Realities** — designed around Pakistan's actual electrical, logistical, and supply-chain conditions, not assumptions imported from elsewhere.
   - **Open Where It Matters** — firmware and documentation lean open-source rather than walled-garden.

5. **Loose Roadmap** (explicitly NOT date-committed, per the founder's direction — avoid fixed years or promised kits):
   - Present as a flowing narrative or a simple undated horizontal progression, e.g.:
   ```
   Today                         Tomorrow                        Someday
   Desktop 3D printing    →      Expanding tools for makers   →  A full local manufacturing
   for individual creators       and small workshops              ecosystem — CNC, laser
                                                                    cutting, shared maker spaces
   ```
   - Copy should emphasize *direction*, not deadlines: *"We're not attaching fixed dates or promised kits to this journey — creation comes first, and the roadmap grows as we do."*

6. **CTA:** "Follow the journey" (newsletter) + "See what people are building" (→ `/projects`).

### 3.10 Community (`/community`)

- WhatsApp community / Discord / Facebook group links (whichever channels the founder actually runs).
- Short framing copy: *"This is where makers help each other — troubleshooting, sharing designs, and celebrating what gets built."*
- Placeholder section for future: meetups, design challenges, open hardware repository links — structured so it's easy to add entries later without a redesign.

### 3.11 Contact (`/contact`)

- Contact form: name, email/phone, subject (dropdown: General / Product Inquiry / Support / Partnership), message. Submits via HTMX to a Go route, sends an email via the configured Mailtrap SMTP (or, for production, whatever transactional SMTP the founder ultimately uses — Mailtrap is fine for a start but note in `techstack.md` that Mailtrap's free tier is a sandbox and a production SMTP relay should replace it before real customer emails need delivery — document this but do not block build on it).
- WhatsApp Business direct link/button as the primary, lowest-friction contact method (very common and expected pattern in Pakistan's e-commerce/small business context).
- Physical location/city (optional, if the founder wants to disclose it) and business hours.

### 3.12 404 Page

- On-brand, not a generic error page: reuse the preloader's logo glow motif, a friendly line like *"Looks like this page hasn't been built yet — even we start with a blank plate sometimes."*, and links back to Home/Products/Learn.

---

## 4. CMS Data Model (PocketBase Collections)

This section is the authoritative schema. Build these as Go migrations (per `techstack.md` §5.2) so the schema is reproducible without manual admin clicking.

### 4.1 Collection: `content`
Used for Learn articles, Blog posts, and Project stories (unified collection, differentiated by `type`) — this keeps the admin panel simple (one place to manage all written content) while still allowing type-specific templates at render time.

| Field | Type | Rules |
|---|---|---|
| `title` | Text | Required |
| `slug` | Text | Required, Unique, URL-safe (lowercase, hyphenated) |
| `type` | Select | Options: `learn`, `blog`, `project`. Required. |
| `track` | Select | Options: `getting-started`, `designing-making`, `going-deeper`. Only relevant/shown when `type = learn`. |
| `excerpt` | Text (short) | Required. Used in index cards and SEO meta description. |
| `body` | Editor / Long Text | Required. Raw Markdown, compiled via goldmark on save. |
| `hero_image` | File (image) | Required. Auto-processed into responsive WebP variants per `techstack.md` §8. |
| `gallery_images` | File (multiple, images) | Optional. Used primarily for `project` type. |
| `attachment_files` | File (multiple, any) | Optional. Used for `project` type "Files" (e.g., STL downloads) or `learn` type reference PDFs. |
| `author_name` | Text | Optional (defaults to "Jiddat3D Team" if blank). |
| `tags` | Select (multiple) or JSON array | Optional, free-form tagging for related-content lookups. |
| `has_math` | Bool | Default: False. Flags whether to load the optional KaTeX assets on render. |
| `featured` | Bool | Default: False. Used to pull items into homepage spotlights. |
| `published` | Bool | Default: False. Only `true` records get statically cached and publicly routed. |
| `published_at` | Date | Auto-set or manually adjustable, used for sort order and display date. |
| `seo_description` | Text | Optional override for meta description; falls back to `excerpt` if blank. |

### 4.2 Collection: `products`

| Field | Type | Rules |
|---|---|---|
| `name` | Text | Required |
| `slug` | Text | Required, Unique |
| `category` | Select | Options: `3d-printer`, `cnc`, `laser`, `accessory`. Required (future-proofs for expansion beyond printers). |
| `tagline` | Text (short) | Required. One-line pitch. |
| `description` | Editor / Long Text | Markdown body for the full product story/description. |
| `specs` | JSON | Structured key-value spec sheet (e.g., `{"Build Volume": "220x220x250mm", "Layer Resolution": "0.05-0.3mm"}`), rendered as a table. |
| `price` | Number | Optional — if blank, display "Contact for pricing." |
| `currency` | Text | Default: `PKR`. |
| `availability` | Select | Options: `in-stock`, `pre-order`, `coming-soon`. Required. |
| `hero_image` | File (image) | Required. |
| `gallery_images` | File (multiple, images) | Optional. |
| `featured` | Bool | Default: False. |
| `active` | Bool | Default: False. Only `true` products are publicly listed/cached. |
| `whatsapp_message_template` | Text | Optional pre-filled inquiry text for the "Order via WhatsApp" deep link; falls back to a generic default referencing the product name. |

### 4.3 Collection: `subscribers`

| Field | Type | Rules |
|---|---|---|
| `email` | Email | Required, Unique |
| `active` | Bool | Default: True |
| `source` | Text | Optional — records which page/form the signup came from (e.g., `footer`, `learn-article`), useful for later analytics without adding a tracking script. |
| `subscribed_at` | Date | Auto-set on creation. |

### 4.4 Collection: `contact_submissions`

| Field | Type | Rules |
|---|---|---|
| `name` | Text | Required |
| `contact_method` | Text | Required (email or phone) |
| `subject` | Select | Options: `general`, `product-inquiry`, `support`, `partnership` |
| `message` | Long Text | Required |
| `related_product` | Relation → `products` | Optional, set if submitted from a product detail page's inquiry form. |
| `status` | Select | Options: `new`, `read`, `responded`. Default: `new`. Lets the founder triage inquiries directly in the Admin UI. |
| `submitted_at` | Date | Auto-set. |

### 4.5 Collection: `site_settings` (Singleton-style, one record)

Holds editable global content so the founder never needs a code change for small copy tweaks:

| Field | Type | Rules |
|---|---|---|
| `hero_headline` | Text | Default: "Build Something Amazing." |
| `hero_subheadline` | Text | Default per §3.1. |
| `whatsapp_business_number` | Text | Required for all WhatsApp CTA links across the site. |
| `instagram_url` / `facebook_url` / `youtube_url` / `discord_url` | Text (each optional) | Populate footer/community social links. |
| `impact_stat_1_label` / `impact_stat_1_value` (repeat x3-4) | Text/Number | Populates the homepage impact band without needing hardcoded template edits. |
| `manifesto_intro` | Long Text | The About page's origin story block, editable without a redeploy. |

---

## 5. Content Authoring Workflow (For the Founder, Documented for the Agent to Support)

1. Log into `/_/` (PocketBase Admin UI).
2. To publish a new Learn/Blog/Project entry: create a record in `content`, write the body in Markdown directly in the editor field, upload a hero image, set `type` and (if applicable) `track`, toggle `published = true`, save.
3. On save, the Go hook (per `techstack.md` §5.3) automatically compiles and caches the page — no further action needed, no "publish/deploy" button beyond the save itself.
4. To take content down: toggle `published = false` (soft) or delete the record (hard — also removes the cached file).
5. To feature something on the homepage: toggle `featured = true` on a `content` or `products` record.
6. To edit global copy (hero text, social links, WhatsApp number): edit the single `site_settings` record.

This workflow must require **zero technical knowledge beyond basic form-filling and Markdown syntax** — the explicit success criterion from the project brief ("all I have to do is... go to the admin panel to change content").

---

## 6. SEO & Metadata Requirements

- Every cached page must render proper `<title>`, `<meta name="description">`, Open Graph tags (`og:title`, `og:description`, `og:image` using the record's `hero_image`), and `twitter:card` tags — all templated from the content record's fields (falling back sensibly if a field like `seo_description` is empty).
- `sitemap.xml` auto-generated from all `published`/`active` records across `content` and `products` collections (Go route, regenerated on the same hooks that trigger page caching).
- Semantic heading structure (one `<h1>` per page) enforced by the templates, not left to content authors to get right in Markdown bodies (the Markdown body should start at `<h2>` level internally since the page template already renders the `<h1>` from the `title` field).

# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

- **Primary:** Engineering students (including Final Year Project / FYP students) and independent makers in Pakistan who need fast, transparent PKR quotes to 3D-print custom CAD parts (`.stl`, `.obj`, `.3mf`) or want to buy a reliable, locally built and repairable 3D printer via WhatsApp.
- **Secondary:** Pakistani schools, colleges, and engineering universities evaluating turn-key STEAM makerspace packages and institutional tenders, plus beginners using free guides to learn 3D printing under local conditions (humidity, grid load-shedding, UPS sizing).

## Product Purpose

Jiddat3D's website is the digital front door, instant 3D-print quotation engine, educational hub, and WhatsApp sales pipeline for Jiddat3D's Karachi engineering lab. Success means a visitor can either upload a 3D model and immediately get an honest PKR print quote to confirm on WhatsApp (`0300-1318112`), evaluate and inquire about an in-house Jiddat3D printer with full confidence in local support, or request a school STEAM lab package without friction.

## Positioning

Unlike resellers that dropship closed-source imported consumer printers (e.g., Creality, Bambu Lab, Anycubic) with no local repair path, Jiddat3D designs, bench-assembles, and print-tests its own specialized FDM/FFM 3D printers in Karachi with a **100% open-source hardware and firmware guarantee** (CAD, BOMs, Marlin/Klipper profiles remain public and repairable with off-the-shelf parts even if Jiddat3D ceased to exist) paired with an instant server-side 3D model slicer quoting directly in Pakistani Rupees.

## Operating Context

- **Geography & infrastructure:** Operated exclusively out of **Karachi, Sindh, Pakistan** (DHA Phase 1), serving customers across Pakistan. Operating realities include summer humidity, grid load-shedding (UPS usage), and mobile-heavy visitors on variable bandwidth.
- **Commercial workflow:** Every order is **100% prepaid via Raast or direct bank transfer (IBFT)** with formal NTN invoices; **Cash on Delivery (COD) is not offered**. Order closing and engineer confirmation happen directly on **WhatsApp (`0300-1318112` / `+92 300 1318112`)** rather than through an online shopping cart.
- **Quotation workflow:** Visitors drop or select `.3mf`, `.obj`, or `.stl` files at `/instant-quote`; the Go server parses geometry, computes print time and material weight using admin-configured `slicer_settings` and `filament_types` (LayerInvoice formulas with internal safety buffers kept admin-only), and prefills a WhatsApp order message.
- **Admin workflow:** Content (blog, learn guides, projects), products, contact inquiries, newsletter subscribers, site settings, and slicer/filament pricing are managed through PocketBase Admin UI (`/_/`).

## Capabilities and Constraints

- **Implemented capabilities:**
  - Server-side instant 3D slicer & PKR quote calculator (`.3mf`, `.obj`, `.stl`) with interactive infill density and filament switching (`/instant-quote`, `/api/instant-slice`, `/api/recalculate-quote`).
  - Specialized printer catalog and detail pages (`/printers`, `/products`, `/printers/{slug}`), print service & FYP clinic (`/print-service`, `/fyp`), schools & tenders (`/schools`), workshops (`/workshops`), filament/spares (`/filament`), gallery (`/gallery`), support & warranty (`/support`), printer finder quiz (`/find-my-printer`), handbook/blog/projects (`/learn`, `/blog`, `/projects`), community (`/community`), and contact/newsletter forms (`/api/contact/submit`, `/api/newsletter/subscribe`).
- **Hard business & technical constraints:**
  - **Karachi-only HQ:** Never reference Lahore or other cities as headquarters.
  - **No imported/dropshipped brands:** Only specialized in-house Jiddat3D machines.
  - **Prepaid only:** Strictly IBFT / Raast prepaid; never promise Cash on Delivery.
  - **Warranty terms:** Strictly **6-Month Working Warranty** and **STEAM Education Warranty** (never claim a generic 12-month warranty).
  - **V1 scope:** No public user accounts, logins, shopping cart, or online card checkout.
  - **Lean stack:** Go + PocketBase SSR + `html/template` + standalone Tailwind CSS v3 CLI + vendored HTMX & Alpine.js; do not introduce Node.js, npm, or JS build bundlers.
  - **Open/roadmap boundaries:** In-house filament extrusion, MSLA resin printing, rapid PCB prototyping, and desktop CNC milling are future R&D / roadmap items and must not be represented as in-stock today.

## Brand Commitments

- **Brand name:** Jiddat3D (*Jiddat* = innovation / novelty in Urdu).
- **Primary contact & community channels:** WhatsApp Business `0300-1318112` (`+92 300 1318112`), email `info@jiddat3d.com`, WhatsApp Community, and Discord.
- **Logo asset:** `ui/static/img/logo/JiddatBrightTrasnparent-removebg-preview.png` (dark violet/magenta mark; binding requirement to keep it clearly readable on a light plinth against dark surfaces).

## Evidence on Hand

- **Available now:**
  - Official Jiddat3D logo (`ui/static/img/logo/JiddatBrightTrasnparent-removebg-preview.png`).
  - Live server-side `.3mf`/`.obj`/`.stl` geometry parser and LayerInvoice PKR quotation engine (`internal/slicer/`).
  - Six seeded Jiddat3D machine definitions (`pb_migrations/1700000004_business_seed.go`) and initial RepRap/Pakistan editorial article (`pb_migrations/1700000002_seed_blog.go`).
- **Explicit absences (do not fabricate):**
  - No real workshop or Jiddat3D machine photography is on hand yet.
  - No real printed-part photos with verified print logs are on hand yet.
  - No real customer testimonials, school/university partner logos, or unit-count/maker-count statistics are on hand yet. Future work must not invent fake numbers, fake customer quotes, or fake institutional endorsements.

## Product Principles

1. **Instant clarity over guesswork:** Give students and makers immediate, transparent PKR numbers and concrete specs before asking them to message sales.
2. **Sovereign, repairable engineering:** Emphasize open-source transparency, local Karachi repairability, and honest limits (prepaid, 6-month warranty, UPS guidance) over hype.
3. **WhatsApp-first conversion without clutter:** Make locking in an order or quote on WhatsApp effortless from every key decision point without repeating competing buttons on the same screen.
4. **Honest proof only:** Show real engineering mechanics and live tools (like the instant slicer) rather than fabricated testimonials, inflated metrics, or stock import imagery.
5. **Fast and lean on Pakistani connections:** Keep pages lightweight, server-rendered, and responsive on mobile devices with zero heavy client bundles.

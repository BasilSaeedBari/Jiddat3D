/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    "./ui/templates/**/*.html",
    "./ui/static/js/**/*.js",
    "./internal/routes/**/*.go"
  ],
  theme: {
    extend: {
      colors: {
        // Charcoal Dark Palette
        workshop: {
          bg: '#131217',
          surface: '#1B1921',
          card: '#23202B',
          hover: '#2C2836',
          border: '#343040',
          subtle: '#17151C',
        },
        // Warm Stone-Ash Light Palette (pairs cleanly with Dusty Purple & Charcoal)
        parchment: {
          bg: '#F1EEEA',
          surface: '#FAF8F5',
          card: '#FAF8F5',
          hover: '#E6E1DB',
          border: '#D1C9C0',
          subtle: '#E8E3DC',
          ink: '#17151C',
          muted: '#5A5363',
        },
        // Dusty Purple (Brand Structural Accent)
        plum: {
          DEFAULT: '#5A4769',
          deep: '#3E314A',
          dark: '#2D2436',
          light: '#8F7CA0',
          mist: '#EAE3F0',
        },
        // Burnished Gold (Primary Focal Accent — replaces orange across all templates)
        filament: {
          DEFAULT: '#D4A843',
          hover: '#E5BC5E',
          dark: '#9A7218',
          light: '#E8C97A',
        },
        primary: {
          DEFAULT: '#D4A843',
          light: '#E5BC5E',
          dark: '#9A7218',
        },
        // Dusty Purple mapped onto secondary accent slots for site-wide coherence
        teal: {
          DEFAULT: '#8F7CA0',
          light: '#A896B8',
          dark: '#5A4769',
        },
        // Restrained WhatsApp Action Green
        whatsapp: {
          DEFAULT: '#1D8F4B',
          dark: '#156E38',
          hover: '#22A356',
        },
        // Supporting Tokens (all aligned to Gold / Dusty Purple / Charcoal)
        accent: {
          gold: '#D4A843',
          teal: '#8F7CA0',
          terracotta: '#9A7218',
          green: '#2A8F56',
          blue: '#5A4769',
        },
        // Canvas & Surface
        bg: {
          DEFAULT: '#131217',
          subtle: '#1B1921',
          card: '#23202B',
        },
        // Typography Colors
        ink: {
          DEFAULT: '#F2EFE9',
          muted: '#A39DAE',
          dim: '#756F80',
          dark: '#17151C',
        },
        border: {
          DEFAULT: '#343040',
          subtle: '#25222E',
          bright: '#4A435A',
          light: '#D1C9C0',
        },
      },
      fontFamily: {
        sans: ['Figtree', 'system-ui', 'sans-serif'],
        display: ['"Bricolage Grotesque"', 'system-ui', 'sans-serif'],
        editorial: ['"Instrument Serif"', 'Georgia', 'serif'],
        serif: ['"Instrument Serif"', 'Georgia', 'serif'],
        calligraphic: ['"Instrument Serif"', 'Georgia', 'serif'],
        numeral: ['"Martian Mono"', 'monospace'],
        mono: ['"Martian Mono"', 'monospace'],
        urdu: ['"Noto Nastaliq Urdu"', 'serif'],
      },
      borderRadius: {
        'nouveau': '10px',
        'nouveau-sm': '6px',
        'nouveau-lg': '14px',
        'nouveau-arch': '10px',
        'none': '0px',
      },
      boxShadow: {
        'card': '0 10px 30px -12px rgba(0, 0, 0, 0.5)',
        'card-hover': '0 16px 36px -14px rgba(0, 0, 0, 0.6)',
        'card-light': '0 8px 24px -12px rgba(23, 21, 28, 0.10)',
        'card-light-hover': '0 14px 32px -12px rgba(23, 21, 28, 0.16)',
      },
      fontSize: {
        'hero': 'clamp(2.6rem, 5.4vw + 0.8rem, 4.85rem)',
        'h1': 'clamp(2.1rem, 3.4vw + 0.8rem, 3.35rem)',
        'h2': 'clamp(1.75rem, 2.4vw + 0.8rem, 2.65rem)',
        'h3': 'clamp(1.25rem, 1.1vw + 0.8rem, 1.55rem)',
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}

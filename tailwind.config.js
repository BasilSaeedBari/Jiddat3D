/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    "./ui/templates/**/*.html",
    "./ui/static/js/**/*.js"
  ],
  theme: {
    extend: {
      colors: {
        // Modern Workshop Dark Palette
        workshop: {
          bg: '#0C0E12',
          surface: '#141820',
          card: '#1B202A',
          hover: '#242B38',
          border: '#2C3545',
          subtle: '#181E27',
        },
        // Artisan Light Palette (Warm Drafting Parchment)
        parchment: {
          bg: '#FAF7F2',
          surface: '#FFFFFF',
          card: '#FFFFFF',
          hover: '#F2EDE4',
          border: '#E2D7C7',
          subtle: '#F5EFEB',
          ink: '#1D1A16',
          muted: '#635D54',
        },
        // Filament Orange (Primary Brand Accent)
        filament: {
          DEFAULT: '#FF6A1A',
          hover: '#FF8038',
          dark: '#E0550B',
          light: '#FFA066',
        },
        primary: {
          DEFAULT: '#FF6A1A',
          light: '#FF8038',
          dark: '#E0550B',
        },
        // Precision Cool Teal / Art Nouveau Jade
        teal: {
          DEFAULT: '#17C3B2',
          light: '#2FE0CE',
          dark: '#118C7F',
        },
        // Official WhatsApp Business Green
        whatsapp: {
          DEFAULT: '#25D366',
          dark: '#128C7E',
          hover: '#20BA5A',
        },
        // UI Accents & Memphis Contrasts
        accent: {
          gold: '#E5A93C',
          teal: '#17C3B2',
          terracotta: '#D9530F',
          magenta: '#E83E8C',
          green: '#3DDC84',
          blue: '#2B59FF',
        },
        // Canvas & Surface
        bg: {
          DEFAULT: '#0C0E12',
          subtle: '#141820',
          card: '#1B202A',
        },
        // Typography Colors
        ink: {
          DEFAULT: '#F5F3EF',
          muted: '#9EA8B6',
          dim: '#6B7687',
          dark: '#1D1A16',
        },
        border: {
          DEFAULT: '#2C3545',
          subtle: '#1E2532',
          bright: '#414D63',
          light: '#E2D7C7',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        display: ['Space Grotesk', 'system-ui', 'sans-serif'],
        serif: ['Space Grotesk', 'system-ui', 'sans-serif'],
        calligraphic: ['"Cinzel Decorative"', '"Cormorant Garamond"', '"Playfair Display"', 'Georgia', 'serif'],
        numeral: ['"Cormorant Garamond"', '"Cinzel Decorative"', '"Playfair Display"', 'Georgia', 'serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      borderRadius: {
        'nouveau': '28px 4px 28px 4px',
        'nouveau-sm': '14px 2px 14px 2px',
        'nouveau-lg': '42px 6px 42px 6px',
        'nouveau-arch': '36px 36px 4px 4px',
        'none': '0px',
      },
      dropShadow: {
        'glow': '0 0 20px rgba(255, 106, 26, 0.45)',
        'glow-teal': '0 0 20px rgba(23, 195, 178, 0.45)',
        'glow-wa': '0 0 20px rgba(37, 211, 102, 0.45)',
      },
      boxShadow: {
        'card': '0 4px 20px rgba(0, 0, 0, 0.4)',
        'card-hover': '0 8px 30px rgba(0, 0, 0, 0.6), 0 0 15px rgba(255, 106, 26, 0.15)',
        'card-light': '0 4px 20px rgba(180, 150, 120, 0.15)',
        'card-light-hover': '0 8px 30px rgba(180, 150, 120, 0.25), 0 0 12px rgba(217, 83, 15, 0.2)',
        'glow-filament': '0 0 25px rgba(255, 106, 26, 0.3)',
      },
      fontSize: {
        'hero': 'clamp(2.25rem, 5vw + 1rem, 4.25rem)',
        'h1': 'clamp(1.875rem, 3vw + 1rem, 3rem)',
        'h2': 'clamp(1.5rem, 2vw + 1rem, 2.25rem)',
        'h3': 'clamp(1.25rem, 1vw + 1rem, 1.5rem)',
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}

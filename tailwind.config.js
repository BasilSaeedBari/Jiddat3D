/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./ui/templates/**/*.html",
    "./ui/static/js/**/*.js"
  ],
  theme: {
    extend: {
      colors: {
        // Modern Workshop Dark Palette
        workshop: {
          bg: '#0C0E11',
          surface: '#14181F',
          card: '#1A1F29',
          hover: '#222936',
          border: '#28303F',
          subtle: '#181D26',
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
        // Precision Cool Teal
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
        // UI Accents
        accent: {
          gold: '#FFB400',
          teal: '#17C3B2',
          terracotta: '#FF6A1A',
          magenta: '#E83E8C',
          green: '#3DDC84',
          blue: '#2B59FF',
        },
        // Canvas & Surface
        bg: {
          DEFAULT: '#0C0E11',
          subtle: '#14181F',
          card: '#1A1F29',
        },
        // Typography Colors (High contrast on dark workshop)
        ink: {
          DEFAULT: '#F4F1EA',
          muted: '#9DA3AF',
          dim: '#656D7D',
          dark: '#0C0E11',
        },
        border: {
          DEFAULT: '#28303F',
          subtle: '#1E2430',
          bright: '#3E495D',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        display: ['Space Grotesk', 'system-ui', 'sans-serif'],
        serif: ['Space Grotesk', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      dropShadow: {
        'glow': '0 0 20px rgba(255, 106, 26, 0.45)',
        'glow-teal': '0 0 20px rgba(23, 195, 178, 0.45)',
        'glow-wa': '0 0 20px rgba(37, 211, 102, 0.45)',
      },
      boxShadow: {
        'card': '0 4px 20px rgba(0, 0, 0, 0.4)',
        'card-hover': '0 8px 30px rgba(0, 0, 0, 0.6), 0 0 15px rgba(255, 106, 26, 0.15)',
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

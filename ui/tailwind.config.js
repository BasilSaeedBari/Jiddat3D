/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "../ui/templates/**/*.html",
    "../ui/static/js/**/*.js"
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#5B2A86',
          light: '#9B5DE5',
          dark: '#3A1A57',
        },
        accent: {
          gold: '#E8A93A',
          teal: '#1E8A8A',
          terracotta: '#D9603B',
          magenta: '#C43F6E',
        },
        bg: {
          DEFAULT: '#FFFFFF',
          subtle: '#FAF8FC',
        },
        ink: {
          DEFAULT: '#211A2B',
          muted: '#6B6275',
        },
        border: '#E8E2EF',
      },
      fontFamily: {
        serif: ['Fraunces', 'serif'],
        sans: ['Inter', 'sans-serif'],
      },
      dropShadow: {
        'glow': '0 0 16px rgba(155, 93, 229, 0.5)',
      },
      fontSize: {
        'hero': 'clamp(2.25rem, 5vw + 1rem, 4.5rem)',
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

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'dark-bg': '#05060A',
        'dark-surface': '#0A0B10',
        'dark-border': '#1F2933',
        'primary': {
          500: '#3B82F6',
          600: '#2563EB',
        },
        'neon-blue': '#00D4FF',
        'neon-blue-hover': '#00A8CC',
      },
    },
  },
  plugins: [
    require('tailwindcss-rtl'),
  ],
}

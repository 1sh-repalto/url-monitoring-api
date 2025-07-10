/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",                // Vite entry
    "./src/**/*.{js,ts,jsx,tsx}",  // <‑‑ include TS/TSX!
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}

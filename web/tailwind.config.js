/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: '#0B1121', // Deep dark blue/black
        surface: '#151E32',    // Lighter blue/slate for cards
        surfaceHighlight: '#1E293B',
        primary: '#3B82F6',    // Neon Blue
        secondary: '#10B981',  // Neon Green
        accent: '#8B5CF6',     // Neon Purple
        danger: '#EF4444',
        warning: '#F59E0B',
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      boxShadow: {
        'glow-blue': '0 0 20px -5px rgba(59, 130, 246, 0.5)',
        'glow-green': '0 0 20px -5px rgba(16, 185, 129, 0.5)',
      }
    },
  },
  plugins: [],
}

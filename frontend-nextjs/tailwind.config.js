/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        cyber: {
          dark: "#030712",        // Deep outer space gray-black
          card: "rgba(17, 24, 39, 0.65)", // Glassmorphic card surface
          orange: "#ea580c",      // Energetic PromptWars orange
          cyan: "#06b6d4",        // Sleek medical digital twin cyan
          lime: "#10b981",        // Healthy success indicators
          gray: "#9ca3af",        // Subdued high-tech labels
          border: "rgba(255, 255, 255, 0.08)" // Micro-thin glass border
        }
      },
      fontFamily: {
        orbitron: ["var(--font-orbitron)", "sans-serif"],
        mono: ["var(--font-mono)", "monospace"],
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'grid-pattern': "linear-gradient(to right, rgba(255, 255, 255, 0.03) 1px, transparent 1px), linear-gradient(to bottom, rgba(255, 255, 255, 0.03) 1px, transparent 1px)",
      },
    },
  },
  plugins: [],
}

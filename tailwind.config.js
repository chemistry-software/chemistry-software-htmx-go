/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class', // Enable dark mode based on class
  content: [
    "./templates/*.html",
    "./static/js/*.js",
    "./main.go", // Or any other Go files that might include class names
  ],
  theme: {
    extend: {
      animation: { // Ensure 'spin' animation is defined
        spin: 'spin 1s linear infinite',
      },
      keyframes: { // Define 'spin' keyframes if not already defined
        spin: {
          from: { transform: 'rotate(0deg)' },
          to: { transform: 'rotate(360deg)' },
        },
      },
    },
  },
  plugins: [],
} 
/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class', // Enable dark mode based on class
  content: [
    "./templates/*.html",
    "./static/js/*.js",
    "./main.go", // Or any other Go files that might include class names
  ],
  theme: {
    extend: {},
  },
  plugins: [],
} 
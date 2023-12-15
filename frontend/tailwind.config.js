/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./**/**/*.{html,js,go,templ}"],
    theme: {
      extend: {
        colors: {
          main: "#e51636",
          black: "#4a5568"
        },
      },
    },
    plugins: [],
  }
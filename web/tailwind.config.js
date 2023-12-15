/* eslint-disable no-undef */
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{ts,vue,css}",
  ],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      "business",
    {
      light: {
        ...require("daisyui/src/theming/themes")["corporate"],
        accent: "oklch(0.67271 0.167726 35.7915)",
      },
    },],
  },
}


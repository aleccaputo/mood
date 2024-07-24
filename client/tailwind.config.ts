import type { Config } from "tailwindcss";

export default {
  content: ["./app/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        main: '#c4a1ff',
        overlay: 'rgba(0,0,0,0.8)',
        // background color overlay for alert dialogs, modals, etc.
        // light mode
        bg: '#daf5f0',
        text: '#000',
        border: '#000',
        // dark mode
        darkBg: '#0f3730',
        darkText: '#eeefe9',
        darkBorder: '#000',
      },
      borderRadius: {
        base: '5px'
      },
      boxShadow: {
        light: '3px 4px 0px 0px #000',
        dark: '3px 4px 0px 0px #000',
      },
      translate: {
        boxShadowX: '3px',
        boxShadowY: '4px',
        reverseBoxShadowX: '-3px',
        reverseBoxShadowY: '-4px',
      },
      fontWeight: {
        base: '500',
        heading: '700',
      },
    },
  },
  darkMode: 'media',
  plugins: [],
} satisfies Config;

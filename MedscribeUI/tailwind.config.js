/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx,mdx}',
    './.storybook/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      // Add custom scrollbar styles
      scrollbar: {
        width: '8px',
        track: 'transparent',
        thumb: '#E5E7EB',
        'thumb-hover': '#D1D5DB',
        radius: '9999px',
      },
    },
  },
  plugins: [
    function({ addUtilities }) {
      const newUtilities = {
        '.scrollbar-thin': {
          scrollbarWidth: 'thin',
          '&::-webkit-scrollbar': {
            width: '8px',
            height: '8px',
          },
        },
        '.scrollbar-thumb-rounded-full': {
          '&::-webkit-scrollbar-thumb': {
            borderRadius: '9999px',
          },
        },
        '.scrollbar-thumb-gray-300': {
          '&::-webkit-scrollbar-thumb': {
            backgroundColor: '#D1D5DB',
          },
        },
        '.scrollbar-thumb-gray-400': {
          '&::-webkit-scrollbar-thumb': {
            backgroundColor: '#9CA3AF',
          },
        },
        '.scrollbar-track-transparent': {
          '&::-webkit-scrollbar-track': {
            backgroundColor: 'transparent',
          },
        },
        '.hover\\:scrollbar-thumb-gray-400:hover': {
          '&::-webkit-scrollbar-thumb': {
            backgroundColor: '#9CA3AF',
          },
        },
      };
      addUtilities(newUtilities, ['responsive', 'hover']);
    },
  ],
};

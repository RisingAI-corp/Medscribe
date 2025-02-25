import type { Preview } from '@storybook/react';
import '../src/index.css'; // must require to expose tailwind styles to stories

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },
};

export default preview;

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import '@mantine/core/styles.css';
import '@mantine/tiptap/styles.css';
import './index.css';

import { MantineProvider } from '@mantine/core';

const rootElement = document.getElementById('root');

if (rootElement) {
  createRoot(rootElement).render(
    <StrictMode>
      <MantineProvider>
        <App />
      </MantineProvider>
    </StrictMode>,
  );
} else {
  console.error('Root element not found');
}

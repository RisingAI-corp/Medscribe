import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import App from './App';
import '@mantine/core/styles.css';
import '@mantine/tiptap/styles.css';
import './index.css';
import '@mantine/notifications/styles.css';
import { Notifications } from '@mantine/notifications';

import { MantineProvider } from '@mantine/core';

const rootElement = document.getElementById('root');
const queryClient = new QueryClient();
if (rootElement) {
  createRoot(rootElement).render(
    <StrictMode>
      <MantineProvider>
        <Notifications />
        <QueryClientProvider client={queryClient}>
          <App />
        </QueryClientProvider>
      </MantineProvider>
    </StrictMode>,
  );
} else {
  console.error('Root element not found');
}

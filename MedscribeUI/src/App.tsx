import { useState } from 'react';
import reactLogo from './assets/react.svg';
import viteLogo from '/vite.svg';
import { Loader } from '@mantine/core';
import './App.css';

import NoteControlsLayout from './components/NoteControls/NotecontrolsLayout';

function App() {
  return (
    <>
      <div>
        <NoteControlsLayout />
      </div>
    </>
  );
}

export default App;

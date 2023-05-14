import DarkMode from './DarkMode';
import React from 'react';
import { ThemeProvider } from '../../shared/contexts/theme-context/ThemeContext';

export default {
  title: 'Components/DarkMode',
  component: DarkMode,
};

const Template = () => {
  return (
    <ThemeProvider>
      <DarkMode />
    </ThemeProvider>
  );
};

export const Default = Template.bind({});
// Default.args = {};

import * as React from 'react';

import { FormControlLabel, FormGroup, Switch } from '@mui/material';

import { ThemeContext } from '../../shared/contexts/theme-context/ThemeContext';
import { useContext } from 'react';

export const DarkMode = () => {
  const { isDarkMode, toggleDarkMode } = useContext(ThemeContext);

  return (
    <FormGroup>
      <FormControlLabel control={<Switch checked={isDarkMode} onChange={toggleDarkMode} size="small" />} label={isDarkMode ? 'Dark' : 'Light'} />
    </FormGroup>
  );
};

export default DarkMode;

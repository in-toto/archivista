/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-empty-function */
import * as React from 'react';

import darkTheme, { lightTheme } from '../../../gatsby-theme-material-ui-top-layout/theme';

import CssBaseline from '@mui/material/CssBaseline';
import { ThemeProvider as MuiThemeProvider } from '@mui/material/styles';
import { SimpleWrapperProps } from '..';
import { createContext } from 'react';
import useLocalStorage from '../../hooks/use-local-storage/useLocalStorage';
import useMediaQuery from '@mui/material/useMediaQuery';

const ThemeContext = createContext({
  isDarkMode: true,
  toggleDarkMode: () => {},
});

const ThemeProvider = ({ children }: SimpleWrapperProps) => {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
  const [isDarkMode, setIsDarkMode] = useLocalStorage('theme.mode', prefersDarkMode);

  const toggleDarkMode = () => {
    // when the user changes their theme, because of SSG we require them to refresh the window.
    // sucks to loose navigation state and reload the app but such is static html/css
    setIsDarkMode(!isDarkMode);
    window.location.reload();
  };

  const theme = isDarkMode ? darkTheme : lightTheme;

  return (
    <ThemeContext.Provider value={{ isDarkMode, toggleDarkMode }}>
      <MuiThemeProvider theme={theme}>{children}</MuiThemeProvider>
    </ThemeContext.Provider>
  );
};

export { ThemeProvider, ThemeContext };

export const useTheme = () => React.useContext(ThemeContext);

import { createTheme } from '@mui/material';

export const lightTheme = createTheme({
  palette: {
    action: {
      active: '#38B3A1',
    },
    primary: {
      main: '#000066',
    },
    secondary: {
      main: '#38B3A1',
    },
    mode: 'light',
  },
  typography: {
    fontFamily: 'Roboto',
  },
});

export const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    secondary: {
      main: '#000066',
    },
    background: {
      default: '#000066',
    },
    primary: {
      main: '#FFFFFF',
    },
    text: {
      primary: '#FFFFFF',
      secondary: '#BDBDBD',
    },
  },
  typography: {
    fontFamily: 'Roboto',
  },
});

export default lightTheme;

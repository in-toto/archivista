import { createTheme } from "@mui/material";

const Theme = createTheme({
  palette: {
    action: {
      active: "#38B3A1",
    },
    primary: {
      main: "#000066",
    },
    secondary: {
      main: "#38B3A1",
    },
  },
  typography: {
    fontFamily: [
      "-apple-system",
      "BlinkMacSystemFont",
      '"Segoe UI"',
      "Roboto",
      '"Helvetica Neue"',
      "Arial",
      "sans-serif",
      '"Apple Color Emoji"',
      '"Segoe UI Emoji"',
      '"Segoe UI Symbol"',
    ].join(","),
  },
});

export { Theme };

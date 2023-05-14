import * as React from "react";

import { App } from "./App";
import { BrowserRouter } from "react-router-dom";
import ReactDOM from "react-dom";

const app = document.getElementById("app");
ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
  app
);

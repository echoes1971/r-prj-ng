import React from 'react';
import ReactDOM from 'react-dom/client';
import './i18n';
import App from './App';
import { ThemeProvider } from "./ThemeContext";
import { app_cfg } from './app.cfg';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import './index.css';

// Set page title from runtime config
document.title = app_cfg.site_title;

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <ThemeProvider>
      <App />
    </ThemeProvider>
  </React.StrictMode>
);

import React, { useContext, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Navbar, Nav, NavDropdown, Container, Button, Dropdown, NavItem } from "react-bootstrap";
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";
import { app_cfg } from "./app.cfg";
import axios from "./axios";

export function AppFooter() {
    // const { dark, toggleTheme } = useContext(ThemeContext);
    // const { t, i18n } = useTranslation();

    return (
      <footer className="py-3 mt-auto">
        <div className="container">
          <div className="row">
            <div className="col-md-6 text-center text-md-start">
              {app_cfg.app_name && <small>Powered by: {app_cfg.app_name} - v. {app_cfg.app_version || '1.0.0'}</small>}
            </div>
            <div className="col-md-6 text-center text-md-end">
              {app_cfg.app_copyright && <small>&copy; 2025 {app_cfg.app_copyright || 'echoes1971'}</small>}
            </div>
          </div>
        </div>
      </footer>
    );
}

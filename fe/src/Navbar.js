import React, { useContext, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Navbar, Nav, NavDropdown, Container, Button, Dropdown } from "react-bootstrap";
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";
import { app_cfg } from "./app.cfg";

function AppNavbar() {
  const navigate = useNavigate();
  const [username, setUsername] = useState(localStorage.getItem("username"));
  const { dark, toggleTheme } = useContext(ThemeContext);
  const { t, i18n } = useTranslation();
  const site_title = app_cfg.site_title;
  const groups = localStorage.getItem("groups") ? JSON.parse(localStorage.getItem("groups")) : [];
  const isAdmin = groups.includes("-2");

  const changeLanguage = (lng) => {
    i18n.changeLanguage(lng);
    localStorage.setItem("lang", lng);
  };
  const flags = {
    it: "🇮🇹",
    en: "🇬🇧",
    fr: "🇫🇷",
    de: "🇩🇪",
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    setUsername(null);
    navigate("/");
  };

  return (
    <Navbar className={dark ? "navbar bg-gradient-dark" : "navbar bg-gradient-light"} bg={dark ? "dark" : "light"} variant={dark ? "dark" : "light"} expand="lg">
      <Container>
        <Navbar.Brand as={Link} to="/">{site_title}</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />

        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="ms-auto">
            
            {username && isAdmin ? (
              <Dropdown className="me-2">
                <Dropdown.Toggle variant="outline-secondary" size="sm">
                  Admin ⚙️
                </Dropdown.Toggle>
                <Dropdown.Menu>
                  <Dropdown.Item as={Link} to="/users">{t("users.users")}</Dropdown.Item>
                  <Dropdown.Item as={Link} to="/groups">{t("groups.groups")}</Dropdown.Item>
                </Dropdown.Menu>
              </Dropdown>
            ) : null}
          
            {/* Switch Language: */}
            <Dropdown className="me-2">
              <Dropdown.Toggle variant="outline-secondary" size="sm">
                {flags[i18n.language] || "🌍"}
              </Dropdown.Toggle>
              <Dropdown.Menu>
                <Dropdown.Item onClick={() => changeLanguage("en")}>
                  🇬🇧 English
                </Dropdown.Item>
                <Dropdown.Item onClick={() => changeLanguage("fr")}>
                  🇫🇷 Français
                </Dropdown.Item>
                <Dropdown.Item onClick={() => changeLanguage("de")}>
                  🇩🇪 Deutsch
                </Dropdown.Item>
                <Dropdown.Item onClick={() => changeLanguage("it")}>
                  🇮🇹 Italiano
                </Dropdown.Item>
              </Dropdown.Menu>
            </Dropdown>

            {/* Bottone toggle tema */}
            <Button
              variant={dark ? "secondary" : "outline-secondary"}
              className="ms-2"
              onClick={toggleTheme}
            >
              {dark ? <i className="bi bi-sun"></i> : <i className="bi bi-moon"></i>}
            </Button>
          </Nav>
        </Navbar.Collapse>

      </Container>
    </Navbar>
  );
}

export default AppNavbar;

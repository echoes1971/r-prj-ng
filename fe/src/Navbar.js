import React, { useContext, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Navbar, Nav, NavDropdown, Container, Button, Dropdown, NavItem } from "react-bootstrap";
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";
import { app_cfg } from "./app.cfg";
import axios from "./axios";

function AppNavbar() {
  const navigate = useNavigate();
  const [username, setUsername] = useState(localStorage.getItem("username"));
  const [searchVisible, setSearchVisible] = useState(false);
  const [searchText, setSearchText] = useState('');
  const { dark, toggleTheme } = useContext(ThemeContext);
  const { t, i18n } = useTranslation();
  const site_title = app_cfg.site_title;
  const site_root = app_cfg.app_home_object_id;
  const [children, setChildren] = useState([]);
  const groups = localStorage.getItem("groups") ? JSON.parse(localStorage.getItem("groups")) : [];
  const isAdmin = groups.includes("-2");
  const isWebmaster = groups.includes(app_cfg.webmaster_group_id);

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

  // Load root children
  useEffect(() => {
    const loadChildren = async () => {
      try {
        const response = await axios.get(`/nav/children/${site_root}`);
        // filter those with metadata DBFolder
        const filteredChildren = (response.data.children || []).filter(child => child.metadata && child.metadata.classname === "DBFolder");
        // alert(JSON.stringify(filteredChildren));
        setChildren(filteredChildren);
      } catch (error) {
        console.error("Error loading root children:", error);
      }
    };
    loadChildren();
  }, [site_root]);

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    setUsername(null);
    navigate("/");
  };

  const handleProfileClick = async () => {
    try {
      const userId = localStorage.getItem("user_id");
      const response = await axios.get(`/users/${userId}/person`);
      const { person_id } = response.data;
      navigate(`/c/${person_id}`);
    } catch (error) {
      console.error("Error fetching person:", error);
    }
  };

  const handleSearchSubmit = (e) => {
    e.preventDefault();
    if (searchText.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchText.trim())}`);
      setSearchText('');
      setSearchVisible(false);
    }
  };

  return (
    <Navbar className={dark ? "navbar bg-gradient-dark" : "navbar bg-gradient-light"} bg={dark ? "dark" : "light"} variant={dark ? "dark" : "light"} expand="lg">
      <Container>
        <Navbar.Brand as={Link} to="/">{site_title}</Navbar.Brand>
        
        <div className="d-flex d-lg-none ms-auto">
          <Button
            variant={dark ? "secondary" : "outline-secondary"}
            size="sm"
            onClick={() => navigate('/search')}
            className="me-2"
          >
            <i className="bi bi-search"></i>
          </Button>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
        </div>

        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="ms-auto">
            {/* Display non public folders when a user is logged in */}
            {children.map(child => (!username || child.data.permissions.slice(-3)==='---') && (
              <Nav.Link as={Link} key={child.data.id} to={`/c/${child.data.id}`}>
                {child.data.name}
              </Nav.Link>
            ))}

            {/* Search toggle and field */}
            {searchVisible ? (
              <form onSubmit={handleSearchSubmit} className="d-flex align-items-center me-2">
                <input
                  type="text"
                  className="form-control form-control-sm"
                  placeholder={t('common.search_placeholder') || 'Search...'}
                  value={searchText}
                  onChange={(e) => setSearchText(e.target.value)}
                  autoFocus
                  style={{ width: '200px' }}
                />
                <Button
                  variant="link"
                  size="sm"
                  onClick={() => {
                    setSearchVisible(false);
                    setSearchText('');
                  }}
                  className="text-secondary ms-1"
                >
                  <i className="bi bi-x-lg"></i>
                </Button>
              </form>
            ) : (
              <Button
                variant={dark ? "secondary" : "outline-secondary"}
                size="sm"
                onClick={() => setSearchVisible(true)}
                className="me-2 d-none d-lg-inline"
              >
                <i className="bi bi-search"></i>
              </Button>
            )}

            {username && isWebmaster ? (
              <NavDropdown title="Webmaster 🛠️" id="webmaster-nav-dropdown" align="end" {...(dark ? { menuVariant: 'dark' } : {})}>
                <NavDropdown.Item as={Link} to="/">{t("common.home")}</NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/folders">{t("folder.folders")}</NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/pages">{t("page.pages")}</NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/news">{t("news.news")}</NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/files">{t("files.files")}</NavDropdown.Item>
              </NavDropdown>
            ) : null}
            {username && (
              <NavDropdown title={t("common.contacts")} id="webmaster-nav-dropdown" align="end" {...(dark ? { menuVariant: 'dark' } : {})}>
                <NavDropdown.Item as={Link} to="/companies">{t("company.companies")}</NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/people">{t("person.people")}</NavDropdown.Item>
              </NavDropdown>
            )}
            {username && isAdmin ? (
              <NavDropdown title="Admin ⚙️" id="admin-nav-dropdown" align="end" {...(dark ? { menuVariant: 'dark' } : {})}>
                <NavDropdown.Item as={Link} to="/dashboard" disabled>{t("common.dashboard")}</NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item as={Link} to="/users">{t("users.users")}</NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/groups">{t("groups.groups")}</NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item as={Link} to="/db" disabled>{t("common.db")}</NavDropdown.Item>
                <NavDropdown.Item as={Link} to="/log" disabled>{t("common.log")}</NavDropdown.Item>
              </NavDropdown>
            ) : null}

          
            {!username ? (
              <Button as={Link} to="/login" variant={dark ? "secondary" : "outline-secondary"}>
                {t("common.login")}
              </Button>
            ) : null}

            {username ? (
              <NavDropdown title={username} id="basic-nav-dropdown" align="end" {...(dark ? { menuVariant: 'dark' } : {})}>
                <NavDropdown.Item onClick={handleProfileClick}>
                  <i className="bi bi-person-circle me-2"></i>{t("users.user_profile")}
                </NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item as={Link} to="/settings" disabled>
                  <i className="bi bi-gear me-2"></i>{t("common.settings")}
                </NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item onClick={handleLogout}>{t("common.logout")}</NavDropdown.Item>
              </NavDropdown>
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

import React, { useContext, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Navbar, Nav, NavDropdown, Container, Button, Dropdown, NavItem } from "react-bootstrap";
import { ThemeContext } from "../ThemeContext";
import { useTranslation } from "react-i18next";
import { app_cfg } from "../app.cfg";
import axios from "../axios";

export function AdminDashboard() {
    const { t } = useTranslation();
    const [errorMessage, setErrorMessage] = useState("");
    const { dark, themeClass } = useContext(ThemeContext);

    const [userCount, setUserCount] = useState(0);
    const [groupCount, setGroupCount] = useState(0);
    const [storageUsage, setStorageUsage] = useState(0);

    useEffect(() => {
        // Fetch users number from the API
        axios.get('/users')
            .then(response => {
                setUserCount(response.data.length);
            })
            .catch(error => {
                console.error("There was an error fetching the users!", error);
            });
        axios.get('/groups')
            .then(response => {
                setGroupCount(response.data.length);
            })
            .catch(error => {
                console.error("There was an error fetching the groups!", error);
            });
        axios.get('/files/storage-usage')
            .then(response => {
                setStorageUsage(response.data.usage);
            })
            .catch(error => {
                console.error("There was an error fetching the storage usage!", error);
            });
    }, []);

    // Fetch users number



    return (
        <div className={`container mt-3 ${themeClass}`}>
          <h2 className={dark ? "text-light" : "text-dark"}>{t("admin.dashboard")}</h2>
            <p>Welcome to the admin dashboard. Here you can manage the application settings and monitor system status.</p>
            <p>{t("admin.total_users")}: {userCount}</p>
            <p>{t("admin.total_groups")}: {groupCount}</p>
        </div>
    );
}

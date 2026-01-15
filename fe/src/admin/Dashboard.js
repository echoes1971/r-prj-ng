import React, { useContext, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Navbar, Nav, NavDropdown, Container, Button, Dropdown, NavItem } from "react-bootstrap";
import { ThemeContext } from "../ThemeContext";
import { useTranslation } from "react-i18next";
import { app_cfg } from "../app.cfg";
import axios from "../axios";
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip, BarChart, Bar, XAxis, YAxis, CartesianGrid } from "recharts";

// Colori per le slice della torta
const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d', '#ffc658', '#ff7c7c'];

export function AdminDashboard() {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const [errorMessage, setErrorMessage] = useState("");
    const { dark, themeClass } = useContext(ThemeContext);

    const [stats, setStats] = useState({});

    useEffect(() => {
        // retrieve statistics from the API
        const fetchStats = async () => {
            const token = localStorage.getItem("token");
            try {
                const res = await axios.get("/admin/dashboard", {
                    headers: { Authorization: `Bearer ${token}` },
                });
                console.log("AdminDashboard: fetched user statistics:", res.data);
                setStats(res.data || {});
            } catch (err) {
                console.log("Error loading user statistics:", err);
            }
        };
        fetchStats();
    }, []);

    // Prepara i dati per il pie chart: trasforma object_stats in array
    const pieData = stats.object_stats ? Object.keys(stats.object_stats).map(className => ({
        name: className,
        value: stats.object_stats[className] !== null && stats.object_stats[className].count || 0
    })).filter(item => item.value > 0) : []; // mostra solo tipi con almeno 1 oggetto

    // Prepara i dati per active users bar chart
    const activeUsersData = stats.users_stats ? [
        { period: t('admin.last_24h'), users: stats.users_stats.active_last_24h || 0 },
        { period: t('admin.last_7_days'), users: stats.users_stats.active_last_7_days || 0 },
        { period: t('admin.last_30_days'), users: stats.users_stats.active_last_30_days || 0 }
    ] : [];

    // Calcola totali per object activity (created, modified, deleted last week)
    const objectActivityData = stats.object_stats ? [
        {
            action: t('admin.created'),
            count: Object.values(stats.object_stats).reduce((sum, obj) => sum + (obj?.created_last_week || 0), 0),
            fill: '#00C49F' // verde
        },
        {
            action: t('admin.modified'),
            count: Object.values(stats.object_stats).reduce((sum, obj) => sum + (obj?.modified_last_week || 0), 0),
            fill: '#0088FE' // blu
        },
        {
            action: t('admin.deleted'),
            count: Object.values(stats.object_stats).reduce((sum, obj) => sum + (obj?.deleted_count || 0), 0),
            fill: '#FF8042' // rosso
        }
    ] : [];

    // Prepara i dati per il pie chart dei gruppi (users per gruppo)
    const groupsPieData = stats.groups_stats ? Object.keys(stats.groups_stats).map(groupName => ({
        name: groupName,
        value: stats.groups_stats[groupName] || 0
    })).filter(item => item.value > 0) : [];

    return (
        <div className={`container mt-3 ${themeClass}`}>
          <h2 className={dark ? "text-light" : "text-dark"}>{t("admin.dashboard")}</h2>
            <p>{t("admin.welcome_message")}</p>

            {/* Active Users Cards */}
            {stats.users_stats && (
                <div className="row mb-4">
                    <div className="col-12 mb-3">
                        <h4 className={dark ? "text-light" : "text-dark"}>{t("admin.active_users")}</h4>
                    </div>
                    <div className="col-md-4 mb-3">
                        <div className={`card text-center ${dark ? 'bg-dark text-light' : 'bg-light'}`}>
                            <div className="card-body">
                                <h6 className="card-subtitle mb-2 text-secondary">{t("admin.last_24h")}</h6>
                                <h2 className="card-title" style={{fontSize: '3rem', color: '#0088FE'}}>
                                    {stats.users_stats.active_last_24h || 0}
                                </h2>
                            </div>
                        </div>
                    </div>
                    <div className="col-md-4 mb-3">
                        <div className={`card text-center ${dark ? 'bg-dark text-light' : 'bg-light'}`}>
                            <div className="card-body">
                                <h6 className="card-subtitle mb-2 text-secondary">{t("admin.last_7_days")}</h6>
                                <h2 className="card-title" style={{fontSize: '3rem', color: '#00C49F'}}>
                                    {stats.users_stats.active_last_7_days || 0}
                                </h2>
                            </div>
                        </div>
                    </div>
                    <div className="col-md-4 mb-3">
                        <div className={`card text-center ${dark ? 'bg-dark text-light' : 'bg-light'}`}>
                            <div className="card-body">
                                <h6 className="card-subtitle mb-2 text-secondary">{t("admin.last_30_days")}</h6>
                                <h2 className="card-title" style={{fontSize: '3rem', color: '#FFBB28'}}>
                                    {stats.users_stats.active_last_30_days || 0}
                                </h2>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {/* Bar Charts Row - Side by side on desktop */}
            <div className="row mb-4">
                {/* Pie Chart per distribuzione oggetti */}
                {pieData.length > 0 && (
                    <div className="col-12 col-lg-6">
                        <div className="col-12">
                            <h4 className={dark ? "text-light" : "text-dark"}>{t("admin.objects_distribution")}</h4>
                            <ResponsiveContainer width="100%" height={400}>
                                <PieChart>
                                    <Pie
                                        data={pieData}
                                        cx="50%"
                                        cy="50%"
                                        labelLine={false}
                                        label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                                        outerRadius={180}
                                        fill="#8884d8"
                                        dataKey="value"
                                    >
                                        {pieData.map((entry, index) => (
                                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                        ))}
                                    </Pie>
                                    <Tooltip />
                                </PieChart>
                            </ResponsiveContainer>
                        </div>
                    </div>
                )}

                {/* Object Activity Bar Chart (Last Week) */}
                {objectActivityData.length > 0 && (
                    <div className="col-12 col-lg-6">
                        <h4 className={dark ? "text-light" : "text-dark"}>{t("admin.object_activity")}</h4>
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={objectActivityData}>
                                <CartesianGrid strokeDasharray="3 3" />
                                <XAxis dataKey="action" />
                                <YAxis />
                                <Tooltip />
                                <Bar dataKey="count" />
                            </BarChart>
                        </ResponsiveContainer>
                    </div>
                )}
            </div>

            <div className="row mb-4">
                {groupsPieData.length > 0 && (
                    <div className="col-12 col-lg-6 mb-4 mb-lg-0">
                        <h4 className={dark ? "text-light" : "text-dark"}>{t("admin.groups_distribution")}</h4>
                        <ResponsiveContainer width="100%" height={400}>
                            <PieChart>
                                <Pie
                                    data={groupsPieData}
                                    cx="50%"
                                    cy="50%"
                                    labelLine={false}
                                    label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                                    outerRadius={180}
                                    fill="#8884d8"
                                    dataKey="value"

                                    onClick={(data) => {
                                        console.log(data.payload);
                                        navigate(`/groups/${data.payload.name}`);
                                    }}
                                >
                                    {groupsPieData.map((entry, index) => (
                                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                    ))}
                                </Pie>
                                <Tooltip />
                            </PieChart>
                        </ResponsiveContainer>
                    </div>
                )}

                {/* Active Users Bar Chart */}
                {activeUsersData.length > 0 && (
                    <div className="col-12 col-lg-6 mb-4 mb-lg-0">
                        <h4 className={dark ? "text-light" : "text-dark"}>{t("admin.active_users")}</h4>
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={activeUsersData}>
                                <CartesianGrid strokeDasharray="3 3" />
                                <XAxis dataKey="period" />
                                <YAxis />
                                <Tooltip />
                                <Bar dataKey="users" fill="#0088FE" />
                            </BarChart>
                        </ResponsiveContainer>
                    </div>
                )}
            </div>

            {/* Groups Distribution Pie Chart */}
            <div className="row mb-4">
                    <div className="col-md-4 mb-3">
                        <div className={`card text-center ${dark ? 'bg-dark text-light' : 'bg-light'}`}>
                            <div className="card-body">
                                <h6 className="card-subtitle mb-2 text-secondary">{t("admin.total_users")}</h6>
                                <h2 className="card-title" style={{fontSize: '3rem', color: '#0088FE'}}>
                                    {stats.users_count || 0}
                                </h2>
                            </div>
                        </div>
                    </div>
                    <div className="col-md-4 mb-3">
                        <div className={`card text-center ${dark ? 'bg-dark text-light' : 'bg-light'}`}>
                            <div className="card-body">
                                <h6 className="card-subtitle mb-2 text-secondary">{t("admin.total_groups")}</h6>
                                <h2 className="card-title" style={{fontSize: '3rem', color: '#00C49F'}}>
                                    {stats.groups_count || 0}
                                </h2>
                            </div>
                        </div>
                    </div>
            </div>
        </div>
    );
}

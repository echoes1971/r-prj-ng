import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AppNavbar from "./Navbar";
import './App.css';

import { app_cfg } from './app.cfg';

import DefaultPage from "./DefaultPage";
import Login from "./Login";
import Users from "./Users";
import Groups from './Groups';
import SiteNavigation from './SiteNavigation';

function App() {
  const token = localStorage.getItem("token");
  const groups = localStorage.getItem("groups") ? JSON.parse(localStorage.getItem("groups")) : [];
  const isAdmin = groups.includes("-2");
  
  return (
    <Router>
      <AppNavbar />
      <Routes>
        <Route path="/" element={<SiteNavigation />} />
        <Route path="/default" element={<DefaultPage />} />
        <Route path="/login" element={<Login />} />

        {/* Site Navigation - content by object ID */}
        <Route path="/c/:objectId" element={<SiteNavigation />} />
        <Route path="/c" element={<SiteNavigation />} />

        {/* Protected routes - only for admins (group -2) */}
        <Route
          path="/users"
          element={token && isAdmin ? <Users /> : <Navigate to="/" />}
        />

        <Route
          path="/groups"
          element={token && isAdmin ? <Groups /> : <Navigate to="/" />}
        />

        {/* Default -> redirect to / */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
}

export default App;

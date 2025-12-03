import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AppNavbar from "./Navbar";
import './App.css';

import { app_cfg } from './app.cfg';

import DefaultPage from "./DefaultPage";
import Login from "./Login";
import Users from "./Users";
import UserProfile from "./UserProfile";
import Groups from './Groups';
import GroupProfile from './GroupProfile';
import SiteNavigation from './SiteNavigation';
import ContentEdit from './ContentEdit';

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

        {/* Content Edit - edit object by ID (requires authentication) */}
        <Route path="/e/:id" element={token ? <ContentEdit /> : <Navigate to="/c/:id" />} />

        {/* User profile - accessible by the user themselves or admins */}
        <Route
          path="/users/:userId"
          element={token ? <UserProfile /> : <Navigate to="/login" />}
        />

        {/* Group profile - only for admins */}
        <Route
          path="/groups/:groupId"
          element={token ? <GroupProfile /> : <Navigate to="/" />}
        />

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
        <Route path="*" element={<Navigate to="/default" />} />
      </Routes>
    </Router>
  );
}

export default App;

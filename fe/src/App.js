import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AppNavbar from "./Navbar";
import './App.css';

import { app_cfg } from './app.cfg';

import DefaultPage from "./DefaultPage";
import Login from "./Login";
import Users from "./Users";
import Groups from './groups';

function App() {
  const token = localStorage.getItem("token");
  const group_ids = localStorage.getItem("group_ids") ? localStorage.getItem("group_ids").split(",") : [];
  return (
    <Router>
      <AppNavbar />
      <Routes>
        <Route path="/" element={<DefaultPage />} />
        <Route path="/login" element={<Login />} />

        {/* Protected route */}
        <Route
          path="/users"
          element={token ? <Users /> : <Navigate to="/" />}
        />

        <Route
          path="/groups"
          element={token ? <Groups /> : <Navigate to="/" />}
        />

        {/* Default -> redirect to / */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
}

export default App;

import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AppNavbar from "./Navbar";
// import logo from './logo.svg';
import './App.css';
// import { ThemeProvider, ThemeContext } from "./ThemeContext";

import Login from "./Login";
import Users from "./Users";

function App() {
  const token = localStorage.getItem("token");
  return (
    <Router>
      <AppNavbar />
      <Routes>
        {/* Rotta login sempre accessibile */}
        <Route path="/login" element={<Login />} />

        {/* Rotta utenti protetta */}
        <Route
          path="/users"
          element={token ? <Users /> : <Navigate to="/login" />}
        />

        {/* Rotta di default → redirect a /login */}
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    </Router>
  );
}

      // <div className={themeClass}>
      //   <button className="btn btn-secondary" onClick={toggleTheme}>
      //     Toggle Tema
      //   </button>
      //   {/* router e pagine */}
      // </div>

export default App;

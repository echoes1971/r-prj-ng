import { createContext, useEffect, useState } from "react";

export const ThemeContext = createContext();

export function ThemeProvider({ children }) {
  // Load theme from localStorage on init
  const [dark, setDark] = useState(() => {
    const savedTheme = localStorage.getItem("theme");
    return savedTheme === "dark";
  });
  
  const toggleTheme = () => setDark(!dark);

  const themeClass = dark ? "bg-dark text-light" : "bg-light text-dark";

  // Save theme to localStorage and update body class
  useEffect(() => {
    document.body.className = dark ? "dark-theme" : "light-theme";
    localStorage.setItem("theme", dark ? "dark" : "light");
  }, [dark]);

  return (
    <ThemeContext.Provider value={{ dark, toggleTheme, themeClass }}>
      {children}
    </ThemeContext.Provider>
  );
}

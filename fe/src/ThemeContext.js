import { createContext, useState } from "react";

export const ThemeContext = createContext();

export function ThemeProvider({ children }) {
  const [dark, setDark] = useState(false);
  const toggleTheme = () => setDark(!dark);

  const themeClass = dark ? "bg-dark text-light" : "bg-light text-dark";

  return (
    <ThemeContext.Provider value={{ dark, toggleTheme, themeClass }}>
      {children}
    </ThemeContext.Provider>
  );
}

import React, { useContext, useState } from "react";
import api from "./axios";
import { useTranslation } from "react-i18next";
import { ThemeContext } from "./ThemeContext";

function Login() {
  const { t } = useTranslation();
  const [login, setLogin] = useState("");
  const [pwd, setPwd] = useState("");
  const { dark, themeClass } = useContext(ThemeContext);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const res = await api.post("/login", { login, pwd });
      localStorage.setItem("token", res.data.access_token);
      // store username for Navbar display
      localStorage.setItem("username", login);
      localStorage.setItem("groups", JSON.stringify(res.data.groups));
      window.location.href = "/users";
    } catch (err) {
      alert("Login failed");
    }
  };

  return (
    <div className={`container mt-3  ${themeClass}`}>
      <form onSubmit={handleSubmit} className="p-3">
        <div className="form-group row">
          <label class="col-md-1 col-form-label">Login</label>
          <div class="col-sm-3">
            <input className="form-control mb-2" placeholder="Login" value={login} onChange={e => setLogin(e.target.value)} />
          </div>
        </div>
        <div class="form-group row">
          <label class="col-md-1 col-form-label">Password</label>
          <div class="col-sm-3">
            <input className="form-control mb-2" type="password" placeholder="Password" value={pwd} onChange={e => setPwd(e.target.value)} />
          </div>
        </div>
        <button className="btn btn-primary">{t("common.login")}</button>
      </form>
    </div>
  );
}

export default Login;

import React, { useContext, useState } from "react";
import api from "./axios";
import { useTranslation } from "react-i18next";
import { ThemeContext } from "./ThemeContext";
import { getErrorMessage } from "./errorHandler";

function Login() {
  const { t } = useTranslation();
  const [login, setLogin] = useState("");
  const [pwd, setPwd] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const { dark, themeClass } = useContext(ThemeContext);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const res = await api.post("/login", { login, pwd });
      localStorage.setItem("token", res.data.access_token);
      localStorage.setItem("expires_at", res.data.expires_at);
      const expiryDate = new Date(Number(res.data.expires_at)*1000);
      const now = new Date();
      const timeDiff = new Date(145*1000); //new Date(expiryDate - now);
      // alert("Expires at: " + expiryDate.toString()+"\nNow: " + now.toString()+"\nTime remaining (ms): " + timeDiff.getUTCHours()+":" + timeDiff.getUTCMinutes()+":" + timeDiff.getUTCSeconds() );
      // store username for Navbar display
      localStorage.setItem("username", login);
      localStorage.setItem("user_id", res.data.user_id);
      localStorage.setItem("groups", JSON.stringify(res.data.groups));
      window.location.href = "/";
    } catch (err) {
      const errorMsg = getErrorMessage(err, t("common.login_failed") || "Login failed");
      setErrorMessage(errorMsg);
    }
  };

  return (
    // Center horizontally and vertically, expand fields to reasonable size
    <div className={`container mt-3 d-flex justify-content-center align-items-center ${themeClass}`}>
      <form onSubmit={handleSubmit} className="p-3">
        {errorMessage && (
          <div className="alert alert-danger" role="alert">
            {errorMessage}
          </div>
        )}
        <div className="form-group row">
          <label className="col-md-4 col-form-label text-md-end">Login</label>
          <div className="col-md-8 col-sm-3">
            <input className="form-control mb-2" placeholder="Login" value={login} onChange={e => setLogin(e.target.value)} />
          </div>
        </div>
        <div className="form-group row">
          <label className="col-md-4 col-form-label text-md-end">Password</label>
          <div className="col-md-8 col-sm-3">
            <input className="form-control mb-2" type="password" placeholder="Password" value={pwd} onChange={e => setPwd(e.target.value)} />
          </div>
        </div>
        <div className="form-group row">
          <div className="col-md-4"></div>
          <div className="col-md-8">
            <button className="btn btn-primary">{t("common.login")}</button>
          </div>
        </div>
      </form>
    </div>
  );
}

export default Login;

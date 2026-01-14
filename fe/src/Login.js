import React, { useContext, useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import api from "./axios";
import { useTranslation } from "react-i18next";
import { ThemeContext } from "./ThemeContext";
import { getErrorMessage } from "./errorHandler";
import { app_cfg } from "./app.cfg";

function Login() {
  const { t } = useTranslation();
  const [login, setLogin] = useState("");
  const [pwd, setPwd] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const { dark, themeClass } = useContext(ThemeContext);
  const navigate = useNavigate();

  useEffect(() => {
    function onMessage(ev) {
      try {
        const data = ev.data;
        if (!data || !data.access_token) return;
        localStorage.setItem('token', data.access_token);
        localStorage.setItem('expires_at', data.expires_at);
        if (data.user_id) localStorage.setItem('user_id', data.user_id);
        if (data.login) localStorage.setItem('username', data.login);
        if (data.groups) localStorage.setItem('groups', JSON.stringify((data.groups+"").split(',')));
        navigate('/');
      } catch (e) {
        // ignore
      }
    }
    window.addEventListener('message', onMessage);
    return () => { window.removeEventListener('message', onMessage); };
  }, [navigate]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const res = await api.post("/login", { login, pwd });
      localStorage.setItem("token", res.data.access_token);
      localStorage.setItem("expires_at", res.data.expires_at);
      const expiryDate = new Date(Number(res.data.expires_at)*1000);
      const now = new Date();
      const timeDiff = new Date(145*1000); //new Date(expiryDate - now);
      // console.log("Expires at: " + expiryDate.toString()+"\nNow: " + now.toString()+"\nTime remaining (ms): " + timeDiff.getUTCHours()+":" + timeDiff.getUTCMinutes()+":" + timeDiff.getUTCSeconds() );
      // store username for Navbar display
      localStorage.setItem("username", login);
      localStorage.setItem("user_id", res.data.user_id);
      localStorage.setItem("groups", JSON.stringify(res.data.groups));
      // go to the main page and refresh to load user-specific data
      navigate("/");
      window.location.reload();
    } catch (err) {
      const errorMsg = getErrorMessage(err, t("common.login_failed") || "Login failed");
      setErrorMessage(errorMsg);
    }
  };

  return (
    // Center horizontally and vertically, expand fields to reasonable size
    <div className={`container mt-2 mt-md-5 d-flex justify-content-center align-items-center ${themeClass}`}>
      <form onSubmit={handleSubmit} className="p-3">
        {errorMessage && (
          <div className="alert alert-danger" role="alert">
            {errorMessage}
          </div>
        )}
        <div className="form-group row">
          <label className="col-md-4 col-form-label text-md-end">Login</label>
          <div className="col-md-8">
            <input className="form-control mb-2" placeholder="Login" value={login} onChange={e => setLogin(e.target.value)} />
          </div>
        </div>
        <div className="form-group row">
          <label className="col-md-4 col-form-label text-md-end">Password</label>
          <div className="col-md-8">
            <input className="form-control mb-2" type="password" placeholder="Password" value={pwd} onChange={e => setPwd(e.target.value)} />
          </div>
        </div>
        <div className="form-group row">
          <div className="col-md-4"></div>
          <div className="col-md-8">
            <button className="btn btn-primary">{t("common.login")}</button>
            </div>
        </div>
            { (app_cfg.enable_google_oauth === 'true' || app_cfg.enable_google_oauth === true || app_cfg.enable_google_oauth === '1') && (
              <div className="form-group row mt-3 m-auto">
                <div className="col-md-12 text-center">
                    <button type="button" className="btn btn-outline-secondary" onClick={() => {
                      const startUrl = `${app_cfg.endpoint}/oauth/google/start`;
                      window.location.href = startUrl;
                    }}>
                      <svg className="me-2" style={{height: '1.5rem'}} xmlns="https://www.w3.org/2000/svg" viewBox="0 0 48 48"><path fill="#4285F4" d="M45.12 24.5c0-1.56-.14-3.06-.4-4.5H24v8.51h11.84c-.51 2.75-2.06 5.08-4.39 6.64v5.52h7.11c4.16-3.83 6.56-9.47 6.56-16.17z"></path><path fill="#34A853" d="M24 46c5.94 0 10.92-1.97 14.56-5.33l-7.11-5.52c-1.97 1.32-4.49 2.1-7.45 2.1-5.73 0-10.58-3.87-12.31-9.07H4.34v5.7C7.96 41.07 15.4 46 24 46z"></path><path fill="#FBBC05" d="M11.69 28.18C11.25 26.86 11 25.45 11 24s.25-2.86.69-4.18v-5.7H4.34C2.85 17.09 2 20.45 2 24c0 3.55.85 6.91 2.34 9.88l7.35-5.7z"></path><path fill="#EA4335" d="M24 10.75c3.23 0 6.13 1.11 8.41 3.29l6.31-6.31C34.91 4.18 29.93 2 24 2 15.4 2 7.96 6.93 4.34 14.12l7.35 5.7c1.73-5.2 6.58-9.07 12.31-9.07z"></path><path fill="none" d="M2 2h44v44H2z"></path></svg>
                      {t("auth.sign_in_with_google") || 'Sign in with Google'}
                    </button>
                </div>
              </div>
            ) }
            { (app_cfg.enable_github_oauth === 'true' || app_cfg.enable_github_oauth === true || app_cfg.enable_github_oauth === '1') && (
              <div className="form-group row mt-3 m-auto">
                <div className="col-md-12 text-center">
                    <button type="button" className="btn btn-outline-secondary" onClick={() => {
                        const startUrl = `${app_cfg.endpoint}/oauth/github/start`;
                        window.location.href = startUrl;
                      }}>
                      <svg className="me-2" xmlns="https://www.w3.org/2000/svg" viewBox="0 0 22 22" style={{height: '1.5rem'}}><path d="M12 1C5.923 1 1 5.923 1 12c0 4.867 3.149 8.979 7.521 10.436.55.096.756-.233.756-.522 0-.262-.013-1.128-.013-2.049-2.764.509-3.479-.674-3.699-1.292-.124-.317-.66-1.293-1.127-1.554-.385-.207-.936-.715-.014-.729.866-.014 1.485.797 1.691 1.128.99 1.663 2.571 1.196 3.204.907.096-.715.385-1.196.701-1.471-2.448-.275-5.005-1.224-5.005-5.432 0-1.196.426-2.186 1.128-2.956-.111-.275-.496-1.402.11-2.915 0 0 .921-.288 3.024 1.128a10.193 10.193 0 0 1 2.75-.371c.936 0 1.871.123 2.75.371 2.104-1.43 3.025-1.128 3.025-1.128.605 1.513.221 2.64.111 2.915.701.77 1.127 1.747 1.127 2.956 0 4.222-2.571 5.157-5.019 5.432.399.344.743 1.004.743 2.035 0 1.471-.014 2.654-.014 3.025 0 .289.206.632.756.522C19.851 20.979 23 16.854 23 12c0-6.077-4.922-11-11-11Z"></path></svg>
                      {t("auth.sign_in_with_github") || 'Sign in with GitHub'}
                    </button>
                </div>
              </div>
            ) }
            { (app_cfg.enable_telegram_oauth === 'true' || app_cfg.enable_telegram_oauth === true || app_cfg.enable_telegram_oauth === '1') && app_cfg.telegram_bot_id && (
              <div className="form-group row mt-3 m-auto">
                <div className="col-md-12 text-center">
                    <button type="button" className="btn btn-outline-secondary" onClick={() => {
                        const callbackUrl = encodeURIComponent(`${window.location.origin}${app_cfg.endpoint}/oauth/telegram/callback`);
                        const telegramAuthUrl = `https://oauth.telegram.org/auth?bot_id=${app_cfg.telegram_bot_id}&origin=${encodeURIComponent(window.location.origin)}&request_access=write&return_to=${callbackUrl}`;
                        window.location.href = telegramAuthUrl;
                      }}>
                      <svg className="me-2" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" style={{height: '1.5rem'}} fill="#0088cc"><path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221l-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.446 1.394c-.14.18-.357.295-.6.295-.002 0-.003 0-.005 0l.213-3.054 5.56-5.022c.24-.213-.054-.334-.373-.121l-6.869 4.326-2.96-.924c-.64-.203-.658-.64.135-.954l11.566-4.458c.538-.196 1.006.128.832.941z"/></svg>
                      {t("auth.sign_in_with_telegram") || 'Sign in with Telegram'}
                    </button>
                </div>
              </div>
            ) }
      </form>
    </div>
  );
}

export default Login;

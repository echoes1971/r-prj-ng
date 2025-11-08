import React, { useState } from "react";
import axios from "axios";

function Login() {
  const [login, setLogin] = useState("");
  const [pwd, setPwd] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      // const res = await axios.post("http://localhost:1971/login", { login, pwd });
      const res = await axios.post("/login", { login, pwd });
      localStorage.setItem("token", res.data.access_token);
      // store username for Navbar display
      localStorage.setItem("username", login);
      // redirect to /users
      window.location.href = "/users";
    } catch (err) {
      alert("Login fallito");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="p-3">
      <input className="form-control mb-2" placeholder="Login" value={login} onChange={e => setLogin(e.target.value)} />
      <input className="form-control mb-2" type="password" placeholder="Password" value={pwd} onChange={e => setPwd(e.target.value)} />
      <button className="btn btn-primary">Login</button>
    </form>
  );
}

export default Login;

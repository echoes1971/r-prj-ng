import axios from "axios";
import { app_cfg } from "./app.cfg";

const endpoint = app_cfg.endpoint;

// crea un'istanza
const api = axios.create({
  baseURL: endpoint, // o la tua API base
});

// interceptor di richiesta - aggiunge il token se presente
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// interceptor di risposta
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("username");
      // setUsername(null);
      // redirect globale
      window.location.href = "/"; // oppure "/login"
    }
    return Promise.reject(error);
  }
);

export default api;

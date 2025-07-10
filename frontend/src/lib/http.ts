import axios from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL, // set in .env.local
  withCredentials: true,                // send/receive cookie
});

export default api;

import { useEffect, useState } from "react";
import { Outlet, Navigate } from "react-router-dom";
import api from "../lib/http";

export default function Protected() {
  const [ok, setOk] = useState<null | boolean>(null);

  useEffect(() => {
    api.get("/auth/refresh")
       .then(() => setOk(true))
       .catch(() => setOk(false));
  }, []);

  if (ok === null) return <div className="p-8">Loading…</div>;
  return ok ? <Outlet /> : <Navigate to="/login" replace />;
}

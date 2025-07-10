import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import api from "../lib/http";

export default function Login() {
  const nav = useNavigate();
  const [email, setEmail] = useState("");
  const [pw, setPw] = useState("");
  const [err, setErr] = useState("");

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post("/auth/login", { email, password: pw });
      nav("/");
    } catch (e: any) {
      setErr(e?.response?.data?.message || "Login failed");
    }
  };

  return (
    <div className="h-screen flex items-center justify-center">
      <form onSubmit={submit} className="w-80 space-y-4">
        <h1 className="text-2xl font-semibold">Sign in</h1>
        {err && <p className="text-red-600">{err}</p>}
        <input
          className="w-full border p-2 rounded"
          placeholder="Email"
          value={email}
          onChange={e => setEmail(e.target.value)}
        />
        <input
          className="w-full border p-2 rounded"
          type="password"
          placeholder="Password"
          value={pw}
          onChange={e => setPw(e.target.value)}
        />
        <button className="w-full bg-blue-600 text-white py-2 rounded">
          Login
        </button>
        <p className="text-sm text-center">
          New? <Link className="text-blue-600" to="/signup">Create account →</Link>
        </p>
      </form>
    </div>
  );
}

import { useState } from "react";
import { Link } from "react-router-dom";
import { useLogin } from "../hooks/authHooks";

export default function Login() {
  const login = useLogin();
  const [email, setEmail] = useState("");
  const [pw, setPw] = useState("");

  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    login.mutate({ email, password: pw });
  };

  return (
    <div className="h-screen flex items-center justify-center">
      <form onSubmit={submit} className="w-80 space-y-4">
        <h1 className="text-2xl font-semibold">Sign in</h1>

        <input
          className="w-full border p-2 rounded"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          className="w-full border p-2 rounded"
          type="password"
          placeholder="Password"
          value={pw}
          onChange={(e) => setPw(e.target.value)}
        />

        <button
          className="w-full bg-blue-600 text-white py-2 rounded disabled:opacity-50"
          disabled={login.isPending}
        >
          {login.isPending ? "…" : "Login"}
        </button>

        <p className="text-sm text-center">
          New? 
          <Link className="text-blue-600" to="/signup">
            Create account →
          </Link>
        </p>
      </form>
    </div>
  );
}

import { useState } from "react";
import { Link } from "react-router-dom";
import { useSignup } from "../hooks/authHooks";

export default function Signup() {
  const signup = useSignup();
  const [email, setEmail] = useState("");
  const [pw, setPw] = useState("");

  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    signup.mutate({ email, password: pw });
  };

  return (
    <div className="h-screen flex items-center justify-center">
      <form onSubmit={submit} className="w-80 space-y-4">
        <h1 className="text-2xl font-semibold">Create account</h1>

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
          disabled={signup.isPending}
        >
          {signup.isPending ? "…" : "Sign up"}
        </button>

        <p className="text-sm text-center">
          Already have an account?{" "}
          <Link className="text-blue-600" to="/login">
            Login →
          </Link>
        </p>
      </form>
    </div>
  );
}

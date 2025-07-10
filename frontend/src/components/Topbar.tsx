import { Link } from "react-router-dom";
import { useLogout } from "../hooks/authHooks";

export default function TopBar() {
  const logout = useLogout();

  return (
    <header className="flex items-center justify-between px-6 py-3 shadow">
      <Link to="/" className="text-xl font-semibold">
        URL Monitor
      </Link>

      <button
        onClick={() => logout.mutate()}
        disabled={logout.isPending}
        className="bg-red-600 text-white px-3 py-1 rounded hover:bg-red-700 disabled:opacity-60"
      >
        {logout.isPending ? "..." : "Logout"}
      </button>
    </header>
  );
}

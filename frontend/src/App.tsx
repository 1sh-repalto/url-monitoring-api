import { BrowserRouter, Routes, Route, Outlet } from "react-router-dom";
import Login from "./pages/Login";
import Signup from "./pages/Signup";
import Dashboard from "./pages/Dashboard";
import Protected from "./components/Protected";
import TopBar from "./components/Topbar";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />

        {/* Protected routes */}
        <Route element={<Protected />}>
          <Route
            element={
              <>
                <TopBar />
                <Outlet /> {/* children pages render here */}
              </>
            }
          >
            <Route path="/" element={<Dashboard />} />
            {/* add other protected routes here */}
          </Route>
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

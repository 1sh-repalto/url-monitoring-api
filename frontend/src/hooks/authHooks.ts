import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import api from "../lib/http";

export function useLogin() {
  const nav = useNavigate();

  return useMutation({
    mutationFn: (payload: { email: string; password: string }) =>
      api.post("/auth/login", payload),
    onSuccess: () => nav("/", { replace: true }),
  });
}

export function useSignup() {
  const nav = useNavigate();

  return useMutation({
    mutationFn: (payload: { email: string; password: string }) =>
      api.post("/auth/signup", payload),
    onSuccess: () => nav("/login", { replace: true }), // or auto‑login
  });
}

export function useLogout() {
  const nav = useNavigate();

  return useMutation({
    mutationFn: () => api.post("/auth/logout"),
    onSuccess: () => {
      // clear React‑Query cache if you like, or just navigate
      nav("/login", { replace: true });
    },
  });
}
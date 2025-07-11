import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { toast } from "react-toastify";
import api from "../lib/http";
import { errMsg } from "../lib/errorMessage";

export function useLogin() {
  const nav = useNavigate();
  const qc = useQueryClient();

  return useMutation({
    mutationFn: (p: { email: string; password: string }) =>
      api.post("/auth/login", p),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["me"] });
      toast.success("Logged in!");
      nav("/", { replace: true });
    },
    onError: (e) => toast.error(errMsg(e)),
  });
}

export function useSignup() {
  const nav = useNavigate();
  const qc = useQueryClient();

  return useMutation({
    mutationFn: (p: { email: string; password: string }) =>
      api.post("/auth/signup", p),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["me"] });
      toast.success("Signup successful — please sign in");
      nav("/login", { replace: true });
    },
    onError: (e) => toast.error(errMsg(e)),
  });
}
export function useLogout() {
  const nav = useNavigate();
  return useMutation({
    mutationFn: () => api.post("/auth/logout"),
    onSuccess: () => {
      toast.success("Logged out");
      nav("/login", { replace: true });
    },
    onError: () => toast.error("Logout failed"),
  });
}

export function useMe() {
  return useQuery({
    queryKey: ["me"],
    queryFn: () => api.get("/auth/me").then((r) => r.data),
    staleTime: 5 * 60 * 1000,
  });
}

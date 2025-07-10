import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "../lib/http";

export function useUrls() {
  return useQuery({
    queryKey: ["urls"],
    queryFn: () => api.get("/urls").then(r => r.data),
  });
}

export function useAddUrl() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (url: string) => api.post("/urls", { url }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["urls"] }),
  });
}

export function useDeleteUrl() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`/urls/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["urls"] }),
  });
}

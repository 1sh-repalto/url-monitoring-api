import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "react-toastify";
import api from "../lib/http";
import { errMsg } from "../lib/errorMessage";

export function useUrls() {
  return useQuery({
    queryKey: ["urls"],
    queryFn: () => api.get("/urls").then((r) => r.data),
  });
}

export function useAddUrl() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (url: string) => api.post("/urls", { url }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["urls"] });
      toast.success("URL added!");
    },
    onError: (e) => toast.error(errMsg(e)),
  });
}

export function useDeleteUrl() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.delete(`/urls/${id}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["urls"] });
      toast.info("URL removed");
    },
    onError: () => toast.error("Delete failed"),
  });
}

export function useUrlLogs(urlId?: string) {
  return useQuery({
    queryKey: ["urlLogs", urlId],
    queryFn: () => api.get(`/urls/${urlId}/logs`).then((res) => res.data),
    enabled: !!urlId,
    staleTime: 1000 * 30,
  });
}

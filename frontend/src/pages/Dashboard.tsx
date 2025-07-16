import { useState } from "react";
import { useUrls, useAddUrl, useDeleteUrl } from "../hooks/useURLs";
import { Trash2 } from "lucide-react";
import LogsModal from "../components/LogsModal";

export default function Dashboard() {
  const { data: urls = [], isLoading } = useUrls();
  const [openUrl, setOpenUrl] = useState<{ id: string; url: string } | null>(
    null
  );
  const addUrl = useAddUrl();
  const delUrl = useDeleteUrl();

  const [input, setInput] = useState("");

  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim()) return;
    addUrl.mutate(input, { onSuccess: () => setInput("") });
  };

  return (
    <div className="p-6 max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Monitored URLs</h1>

      {/* add‑url input */}
      <form onSubmit={submit} className="flex gap-2 mb-6">
        <input
          className="flex-1 border p-2 rounded"
          placeholder="https://example.com"
          value={input}
          onChange={(e) => setInput(e.target.value)}
        />
        <button
          className="bg-blue-600 text-white px-4 rounded disabled:opacity-50"
          disabled={addUrl.isPending}
        >
          {addUrl.isPending ? "Adding…" : "Add"}
        </button>
      </form>

      {isLoading ? (
        <p>Loading…</p>
      ) : urls.length === 0 ? (
        <p>No URLs yet.</p>
      ) : (
        <ul className="space-y-2">
          {urls.map((u: any) => (
            <li
              key={u.id}
              className="flex items-center justify-between p-3 border rounded-lg"
            >
              <button
                onClick={() => setOpenUrl({ id: u.id, url: u.url })}
                className="text-left flex-1 break-all hover:underline"
              >
                {u.url}
              </button>

              <button
                onClick={() => delUrl.mutate(u.id)}
                className="text-gray-500 hover:text-red-600"
              >
                <Trash2 size={18} />
              </button>
            </li>
          ))}
        </ul>
      )}
      {openUrl && (
        <LogsModal
          urlId={openUrl.id}
          url={openUrl.url}
          onClose={() => setOpenUrl(null)}
        />
      )}
    </div>
  );
}

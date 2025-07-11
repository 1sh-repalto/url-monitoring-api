import {
  Dialog,
  DialogBackdrop,
  DialogPanel,
  DialogTitle,
} from "@headlessui/react";
import { useUrlLogs } from "../hooks/useURLs";

type Props = { urlId: string; url: string; onClose: () => void };

export default function LogsModal({ urlId, url, onClose }: Props) {
  const { data = [], isLoading, error } = useUrlLogs(urlId);
  console.log(data);
  return (
    <Dialog open onClose={onClose} className="relative z-50">
      <DialogBackdrop className="fixed inset-0 bg-black/30" />

      <div className="fixed inset-0 flex w-screen items-center justify-center p-4">
        <DialogPanel className="w-full max-w-xl rounded bg-white p-6">
          <DialogTitle className="mb-4 text-lg font-semibold">
            Logs for {url}
          </DialogTitle>

          {isLoading ? (
            <p>Loading…</p>
          ) : error ? (
            <p className="text-red-600">Failed to load logs</p>
          ) : data.length === 0 ? (
            <p>No logs yet.</p>
          ) : (
            <ul className="max-h-80 space-y-2 overflow-y-auto text-sm">
              {data.map((l: any) => (
                <li
                  key={l.id}
                  className="grid grid-cols-[200px_90px_1fr] gap-3 border-b py-1"
                >
                  <span>{new Date(l.checked_at).toLocaleString()}</span>
                  <span>{l.status_code}</span>
                  <span>{l.response_time_ms} ms</span>
                </li>
              ))}
            </ul>
          )}

          <button
            onClick={onClose}
            className="mt-6 rounded bg-blue-600 px-4 py-2 text-white"
          >
            Close
          </button>
        </DialogPanel>
      </div>
    </Dialog>
  );
}

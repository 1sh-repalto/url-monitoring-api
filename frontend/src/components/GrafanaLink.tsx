import { useMe } from "../hooks/authHooks";

export default function GrafanaLink() {
  const { data: me, isLoading, error } = useMe();

  if (isLoading || error || !me) return null;

  const url =
    `http://localhost:3001/d/891dd564-9d27-4ea3-b860-11e9fe616176/url-monitoring-api-dashboard
?orgId=1&from=now-15m&to=now&timezone=browser&var-user_id=${me.user_id}&refresh=1m`.replace(
      /\s+/g,
      ""
    );

  return (
    <a
      href={url}
      target="_blank"
      rel="noopener noreferrer"
      className="inline-block bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700"
    >
      View my Grafana dashboard
    </a>
  );
}

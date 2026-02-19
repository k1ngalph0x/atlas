import { useState, useEffect } from "react";
import { alerts } from "../../api/client";

export default function AlertsList({ projectId }) {
  const [alertsList, setAlertsList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    loadAlerts();
  }, [projectId]);

  const loadAlerts = async () => {
    try {
      const { data } = await alerts.getAlerts(projectId);
      setAlertsList(data.alerts || []);
    } catch (err) {
      setError("Failed to load alerts", err);
    } finally {
      setLoading(false);
    }
  };

  const handleAcknowledge = async (alertId) => {
    try {
      await alerts.acknowledgeAlert(alertId);
      setAlertsList(
        alertsList.map((alert) =>
          alert.id === alertId ? { ...alert, acknowledged: true } : alert,
        ),
      );
    } catch (err) {
      console.error("Failed to acknowledge alert:", err);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <p className="text-gray-600">Loading alerts...</p>
      </div>
    );
  }

  if (error) {
    return <div className="bg-red-50 text-red-600 p-4 rounded">{error}</div>;
  }

  if (alertsList.length === 0) {
    return (
      <div className="bg-white shadow rounded-lg p-6 text-center">
        <p className="text-gray-600">No alerts fired yet.</p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow rounded-lg overflow-hidden">
      <ul className="divide-y divide-gray-200">
        {alertsList.map((alert) => (
          <li
            key={alert.id}
            className={`px-6 py-4 ${
              alert.acknowledged ? "bg-white" : "bg-blue-50"
            }`}
          >
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center space-x-2 mb-1">
                  {!alert.acknowledged && (
                    <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
                      New
                    </span>
                  )}
                  <span className="text-xs text-gray-500">
                    {new Date(alert.fired_at).toLocaleString()}
                  </span>
                </div>
                <p className="text-sm text-gray-900">{alert.message}</p>
                <p className="text-xs text-gray-500 mt-1">
                  Issue ID: {alert.issue_id}
                </p>
              </div>
              {!alert.acknowledged && (
                <button
                  onClick={() => handleAcknowledge(alert.id)}
                  className="ml-4 text-xs text-blue-600 hover:text-blue-700"
                >
                  Acknowledge
                </button>
              )}
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

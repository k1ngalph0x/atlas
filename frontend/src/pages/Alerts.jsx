import { useParams, useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import { alerts } from "../api/client";
import AlertsList from "../components/alerts/AlertsList";
import AlertRules from "../components/alerts/AlertRules";

export default function Alerts() {
  const { projectId } = useParams();
  const navigate = useNavigate();
  const [unreadCount, setUnreadCount] = useState(0);
  const [activeTab, setActiveTab] = useState("alerts");

  useEffect(() => {
    alerts
      .getUnreadCount(projectId)
      .then(({ data }) => setUnreadCount(data.unread_count))
      .catch((err) => console.error("Failed to load unread count:", err));
  }, [projectId]);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Alerts</h1>
              <p className="text-sm text-gray-600 mt-1">
                Alert notifications and rules for this project
              </p>
            </div>
            <button
              onClick={() => navigate(`/projects/${projectId}`)}
              className="text-sm text-blue-600 hover:text-blue-700"
            >
              ← Back to issues
            </button>
          </div>

          <div className="border-b border-gray-200 mb-6">
            <nav className="-mb-px flex space-x-8">
              <button
                onClick={() => setActiveTab("alerts")}
                className={`${
                  activeTab === "alerts"
                    ? "border-blue-500 text-blue-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center space-x-2`}
              >
                <span>Fired Alerts</span>
                {unreadCount > 0 && (
                  <span className="bg-blue-100 text-blue-800 px-2 py-0.5 rounded-full text-xs font-medium">
                    {unreadCount}
                  </span>
                )}
              </button>
              <button
                onClick={() => setActiveTab("rules")}
                className={`${
                  activeTab === "rules"
                    ? "border-blue-500 text-blue-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm`}
              >
                Rules
              </button>
            </nav>
          </div>

          {activeTab === "alerts" && <AlertsList projectId={projectId} />}
          {activeTab === "rules" && <AlertRules projectId={projectId} />}
        </div>
      </div>
    </div>
  );
}

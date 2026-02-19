import { useParams, useNavigate, useLocation } from "react-router-dom";
import { useState, useEffect } from "react";
import { issues, alerts } from "../api/client";
import IssuesList from "../components/issues/IssuesList";

export default function ProjectIssues() {
  const { projectId } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const [overview, setOverview] = useState(null);
  const [unreadCount, setUnreadCount] = useState(0);

  const activeTab = location.pathname.includes("/alerts") ? "alerts" : "issues";

  useEffect(() => {
    const loadData = async () => {
      try {
        const [overviewRes, unreadRes] = await Promise.all([
          issues.getOverview(projectId),
          alerts.getUnreadCount(projectId),
        ]);

        setOverview(overviewRes.data);
        setUnreadCount(unreadRes.data.unread_count);
      } catch (err) {
        console.error("Failed to load project data:", err);
      }
    };

    if (projectId) {
      loadData();
    }
  }, [projectId]);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="flex items-center justify-between mb-6">
            <h1 className="text-2xl font-bold text-gray-900">
              Project Dashboard
            </h1>
            <button
              onClick={() => navigate("/dashboard")}
              className="text-sm text-blue-600 hover:text-blue-700"
            >
              ← Back to projects
            </button>
          </div>

          {/* Tabs */}
          <div className="border-b border-gray-200 mb-6">
            <nav className="-mb-px flex space-x-8">
              <button
                onClick={() => navigate(`/projects/${projectId}`)}
                className={`${
                  activeTab === "issues"
                    ? "border-blue-500 text-blue-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm`}
              >
                Issues
              </button>
              <button
                onClick={() => navigate(`/projects/${projectId}/alerts`)}
                className={`${
                  activeTab === "alerts"
                    ? "border-blue-500 text-blue-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center space-x-2`}
              >
                <span>Alerts</span>
                {unreadCount > 0 && (
                  <span className="bg-blue-100 text-blue-800 px-2 py-0.5 rounded-full text-xs font-medium">
                    {unreadCount}
                  </span>
                )}
              </button>
            </nav>
          </div>

          {overview && (
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
              <div className="bg-white shadow rounded-lg p-4">
                <div className="text-sm text-gray-500">Total Issues</div>
                <div className="text-2xl font-bold text-gray-900">
                  {overview.total_issues}
                </div>
              </div>
              <div className="bg-white shadow rounded-lg p-4">
                <div className="text-sm text-gray-500">Open</div>
                <div className="text-2xl font-bold text-blue-600">
                  {overview.open_issues}
                </div>
              </div>
              <div className="bg-white shadow rounded-lg p-4">
                <div className="text-sm text-gray-500">Critical</div>
                <div className="text-2xl font-bold text-red-600">
                  {overview.critical_count}
                </div>
              </div>
              <div className="bg-white shadow rounded-lg p-4">
                <div className="text-sm text-gray-500">Errors</div>
                <div className="text-2xl font-bold text-orange-600">
                  {overview.error_count}
                </div>
              </div>
            </div>
          )}

          <IssuesList projectId={projectId} />
        </div>
      </div>
    </div>
  );
}

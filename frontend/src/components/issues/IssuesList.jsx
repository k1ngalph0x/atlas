import { useState, useEffect } from "react";
import { issues } from "../../api/client";
import { useNavigate } from "react-router-dom";

export default function IssuesList({ projectId }) {
  const [issuesList, setIssuesList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    loadIssues();
  }, [projectId]);

  const loadIssues = async () => {
    try {
      const { data } = await issues.getIssues(projectId);
      setIssuesList(data.issues || []);
    } catch (err) {
      setError("Failed to load issues", err);
    } finally {
      setLoading(false);
    }
  };

  const getLevelColor = (level) => {
    switch (level) {
      case "critical":
        return "bg-red-100 text-red-800";
      case "error":
        return "bg-orange-100 text-orange-800";
      case "warning":
        return "bg-yellow-100 text-yellow-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const getStatusColor = (status) => {
    return status === "open"
      ? "bg-blue-100 text-blue-800"
      : "bg-green-100 text-green-800";
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <p className="text-gray-600">Loading issues...</p>
      </div>
    );
  }

  if (error) {
    return <div className="bg-red-50 text-red-600 p-4 rounded">{error}</div>;
  }

  if (issuesList.length === 0) {
    return (
      <div className="bg-white shadow rounded-lg p-6 text-center">
        <p className="text-gray-600">
          No issues yet. Send your first error event to get started.
        </p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow rounded-lg overflow-hidden">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Issue
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Level
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Count
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Last Seen
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {issuesList.map((issue) => (
              <tr
                key={issue.id}
                onClick={() =>
                  navigate(`/projects/${projectId}/issues/${issue.id}`)
                }
                className="hover:bg-gray-50 cursor-pointer"
              >
                <td className="px-6 py-4">
                  <div className="text-sm font-medium text-gray-900 truncate max-w-md">
                    {issue.title}
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    First seen: {new Date(issue.first_seen).toLocaleString()}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getLevelColor(issue.level)}`}
                  >
                    {issue.level}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {issue.count}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(issue.status)}`}
                  >
                    {issue.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(issue.last_seen).toLocaleString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

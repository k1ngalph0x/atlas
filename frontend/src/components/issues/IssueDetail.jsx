import { useState, useEffect, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { issues, insights } from "../../api/client";

const POLL_INTERVAL = 5000;
const MAX_POLLS = 12;

function InsightPanel({ issueId }) {
  const [insight, setInsight] = useState(null);
  const [gaveUp, setGaveUp] = useState(false);
  const pollRef = useRef(null);
  const pollCount = useRef(0);

  useEffect(() => {
    if (!issueId) return;

    let isMounted = true;

    const fetchInsight = async () => {
      try {
        const res = await insights.getIssueInsight(issueId);

        if (res.status === 200 && isMounted) {
          setInsight(res.data.insight);
          return;
        }
      } catch {
        startPolling();
      }
    };

    const startPolling = () => {
      pollCount.current = 0;

      pollRef.current = setInterval(async () => {
        pollCount.current++;

        if (pollCount.current > MAX_POLLS) {
          clearInterval(pollRef.current);
          if (isMounted) setGaveUp(true);
          return;
        }

        try {
          const res = await insights.getIssueInsight(issueId);

          if (res.status === 200) {
            if (isMounted) setInsight(res.data.insight);
            clearInterval(pollRef.current);
          }
        } catch (err) {
          console.log(err);
        }
      }, POLL_INTERVAL);
    };

    fetchInsight();

    return () => {
      isMounted = false;
      if (pollRef.current) {
        clearInterval(pollRef.current);
      }
    };
  }, [issueId]);

  if (gaveUp) {
    return (
      <div className="bg-gray-50 rounded-lg p-6 text-center">
        <p className="text-sm text-gray-600">
          AI analysis not available for this issue.
        </p>
      </div>
    );
  }

  if (!insight) {
    return (
      <div className="bg-white shadow rounded-lg p-6">
        <h2 className="text-lg font-medium text-gray-900 mb-4">AI Analysis</h2>
        <div className="flex items-center space-x-3 text-sm text-gray-500">
          <svg
            className="animate-spin h-4 w-4 text-blue-500 shrink-0"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8v8z"
            />
          </svg>
          <span>Analyzing with AI — this may take up to a minute...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white shadow rounded-lg p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-medium text-gray-900">AI Analysis</h2>
        <span className="text-xs text-gray-500">
          Powered by {insight.model_used}
        </span>
      </div>
      <div className="space-y-4">
        <div>
          <h3 className="text-sm font-medium text-gray-700 mb-2">Summary</h3>
          <p className="text-sm text-gray-900">{insight.summary}</p>
        </div>
        <div>
          <h3 className="text-sm font-medium text-gray-700 mb-2">Root Cause</h3>
          <p className="text-sm text-gray-900">{insight.root_cause}</p>
        </div>
        <div>
          <h3 className="text-sm font-medium text-gray-700 mb-2">
            Recommended Fix
          </h3>
          <p className="text-sm text-gray-900">{insight.remediation}</p>
        </div>
      </div>
    </div>
  );
}

export default function IssueDetail() {
  const { projectId, issueId } = useParams();
  const navigate = useNavigate();
  const [issue, setIssue] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    issues
      .getIssueDetail(projectId, issueId)
      .then((res) => setIssue(res.data.issue))
      .catch((err) => console.error("Failed to load issue:", err))
      .finally(() => setLoading(false));
  }, [issueId]);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <p className="text-gray-600">Loading...</p>
      </div>
    );
  }

  if (!issue) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <p className="text-gray-600">Issue not found</p>
      </div>
    );
  }

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

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <button
            onClick={() => navigate(`/projects/${projectId}`)}
            className="text-sm text-blue-600 hover:text-blue-700 mb-4"
          >
            ← Back to issues
          </button>

          <div className="bg-white shadow rounded-lg p-6 mb-6">
            <div className="flex items-start justify-between mb-4">
              <div className="flex-1">
                <div className="flex items-center space-x-2 mb-2">
                  <span
                    className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getLevelColor(issue.level)}`}
                  >
                    {issue.level}
                  </span>
                  <span className="text-sm text-gray-500">
                    {issue.count} occurrences
                  </span>
                </div>
                <h1 className="text-2xl font-bold text-gray-900">
                  {issue.title}
                </h1>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-gray-500">First seen:</span>
                <span className="ml-2 text-gray-900">
                  {new Date(issue.first_seen).toLocaleString()}
                </span>
              </div>
              <div>
                <span className="text-gray-500">Last seen:</span>
                <span className="ml-2 text-gray-900">
                  {new Date(issue.last_seen).toLocaleString()}
                </span>
              </div>
            </div>
          </div>

          {issue.stack_trace && (
            <div className="bg-white shadow rounded-lg p-6 mb-6">
              <h2 className="text-lg font-medium text-gray-900 mb-4">
                Stack Trace
              </h2>
              <pre className="bg-gray-900 text-gray-100 p-4 rounded text-sm overflow-x-auto">
                {issue.stack_trace}
              </pre>
            </div>
          )}

          <InsightPanel issueId={issueId} />
        </div>
      </div>
    </div>
  );
}

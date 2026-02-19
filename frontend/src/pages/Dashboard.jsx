import { useState, useEffect } from "react";
import { projects as projectsAPI } from "../api/client";
import { useNavigate } from "react-router-dom";
import CreateOrganization from "../components/projects/CreateOrganization";
import CreateProject from "../components/projects/CreateProject";
import ProjectList from "../components/projects/ProjectList";

function APIKeyModal({ apiKey, onClose }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(apiKey);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4 shadow-xl">
        <h3 className="text-lg font-bold text-gray-900 mb-2">Your API Key</h3>
        <div className="flex items-start space-x-2 bg-yellow-50 border border-yellow-200 rounded p-3 mb-4">
          <span className="text-yellow-600 text-sm">⚠️</span>
          <p className="text-sm text-yellow-700">
            Copy this key now — it will <strong>never be shown again</strong>.
            Store it somewhere safe like a password manager or environment
            variable.
          </p>
        </div>
        <div className="bg-gray-100 rounded p-3 font-mono text-sm break-all text-gray-800 mb-4 select-all">
          {apiKey}
        </div>
        <button
          onClick={handleCopy}
          className="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 rounded text-sm font-medium mb-2 transition-colors"
        >
          {copied ? "✓ Copied!" : "Copy to clipboard"}
        </button>
        <button
          onClick={onClose}
          className="w-full border border-gray-300 hover:bg-gray-50 py-2 rounded text-sm text-gray-700 transition-colors"
        >
          I've saved it, close
        </button>
      </div>
    </div>
  );
}

export default function Dashboard() {
  const navigate = useNavigate();
  const email = localStorage.getItem("email");

  const [organizations, setOrganizations] = useState([]);
  const [projects, setProjects] = useState([]);
  const [activeTab, setActiveTab] = useState("projects");
  const [loading, setLoading] = useState(true);
  const [revealedKey, setRevealedKey] = useState(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [orgsRes, projectsRes] = await Promise.all([
        projectsAPI.getOrganizations(),
        projectsAPI.getProjects(),
      ]);
      setOrganizations(orgsRes.data.organizations);
      setProjects(projectsRes.data.projects);
    } catch (err) {
      console.error("Failed to load data:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("email");
    navigate("/signin");
  };

  const handleOrganizationCreated = (org) => {
    setOrganizations([...organizations, org]);
  };

  const handleProjectCreated = (project) => {
    const { api_key, ...projectWithoutKey } = project;
    setProjects([...projects, projectWithoutKey]);
    setActiveTab("projects");
    setRevealedKey(api_key);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <p className="text-gray-600">Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {revealedKey && (
        <APIKeyModal
          apiKey={revealedKey}
          onClose={() => setRevealedKey(null)}
        />
      )}

      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold">Atlas</h1>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-700">{email}</span>
              <button
                onClick={handleLogout}
                className="text-sm text-gray-700 hover:text-gray-900"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-6">
            <h2 className="text-2xl font-bold text-gray-900">Dashboard</h2>
            <p className="mt-1 text-sm text-gray-600">
              Manage your organizations and projects
            </p>
          </div>

          <div className="border-b border-gray-200 mb-6">
            <nav className="-mb-px flex space-x-8">
              <button
                onClick={() => setActiveTab("projects")}
                className={`${
                  activeTab === "projects"
                    ? "border-blue-500 text-blue-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm`}
              >
                Projects
              </button>
              <button
                onClick={() => setActiveTab("setup")}
                className={`${
                  activeTab === "setup"
                    ? "border-blue-500 text-blue-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm`}
              >
                Setup
              </button>
            </nav>
          </div>

          {activeTab === "projects" && (
            <div className="space-y-6">
              <ProjectList projects={projects} />
            </div>
          )}

          {activeTab === "setup" && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <CreateOrganization onSuccess={handleOrganizationCreated} />
              <CreateProject
                organizations={organizations}
                onSuccess={handleProjectCreated}
              />
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

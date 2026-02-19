import { useState } from "react";
import { projects } from "../../api/client";

export default function CreateProject({ organizations, onSuccess }) {
  const [orgId, setOrgId] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const { data } = await projects.createProject(orgId, name);
      setName("");
      setOrgId("");
      onSuccess(data.project);
    } catch (err) {
      setError(err.response?.data?.error || "Failed to create project");
    } finally {
      setLoading(false);
    }
  };

  if (organizations.length === 0) {
    return (
      <div className="bg-white shadow rounded-lg p-6">
        <p className="text-sm text-gray-600">
          Create an organization first before creating projects.
        </p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow rounded-lg p-6">
      <h3 className="text-lg font-medium text-gray-900 mb-4">Create Project</h3>

      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <div className="bg-red-50 text-red-600 p-3 rounded text-sm">
            {error}
          </div>
        )}

        <div>
          <label className="block text-sm font-medium text-gray-700">
            Organization
          </label>
          <select
            value={orgId}
            onChange={(e) => setOrgId(e.target.value)}
            required
            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">Select organization</option>
            {organizations.map((org) => (
              <option key={org.id} value={org.id}>
                {org.organization_name}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">
            Project Name
          </label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            placeholder="e.g., production-api"
            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          {loading ? "Creating..." : "Create Project"}
        </button>
      </form>
    </div>
  );
}

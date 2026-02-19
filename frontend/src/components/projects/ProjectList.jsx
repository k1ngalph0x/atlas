import { useNavigate } from "react-router-dom";

export default function ProjectList({ projects }) {
  const navigate = useNavigate();

  if (projects.length === 0) {
    return (
      <div className="bg-white shadow rounded-lg p-6">
        <p className="text-sm text-gray-600">No projects yet.</p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow rounded-lg overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-200">
        <h3 className="text-lg font-medium text-gray-900">Your Projects</h3>
      </div>
      <ul className="divide-y divide-gray-200">
        {projects.map((project) => (
          <li key={project.id} className="px-6 py-4 hover:bg-gray-50">
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <button
                  onClick={() => navigate(`/projects/${project.id}`)}
                  className="text-sm font-medium text-blue-600 hover:text-blue-700"
                >
                  {project.project_name}
                </button>
                <p className="text-xs text-gray-500 mt-1">
                  Created {new Date(project.created_at).toLocaleDateString()}
                </p>
              </div>

              <div className="flex items-center space-x-2">
                <code className="px-2 py-1 bg-gray-100 text-xs text-gray-500 rounded">
                  atlas_••••••••••••••••
                </code>
              </div>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

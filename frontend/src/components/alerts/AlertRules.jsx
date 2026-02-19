import { useState, useEffect } from "react";
import { alerts } from "../../api/client";

const CONDITIONS = [
  {
    value: "new_issue",
    label: "New Issue",
    desc: "Fires when a brand new issue is detected",
  },
  {
    value: "critical_error",
    label: "Critical / Error",
    desc: "Fires when an error or critical level event occurs",
  },
  {
    value: "count_threshold",
    label: "Count Threshold",
    desc: "Fires when an issue exceeds a set occurrence count",
  },
];

export default function AlertRules({ projectId }) {
  const [rules, setRules] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);

  const [name, setName] = useState("");
  const [condition, setCondition] = useState("new_issue");
  const [threshold, setThreshold] = useState(5);
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState("");

  useEffect(() => {
    loadRules();
  }, [projectId]);

  const loadRules = async () => {
    try {
      const { data } = await alerts.getRules(projectId);
      setRules(data.rules || []);
    } catch {
      setError("Failed to load alert rules");
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (e) => {
    e.preventDefault();
    setFormError("");
    setSaving(true);
    try {
      const { data } = await alerts.createRule(projectId, {
        name,
        condition,
        threshold: condition === "count_threshold" ? Number(threshold) : 0,
      });
      setRules([data.rule, ...rules]);
      setName("");
      setCondition("new_issue");
      setThreshold(5);
      setShowForm(false);
    } catch (err) {
      setFormError(err.response?.data?.error || "Failed to create rule");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (ruleId) => {
    try {
      await alerts.deleteRule(projectId, ruleId);
      setRules(rules.filter((r) => r.id !== ruleId));
    } catch {
      console.error("Failed to delete rule");
    }
  };

  if (loading)
    return <p className="text-sm text-gray-500 py-4">Loading rules...</p>;

  return (
    <div className="bg-white shadow rounded-lg p-6">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h2 className="text-lg font-medium text-gray-900">Alert Rules</h2>
          <p className="text-xs text-gray-500 mt-0.5">
            Define when Atlas should fire an alert
          </p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="text-sm bg-blue-600 hover:bg-blue-700 text-white px-3 py-1.5 rounded"
        >
          {showForm ? "Cancel" : "+ New Rule"}
        </button>
      </div>

      {error && <p className="text-sm text-red-600 mb-4">{error}</p>}

      {showForm && (
        <form
          onSubmit={handleCreate}
          className="bg-gray-50 rounded-lg p-4 mb-4 space-y-3"
        >
          {formError && <p className="text-sm text-red-600">{formError}</p>}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Rule Name
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              placeholder="e.g. Alert on new errors"
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Condition
            </label>
            <select
              value={condition}
              onChange={(e) => setCondition(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              {CONDITIONS.map((c) => (
                <option key={c.value} value={c.value}>
                  {c.label}
                </option>
              ))}
            </select>
            <p className="text-xs text-gray-500 mt-1">
              {CONDITIONS.find((c) => c.value === condition)?.desc}
            </p>
          </div>

          {condition === "count_threshold" && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Threshold (fire when count exceeds)
              </label>
              <input
                type="number"
                min={1}
                value={threshold}
                onChange={(e) => setThreshold(e.target.value)}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          )}

          <button
            type="submit"
            disabled={saving}
            className="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 rounded text-sm font-medium disabled:opacity-50"
          >
            {saving ? "Creating..." : "Create Rule"}
          </button>
        </form>
      )}

      {rules.length === 0 ? (
        <p className="text-sm text-gray-500 text-center py-4">
          No rules yet — create one to start receiving alerts.
        </p>
      ) : (
        <ul className="divide-y divide-gray-100">
          {rules.map((rule) => (
            <li
              key={rule.id}
              className="py-3 flex items-center justify-between"
            >
              <div>
                <p className="text-sm font-medium text-gray-900">{rule.name}</p>
                <p className="text-xs text-gray-500 mt-0.5">
                  {CONDITIONS.find((c) => c.value === rule.condition)?.label}
                  {rule.condition === "count_threshold" &&
                    ` › ${rule.threshold}`}
                </p>
              </div>
              <div className="flex items-center space-x-3">
                <span
                  className={`text-xs px-2 py-0.5 rounded-full font-medium ${
                    rule.is_active
                      ? "bg-green-100 text-green-700"
                      : "bg-gray-100 text-gray-500"
                  }`}
                >
                  {rule.is_active ? "active" : "inactive"}
                </span>
                <button
                  onClick={() => handleDelete(rule.id)}
                  className="text-xs text-red-500 hover:text-red-700"
                >
                  Delete
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

import axios from "axios";

const createClient = (baseURL) => {
  const instance = axios.create({
    baseURL,
    headers: { "Content-Type": "application/json" },
    validateStatus: (status) => status < 300 || status === 202,
  });

  instance.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) config.headers.Authorization = `Bearer ${token}`;
    return config;
  });

  return instance;
};

const identityClient = createClient(import.meta.env.VITE_IDENTITY_URL);
const processingClient = createClient(import.meta.env.VITE_PROCESSING_URL);
const intelligenceClient = createClient(import.meta.env.VITE_INTELLIGENCE_URL);
const alertClient = createClient(import.meta.env.VITE_ALERT_URL);

export const auth = {
  signUp: (email, password) =>
    identityClient.post("/auth/signup", { email, password }),
  signIn: (email, password) =>
    identityClient.post("/auth/signin", { email, password }),
};

export const projects = {
  createOrganization: (organization_name) =>
    identityClient.post("/project/create-organization", { organization_name }),
  createProject: (organization_id, project_name) =>
    identityClient.post("/project/create-project", {
      organization_id,
      project_name,
    }),
  getOrganizations: () => identityClient.get("/project/organizations"),
  getProjects: () => identityClient.get("/project/projects"),
};

export const issues = {
  getIssues: (projectId) =>
    processingClient.get(`/projects/${projectId}/issues`),
  getIssueDetail: (projectId, issueId) =>
    processingClient.get(`/projects/${projectId}/issues/${issueId}`),
  getOverview: (projectId) =>
    processingClient.get(`/projects/${projectId}/overview`),
};

export const insights = {
  getIssueInsight: (issueId) =>
    intelligenceClient.get(`/issues/${issueId}/insight`),
};

export const alerts = {
  getAlerts: (projectId) => alertClient.get(`/projects/${projectId}/alerts`),
  getUnreadCount: (projectId) =>
    alertClient.get(`/projects/${projectId}/alerts/unread`),
  acknowledgeAlert: (alertId) =>
    alertClient.post(`/alerts/${alertId}/acknowledge`),

  getRules: (projectId) => alertClient.get(`/projects/${projectId}/rules`),
  createRule: (projectId, rule) =>
    alertClient.post(`/projects/${projectId}/rules`, rule),
  deleteRule: (projectId, ruleId) =>
    alertClient.delete(`/projects/${projectId}/rules/${ruleId}`),
};

export default identityClient;

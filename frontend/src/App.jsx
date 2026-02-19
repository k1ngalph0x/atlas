import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import SignIn from "./components/auth/SignIn";
import SignUp from "./components/auth/SignUp";
import Dashboard from "./pages/Dashboard";
import ProjectIssues from "./pages/ProjectIssues";
import Alerts from "./pages/Alerts";
import IssueDetail from "./components/issues/IssueDetail";
import ProtectedRoute from "./components/ProtectedRoute";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/signin" replace />} />
        <Route path="/signin" element={<SignIn />} />
        <Route path="/signup" element={<SignUp />} />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/projects/:projectId"
          element={
            <ProtectedRoute>
              <ProjectIssues />
            </ProtectedRoute>
          }
        />
        <Route
          path="/projects/:projectId/alerts"
          element={
            <ProtectedRoute>
              <Alerts />
            </ProtectedRoute>
          }
        />
        <Route
          path="/projects/:projectId/issues/:issueId"
          element={
            <ProtectedRoute>
              <IssueDetail />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;

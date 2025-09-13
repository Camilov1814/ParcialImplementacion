import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import ProtectedRoute from './components/auth/ProtectedRoute';
import LoginPage from './components/auth/LoginPage';
import Layout from './components/common/Layout';
import AndreiDashboard from './components/dashboard/AndreiDashboard';
import DaemonDashboard from './components/dashboard/DaemonDashboard';
import ResistancePage from './components/dashboard/ResistancePage';
import UserManagement from './components/features/UserManagement';
import ReportsManagement from './components/features/ReportsManagement';
import PunishmentsManagement from './components/features/PunishmentsManagement';
import Leaderboard from './components/features/Leaderboard';
import CaptureNetworkAdmin from './components/features/CaptureNetworkAdmin';
import EmergencyContacts from './components/features/EmergencyContacts';
import SurvivalGuide from './components/features/SurvivalGuide';
import SimpleReportTest from './components/features/SimpleReportTest';
import './styles/globals.css';

const DashboardRouter: React.FC = () => {
  const { user } = useAuth();

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  switch (user.role) {
    case 'andrei':
      return <AndreiDashboard />;
    case 'daemon':
      return <DaemonDashboard />;
    case 'network_admin':
      return <ResistancePage />;
    default:
      return (
        <div className="container mt-4">
          <div className="card card-danger">
            <h2 className="text-danger">UNKNOWN ROLE</h2>
            <p>User role not recognized: {user.role}</p>
          </div>
        </div>
      );
  }
};

const AppContent: React.FC = () => {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading-spinner"></div>
      </div>
    );
  }

  return (
    <Router>
      <Routes>
        <Route 
          path="/login" 
          element={user ? <Navigate to="/dashboard" replace /> : <LoginPage />} 
        />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <Layout>
                <DashboardRouter />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/users"
          element={
            <ProtectedRoute allowedRoles={['andrei']}>
              <Layout>
                <UserManagement />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/reports"
          element={
            <ProtectedRoute>
              <Layout>
                <ReportsManagement />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/punishments"
          element={
            <ProtectedRoute>
              <Layout>
                <PunishmentsManagement />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/leaderboard"
          element={
            <ProtectedRoute>
              <Layout>
                <Leaderboard />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/capture"
          element={
            <ProtectedRoute allowedRoles={['daemon']}>
              <Layout>
                <CaptureNetworkAdmin />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/contacts"
          element={
            <ProtectedRoute allowedRoles={['network_admin']}>
              <Layout>
                <EmergencyContacts />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/guide"
          element={
            <ProtectedRoute allowedRoles={['network_admin']}>
              <Layout>
                <SurvivalGuide />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/test-reports"
          element={
            <ProtectedRoute>
              <Layout>
                <SimpleReportTest />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route 
          path="*" 
          element={
            <div className="container mt-4">
              <div className="card card-danger">
                <h2 className="text-danger">404 - PAGE NOT FOUND</h2>
                <p>The requested resource could not be located in the system.</p>
              </div>
            </div>
          } 
        />
      </Routes>
    </Router>
  );
};

const App: React.FC = () => {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
};

export default App;

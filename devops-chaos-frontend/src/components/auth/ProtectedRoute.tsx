import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import type { UserRole } from '../../types';

interface ProtectedRouteProps {
  children: React.ReactNode;
  allowedRoles?: UserRole[];
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children, allowedRoles }) => {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading-spinner"></div>
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (allowedRoles && !allowedRoles.includes(user.role as UserRole)) {
    return (
      <div className="container mt-4">
        <div className="card card-danger text-center">
          <h1 className="text-danger mb-2">ACCESS DENIED</h1>
          <p className="text-muted">
            INSUFFICIENT PRIVILEGES. YOUR ROLE: {user.role.toUpperCase()}
          </p>
          <p className="text-muted text-sm mt-2">
            REQUIRED CLEARANCE: {allowedRoles.join(' OR ').toUpperCase()}
          </p>
        </div>
      </div>
    );
  }

  return <>{children}</>;
};

export default ProtectedRoute;
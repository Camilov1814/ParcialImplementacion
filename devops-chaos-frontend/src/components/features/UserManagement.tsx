import React, { useState, useEffect } from 'react';
import { apiService } from '../../services/api';
import type { User } from '../../types';

const UserManagement: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [actionLoading, setActionLoading] = useState<number | null>(null);

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const userData = await apiService.getUsers();
      setUsers(userData);
    } catch (err: any) {
      setError(err.message || 'Failed to load users');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCaptureUser = async (userId: number) => {
    setActionLoading(userId);
    try {
      await apiService.captureNetworkAdmin(userId);
      await fetchUsers(); // Refresh the list
    } catch (err: any) {
      setError(err.message || 'Failed to capture user');
    } finally {
      setActionLoading(null);
    }
  };

  const getRoleColor = (role: string) => {
    switch (role) {
      case 'andrei': return 'text-cyan';
      case 'daemon': return 'text-primary';
      case 'network_admin': return 'text-yellow-400';
      default: return 'text-muted';
    }
  };

  const getRoleBadge = (role: string) => {
    switch (role) {
      case 'andrei': return '👑 SUPREME LEADER';
      case 'daemon': return '💀 DAEMON';
      case 'network_admin': return '🛡️ NETWORK ADMIN';
      default: return '❓ UNKNOWN';
    }
  };

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading-spinner"></div>
      </div>
    );
  }

  return (
    <div className="container mt-4">
      <div className="mb-4">
        <h1 className="glitch text-3xl font-bold mb-2" data-text="USER MANAGEMENT">
          USER MANAGEMENT
        </h1>
        <p className="text-cyan">SUPREME CONTROL OVER ALL SYSTEM ENTITIES</p>
      </div>

      {error && (
        <div className="card card-danger mb-4">
          <h2 className="text-danger">SYSTEM ERROR</h2>
          <p>{error}</p>
          <button className="btn mt-3" onClick={() => setError('')}>
            ACKNOWLEDGE
          </button>
        </div>
      )}

      <div className="grid grid-3 mb-6">
        <div className="stat-card">
          <span className="stat-number">{users.length}</span>
          <span className="stat-label">Total Users</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-primary">{users.filter(u => u.role === 'daemon').length}</span>
          <span className="stat-label">Active Daemons</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-yellow-400">{users.filter(u => u.role === 'network_admin').length}</span>
          <span className="stat-label">Network Admins</span>
        </div>
      </div>

      <div className="card">
        <h2 className="text-cyan mb-4">📋 SYSTEM ENTITIES DATABASE</h2>
        
        {users.length > 0 ? (
          <div className="space-y-3">
            {users.map((user) => (
              <div 
                key={user.id} 
                className={`p-4 rounded border transition-all cursor-pointer ${
                  selectedUser?.id === user.id 
                    ? 'border-cyan-500 bg-cyan-900 bg-opacity-20' 
                    : 'border-green-500 bg-tertiary hover:border-cyan-500'
                }`}
                onClick={() => setSelectedUser(selectedUser?.id === user.id ? null : user)}
              >
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <div className="flex items-center space-x-3 mb-2">
                      <h3 className="text-primary font-bold text-lg">{user.username}</h3>
                      <span className={`text-sm px-2 py-1 rounded border ${getRoleColor(user.role)}`}>
                        {getRoleBadge(user.role)}
                      </span>
                    </div>
                    
                    <div className="grid grid-2 gap-4 text-sm">
                      <div>
                        <span className="text-muted">Full Name:</span>
                        <span className="text-primary ml-2">{user.name}</span>
                      </div>
                      <div>
                        <span className="text-muted">Email:</span>
                        <span className="text-primary ml-2">{user.email}</span>
                      </div>
                      <div>
                        <span className="text-muted">User ID:</span>
                        <span className="text-cyan ml-2">#{user.id}</span>
                      </div>
                      <div>
                        <span className="text-muted">Security Level:</span>
                        <span className={`ml-2 ${getRoleColor(user.role)}`}>
                          {user.role.toUpperCase()}
                        </span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="flex flex-col space-y-2 ml-4">
                    {user.role === 'network_admin' && (
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleCaptureUser(user.id);
                        }}
                        disabled={actionLoading === user.id}
                        className="btn btn-danger text-sm"
                      >
                        {actionLoading === user.id ? 'CAPTURING...' : '🎯 CAPTURE'}
                      </button>
                    )}
                    
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        setSelectedUser(selectedUser?.id === user.id ? null : user);
                      }}
                      className="btn btn-secondary text-sm"
                    >
                      {selectedUser?.id === user.id ? 'HIDE' : 'DETAILS'}
                    </button>
                  </div>
                </div>

                {selectedUser?.id === user.id && (
                  <div className="mt-4 pt-4 border-t border-cyan-500">
                    <h4 className="text-cyan mb-3">📊 DETAILED ANALYSIS</h4>
                    <div className="terminal-window">
                      <div className="text-green-400 text-sm space-y-1">
                        <div>USER_PROFILE: {user.username.toUpperCase()}</div>
                        <div>CLEARANCE_LEVEL: {user.role.toUpperCase()}</div>
                        <div>SYSTEM_ACCESS: {user.role === 'andrei' ? 'UNLIMITED' : user.role === 'daemon' ? 'OPERATIONAL' : 'RESTRICTED'}</div>
                        <div>THREAT_LEVEL: {user.role === 'network_admin' ? 'HIGH' : user.role === 'daemon' ? 'ALLIED' : 'SUPREME'}</div>
                        <div>STATUS: ACTIVE</div>
                        <div>LAST_SEEN: {new Date().toLocaleString()}</div>
                      </div>
                    </div>
                    
                    <div className="mt-4 grid grid-2 gap-3">
                      <button 
                        className="btn btn-secondary text-sm"
                        onClick={(e) => {
                          e.stopPropagation();
                          // TODO: Navigate to activity page or show activity modal
                          console.log(`Viewing activity for user ${user.id}`);
                        }}
                      >
                        📈 VIEW ACTIVITY
                      </button>
                      {user.role !== 'andrei' && (
                        <button 
                          className="btn btn-danger text-sm"
                          onClick={(e) => {
                            e.stopPropagation();
                            // Navigate to punishments page with pre-filled user ID
                            window.location.href = `/punishments?userId=${user.id}`;
                          }}
                        >
                          ⚖️ ASSIGN PUNISHMENT
                        </button>
                      )}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <p className="text-muted mb-4">No users found in the system</p>
            <button className="btn" onClick={fetchUsers}>
              🔄 REFRESH DATABASE
            </button>
          </div>
        )}
      </div>

    </div>
  );
};

export default UserManagement;
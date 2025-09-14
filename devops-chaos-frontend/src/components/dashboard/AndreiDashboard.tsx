import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { DashboardStats, LeaderboardEntry, Report, CaptureDetail, DashboardResponse } from '../../types';

const AndreiDashboard: React.FC = () => {
  const navigate = useNavigate();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [recentReports, setRecentReports] = useState<Report[]>([]);
  const [allCaptures, setAllCaptures] = useState<CaptureDetail[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const dashboardData: DashboardResponse = await apiService.getAndreiDashboard();

        setStats(dashboardData.system_stats);
        setLeaderboard(dashboardData.top_daemons || []);
        setRecentReports(dashboardData.recent_reports || []);
        setAllCaptures(dashboardData.all_captures || []);
      } catch (err: any) {
        setError('Failed to load dashboard data');
        console.error('Dashboard error:', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading-spinner"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mt-4">
        <div className="card card-danger">
          <h2 className="text-danger">SYSTEM ERROR</h2>
          <p>{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mt-4">
      <div className="mb-4">
        <h1 className="glitch text-3xl font-bold mb-2" data-text="ANDREI COMMAND CENTER">
          ANDREI COMMAND CENTER
        </h1>
        <p className="text-cyan">SUPREME LEADER CONTROL PANEL</p>
      </div>

      {stats && (
        <div className="stats-grid">
          <div className="stat-card">
            <span className="stat-number">{stats.total_users}</span>
            <span className="stat-label">Total Users</span>
          </div>
          <div className="stat-card">
            <span className="stat-number">{stats.total_daemons}</span>
            <span className="stat-label">Active Daemons</span>
          </div>
          <div className="stat-card">
            <span className="stat-number text-danger">{stats.captured_admins}</span>
            <span className="stat-label">Captured Admins</span>
          </div>
          <div className="stat-card">
            <span className="stat-number">{stats.total_reports}</span>
            <span className="stat-label">Total Reports</span>
          </div>
        </div>
      )}

      <div className="grid grid-2 mt-6">
        <div className="card">
          <h2 className="text-cyan mb-4">🏆 TOP PERFORMING DAEMONS</h2>
          {leaderboard.length > 0 ? (
            <div className="space-y-3">
              {leaderboard.map((daemon, index) => (
                <div key={daemon.id} className="flex justify-between items-center p-2 bg-tertiary rounded">
                  <div className="flex items-center space-x-3">
                    <span className="text-primary font-bold">#{index + 1}</span>
                    <div>
                      <div className="text-primary font-semibold">{daemon.username}</div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-primary font-bold">{daemon.points} pts</div>
                    <div className="text-cyan text-sm">{daemon.captures_count} captures</div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-muted">No daemon data available</p>
          )}
        </div>

        <div className="card">
          <h2 className="text-cyan mb-4">📊 RECENT INTEL REPORTS</h2>
          {recentReports.length > 0 ? (
            <div className="space-y-3">
              {recentReports.map((report) => (
                <div key={report.id} className="p-2 bg-tertiary rounded">
                  <div className="text-primary font-semibold text-sm">{report.title}</div>
                  <div className="text-muted text-xs">{report.description}</div>
                  <div className="flex justify-between items-center mt-2">
                    <span className={`text-xs px-2 py-1 rounded ${
                      report.status === 'active' ? 'bg-green-900 text-green-300' : 'bg-gray-700 text-gray-300'
                    }`}>
                      {report.status}
                    </span>
                    <span className="text-muted text-xs">
                      {report.created_at ? new Date(report.created_at).toLocaleDateString() : 'Unknown date'}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-muted">No recent reports</p>
          )}
        </div>
      </div>

      {allCaptures.length > 0 && (
        <div className="card mt-6">
          <h2 className="text-danger mb-4">👤 ALL NETWORK ADMIN CAPTURES</h2>
          <div className="space-y-3">
            {allCaptures.map((capture) => (
              <div key={capture.id} className="p-3 bg-red-900 bg-opacity-20 rounded border border-red-500">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <div className="flex items-center space-x-3 mb-2">
                      <div className="text-danger font-bold text-lg">{capture.target_name}</div>
                      <span className="text-muted text-sm">captured by</span>
                      <div className="text-cyan font-bold">{capture.daemon_name}</div>
                    </div>
                    <div className="text-muted text-sm">
                      {new Date(capture.capture_date).toLocaleDateString()} - {new Date(capture.capture_date).toLocaleTimeString()}
                    </div>
                    <div className="flex items-center space-x-2 mt-2">
                      <span className={`text-xs px-2 py-1 rounded ${
                        capture.difficulty === 'hard' ? 'bg-red-900 text-red-300' :
                        capture.difficulty === 'medium' ? 'bg-yellow-900 text-yellow-300' :
                        'bg-green-900 text-green-300'
                      }`}>
                        {capture.difficulty?.toUpperCase()}
                      </span>
                      <span className="text-xs px-2 py-1 rounded bg-blue-900 text-blue-300">
                        {capture.method}
                      </span>
                      <span className={`text-xs px-2 py-1 rounded ${
                        capture.status === 'captured' ? 'bg-red-900 text-red-300' :
                        capture.status === 'released' ? 'bg-yellow-900 text-yellow-300' :
                        'bg-gray-900 text-gray-300'
                      }`}>
                        {capture.status?.toUpperCase()}
                      </span>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-cyan font-bold">+{capture.points}</div>
                    <div className="text-muted text-xs">points awarded</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="card mt-6">
        <h2 className="text-cyan mb-4">⚡ COMMAND OPERATIONS</h2>
        <div className="grid grid-4 gap-4">
          <button className="btn" onClick={() => navigate('/users')}>
            👥 MANAGE USERS
          </button>
          <button className="btn" onClick={() => navigate('/reports')}>
            📋 VIEW ALL REPORTS
          </button>
          <button className="btn btn-danger" onClick={() => navigate('/punishments')}>
            ⚖️ MANAGE PUNISHMENTS
          </button>
          <button className="btn btn-secondary" onClick={() => navigate('/leaderboard')}>
            📈 LEADERBOARD
          </button>
        </div>
      </div>

      <div className="card mt-6">
        <h2 className="text-cyan mb-4">🔴 RECENT SYSTEM ACTIVITY</h2>
        <div className="terminal-window">
          <div className="text-green-400 text-sm space-y-1">
            <div>[{new Date().toLocaleTimeString()}] SYSTEM_STATUS: OPERATIONAL</div>
            <div>[{new Date().toLocaleTimeString()}] CHAOS_PROTOCOL: ACTIVE</div>
            <div>[{new Date().toLocaleTimeString()}] DAEMON_NETWORK: {stats?.total_daemons} NODES CONNECTED</div>
            <div>[{new Date().toLocaleTimeString()}] RESISTANCE_THREAT: MODERATE</div>
            <div>[{new Date().toLocaleTimeString()}] ANDREI_STATUS: SUPREME_LEADER_ONLINE</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AndreiDashboard;
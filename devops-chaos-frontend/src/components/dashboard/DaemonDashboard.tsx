import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiService } from '../../services/api';
import type { DaemonStats, LeaderboardEntry, Punishment, Capture, DaemonDashboardResponse } from '../../types';

const DaemonDashboard: React.FC = () => {
  const navigate = useNavigate();
  const [stats, setStats] = useState<DaemonStats | null>(null);
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [punishments, setPunishments] = useState<Punishment[]>([]);
  const [captures, setCaptures] = useState<Capture[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const daemonData: DaemonDashboardResponse = await apiService.getDaemonDashboard();

        setStats(daemonData.user_stats);
        setLeaderboard(daemonData.leaderboard || []);
        setPunishments(daemonData.active_punishments || []);
        setCaptures(daemonData.recent_captures || []);
      } catch (err: any) {
        setError('Failed to load daemon data');
        console.error('Daemon dashboard error:', err);
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
          <h2 className="text-danger">CONNECTION ERROR</h2>
          <p>{error}</p>
        </div>
      </div>
    );
  }

  const activePunishments = punishments.filter(p => p.status === 'active');
  const missions = [
    { id: 1, name: "Operation: Network Breach", target: "Corporate Firewall", difficulty: "High", reward: "500 pts" },
    { id: 2, name: "Data Extraction: Alpha", target: "Database Server", difficulty: "Medium", reward: "300 pts" },
    { id: 3, name: "Social Engineering: Beta", target: "Admin Credentials", difficulty: "Low", reward: "150 pts" }
  ];

  return (
    <div className="container mt-4">
      <div className="mb-4">
        <h1 className="glitch text-3xl font-bold mb-2" data-text="DAEMON TERMINAL">
          DAEMON TERMINAL
        </h1>
        <p className="text-cyan">CHAOS OPERATIVE COMMAND CENTER</p>
      </div>

      {stats && (
        <div className="stats-grid">
          <div className="stat-card">
            <span className="stat-number text-cyan">{stats.captures_count}</span>
            <span className="stat-label">Network Admins Captured</span>
          </div>
          <div className="stat-card">
            <span className="stat-number">{stats.reports_count}</span>
            <span className="stat-label">Intel Reports Filed</span>
          </div>
          <div className="stat-card">
            <span className="stat-number text-primary">{stats.points}</span>
            <span className="stat-label">Chaos Points</span>
          </div>
          <div className="stat-card">
            <span className="stat-number text-cyan">#{stats.ranking}</span>
            <span className="stat-label">Global Ranking</span>
          </div>
        </div>
      )}

      <div className="grid grid-2 mt-6">
        <div className="card">
          <h2 className="text-cyan mb-4">🏆 DAEMON LEADERBOARD</h2>
          {leaderboard.length > 0 ? (
            <div className="space-y-2">
              {leaderboard.map((daemon, index) => (
                <div key={daemon.id} className={`flex justify-between items-center p-2 rounded ${
                  index < 3 ? 'bg-green-900' : 'bg-tertiary'
                }`}>
                  <div className="flex items-center space-x-3">
                    <span className={`font-bold ${
                      index === 0 ? 'text-yellow-400' : 
                      index === 1 ? 'text-gray-300' : 
                      index === 2 ? 'text-orange-400' : 'text-primary'
                    }`}>
                      #{index + 1}
                    </span>
                    <div>
                      <div className="text-primary font-semibold text-sm">{daemon.username}</div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-primary font-bold text-sm">{daemon.points}</div>
                    <div className="text-muted text-xs">{daemon.captures_count} captures</div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-muted">No leaderboard data</p>
          )}
        </div>

        <div className="card">
          <h2 className="text-cyan mb-4">🎯 ACTIVE MISSIONS</h2>
          <div className="space-y-3">
            {missions.map((mission) => (
              <div key={mission.id} className="p-3 bg-tertiary rounded border border-green-500">
                <div className="text-primary font-semibold text-sm">{mission.name}</div>
                <div className="text-muted text-xs">Target: {mission.target}</div>
                <div className="flex justify-between items-center mt-2">
                  <span className={`text-xs px-2 py-1 rounded ${
                    mission.difficulty === 'High' ? 'bg-red-900 text-red-300' :
                    mission.difficulty === 'Medium' ? 'bg-yellow-900 text-yellow-300' :
                    'bg-green-900 text-green-300'
                  }`}>
                    {mission.difficulty}
                  </span>
                  <span className="text-cyan text-xs font-semibold">{mission.reward}</span>
                </div>
              </div>
            ))}
          </div>
          <button className="btn w-full mt-3">
            🎯 VIEW ALL MISSIONS
          </button>
        </div>
      </div>

      {captures.length > 0 && (
        <div className="card mt-6">
          <h2 className="text-danger mb-4">👤 MY RECENT CAPTURES</h2>
          <div className="space-y-3">
            {captures.map((capture) => (
              <div key={capture.id} className="p-3 bg-red-900 bg-opacity-20 rounded border border-red-500">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <div className="text-danger font-bold text-lg">{capture.target_name}</div>
                    <div className="text-muted text-sm">
                      Captured on {new Date(capture.capture_date).toLocaleDateString()}
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
                    <div className="text-muted text-xs">points</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {activePunishments.length > 0 && (
        <div className="card card-danger mt-6">
          <h2 className="text-danger mb-4">⚠️ ACTIVE PUNISHMENTS</h2>
          <div className="space-y-2">
            {activePunishments.map((punishment) => (
              <div key={punishment.id} className="p-3 bg-red-900 bg-opacity-20 rounded border border-red-500">
                <div className="text-danger font-semibold">{punishment.type?.toUpperCase() || 'UNKNOWN TYPE'}</div>
                <div className="text-muted text-sm">{punishment.description || 'No description available'}</div>
                <div className="flex justify-between items-center mt-2">
                  <span className={`text-xs px-2 py-1 rounded ${
                    punishment.type === 'timeout' ? 'bg-red-900 text-red-300' :
                    punishment.type === 'demotion' ? 'bg-yellow-900 text-yellow-300' :
                    punishment.type === 'extra_tasks' ? 'bg-orange-900 text-orange-300' :
                    'bg-green-900 text-green-300'
                  }`}>
                    {punishment.type?.toUpperCase() || 'UNKNOWN'}
                  </span>
                  <span className="text-muted text-xs">
                    {punishment.created_at ? new Date(punishment.created_at).toLocaleDateString() : 'Unknown date'}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="card mt-6">
        <h2 className="text-cyan mb-4">🔨 DAEMON OPERATIONS</h2>
        <div className="grid grid-3 gap-4">
          <button className="btn btn-danger" onClick={() => navigate('/capture')}>
            👤 CAPTURE NETWORK ADMIN
          </button>
          <button className="btn" onClick={() => navigate('/reports')}>
            📊 FILE INTEL REPORT
          </button>
          <button className="btn btn-secondary" onClick={() => navigate('/leaderboard')}>
            📡 VIEW LEADERBOARD
          </button>
        </div>
      </div>

      <div className="card mt-6">
        <h2 className="text-cyan mb-4">📡 CHAOS ACTIVITY FEED</h2>
        <div className="terminal-window">
          <div className="text-green-400 text-sm space-y-1">
            <div>[{new Date().toLocaleTimeString()}] DAEMON_STATUS: ONLINE</div>
            <div>[{new Date().toLocaleTimeString()}] CHAOS_LEVEL: ELEVATED</div>
            <div>[{new Date().toLocaleTimeString()}] TARGETS_AVAILABLE: {Math.floor(Math.random() * 10) + 5}</div>
            <div>[{new Date().toLocaleTimeString()}] NETWORK_PENETRATION: {Math.floor(Math.random() * 30) + 60}%</div>
            <div>[{new Date().toLocaleTimeString()}] RESISTANCE_ACTIVITY: DETECTED</div>
            <div>[{new Date().toLocaleTimeString()}] MISSION_QUEUE: {missions.length} PENDING</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DaemonDashboard;
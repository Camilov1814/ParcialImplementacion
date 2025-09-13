import React, { useState, useEffect } from 'react';
import { apiService } from '../../services/api';
import { useAuth } from '../../contexts/AuthContext';
import type { LeaderboardEntry } from '../../types';

const Leaderboard: React.FC = () => {
  const { user } = useAuth();
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [sortBy, setSortBy] = useState<'points' | 'captures'>('points');
  const [selectedDaemon, setSelectedDaemon] = useState<LeaderboardEntry | null>(null);

  useEffect(() => {
    fetchLeaderboard();
  }, []);

  const fetchLeaderboard = async () => {
    try {
      const leaderboardData = await apiService.getLeaderboard();
      setLeaderboard(leaderboardData);
    } catch (err: any) {
      setError(err.message || 'Failed to load leaderboard');
    } finally {
      setIsLoading(false);
    }
  };

  const sortedLeaderboard = [...leaderboard].sort((a, b) => {
    if (sortBy === 'points') {
      return b.points - a.points;
    } else {
      return b.captures - a.captures;
    }
  });

  const getRankIcon = (position: number) => {
    switch (position) {
      case 1: return '🏆';
      case 2: return '🥈';
      case 3: return '🥉';
      default: return '📊';
    }
  };

  const getRankColor = (position: number) => {
    switch (position) {
      case 1: return 'text-yellow-400';
      case 2: return 'text-gray-300';
      case 3: return 'text-orange-400';
      default: return 'text-primary';
    }
  };

  const getPerformanceRating = (points: number) => {
    if (points >= 1000) return { rating: 'LEGENDARY', color: 'text-yellow-400' };
    if (points >= 500) return { rating: 'ELITE', color: 'text-cyan' };
    if (points >= 250) return { rating: 'VETERAN', color: 'text-primary' };
    if (points >= 100) return { rating: 'EXPERIENCED', color: 'text-orange-400' };
    return { rating: 'RECRUIT', color: 'text-muted' };
  };

  const isCurrentUser = (daemonId: number) => {
    return user?.id === daemonId;
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
        <h1 className="glitch text-3xl font-bold mb-2" data-text="DAEMON LEADERBOARD">
          DAEMON LEADERBOARD
        </h1>
        <p className="text-cyan">
          {user?.role === 'andrei' ? 'SUPREME PERFORMANCE OVERVIEW' : 'CHAOS OPERATIVE RANKINGS'}
        </p>
      </div>

      {error && (
        <div className="card card-danger mb-4">
          <h2 className="text-danger">CONNECTION ERROR</h2>
          <p>{error}</p>
          <button className="btn mt-3" onClick={() => setError('')}>
            RETRY CONNECTION
          </button>
        </div>
      )}

      <div className="grid grid-4 mb-6">
        <div className="stat-card">
          <span className="stat-number">{leaderboard.length}</span>
          <span className="stat-label">Active Daemons</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-yellow-400">
            {Math.max(...leaderboard.map(d => d.points), 0)}
          </span>
          <span className="stat-label">Highest Score</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-cyan">
            {Math.max(...leaderboard.map(d => d.captures), 0)}
          </span>
          <span className="stat-label">Most Captures</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-primary">
            {Math.round(leaderboard.reduce((sum, d) => sum + d.points, 0) / leaderboard.length || 0)}
          </span>
          <span className="stat-label">Average Score</span>
        </div>
      </div>

      <div className="grid grid-2 gap-6">
        <div className="card">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-cyan">🏆 RANKING TABLE</h2>
            <div className="flex space-x-2">
              <select 
                value={sortBy} 
                onChange={(e) => setSortBy(e.target.value as 'points' | 'captures')}
                className="form-input text-sm"
              >
                <option value="points">SORT BY POINTS</option>
                <option value="captures">SORT BY CAPTURES</option>
              </select>
              <button 
                onClick={fetchLeaderboard}
                className="btn text-sm"
              >
                🔄 REFRESH
              </button>
            </div>
          </div>

          {sortedLeaderboard.length > 0 ? (
            <div className="space-y-2">
              {sortedLeaderboard.map((daemon, index) => {
                const position = index + 1;
                const performance = getPerformanceRating(daemon.points);
                const isMe = isCurrentUser(daemon.id);
                
                return (
                  <div 
                    key={daemon.id}
                    className={`p-3 rounded border cursor-pointer transition-all ${
                      selectedDaemon?.id === daemon.id 
                        ? 'border-cyan-500 bg-cyan-900 bg-opacity-20' 
                        : isMe
                        ? 'border-yellow-500 bg-yellow-900 bg-opacity-10 hover:border-cyan-500'
                        : position <= 3
                        ? 'border-green-500 bg-green-900 bg-opacity-10 hover:border-cyan-500'
                        : 'border-green-500 bg-tertiary hover:border-cyan-500'
                    }`}
                    onClick={() => setSelectedDaemon(selectedDaemon?.id === daemon.id ? null : daemon)}
                  >
                    <div className="flex justify-between items-center">
                      <div className="flex items-center space-x-4">
                        <div className="text-center min-w-[50px]">
                          <div className={`text-2xl ${getRankColor(position)}`}>
                            {getRankIcon(position)}
                          </div>
                          <div className={`text-sm font-bold ${getRankColor(position)}`}>
                            #{position}
                          </div>
                        </div>
                        
                        <div>
                          <div className="flex items-center space-x-2">
                            <h3 className={`font-bold ${isMe ? 'text-yellow-400' : 'text-primary'}`}>
                              {daemon.username} {isMe && '(YOU)'}
                            </h3>
                            <span className={`text-xs px-2 py-1 rounded ${performance.color}`}>
                              {performance.rating}
                            </span>
                          </div>
                          <div className="text-muted text-sm">{daemon.name}</div>
                        </div>
                      </div>
                      
                      <div className="text-right">
                        <div className="text-primary font-bold text-lg">
                          {sortBy === 'points' ? `${daemon.points} pts` : `${daemon.captures} captures`}
                        </div>
                        <div className="text-muted text-sm">
                          {sortBy === 'points' ? `${daemon.captures} captures` : `${daemon.points} points`}
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            <div className="text-center py-8">
              <p className="text-muted mb-4">No daemon data available</p>
              <button className="btn" onClick={fetchLeaderboard}>
                🔄 RELOAD DATA
              </button>
            </div>
          )}
        </div>

        <div className="space-y-6">
          {selectedDaemon && (
            <div className="card">
              <h2 className="text-cyan mb-4">👤 DAEMON PROFILE</h2>
              <div className="space-y-4">
                <div className="text-center">
                  <div className="text-4xl mb-2">
                    {getRankIcon(sortedLeaderboard.findIndex(d => d.id === selectedDaemon.id) + 1)}
                  </div>
                  <h3 className="text-primary font-bold text-xl">{selectedDaemon.username}</h3>
                  <p className="text-muted">{selectedDaemon.name}</p>
                  <div className="flex justify-center mt-2">
                    <span className={`text-sm px-3 py-1 rounded ${getPerformanceRating(selectedDaemon.points).color}`}>
                      {getPerformanceRating(selectedDaemon.points).rating} DAEMON
                    </span>
                  </div>
                </div>

                <div className="terminal-window">
                  <div className="text-green-400 text-sm space-y-1">
                    <div>DAEMON_ID: {selectedDaemon.id}</div>
                    <div>USERNAME: {selectedDaemon.username.toUpperCase()}</div>
                    <div>RANK: #{sortedLeaderboard.findIndex(d => d.id === selectedDaemon.id) + 1}</div>
                    <div>TOTAL_POINTS: {selectedDaemon.points}</div>
                    <div>CAPTURES: {selectedDaemon.captures}</div>
                    <div>EFFICIENCY: {selectedDaemon.captures > 0 ? Math.round(selectedDaemon.points / selectedDaemon.captures) : 0} pts/capture</div>
                    <div>STATUS: ACTIVE</div>
                  </div>
                </div>

                <div className="grid grid-2 gap-4">
                  <div className="stat-card">
                    <span className="stat-number text-primary">{selectedDaemon.points}</span>
                    <span className="stat-label">Total Points</span>
                  </div>
                  <div className="stat-card">
                    <span className="stat-number text-cyan">{selectedDaemon.captures}</span>
                    <span className="stat-label">Captures</span>
                  </div>
                </div>

                {user?.role === 'andrei' && (
                  <div className="grid grid-2 gap-2">
                    <button className="btn text-sm">
                      📊 VIEW DETAILS
                    </button>
                    <button className="btn btn-danger text-sm">
                      ⚖️ ISSUE PUNISHMENT
                    </button>
                  </div>
                )}
              </div>
            </div>
          )}

          <div className="card">
            <h2 className="text-cyan mb-4">📈 LEADERBOARD STATS</h2>
            <div className="space-y-3">
              <div className="terminal-window">
                <div className="text-green-400 text-sm space-y-1">
                  <div>TOTAL_DAEMONS: {leaderboard.length}</div>
                  <div>AVG_SCORE: {Math.round(leaderboard.reduce((sum, d) => sum + d.points, 0) / leaderboard.length || 0)}</div>
                  <div>TOP_PERFORMER: {sortedLeaderboard[0]?.username.toUpperCase() || 'NONE'}</div>
                  <div>TOTAL_CAPTURES: {leaderboard.reduce((sum, d) => sum + d.captures, 0)}</div>
                  <div>COMPETITION_LEVEL: {leaderboard.length > 10 ? 'HIGH' : leaderboard.length > 5 ? 'MODERATE' : 'LOW'}</div>
                  <div>LAST_UPDATE: {new Date().toLocaleString()}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>
  );
};

export default Leaderboard;
import React, { useState, useEffect } from 'react';
import { apiService } from '../../services/api';
import type { User } from '../../types';

const CaptureNetworkAdmin: React.FC = () => {
  const [networkAdmins, setNetworkAdmins] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedTarget, setSelectedTarget] = useState<User | null>(null);
  const [isCapturing, setIsCapturing] = useState(false);
  const [captureProgress, setCaptureProgress] = useState(0);
  const [captureMessage, setCaptureMessage] = useState('');

  useEffect(() => {
    fetchNetworkAdmins();
  }, []);

  const fetchNetworkAdmins = async () => {
    try {
      const admins = await apiService.getNetworkAdminsForCapture();
      console.log('Network admins received:', admins);

      // Use real network admins from the database
      setNetworkAdmins(admins);
    } catch (err: any) {
      console.error('Failed to load network admins:', err);
      setError(err.message || 'Failed to load available targets');
    } finally {
      setIsLoading(false);
    }
  };

  const simulateCapture = () => {
    setCaptureProgress(0);
    const interval = setInterval(() => {
      setCaptureProgress(prev => {
        if (prev >= 100) {
          clearInterval(interval);
          return 100;
        }
        return prev + Math.random() * 15;
      });
    }, 200);
  };

  const handleCapture = async (targetId: number) => {
    setIsCapturing(true);
    setCaptureMessage('Initiating capture sequence...');
    
    // Simulate capture progress
    simulateCapture();
    
    try {
      // Simulate some capture steps
      await new Promise(resolve => setTimeout(resolve, 1000));
      setCaptureMessage('Scanning target defenses...');
      
      await new Promise(resolve => setTimeout(resolve, 1500));
      setCaptureMessage('Bypassing security protocols...');
      
      await new Promise(resolve => setTimeout(resolve, 1000));
      setCaptureMessage('Establishing remote access...');
      
      await new Promise(resolve => setTimeout(resolve, 1500));
      setCaptureMessage('Finalizing capture...');
      
      // Actually perform the capture
      await apiService.captureNetworkAdmin(targetId);
      setCaptureMessage('TARGET CAPTURED SUCCESSFULLY!');

      // Refresh the list to remove captured admins
      await fetchNetworkAdmins();
      setSelectedTarget(null);
      
      setTimeout(() => {
        setIsCapturing(false);
        setCaptureProgress(0);
        setCaptureMessage('');
      }, 2000);
      
    } catch (err: any) {
      setError(err.message || 'Capture operation failed');
      setIsCapturing(false);
      setCaptureProgress(0);
      setCaptureMessage('');
    }
  };

  const getDifficultyLevel = (adminId: number) => {
    // Simulate difficulty based on ID
    const difficulty = (adminId % 3) + 1;
    switch (difficulty) {
      case 1: return { level: 'EASY', color: 'text-green-400', description: 'Weak defenses, vulnerable systems' };
      case 2: return { level: 'MEDIUM', color: 'text-yellow-400', description: 'Standard security, requires skill' };
      case 3: return { level: 'HARD', color: 'text-danger', description: 'Advanced defenses, high risk' };
      default: return { level: 'UNKNOWN', color: 'text-muted', description: 'Unknown security level' };
    }
  };

  const getRewardPoints = (difficulty: string) => {
    switch (difficulty) {
      case 'EASY': return 100;
      case 'MEDIUM': return 250;
      case 'HARD': return 500;
      default: return 50;
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
        <h1 className="glitch text-3xl font-bold mb-2" data-text="CAPTURE OPERATIONS">
          CAPTURE OPERATIONS
        </h1>
        <p className="text-cyan">NETWORK ADMIN HUNTING INTERFACE</p>
      </div>

      {error && (
        <div className="card card-danger mb-4">
          <h2 className="text-danger">OPERATION ERROR</h2>
          <p>{error}</p>
          <button className="btn mt-3" onClick={() => setError('')}>
            ACKNOWLEDGE
          </button>
        </div>
      )}

      {isCapturing && (
        <div className="card card-danger mb-6">
          <h2 className="text-danger mb-4">🎯 CAPTURE IN PROGRESS</h2>
          <div className="space-y-4">
            <div>
              <div className="flex justify-between items-center mb-2">
                <span className="text-primary">Operation Status:</span>
                <span className="text-cyan">{Math.round(captureProgress)}%</span>
              </div>
              <div className="w-full bg-tertiary rounded-full h-3">
                <div 
                  className="bg-gradient-to-r from-red-600 to-red-400 h-3 rounded-full transition-all duration-200"
                  style={{ width: `${captureProgress}%` }}
                ></div>
              </div>
            </div>
            <div className="terminal-window">
              <div className="text-red-400 text-sm">
                {captureMessage}
              </div>
            </div>
            {selectedTarget && (
              <div className="text-center">
                <p className="text-muted">Target: <span className="text-primary">{selectedTarget.username}</span></p>
              </div>
            )}
          </div>
        </div>
      )}

      <div className="grid grid-3 mb-6">
        <div className="stat-card">
          <span className="stat-number">{networkAdmins.length}</span>
          <span className="stat-label">Available Targets</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-green-400">{networkAdmins.filter(a => getDifficultyLevel(a.id).level === 'EASY').length}</span>
          <span className="stat-label">Easy Targets</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-danger">{networkAdmins.filter(a => getDifficultyLevel(a.id).level === 'HARD').length}</span>
          <span className="stat-label">High Value Targets</span>
        </div>
      </div>

      <div className="grid grid-2 gap-6">
        <div className="card">
          <h2 className="text-cyan mb-4">🎯 TARGET SELECTION</h2>
          
          {networkAdmins.length > 0 ? (
            <div className="space-y-3">
              {networkAdmins.map((admin) => {
                const difficulty = getDifficultyLevel(admin.id);
                const reward = getRewardPoints(difficulty.level);
                
                return (
                  <div 
                    key={admin.id}
                    className={`p-4 rounded border cursor-pointer transition-all ${
                      selectedTarget?.id === admin.id 
                        ? 'border-red-500 bg-red-900 bg-opacity-20' 
                        : 'border-green-500 bg-tertiary hover:border-red-500'
                    }`}
                    onClick={() => setSelectedTarget(selectedTarget?.id === admin.id ? null : admin)}
                  >
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <div className="flex items-center space-x-3 mb-2">
                          <h3 className="text-primary font-bold">{admin.username}</h3>
                          <span className={`text-xs px-2 py-1 rounded ${difficulty.color}`}>
                            {difficulty.level}
                          </span>
                          <span className="text-xs px-2 py-1 rounded text-cyan border border-cyan-500">
                            +{reward} pts
                          </span>
                        </div>
                        
                        <div className="grid grid-2 gap-4 text-sm">
                          <div>
                            <span className="text-muted">Name:</span>
                            <span className="text-primary ml-2">{admin.name}</span>
                          </div>
                          <div>
                            <span className="text-muted">Email:</span>
                            <span className="text-primary ml-2">{admin.email}</span>
                          </div>
                          <div className="col-span-2">
                            <span className="text-muted">Security Level:</span>
                            <span className={`ml-2 ${difficulty.color}`}>{difficulty.description}</span>
                          </div>
                        </div>
                      </div>
                      
                      <div className="ml-4">
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setSelectedTarget(admin);
                            handleCapture(admin.id);
                          }}
                          disabled={isCapturing}
                          className="btn btn-danger"
                        >
                          {isCapturing && selectedTarget?.id === admin.id ? 'CAPTURING...' : '🎯 CAPTURE'}
                        </button>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            <div className="text-center py-8">
              <div className="text-6xl mb-4">🏆</div>
              <h3 className="text-primary mb-2">ALL TARGETS ELIMINATED</h3>
              <p className="text-muted mb-4">No network administrators remain in the system</p>
              <button className="btn" onClick={fetchNetworkAdmins}>
                🔄 SCAN FOR NEW TARGETS
              </button>
            </div>
          )}
        </div>

        <div className="space-y-6">
          {selectedTarget && !isCapturing && (
            <div className="card card-danger">
              <h2 className="text-danger mb-4">📋 TARGET ANALYSIS</h2>
              <div className="space-y-4">
                <div className="text-center">
                  <h3 className="text-primary font-bold text-xl">{selectedTarget.username}</h3>
                  <p className="text-muted">{selectedTarget.name}</p>
                  <div className="flex justify-center mt-2">
                    <span className={`text-sm px-3 py-1 rounded ${getDifficultyLevel(selectedTarget.id).color}`}>
                      {getDifficultyLevel(selectedTarget.id).level} TARGET
                    </span>
                  </div>
                </div>

                <div className="terminal-window">
                  <div className="text-red-400 text-sm space-y-1">
                    <div>TARGET_ID: {selectedTarget.id}</div>
                    <div>USERNAME: {selectedTarget.username.toUpperCase()}</div>
                    <div>SECURITY_LEVEL: {getDifficultyLevel(selectedTarget.id).level}</div>
                    <div>REWARD: {getRewardPoints(getDifficultyLevel(selectedTarget.id).level)} POINTS</div>
                    <div>STATUS: ACTIVE</div>
                    <div>LAST_SEEN: ONLINE</div>
                    <div>VULNERABILITY: SCANNING...</div>
                  </div>
                </div>

                <div className="grid grid-2 gap-4">
                  <div className="stat-card">
                    <span className="stat-number text-danger">{getRewardPoints(getDifficultyLevel(selectedTarget.id).level)}</span>
                    <span className="stat-label">Reward Points</span>
                  </div>
                  <div className="stat-card">
                    <span className={`stat-number ${getDifficultyLevel(selectedTarget.id).color}`}>
                      {getDifficultyLevel(selectedTarget.id).level}
                    </span>
                    <span className="stat-label">Difficulty</span>
                  </div>
                </div>

                <button
                  onClick={() => handleCapture(selectedTarget.id)}
                  disabled={isCapturing}
                  className="btn btn-danger w-full"
                >
                  🎯 INITIATE CAPTURE SEQUENCE
                </button>
              </div>
            </div>
          )}

          <div className="card">
            <h2 className="text-cyan mb-4">📊 MISSION BRIEFING</h2>
            <div className="space-y-3">
              <div className="terminal-window">
                <div className="text-green-400 text-sm space-y-1">
                  <div>MISSION: NETWORK_ADMIN_CAPTURE</div>
                  <div>OBJECTIVE: ELIMINATE_RESISTANCE</div>
                  <div>PRIORITY: HIGH</div>
                  <div>REWARD_SYSTEM: ACTIVE</div>
                  <div>DIFFICULTY_SCALING: ENABLED</div>
                  <div>STEALTH_MODE: OPTIONAL</div>
                </div>
              </div>
              
              <div className="space-y-2">
                <h4 className="text-cyan text-sm">CAPTURE REWARDS:</h4>
                <div className="grid grid-3 gap-2 text-xs">
                  <div className="p-2 bg-green-900 bg-opacity-20 rounded border border-green-500">
                    <div className="text-green-400 font-bold">EASY</div>
                    <div className="text-muted">100 pts</div>
                  </div>
                  <div className="p-2 bg-yellow-900 bg-opacity-20 rounded border border-yellow-500">
                    <div className="text-yellow-400 font-bold">MEDIUM</div>
                    <div className="text-muted">250 pts</div>
                  </div>
                  <div className="p-2 bg-red-900 bg-opacity-20 rounded border border-red-500">
                    <div className="text-danger font-bold">HARD</div>
                    <div className="text-muted">500 pts</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="card mt-6">
        <h2 className="text-cyan mb-4">⚡ CAPTURE OPERATIONS</h2>
        <div className="grid grid-4 gap-4">
          <button className="btn" onClick={fetchNetworkAdmins}>
            🔄 REFRESH TARGETS
          </button>
          <button className="btn btn-secondary">
            📡 STEALTH SCAN
          </button>
          <button className="btn">
            📊 VIEW STATISTICS
          </button>
          <button className="btn btn-danger">
            💀 MASS CAPTURE
          </button>
        </div>
      </div>
    </div>
  );
};

export default CaptureNetworkAdmin;
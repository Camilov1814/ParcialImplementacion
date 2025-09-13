import React, { useState, useEffect } from 'react';
import { apiService } from '../../services/api';
import { useAuth } from '../../contexts/AuthContext';
import type { Punishment, User } from '../../types';

const PunishmentsManagement: React.FC = () => {
  const { user } = useAuth();
  const [punishments, setPunishments] = useState<Punishment[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedPunishment, setSelectedPunishment] = useState<Punishment | null>(null);
  const [detailedPunishment, setDetailedPunishment] = useState<Punishment | null>(null);
  const [loadingDetails, setLoadingDetails] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [newPunishment, setNewPunishment] = useState({
    target_id: 0,
    type: '',
    description: '',
    expires_at: ''
  });
  const [filterSeverity, setFilterSeverity] = useState<string>('all');
  const [statusUpdateLoading, setStatusUpdateLoading] = useState<number | null>(null);
  const [isModifying, setIsModifying] = useState(false);
  const [modifyForm, setModifyForm] = useState({
    description: '',
    expires_at: ''
  });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [punishmentsData, usersData] = await Promise.all([
        apiService.getPunishments(),
        user?.role === 'andrei' ? apiService.getUsers() : Promise.resolve([])
      ]);
      console.log('Punishments data received:', punishmentsData);
      
      // Load detailed data for each punishment to get descriptions
      const punishmentsWithDetails = await Promise.all(
        punishmentsData.map(async (punishment) => {
          try {
            const details = await apiService.getPunishmentDetails(punishment.id);
            return { ...punishment, description: details.description };
          } catch (err) {
            console.warn(`Failed to load details for punishment ${punishment.id}:`, err);
            return { ...punishment, description: 'Failed to load description' };
          }
        })
      );
      
      setPunishments(punishmentsWithDetails);
      setUsers(usersData);
    } catch (err: any) {
      setError(err.message || 'Failed to load punishments data');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreatePunishment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newPunishment.type.trim() || !newPunishment.description.trim() || !newPunishment.target_id) return;

    try {
      // Convert datetime-local format to RFC3339 if provided
      const punishmentData = {
        ...newPunishment,
        expires_at: newPunishment.expires_at ? new Date(newPunishment.expires_at).toISOString() : undefined
      };
      
      // Remove empty expires_at field
      if (!punishmentData.expires_at) {
        delete punishmentData.expires_at;
      }
      
      console.log('Creating punishment with data:', punishmentData);
      await apiService.createPunishment(punishmentData);
      setNewPunishment({
        target_id: 0,
        type: '',
        description: '',
        expires_at: ''
      });
      setIsCreating(false);
      await fetchData();
    } catch (err: any) {
      console.error('Punishment creation error:', err);
      console.error('Error response:', err.response?.data);
      setError(err.response?.data?.message || err.message || 'Failed to create punishment');
    }
  };

  const handleSelectPunishment = async (punishment: Punishment) => {
    if (selectedPunishment?.id === punishment.id) {
      // Deselect if clicking the same punishment
      setSelectedPunishment(null);
      setDetailedPunishment(null);
      return;
    }

    setSelectedPunishment(punishment);
    setLoadingDetails(true);
    
    try {
      console.log('Loading details for punishment ID:', punishment.id);
      const details = await apiService.getPunishmentDetails(punishment.id);
      console.log('Punishment details loaded:', details);
      console.log('Details description:', details.description);
      console.log('Details type:', typeof details.description);
      setDetailedPunishment(details);
    } catch (err: any) {
      console.error('Failed to load punishment details:', err);
      setError('Failed to load punishment details');
    } finally {
      setLoadingDetails(false);
    }
  };

  const handleMarkCompleted = async (punishmentId: number) => {
    setStatusUpdateLoading(punishmentId);
    try {
      await apiService.updatePunishment(punishmentId, { status: 'completed' });
      await fetchData(); // Refresh the list
      // Update selected punishment if it's the one being updated
      if (selectedPunishment?.id === punishmentId) {
        setSelectedPunishment({ ...selectedPunishment, status: 'completed' });
      }
    } catch (err: any) {
      console.error('Status update error:', err);
      setError(err.response?.data?.message || err.message || 'Failed to update punishment status');
    } finally {
      setStatusUpdateLoading(null);
    }
  };

  const handleStartModify = (punishment: Punishment) => {
    setModifyForm({
      description: punishment.description || '',
      expires_at: punishment.expires_at ? new Date(punishment.expires_at).toISOString().slice(0, 16) : ''
    });
    setIsModifying(true);
  };

  const handleModifySubmit = async (punishmentId: number) => {
    setStatusUpdateLoading(punishmentId);
    try {
      const updates: { description?: string; expires_at?: string } = {};
      
      if (modifyForm.description.trim()) {
        updates.description = modifyForm.description.trim();
      }
      
      if (modifyForm.expires_at) {
        updates.expires_at = new Date(modifyForm.expires_at).toISOString();
      }
      
      await apiService.updatePunishment(punishmentId, updates);
      await fetchData(); // Refresh the list
      setIsModifying(false);
      setModifyForm({ description: '', expires_at: '' });
    } catch (err: any) {
      console.error('Modify error:', err);
      setError(err.response?.data?.message || err.message || 'Failed to modify punishment');
    } finally {
      setStatusUpdateLoading(null);
    }
  };

  const filteredPunishments = punishments.filter(punishment => 
    filterSeverity === 'all' || punishment.type === filterSeverity
  );

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'timeout': return 'text-orange-400';
      case 'demotion': return 'text-danger';
      case 'extra_tasks': return 'text-cyan';
      case 'reward': return 'text-green-400';
      default: return 'text-muted';
    }
  };

  const getTypeBadge = (type: string) => {
    switch (type) {
      case 'timeout': return '⏰ TIMEOUT';
      case 'demotion': return '📉 DEMOTION';
      case 'extra_tasks': return '📝 EXTRA TASKS';
      case 'reward': return '🏆 REWARD';
      default: return '❓ UNKNOWN';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'text-danger';
      case 'completed': return 'text-cyan';
      case 'cancelled': return 'text-muted';
      default: return 'text-muted';
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'active': return '🔴 ACTIVE';
      case 'completed': return '✅ COMPLETED';
      case 'cancelled': return '🚫 CANCELLED';
      default: return '❓ UNKNOWN';
    }
  };

  const daemons = users.filter(u => u.role === 'daemon');

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
        <h1 className="glitch text-3xl font-bold mb-2" data-text="PUNISHMENT SYSTEM">
          PUNISHMENT SYSTEM
        </h1>
        <p className="text-cyan">
          {user?.role === 'andrei' ? 'SUPREME DISCIPLINARY CONTROL' : 'ACTIVE PUNISHMENT RECORDS'}
        </p>
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

      <div className="grid grid-4 mb-6">
        <div className="stat-card">
          <span className="stat-number">{punishments.length}</span>
          <span className="stat-label">Total Punishments</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-danger">{punishments.filter(p => p.status === 'active').length}</span>
          <span className="stat-label">Active</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-cyan">{punishments.filter(p => p.status === 'completed').length}</span>
          <span className="stat-label">Completed</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-danger">{punishments.filter(p => p.type === 'timeout' || p.type === 'demotion').length}</span>
          <span className="stat-label">Severe Punishments</span>
        </div>
      </div>

      <div className="grid grid-2 gap-6">
        <div className="card">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-cyan">⚖️ DISCIPLINARY RECORDS</h2>
            <div className="flex space-x-2">
              <select 
                value={filterSeverity} 
                onChange={(e) => setFilterSeverity(e.target.value)}
                className="form-input text-sm"
              >
                <option value="all">ALL TYPES</option>
                <option value="timeout">TIMEOUT</option>
                <option value="demotion">DEMOTION</option>
                <option value="extra_tasks">EXTRA TASKS</option>
                <option value="reward">REWARD</option>
              </select>
              {user?.role === 'andrei' && (
                <button 
                  onClick={() => setIsCreating(!isCreating)}
                  className="btn btn-danger text-sm"
                >
                  {isCreating ? 'CANCEL' : '⚖️ NEW PUNISHMENT'}
                </button>
              )}
            </div>
          </div>

          {filteredPunishments.length > 0 ? (
            <div className="space-y-3">
              {filteredPunishments.map((punishment) => (
                <div 
                  key={punishment.id}
                  className={`p-3 rounded border cursor-pointer transition-all ${
                    selectedPunishment?.id === punishment.id 
                      ? 'border-red-500 bg-red-900 bg-opacity-20' 
                      : 'border-orange-500 bg-tertiary hover:border-red-500'
                  }`}
                  onClick={() => handleSelectPunishment(punishment)}
                >
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <div className="flex items-center space-x-2 mb-2">
                        <h3 className="text-primary font-semibold">{getTypeBadge(punishment.type)}</h3>
                        <span className={`text-xs px-2 py-1 rounded ${getTypeColor(punishment.type)}`}>
                          {punishment.type.toUpperCase()}
                        </span>
                        <span className={`text-xs px-2 py-1 rounded ${getStatusColor(punishment.status)}`}>
                          {getStatusBadge(punishment.status)}
                        </span>
                      </div>
                      <p className="text-muted text-sm">
                        {punishment.description ? 
                          (punishment.description.length > 80 ? 
                            punishment.description.substring(0, 80) + '...' : 
                            punishment.description
                          ) : 
                          'No description available'
                        }
                      </p>
                      <div className="flex justify-between items-center mt-2">
                        <span className="text-muted text-xs">
                          Target: {punishment.target?.username || punishment.target_name || 'Unknown'}
                        </span>
                        <span className="text-muted text-xs">
                          {punishment.created_at ? new Date(punishment.created_at).toLocaleDateString() : 'Unknown date'}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <p className="text-muted mb-4">No punishments match the current filter</p>
              <button className="btn" onClick={fetchData}>
                🔄 REFRESH DATA
              </button>
            </div>
          )}
        </div>

        <div className="space-y-6">
          {isCreating && user?.role === 'andrei' && (
            <div className="card card-danger">
              <h2 className="text-danger mb-4">⚖️ ISSUE PUNISHMENT</h2>
              <form onSubmit={handleCreatePunishment} className="space-y-4">
                <div className="form-group">
                  <label className="form-label">TARGET_DAEMON:</label>
                  <select
                    value={newPunishment.target_id}
                    onChange={(e) => setNewPunishment({...newPunishment, target_id: parseInt(e.target.value)})}
                    className="form-input"
                    required
                  >
                    <option value={0}>Select a daemon...</option>
                    {daemons.map((daemon) => (
                      <option key={daemon.id} value={daemon.id}>
                        {daemon.username} - {daemon.name}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="form-group">
                  <label className="form-label">PUNISHMENT_TYPE:</label>
                  <select
                    value={newPunishment.type}
                    onChange={(e) => setNewPunishment({...newPunishment, type: e.target.value})}
                    className="form-input"
                    required
                  >
                    <option value="">Select punishment type...</option>
                    <option value="timeout">⏰ TIMEOUT - Temporary suspension</option>
                    <option value="demotion">📉 DEMOTION - Rank reduction</option>
                    <option value="extra_tasks">📝 EXTRA TASKS - Additional duties</option>
                    <option value="reward">🏆 REWARD - Positive reinforcement</option>
                  </select>
                </div>
                <div className="form-group">
                  <label className="form-label">DETAILED_DESCRIPTION:</label>
                  <textarea
                    value={newPunishment.description}
                    onChange={(e) => setNewPunishment({...newPunishment, description: e.target.value})}
                    className="form-input form-textarea"
                    placeholder="Detailed punishment description and requirements..."
                    required
                    rows={3}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">EXPIRES_AT (Optional):</label>
                  <input
                    type="datetime-local"
                    value={newPunishment.expires_at}
                    onChange={(e) => setNewPunishment({...newPunishment, expires_at: e.target.value})}
                    className="form-input"
                    placeholder="Leave empty for permanent punishment"
                  />
                </div>
                <div className="flex space-x-2">
                  <button type="submit" className="btn btn-danger flex-1">
                    ⚖️ ISSUE PUNISHMENT
                  </button>
                  <button 
                    type="button" 
                    onClick={() => setIsCreating(false)}
                    className="btn btn-secondary"
                  >
                    CANCEL
                  </button>
                </div>
              </form>
            </div>
          )}

          {selectedPunishment && (
            <div className="card card-danger">
              <h2 className="text-danger mb-4">📋 PUNISHMENT DETAILS</h2>
              {loadingDetails ? (
                <div className="text-center py-4">
                  <div className="loading-spinner mx-auto mb-2"></div>
                  <p className="text-cyan">Loading details...</p>
                </div>
              ) : detailedPunishment ? (
                <div className="space-y-4">
                  <div>
                    <h3 className="text-primary font-semibold mb-2">{getTypeBadge(detailedPunishment?.type || '')}</h3>
                    <div className="flex items-center space-x-2 mb-3">
                      <span className={`text-sm px-2 py-1 rounded ${getTypeColor(detailedPunishment?.type || '')}`}>
                        {detailedPunishment?.type?.toUpperCase() || 'UNKNOWN'}
                      </span>
                      <span className={`text-sm px-2 py-1 rounded ${getStatusColor(detailedPunishment?.status || '')}`}>
                        {getStatusBadge(detailedPunishment?.status || '')}
                      </span>
                    </div>
                  </div>

                  <div className="terminal-window">
                    <div className="text-red-400 text-sm space-y-1">
                      <div>PUNISHMENT_ID: {detailedPunishment?.id || 'N/A'}</div>
                      <div>TARGET: {detailedPunishment?.target?.username || 'UNKNOWN'}</div>
                      <div>TYPE: {detailedPunishment?.type?.toUpperCase() || 'UNKNOWN'}</div>
                      <div>STATUS: {detailedPunishment?.status?.toUpperCase() || 'UNKNOWN'}</div>
                      <div>ISSUED: {detailedPunishment?.created_at ? new Date(detailedPunishment.created_at).toLocaleString() : 'Unknown'}</div>
                      <div>ASSIGNED_BY: {detailedPunishment?.assigner?.username || 'SYSTEM'}</div>
                      {detailedPunishment?.expires_at && (
                        <div>EXPIRES: {new Date(detailedPunishment?.expires_at || '').toLocaleString()}</div>
                      )}
                    </div>
                  </div>

                  <div>
                    <h4 className="text-danger mb-2">PUNISHMENT DETAILS:</h4>
                    <p className="text-sm bg-red-900 bg-opacity-20 p-3 rounded border border-red-500">
                      {detailedPunishment?.description || 'No detailed description available'}
                    </p>
                  </div>

                  {user?.role === 'andrei' && detailedPunishment?.status === 'active' && !isModifying && (
                    <div className="grid grid-2 gap-2">
                      <button 
                        className="btn btn-secondary text-sm"
                        disabled={statusUpdateLoading === detailedPunishment?.id}
                        onClick={() => handleMarkCompleted(detailedPunishment?.id || 0)}
                      >
                        {statusUpdateLoading === detailedPunishment?.id ? 'UPDATING...' : '✅ MARK COMPLETED'}
                      </button>
                      <button 
                        className="btn btn-danger text-sm"
                        onClick={() => detailedPunishment && handleStartModify(detailedPunishment)}
                      >
                        📝 MODIFY
                      </button>
                    </div>
                  )}

                  {user?.role === 'andrei' && isModifying && detailedPunishment?.status === 'active' && (
                  <div className="space-y-4">
                    <div className="form-group">
                      <label className="form-label">MODIFY_DESCRIPTION:</label>
                      <textarea
                        value={modifyForm.description}
                        onChange={(e) => setModifyForm({...modifyForm, description: e.target.value})}
                        className="form-input form-textarea"
                        placeholder="Update punishment description..."
                        rows={3}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">MODIFY_EXPIRES_AT:</label>
                      <input
                        type="datetime-local"
                        value={modifyForm.expires_at}
                        onChange={(e) => setModifyForm({...modifyForm, expires_at: e.target.value})}
                        className="form-input"
                      />
                    </div>
                    <div className="grid grid-2 gap-2">
                      <button 
                        className="btn btn-secondary text-sm"
                        disabled={statusUpdateLoading === detailedPunishment?.id}
                        onClick={() => handleModifySubmit(detailedPunishment?.id || 0)}
                      >
                        {statusUpdateLoading === detailedPunishment?.id ? 'SAVING...' : '💾 SAVE CHANGES'}
                      </button>
                      <button 
                        className="btn btn-danger text-sm"
                        onClick={() => {
                          setIsModifying(false);
                          setModifyForm({ description: '', expires_at: '' });
                        }}
                      >
                        ❌ CANCEL
                      </button>
                    </div>
                  </div>
                  )}

                </div>
              ) : (
                <div className="text-center py-4">
                  <p className="text-muted">Select a punishment to view details</p>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

    </div>
  );
};

export default PunishmentsManagement;
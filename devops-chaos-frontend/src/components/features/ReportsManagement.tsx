import React, { useState, useEffect } from 'react';
import { apiService } from '../../services/api';
import { useAuth } from '../../contexts/AuthContext';
import type { Report } from '../../types';

const ReportsManagement: React.FC = () => {
  const { user } = useAuth();
  const [reports, setReports] = useState<Report[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedReport, setSelectedReport] = useState<Report | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [newReport, setNewReport] = useState({
    title: '',
    description: '',
    status: 'pending',
    priority: 'normal'
  });
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [statusUpdateLoading, setStatusUpdateLoading] = useState<number | null>(null);

  useEffect(() => {
    fetchReports();
  }, []);

  const fetchReports = async () => {
    try {
      const reportsData = await apiService.getReports();
      setReports(reportsData);
    } catch (err: any) {
      setError(err.message || 'Failed to load reports');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateReport = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newReport.title.trim() || !newReport.description.trim()) return;

    try {
      console.log('Creating report with data:', newReport);
      
      // Map user roles to report types
      const getReportType = (userRole: string) => {
        switch (userRole) {
          case 'network_admin': return 'resistance';
          case 'daemon': return 'capture';
          case 'andrei': return 'anonymous'; // Andrei can create any type, defaulting to anonymous
          default: return 'anonymous';
        }
      };

      // Only send the required fields that the backend expects
      const reportToSend = {
        title: newReport.title.trim(),
        description: newReport.description.trim(),
        type: getReportType(user?.role || 'anonymous')
      };
      
      console.log('Cleaned report data:', reportToSend);
      await apiService.createReport(reportToSend);
      setNewReport({ title: '', description: '', status: 'pending', priority: 'normal' });
      setIsCreating(false);
      await fetchReports();
    } catch (err: any) {
      console.error('Report creation error:', err);
      console.error('Error response:', err.response?.data);
      setError(err.response?.data?.message || err.message || 'Failed to create report');
    }
  };

  const handleStatusUpdate = async (reportId: number, status: string) => {
    setStatusUpdateLoading(reportId);
    try {
      await apiService.updateReportStatus(reportId, status);
      await fetchReports(); // Refresh the list
      // If this was the selected report, update it
      if (selectedReport?.id === reportId) {
        const updatedReport = reports.find(r => r.id === reportId);
        if (updatedReport) {
          setSelectedReport({ ...updatedReport, status });
        }
      }
    } catch (err: any) {
      console.error('Status update error:', err);
      setError(err.response?.data?.message || err.message || 'Failed to update report status');
    } finally {
      setStatusUpdateLoading(null);
    }
  };

  const filteredReports = reports.filter(report => 
    filterStatus === 'all' || report.status === filterStatus
  );

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return 'text-yellow-400';
      case 'approved': return 'text-cyan';
      case 'rejected': return 'text-danger';
      default: return 'text-muted';
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'pending': return '⏳ PENDING';
      case 'approved': return '✅ APPROVED';
      case 'rejected': return '❌ REJECTED';
      default: return '❓ UNKNOWN';
    }
  };

  const getReportIcon = (type?: string) => {
    if (!type) return '📊';
    switch (type) {
      case 'capture': return '💀';
      case 'resistance': return '🛡️';
      case 'anonymous': return '👑';
      default: return '📊';
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
        <h1 className="glitch text-3xl font-bold mb-2" data-text="INTEL REPORTS">
          INTEL REPORTS
        </h1>
        <p className="text-cyan">
          {user?.role === 'andrei' ? 'SUPREME INTELLIGENCE OVERVIEW' : 
           user?.role === 'daemon' ? 'DAEMON INTELLIGENCE NETWORK' : 
           'RESISTANCE INTELLIGENCE HUB'}
        </p>
      </div>

      {error && (
        <div className="card card-danger mb-4">
          <h2 className="text-danger">TRANSMISSION ERROR</h2>
          <p>{error}</p>
          <button className="btn mt-3" onClick={() => setError('')}>
            CLEAR ERROR
          </button>
        </div>
      )}

      <div className="grid grid-4 mb-6">
        <div className="stat-card">
          <span className="stat-number">{reports.length}</span>
          <span className="stat-label">Total Reports</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-yellow-400">{reports.filter(r => r.status === 'pending').length}</span>
          <span className="stat-label">Pending Review</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-cyan">{reports.filter(r => r.status === 'approved').length}</span>
          <span className="stat-label">Approved</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-danger">{reports.filter(r => r.status === 'rejected').length}</span>
          <span className="stat-label">Rejected</span>
        </div>
      </div>

      <div className="grid grid-2 gap-6">
        <div className="card">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-cyan">📋 INTELLIGENCE DATABASE</h2>
            <div className="flex space-x-2">
              <select 
                value={filterStatus} 
                onChange={(e) => setFilterStatus(e.target.value)}
                className="form-input text-sm"
              >
                <option value="all">ALL STATUS</option>
                <option value="pending">PENDING</option>
                <option value="approved">APPROVED</option>
                <option value="rejected">REJECTED</option>
              </select>
              <button 
                onClick={() => setIsCreating(!isCreating)}
                className="btn text-sm"
              >
                {isCreating ? 'CANCEL' : '📝 NEW REPORT'}
              </button>
            </div>
          </div>

          {filteredReports.length > 0 ? (
            <div className="space-y-3">
              {filteredReports.map((report) => (
                <div 
                  key={report.id}
                  className={`p-3 rounded border cursor-pointer transition-all ${
                    selectedReport?.id === report.id 
                      ? 'border-cyan-500 bg-cyan-900 bg-opacity-20' 
                      : 'border-green-500 bg-tertiary hover:border-cyan-500'
                  }`}
                  onClick={() => setSelectedReport(selectedReport?.id === report.id ? null : report)}
                >
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <div className="flex items-center space-x-2 mb-2">
                        <span className="text-lg">{getReportIcon(report.type)}</span>
                        <h3 className="text-primary font-semibold">{report.title}</h3>
                        <span className={`text-xs px-2 py-1 rounded ${getStatusColor(report.status)}`}>
                          {getStatusBadge(report.status)}
                        </span>
                      </div>
                      <p className="text-muted text-sm">{report.description ? report.description.substring(0, 100) + '...' : 'No description available'}</p>
                      <div className="flex justify-between items-center mt-2">
                        <span className="text-muted text-xs">ID: #{report.id}</span>
                        <span className="text-muted text-xs">
                          {report.createdAt ? new Date(report.createdAt).toLocaleDateString() : 'Unknown date'}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <p className="text-muted mb-4">No reports match the current filter</p>
              <button className="btn" onClick={fetchReports}>
                🔄 REFRESH DATA
              </button>
            </div>
          )}
        </div>

        <div className="space-y-6">
          {isCreating && (
            <div className="card">
              <h2 className="text-cyan mb-4">📝 CREATE INTEL REPORT</h2>
              <form onSubmit={handleCreateReport} className="space-y-4">
                <div className="form-group">
                  <label className="form-label">REPORT_TITLE:</label>
                  <input
                    type="text"
                    value={newReport.title}
                    onChange={(e) => setNewReport({...newReport, title: e.target.value})}
                    className="form-input"
                    placeholder="Brief description of the incident..."
                    required
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">DETAILED_REPORT:</label>
                  <textarea
                    value={newReport.description}
                    onChange={(e) => setNewReport({...newReport, description: e.target.value})}
                    className="form-input form-textarea"
                    placeholder="Provide detailed intelligence information..."
                    required
                    rows={4}
                  />
                </div>
                
                {user?.role === 'andrei' && (
                  <>
                    <div className="form-group">
                      <label className="form-label">STATUS:</label>
                      <select
                        value={newReport.status}
                        onChange={(e) => setNewReport({...newReport, status: e.target.value})}
                        className="form-input"
                      >
                        <option value="pending">PENDING</option>
                        <option value="active">ACTIVE</option>
                        <option value="completed">COMPLETED</option>
                        <option value="rejected">REJECTED</option>
                      </select>
                    </div>
                    <div className="form-group">
                      <label className="form-label">PRIORITY:</label>
                      <select
                        value={newReport.priority}
                        onChange={(e) => setNewReport({...newReport, priority: e.target.value})}
                        className="form-input"
                      >
                        <option value="low">LOW</option>
                        <option value="normal">NORMAL</option>
                        <option value="high">HIGH</option>
                        <option value="critical">CRITICAL</option>
                      </select>
                    </div>
                  </>
                )}
                
                <div className="flex space-x-2">
                  <button type="submit" className="btn flex-1">
                    📡 TRANSMIT REPORT
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

          {selectedReport && (
            <div className="card">
              <h2 className="text-cyan mb-4">📊 REPORT ANALYSIS</h2>
              <div className="space-y-4">
                <div>
                  <h3 className="text-primary font-semibold mb-2">{selectedReport.title}</h3>
                  <div className="flex items-center space-x-2 mb-3">
                    <span className={`text-sm px-2 py-1 rounded ${getStatusColor(selectedReport.status)}`}>
                      {getStatusBadge(selectedReport.status)}
                    </span>
                    <span className="text-muted text-sm">ID: #{selectedReport.id}</span>
                  </div>
                </div>

                <div className="terminal-window">
                  <div className="text-green-400 text-sm space-y-1">
                    <div>REPORT_ID: {selectedReport.id}</div>
                    <div>STATUS: {selectedReport.status.toUpperCase()}</div>
                    <div>CREATED: {selectedReport.createdAt ? new Date(selectedReport.createdAt).toLocaleString() : 'Unknown'}</div>
                    <div>TYPE: {selectedReport.type?.toUpperCase() || 'UNKNOWN'}</div>
                    <div>AUTHOR: {selectedReport.author?.username || 'CLASSIFIED'}</div>
                  </div>
                </div>

                <div>
                  <h4 className="text-cyan mb-2">DETAILED DESCRIPTION:</h4>
                  <p className="text-sm bg-tertiary p-3 rounded border border-green-500">
                    {selectedReport.description || 'No detailed description available'}
                  </p>
                </div>

                {user?.role === 'andrei' && selectedReport.status === 'pending' && (
                  <div className="grid grid-2 gap-2">
                    <button 
                      className="btn btn-secondary text-sm"
                      disabled={statusUpdateLoading === selectedReport.id}
                      onClick={() => handleStatusUpdate(selectedReport.id, 'approved')}
                    >
                      {statusUpdateLoading === selectedReport.id ? 'UPDATING...' : '✅ APPROVE'}
                    </button>
                    <button 
                      className="btn btn-danger text-sm"
                      disabled={statusUpdateLoading === selectedReport.id}
                      onClick={() => handleStatusUpdate(selectedReport.id, 'rejected')}
                    >
                      {statusUpdateLoading === selectedReport.id ? 'UPDATING...' : '❌ REJECT'}
                    </button>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </div>

    </div>
  );
};

export default ReportsManagement;
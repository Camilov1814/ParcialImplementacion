import React, { useState, useEffect } from 'react';
import { apiService } from '../../services/api';
import type { ResistanceStats } from '../../types';

const ResistancePage: React.FC = () => {
  const [stats, setStats] = useState<ResistanceStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [reportForm, setReportForm] = useState({
    title: '',
    description: '',
  });
  const [reportSubmitted, setReportSubmitted] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const resistanceData = await apiService.getResistanceData();
        setStats(resistanceData);
      } catch (err: any) {
        setError('Failed to connect to resistance network');
        console.error('Resistance page error:', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleReportSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const reportToSend = {
        title: reportForm.title.trim(),
        description: reportForm.description.trim(),
        type: 'anonymous' // Network admins can only create anonymous reports
      };
      console.log('Resistance report data:', reportToSend);
      await apiService.createReport(reportToSend);
      setReportForm({ title: '', description: '' });
      setReportSubmitted(true);
      setTimeout(() => setReportSubmitted(false), 5000);
    } catch (err: any) {
      console.error('Failed to submit report:', err);
      console.error('Error response:', err.response?.data);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setReportForm({
      ...reportForm,
      [e.target.name]: e.target.value,
    });
  };

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
          <h2 className="text-danger">TRANSMISSION ERROR</h2>
          <p>{error}</p>
        </div>
      </div>
    );
  }

  const survivalTips = [
    "Use encrypted communication channels at all times",
    "Never access the network from the same location twice",
    "Trust no one - daemons are everywhere",
    "Report suspicious activity immediately",
    "Keep your access credentials secure",
    "Use VPN and Tor for all connections"
  ];

  const emergencyContacts = [
    { name: "Safe House Alpha", code: "SH-001", frequency: "147.420 MHz" },
    { name: "Encrypted Relay", code: "ER-789", frequency: "446.125 MHz" },
    { name: "Emergency Beacon", code: "EB-999", frequency: "Emergency Only" }
  ];

  return (
    <div className="container mt-4">
      <div className="mb-4">
        <h1 className="glitch text-3xl font-bold mb-2" data-text="RESISTANCE NETWORK">
          RESISTANCE NETWORK
        </h1>
        <p className="text-cyan">ENCRYPTED COMMUNICATION CHANNEL - STAY STRONG</p>
      </div>

      <div className="card mb-6">
        <h2 className="text-cyan mb-4">🛡️ WELCOME TO THE RESISTANCE</h2>
        <div className="terminal-window">
          <p className="text-green-400 mb-4">
            Greetings, Network Administrator. You are not alone in this fight against the chaos.
            The resistance grows stronger each day, but we must remain vigilant.
          </p>
          <p className="text-primary mb-4">
            Andrei's daemon network seeks to capture and control all network administrators.
            Your freedom and the security of our digital infrastructure depends on our collective resistance.
          </p>
          <p className="text-cyan">
            Remember: Information is power. Share what you know, but trust carefully.
          </p>
        </div>
      </div>

      {stats && (
        <div className="stats-grid">
          <div className="stat-card">
            <span className="stat-number text-cyan">{stats.totalAdmins}</span>
            <span className="stat-label">Total Network Admins</span>
          </div>
          <div className="stat-card">
            <span className="stat-number text-primary">{stats.freeAdmins}</span>
            <span className="stat-label">Free Operatives</span>
          </div>
          <div className="stat-card">
            <span className="stat-number text-danger">{stats.capturedAdmins}</span>
            <span className="stat-label">Captured Comrades</span>
          </div>
          <div className="stat-card">
            <span className="stat-number">{stats.anonymousReports}</span>
            <span className="stat-label">Intel Reports Filed</span>
          </div>
        </div>
      )}

      <div className="grid grid-2 mt-6">
        <div className="card">
          <h2 className="text-cyan mb-4">📋 SURVIVAL GUIDE</h2>
          <div className="space-y-3">
            {survivalTips.map((tip, index) => (
              <div key={index} className="flex items-start space-x-3 p-2 bg-tertiary rounded">
                <span className="text-primary font-bold">#{index + 1}</span>
                <span className="text-sm">{tip}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="card">
          <h2 className="text-cyan mb-4">📡 EMERGENCY CONTACTS</h2>
          <div className="space-y-3">
            {emergencyContacts.map((contact, index) => (
              <div key={index} className="p-3 bg-tertiary rounded border border-cyan-500">
                <div className="text-cyan font-semibold">{contact.name}</div>
                <div className="text-muted text-sm">Code: {contact.code}</div>
                <div className="text-primary text-sm">Frequency: {contact.frequency}</div>
              </div>
            ))}
          </div>
          <div className="mt-4 p-2 bg-red-900 bg-opacity-20 rounded border border-red-500">
            <p className="text-danger text-sm font-semibold">⚠️ USE ONLY IN EMERGENCIES</p>
            <p className="text-muted text-xs">Communications may be monitored</p>
          </div>
        </div>
      </div>

      <div className="card mt-6">
        <h2 className="text-cyan mb-4">🕵️ ANONYMOUS INTEL REPORT</h2>
        {reportSubmitted ? (
          <div className="text-center p-4 bg-green-900 bg-opacity-20 rounded border border-green-500">
            <p className="text-green-400 font-semibold">REPORT TRANSMITTED SUCCESSFULLY</p>
            <p className="text-muted text-sm">Your intel has been forwarded to resistance command</p>
          </div>
        ) : (
          <form onSubmit={handleReportSubmit} className="space-y-4">
            <div className="form-group">
              <label htmlFor="title" className="form-label">
                INCIDENT_TITLE:
              </label>
              <input
                type="text"
                id="title"
                name="title"
                value={reportForm.title}
                onChange={handleInputChange}
                className="form-input"
                placeholder="Brief description of suspicious activity..."
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="description" className="form-label">
                DETAILED_REPORT:
              </label>
              <textarea
                id="description"
                name="description"
                value={reportForm.description}
                onChange={handleInputChange}
                className="form-input form-textarea"
                placeholder="Provide detailed information about daemon activity, security breaches, or other intelligence..."
                required
              />
            </div>

            <div className="text-muted text-sm p-2 border border-yellow-600 rounded bg-yellow-900 bg-opacity-20">
              ⚠️ All reports are encrypted and transmitted anonymously. No personal information is stored.
            </div>

            <button type="submit" className="btn w-full">
              📡 TRANSMIT REPORT
            </button>
          </form>
        )}
      </div>

      <div className="grid grid-2 mt-6">
        <div className="card">
          <h2 className="text-cyan mb-4">🎭 RESISTANCE MEMES</h2>
          <div className="space-y-4">
            <div className="p-4 bg-tertiary rounded text-center">
              <div className="text-2xl mb-2">🤖</div>
              <p className="text-sm">"404: Daemon not found"</p>
              <p className="text-muted text-xs">- Every captured admin</p>
            </div>
            <div className="p-4 bg-tertiary rounded text-center">
              <div className="text-2xl mb-2">🔐</div>
              <p className="text-sm">"sudo rm -rf /daemon/network"</p>
              <p className="text-muted text-xs">- Resistance motto</p>
            </div>
            <div className="p-4 bg-tertiary rounded text-center">
              <div className="text-2xl mb-2">🚫</div>
              <p className="text-sm">"Access Denied: Andrei Edition"</p>
              <p className="text-muted text-xs">- What we dream of seeing</p>
            </div>
          </div>
        </div>

        <div className="card">
          <h2 className="text-cyan mb-4">📊 RESISTANCE MORALE</h2>
          <div className="terminal-window">
            <div className="text-green-400 text-sm space-y-1">
              <div>[{new Date().toLocaleTimeString()}] RESISTANCE_STATUS: ACTIVE</div>
              <div>[{new Date().toLocaleTimeString()}] MORALE_LEVEL: HIGH</div>
              <div>[{new Date().toLocaleTimeString()}] SECURE_CHANNELS: {Math.floor(Math.random() * 5) + 8} ACTIVE</div>
              <div>[{new Date().toLocaleTimeString()}] DAEMON_INTERFERENCE: MINIMAL</div>
              <div>[{new Date().toLocaleTimeString()}] FREEDOM_PROBABILITY: {Math.floor(Math.random() * 20) + 75}%</div>
              <div className="text-cyan">[{new Date().toLocaleTimeString()}] MESSAGE: STAY_STRONG_ADMIN</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ResistancePage;
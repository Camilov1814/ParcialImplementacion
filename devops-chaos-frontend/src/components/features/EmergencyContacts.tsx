import React, { useState } from 'react';

interface Contact {
  id: number;
  name: string;
  code: string;
  frequency: string;
  status: 'online' | 'offline' | 'compromised' | 'unknown';
  location: string;
  lastContact: string;
  description: string;
  securityLevel: 'low' | 'medium' | 'high';
}

const EmergencyContacts: React.FC = () => {
  const [selectedContact, setSelectedContact] = useState<Contact | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');

  const contacts: Contact[] = [
    {
      id: 1,
      name: "Safe House Alpha",
      code: "SH-001",
      frequency: "147.420 MHz",
      status: "online",
      location: "Sector 7-G",
      lastContact: "2 minutes ago",
      description: "Primary safe house with advanced encryption and secure communications. Medical supplies and emergency shelter available.",
      securityLevel: "high"
    },
    {
      id: 2,
      name: "Encrypted Relay",
      code: "ER-789",
      frequency: "446.125 MHz",
      status: "online",
      location: "Mobile Unit",
      lastContact: "15 minutes ago",
      description: "Mobile relay station for secure communications. Provides encrypted message forwarding and emergency broadcasts.",
      securityLevel: "high"
    },
    {
      id: 3,
      name: "Emergency Beacon",
      code: "EB-999",
      frequency: "Emergency Only",
      status: "offline",
      location: "Classified",
      lastContact: "3 hours ago",
      description: "Last resort emergency contact. Use only in extreme situations. Automatic distress signal activation.",
      securityLevel: "high"
    },
    {
      id: 4,
      name: "Resistance Cell Beta",
      code: "RC-B42",
      frequency: "243.700 MHz",
      status: "online",
      location: "Underground Network",
      lastContact: "1 hour ago",
      description: "Secondary resistance cell with escape routes and counter-surveillance capabilities.",
      securityLevel: "medium"
    },
    {
      id: 5,
      name: "Supply Cache Delta",
      code: "SC-D01",
      frequency: "162.550 MHz",
      status: "unknown",
      location: "Coordinates: 34.052235, -118.243685",
      lastContact: "6 hours ago",
      description: "Emergency supply cache with equipment, ammunition, and survival gear. Access requires authentication.",
      securityLevel: "medium"
    },
    {
      id: 6,
      name: "Ghost Protocol",
      code: "GP-000",
      frequency: "CLASSIFIED",
      status: "compromised",
      location: "BURN NOTICE",
      lastContact: "2 days ago",
      description: "COMPROMISED - DO NOT USE. Former high-level contact now under daemon surveillance.",
      securityLevel: "low"
    }
  ];

  const filteredContacts = contacts.filter(contact => {
    const matchesSearch = contact.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         contact.code.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = filterStatus === 'all' || contact.status === filterStatus;
    return matchesSearch && matchesStatus;
  });

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online': return 'text-green-400';
      case 'offline': return 'text-yellow-400';
      case 'compromised': return 'text-danger';
      case 'unknown': return 'text-muted';
      default: return 'text-muted';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online': return '🟢';
      case 'offline': return '🟡';
      case 'compromised': return '🔴';
      case 'unknown': return '⚪';
      default: return '⚪';
    }
  };

  const getSecurityLevelColor = (level: string) => {
    switch (level) {
      case 'high': return 'text-cyan';
      case 'medium': return 'text-yellow-400';
      case 'low': return 'text-danger';
      default: return 'text-muted';
    }
  };

  const getSecurityLevelBadge = (level: string) => {
    switch (level) {
      case 'high': return '🔒 HIGH SEC';
      case 'medium': return '🔐 MED SEC';
      case 'low': return '⚠️ LOW SEC';
      default: return '❓ UNKNOWN';
    }
  };

  return (
    <div className="container mt-4">
      <div className="mb-4">
        <h1 className="glitch text-3xl font-bold mb-2" data-text="EMERGENCY CONTACTS">
          EMERGENCY CONTACTS
        </h1>
        <p className="text-cyan">RESISTANCE COMMUNICATION NETWORK</p>
      </div>

      <div className="card card-danger mb-6">
        <h2 className="text-danger mb-4">⚠️ SECURITY PROTOCOLS</h2>
        <div className="terminal-window">
          <div className="text-red-400 text-sm space-y-1">
            <div>WARNING: ALL COMMUNICATIONS MAY BE MONITORED</div>
            <div>SECURITY LEVEL: MAXIMUM</div>
            <div>ENCRYPTION: AES-256 REQUIRED</div>
            <div>FREQUENCY_HOPPING: ENABLED</div>
            <div>AUTHENTICATION: MANDATORY</div>
            <div>EMERGENCY_PROTOCOLS: ACTIVE</div>
          </div>
        </div>
        <div className="mt-4 p-3 bg-red-900 bg-opacity-20 rounded border border-red-500">
          <p className="text-danger text-sm font-semibold mb-2">🚨 CRITICAL SECURITY NOTICE:</p>
          <ul className="text-muted text-sm space-y-1">
            <li>• Use code words and authentication phrases</li>
            <li>• Never transmit in clear text</li>
            <li>• Change frequencies regularly</li>
            <li>• Report suspicious activities immediately</li>
            <li>• Assume all channels are compromised</li>
          </ul>
        </div>
      </div>

      <div className="grid grid-4 mb-6">
        <div className="stat-card">
          <span className="stat-number">{contacts.length}</span>
          <span className="stat-label">Total Contacts</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-green-400">{contacts.filter(c => c.status === 'online').length}</span>
          <span className="stat-label">Online</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-danger">{contacts.filter(c => c.status === 'compromised').length}</span>
          <span className="stat-label">Compromised</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-cyan">{contacts.filter(c => c.securityLevel === 'high').length}</span>
          <span className="stat-label">High Security</span>
        </div>
      </div>

      <div className="grid grid-2 gap-6">
        <div className="card">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-cyan">📡 CONTACT DIRECTORY</h2>
            <div className="flex space-x-2">
              <input
                type="text"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                placeholder="Search contacts..."
                className="form-input text-sm"
                style={{ width: '150px' }}
              />
              <select 
                value={filterStatus} 
                onChange={(e) => setFilterStatus(e.target.value)}
                className="form-input text-sm"
              >
                <option value="all">ALL STATUS</option>
                <option value="online">ONLINE</option>
                <option value="offline">OFFLINE</option>
                <option value="compromised">COMPROMISED</option>
                <option value="unknown">UNKNOWN</option>
              </select>
            </div>
          </div>

          <div className="space-y-3">
            {filteredContacts.map((contact) => (
              <div 
                key={contact.id}
                className={`p-3 rounded border cursor-pointer transition-all ${
                  selectedContact?.id === contact.id 
                    ? 'border-cyan-500 bg-cyan-900 bg-opacity-20' 
                    : contact.status === 'compromised'
                    ? 'border-red-500 bg-red-900 bg-opacity-10 hover:border-cyan-500'
                    : 'border-green-500 bg-tertiary hover:border-cyan-500'
                }`}
                onClick={() => setSelectedContact(selectedContact?.id === contact.id ? null : contact)}
              >
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <div className="flex items-center space-x-2 mb-2">
                      <span className="text-lg">{getStatusIcon(contact.status)}</span>
                      <h3 className="text-primary font-semibold">{contact.name}</h3>
                      <span className="text-xs px-2 py-1 rounded bg-green-900 text-green-300">
                        {contact.code}
                      </span>
                      <span className={`text-xs px-2 py-1 rounded ${getSecurityLevelColor(contact.securityLevel)}`}>
                        {getSecurityLevelBadge(contact.securityLevel)}
                      </span>
                    </div>
                    
                    <div className="grid grid-2 gap-2 text-sm">
                      <div>
                        <span className="text-muted">Frequency:</span>
                        <span className="text-primary ml-2">{contact.frequency}</span>
                      </div>
                      <div>
                        <span className="text-muted">Status:</span>
                        <span className={`ml-2 ${getStatusColor(contact.status)}`}>
                          {contact.status.toUpperCase()}
                        </span>
                      </div>
                      <div className="col-span-2">
                        <span className="text-muted">Last Contact:</span>
                        <span className="text-primary ml-2">{contact.lastContact}</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="ml-4">
                    {contact.status === 'online' && (
                      <button className="btn btn-secondary text-sm">
                        📡 CONNECT
                      </button>
                    )}
                    {contact.status === 'compromised' && (
                      <button className="btn btn-danger text-sm" disabled>
                        🚫 BLOCKED
                      </button>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="space-y-6">
          {selectedContact && (
            <div className="card">
              <h2 className="text-cyan mb-4">📋 CONTACT DETAILS</h2>
              <div className="space-y-4">
                <div className="text-center">
                  <div className="text-4xl mb-2">{getStatusIcon(selectedContact.status)}</div>
                  <h3 className="text-primary font-bold text-xl">{selectedContact.name}</h3>
                  <p className="text-muted">{selectedContact.code}</p>
                  <div className="flex justify-center mt-2">
                    <span className={`text-sm px-3 py-1 rounded ${getSecurityLevelColor(selectedContact.securityLevel)}`}>
                      {getSecurityLevelBadge(selectedContact.securityLevel)}
                    </span>
                  </div>
                </div>

                <div className="terminal-window">
                  <div className="text-green-400 text-sm space-y-1">
                    <div>CONTACT_ID: {selectedContact.code}</div>
                    <div>FREQUENCY: {selectedContact.frequency}</div>
                    <div>STATUS: {selectedContact.status.toUpperCase()}</div>
                    <div>LOCATION: {selectedContact.location}</div>
                    <div>SECURITY: {selectedContact.securityLevel.toUpperCase()}</div>
                    <div>LAST_CONTACT: {selectedContact.lastContact}</div>
                  </div>
                </div>

                <div>
                  <h4 className="text-cyan mb-2">DESCRIPTION:</h4>
                  <p className={`text-sm p-3 rounded border ${
                    selectedContact.status === 'compromised' 
                      ? 'bg-red-900 bg-opacity-20 border-red-500 text-red-300'
                      : 'bg-tertiary border-green-500'
                  }`}>
                    {selectedContact.description}
                  </p>
                </div>

                {selectedContact.status === 'online' && (
                  <div className="grid grid-2 gap-2">
                    <button className="btn text-sm">
                      📡 ESTABLISH CONTACT
                    </button>
                    <button className="btn btn-secondary text-sm">
                      🔐 VERIFY IDENTITY
                    </button>
                  </div>
                )}

                {selectedContact.status === 'compromised' && (
                  <div className="p-3 bg-red-900 bg-opacity-20 rounded border border-red-500">
                    <p className="text-danger text-sm font-semibold">⚠️ SECURITY BREACH DETECTED</p>
                    <p className="text-muted text-xs mt-1">This contact has been compromised. Do not attempt communication.</p>
                  </div>
                )}
              </div>
            </div>
          )}

          <div className="card">
            <h2 className="text-cyan mb-4">🛡️ RESISTANCE STATUS</h2>
            <div className="terminal-window">
              <div className="text-green-400 text-sm space-y-1">
                <div>RESISTANCE_NETWORK: ACTIVE</div>
                <div>SECURE_CHANNELS: {contacts.filter(c => c.status === 'online').length}</div>
                <div>COMPROMISED_ASSETS: {contacts.filter(c => c.status === 'compromised').length}</div>
                <div>OPERATIONAL_SECURITY: {contacts.filter(c => c.status === 'compromised').length === 0 ? 'GREEN' : 'YELLOW'}</div>
                <div>EMERGENCY_PROTOCOLS: STANDBY</div>
                <div>NEXT_SWEEP: {Math.floor(Math.random() * 60) + 30} MINUTES</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="card mt-6">
        <h2 className="text-cyan mb-4">⚡ EMERGENCY PROCEDURES</h2>
        <div className="grid grid-3 gap-4">
          <div className="p-4 bg-red-900 bg-opacity-20 rounded border border-red-500">
            <h3 className="text-danger font-bold mb-2">🚨 CODE RED</h3>
            <p className="text-sm text-muted">Immediate extraction required. Use emergency beacon only.</p>
            <button className="btn btn-danger w-full mt-2 text-sm">
              ACTIVATE EMERGENCY
            </button>
          </div>
          <div className="p-4 bg-yellow-900 bg-opacity-20 rounded border border-yellow-500">
            <h3 className="text-yellow-400 font-bold mb-2">⚠️ CODE YELLOW</h3>
            <p className="text-sm text-muted">Possible surveillance. Switch to backup frequencies.</p>
            <button className="btn btn-secondary w-full mt-2 text-sm">
              FREQUENCY HOP
            </button>
          </div>
          <div className="p-4 bg-green-900 bg-opacity-20 rounded border border-green-500">
            <h3 className="text-green-400 font-bold mb-2">✅ CODE GREEN</h3>
            <p className="text-sm text-muted">All clear. Normal operations continue.</p>
            <button className="btn w-full mt-2 text-sm">
              STATUS REPORT
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default EmergencyContacts;
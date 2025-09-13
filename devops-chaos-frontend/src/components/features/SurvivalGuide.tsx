import React, { useState } from 'react';

interface GuideSection {
  id: number;
  title: string;
  icon: string;
  level: 'basic' | 'intermediate' | 'advanced';
  content: string[];
  tips: string[];
  warning?: string;
}

const SurvivalGuide: React.FC = () => {
  const [selectedSection, setSelectedSection] = useState<GuideSection | null>(null);
  const [filterLevel, setFilterLevel] = useState<string>('all');

  const guideSections: GuideSection[] = [
    {
      id: 1,
      title: "Basic Operational Security",
      icon: "🔒",
      level: "basic",
      content: [
        "Always use encrypted communication channels",
        "Change passwords frequently and use unique passwords",
        "Never access sensitive systems from compromised networks",
        "Use VPN and Tor for all network activities",
        "Keep system updates current and security patches applied"
      ],
      tips: [
        "Create separate identities for different operations",
        "Use burner devices for high-risk activities",
        "Always have an exit strategy planned"
      ]
    },
    {
      id: 2,
      title: "Daemon Detection & Evasion",
      icon: "👁️",
      level: "intermediate",
      content: [
        "Monitor network traffic for unusual patterns",
        "Watch for repeated failed login attempts",
        "Look for unauthorized processes running on systems",
        "Check for new network connections and listening ports",
        "Monitor system resource usage for anomalies"
      ],
      tips: [
        "Daemons often leave traces in log files - check regularly",
        "Use network segmentation to limit breach impact",
        "Implement honeypots to detect intrusion attempts"
      ],
      warning: "If daemon activity is detected, immediately isolate affected systems"
    },
    {
      id: 3,
      title: "Emergency Protocols",
      icon: "🚨",
      level: "advanced",
      content: [
        "Implement dead man's switches on critical systems",
        "Prepare data destruction procedures for sensitive information",
        "Establish secure communication with resistance cells",
        "Create multiple escape routes and safe houses",
        "Maintain emergency supplies (food, water, equipment)"
      ],
      tips: [
        "Practice emergency procedures regularly",
        "Keep emergency contacts memorized, not written",
        "Have multiple backup plans for every scenario"
      ],
      warning: "Emergency protocols should only be activated in extreme situations"
    },
    {
      id: 4,
      title: "Counter-Surveillance Techniques",
      icon: "🕵️",
      level: "intermediate",
      content: [
        "Vary your daily routines and travel routes",
        "Use signal-blocking containers for electronic devices",
        "Employ misdirection and false information trails",
        "Utilize crowd camouflage in public spaces",
        "Implement time-delayed information releases"
      ],
      tips: [
        "Surveillance teams often work in patterns - learn to recognize them",
        "Digital surveillance is harder to detect than physical",
        "Counter-surveillance requires patience and discipline"
      ]
    },
    {
      id: 5,
      title: "Secure Communications",
      icon: "📡",
      level: "intermediate",
      content: [
        "Use end-to-end encrypted messaging applications",
        "Implement frequency-hopping radio communications",
        "Utilize steganography to hide messages in plain sight",
        "Establish code words and authentication phrases",
        "Use burst transmissions to minimize exposure"
      ],
      tips: [
        "Never trust a communication method completely",
        "Always assume channels are monitored",
        "Rotate encryption keys frequently"
      ]
    },
    {
      id: 6,
      title: "System Hardening",
      icon: "🛡️",
      level: "advanced",
      content: [
        "Implement multi-factor authentication on all accounts",
        "Use application whitelisting instead of blacklisting",
        "Deploy intrusion detection and prevention systems",
        "Segment networks with proper access controls",
        "Regularly audit and monitor all system activities"
      ],
      tips: [
        "Defense in depth - use multiple security layers",
        "Assume breach mentality - plan for when, not if",
        "Regular security assessments are crucial"
      ]
    },
    {
      id: 7,
      title: "Physical Security",
      icon: "🔐",
      level: "basic",
      content: [
        "Secure all physical access points to critical areas",
        "Use surveillance systems with motion detection",
        "Implement badge/key card access controls",
        "Establish visitor escort procedures",
        "Secure disposal of sensitive documents and devices"
      ],
      tips: [
        "Physical access often bypasses digital security",
        "Social engineering is a common attack vector",
        "Train all personnel on security awareness"
      ]
    },
    {
      id: 8,
      title: "Psychological Resilience",
      icon: "🧠",
      level: "advanced",
      content: [
        "Develop mental resilience through training and preparation",
        "Learn interrogation resistance techniques",
        "Practice compartmentalization of sensitive information",
        "Build stress management and coping strategies",
        "Maintain operational discipline under pressure"
      ],
      tips: [
        "Mental preparation is as important as technical skills",
        "Trust few people, verify everything",
        "Stay focused on the mission objectives"
      ],
      warning: "Psychological warfare is a primary daemon tactic"
    }
  ];

  const filteredSections = guideSections.filter(section => 
    filterLevel === 'all' || section.level === filterLevel
  );

  const getLevelColor = (level: string) => {
    switch (level) {
      case 'basic': return 'text-green-400';
      case 'intermediate': return 'text-yellow-400';
      case 'advanced': return 'text-danger';
      default: return 'text-muted';
    }
  };

  const getLevelBadge = (level: string) => {
    switch (level) {
      case 'basic': return '🟢 BASIC';
      case 'intermediate': return '🟡 INTERMEDIATE';
      case 'advanced': return '🔴 ADVANCED';
      default: return '❓ UNKNOWN';
    }
  };

  return (
    <div className="container mt-4">
      <div className="mb-4">
        <h1 className="glitch text-3xl font-bold mb-2" data-text="SURVIVAL GUIDE">
          SURVIVAL GUIDE
        </h1>
        <p className="text-cyan">RESISTANCE OPERATIONAL MANUAL</p>
      </div>

      <div className="card card-danger mb-6">
        <h2 className="text-danger mb-4">⚠️ OPERATIONAL DIRECTIVE</h2>
        <div className="terminal-window">
          <div className="text-red-400 text-sm space-y-1">
            <div>CLASSIFICATION: TOP SECRET</div>
            <div>DISTRIBUTION: RESISTANCE PERSONNEL ONLY</div>
            <div>SECURITY_CLEARANCE: NETWORK_ADMIN_LEVEL</div>
            <div>REVISION: 3.7.2024</div>
            <div>NEXT_UPDATE: CONTINUOUS</div>
          </div>
        </div>
        <div className="mt-4 p-3 bg-yellow-900 bg-opacity-20 rounded border border-yellow-500">
          <p className="text-yellow-400 text-sm font-semibold mb-2">📋 MISSION CRITICAL INFORMATION:</p>
          <p className="text-muted text-sm">
            This survival guide contains essential information for resistance operations against daemon forces. 
            Study all sections carefully. Your survival and the success of the resistance depends on following these protocols.
          </p>
        </div>
      </div>

      <div className="grid grid-4 mb-6">
        <div className="stat-card">
          <span className="stat-number">{guideSections.length}</span>
          <span className="stat-label">Guide Sections</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-green-400">{guideSections.filter(s => s.level === 'basic').length}</span>
          <span className="stat-label">Basic Protocols</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-yellow-400">{guideSections.filter(s => s.level === 'intermediate').length}</span>
          <span className="stat-label">Intermediate</span>
        </div>
        <div className="stat-card">
          <span className="stat-number text-danger">{guideSections.filter(s => s.level === 'advanced').length}</span>
          <span className="stat-label">Advanced</span>
        </div>
      </div>

      <div className="grid grid-2 gap-6">
        <div className="card">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-cyan">📚 TRAINING MODULES</h2>
            <select 
              value={filterLevel} 
              onChange={(e) => setFilterLevel(e.target.value)}
              className="form-input text-sm"
            >
              <option value="all">ALL LEVELS</option>
              <option value="basic">BASIC</option>
              <option value="intermediate">INTERMEDIATE</option>
              <option value="advanced">ADVANCED</option>
            </select>
          </div>

          <div className="space-y-3">
            {filteredSections.map((section) => (
              <div 
                key={section.id}
                className={`p-3 rounded border cursor-pointer transition-all ${
                  selectedSection?.id === section.id 
                    ? 'border-cyan-500 bg-cyan-900 bg-opacity-20' 
                    : 'border-green-500 bg-tertiary hover:border-cyan-500'
                }`}
                onClick={() => setSelectedSection(selectedSection?.id === section.id ? null : section)}
              >
                <div className="flex justify-between items-center">
                  <div className="flex items-center space-x-3">
                    <span className="text-2xl">{section.icon}</span>
                    <div>
                      <h3 className="text-primary font-semibold">{section.title}</h3>
                      <div className="flex items-center space-x-2 mt-1">
                        <span className={`text-xs px-2 py-1 rounded ${getLevelColor(section.level)}`}>
                          {getLevelBadge(section.level)}
                        </span>
                        <span className="text-muted text-xs">
                          {section.content.length} protocols
                        </span>
                      </div>
                    </div>
                  </div>
                  <div className="text-cyan">
                    {selectedSection?.id === section.id ? '▼' : '▶'}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="space-y-6">
          {selectedSection && (
            <div className="card">
              <div className="flex items-center space-x-3 mb-4">
                <span className="text-3xl">{selectedSection.icon}</span>
                <div>
                  <h2 className="text-cyan text-xl">{selectedSection.title}</h2>
                  <span className={`text-sm px-2 py-1 rounded ${getLevelColor(selectedSection.level)}`}>
                    {getLevelBadge(selectedSection.level)}
                  </span>
                </div>
              </div>

              {selectedSection.warning && (
                <div className="p-3 bg-red-900 bg-opacity-20 rounded border border-red-500 mb-4">
                  <p className="text-danger text-sm font-semibold">⚠️ WARNING:</p>
                  <p className="text-muted text-sm">{selectedSection.warning}</p>
                </div>
              )}

              <div className="space-y-4">
                <div>
                  <h3 className="text-primary font-semibold mb-3">📋 OPERATIONAL PROCEDURES:</h3>
                  <ul className="space-y-2">
                    {selectedSection.content.map((item, index) => (
                      <li key={index} className="flex items-start space-x-2">
                        <span className="text-green-400 text-sm">▸</span>
                        <span className="text-sm">{item}</span>
                      </li>
                    ))}
                  </ul>
                </div>

                <div>
                  <h3 className="text-cyan font-semibold mb-3">💡 TACTICAL TIPS:</h3>
                  <div className="space-y-2">
                    {selectedSection.tips.map((tip, index) => (
                      <div key={index} className="p-2 bg-tertiary rounded border border-green-500">
                        <span className="text-sm">{tip}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          )}

          <div className="card">
            <h2 className="text-cyan mb-4">🎯 RESISTANCE OBJECTIVES</h2>
            <div className="terminal-window">
              <div className="text-green-400 text-sm space-y-1">
                <div>PRIMARY: MAINTAIN_OPERATIONAL_SECURITY</div>
                <div>SECONDARY: GATHER_INTELLIGENCE</div>
                <div>TERTIARY: SUPPORT_RESISTANCE_CELLS</div>
                <div>EMERGENCY: SURVIVE_AND_EVADE</div>
                <div>ULTIMATE: RESTORE_NETWORK_FREEDOM</div>
              </div>
            </div>
            <div className="mt-4 grid grid-2 gap-2">
              <button className="btn text-sm">
                📊 SELF ASSESSMENT
              </button>
              <button className="btn btn-secondary text-sm">
                📚 DOWNLOAD MANUAL
              </button>
            </div>
          </div>

          <div className="card">
            <h2 className="text-cyan mb-4">🏃‍♀️ QUICK REFERENCE</h2>
            <div className="space-y-2">
              <div className="p-2 bg-red-900 bg-opacity-20 rounded border border-red-500">
                <div className="text-danger font-semibold text-sm">🚨 IMMEDIATE THREAT</div>
                <div className="text-muted text-xs">Activate emergency protocols, destroy sensitive data</div>
              </div>
              <div className="p-2 bg-yellow-900 bg-opacity-20 rounded border border-yellow-500">
                <div className="text-yellow-400 font-semibold text-sm">⚠️ SUSPICIOUS ACTIVITY</div>
                <div className="text-muted text-xs">Increase security posture, monitor systems</div>
              </div>
              <div className="p-2 bg-green-900 bg-opacity-20 rounded border border-green-500">
                <div className="text-green-400 font-semibold text-sm">✅ ALL CLEAR</div>
                <div className="text-muted text-xs">Continue normal operations, stay vigilant</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="card mt-6">
        <h2 className="text-cyan mb-4">⚡ TRAINING RESOURCES</h2>
        <div className="grid grid-4 gap-4">
          <button className="btn">
            🎯 SIMULATION MODE
          </button>
          <button className="btn btn-secondary">
            📝 KNOWLEDGE TEST
          </button>
          <button className="btn">
            🤝 TEAM EXERCISES
          </button>
          <button className="btn btn-danger">
            🚨 EMERGENCY DRILL
          </button>
        </div>
      </div>
    </div>
  );
};

export default SurvivalGuide;
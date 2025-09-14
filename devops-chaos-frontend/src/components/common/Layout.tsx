import React from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  if (!user) {
    return <>{children}</>;
  }

  const roleConfig = {
    andrei: {
      title: 'ANDREI SUPREME COMMAND',
      color: 'text-cyan',
      menuItems: [
        { label: 'Command Center', path: '/dashboard', icon: '🎛️' },
        { label: 'Manage Users', path: '/users', icon: '👥' },
        { label: 'View Reports', path: '/reports', icon: '📊' },
        { label: 'Punishments', path: '/punishments', icon: '⚖️' },
        { label: 'Leaderboard', path: '/leaderboard', icon: '🏆' },
      ]
    },
    daemon: {
      title: 'DAEMON OPERATIONS',
      color: 'text-primary',
      menuItems: [
        { label: 'Daemon Terminal', path: '/dashboard', icon: '💻' },
        { label: 'Capture Mission', path: '/capture', icon: '🎯' },
        { label: 'File Report', path: '/reports', icon: '📝' },
        { label: 'Leaderboard', path: '/leaderboard', icon: '🏆' },
        { label: 'My Punishments', path: '/punishments', icon: '⚠️' },
      ]
    },
    network_admin: {
      title: 'RESISTANCE NETWORK',
      color: 'text-cyan',
      menuItems: [
        { label: 'Resistance Hub', path: '/dashboard', icon: '🛡️' },
        { label: 'File Intel Report', path: '/reports', icon: '🕵️' },
      ]
    }
  };

  const config = roleConfig[user.role as keyof typeof roleConfig];

  return (
    <div className="min-h-screen flex">
      <aside className="sidebar">
        <div className="sidebar-header">
          <h2 className={`text-lg font-bold ${config.color}`}>
            {config.title}
          </h2>
          <div className="text-sm text-muted mt-1">
            <div>User: <span className="text-primary">{user.username}</span></div>
            <div>Role: <span className={config.color}>{user.role.toUpperCase()}</span></div>
          </div>
        </div>

        <nav className="sidebar-nav">
          <ul className="space-y-2">
            {config.menuItems.map((item: any, index: number) => (
              <li key={index}>
                <button
                  onClick={() => navigate(item.path)}
                  className="nav-button"
                >
                  <span className="text-lg">{item.icon}</span>
                  <span className="text-sm font-semibold">{item.label}</span>
                </button>
              </li>
            ))}
          </ul>
        </nav>

        <div className="sidebar-footer">
          <div className="terminal-window mb-3">
            <div className="text-primary text-xs space-y-1">
              <div>STATUS: <span className="text-cyan">ONLINE</span></div>
              <div>UPTIME: <span className="text-primary">{Math.floor(Math.random() * 24)}h {Math.floor(Math.random() * 60)}m</span></div>
              <div>SECURITY: <span className="text-cyan">ENCRYPTED</span></div>
            </div>
          </div>
          
          <button
            onClick={handleLogout}
            className="btn btn-danger w-full text-sm"
          >
            🚪 DISCONNECT
          </button>
        </div>
      </aside>

      <main className="flex-1 flex flex-col">
        <header className="main-header">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-xl font-bold text-primary">DevOps Chaos System</h1>
              <div className="text-sm text-muted">
                Logged in as: <span className={config.color}>{user.name}</span>
              </div>
            </div>
            <div className="text-right">
              <div className="text-sm text-muted">
                {new Date().toLocaleDateString()} {new Date().toLocaleTimeString()}
              </div>
              <div className="text-xs text-primary">
                System: <span className="text-cyan">OPERATIONAL</span>
              </div>
            </div>
          </div>
        </header>

        <div className="flex-1 overflow-auto">
          {children}
        </div>
      </main>
    </div>
  );
};

export default Layout;
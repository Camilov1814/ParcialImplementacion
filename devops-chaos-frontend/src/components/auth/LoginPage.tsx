import React, { useState } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';

const LoginPage: React.FC = () => {
  const [credentials, setCredentials] = useState({ username: '', password: '' });
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login(credentials);
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.message || err.response?.data?.message || 'LOGIN FAILED. ACCESS DENIED.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCredentials({
      ...credentials,
      [e.target.name]: e.target.value,
    });
  };

  const quickLogin = (username: string, password: string) => {
    setCredentials({ username, password });
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-primary">
      <div className="container">
        <div className="max-w-md mx-auto">
          <div className="terminal-window mb-4">
            <div className="text-center mb-6">
              <h1 className="glitch text-4xl font-bold mb-2" data-text="DEVOPS CHAOS">
                DEVOPS CHAOS
              </h1>
              <p className="text-cyan text-sm">SYSTEM ACCESS TERMINAL</p>
              <p className="text-muted text-xs mt-2">UNAUTHORIZED ACCESS WILL BE TERMINATED</p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="form-group">
                <label htmlFor="username" className="form-label">
                  USER_ID:
                </label>
                <input
                  type="text"
                  id="username"
                  name="username"
                  value={credentials.username}
                  onChange={handleChange}
                  className="form-input"
                  placeholder="Enter username..."
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="password" className="form-label">
                  ACCESS_KEY:
                </label>
                <input
                  type="password"
                  id="password"
                  name="password"
                  value={credentials.password}
                  onChange={handleChange}
                  className="form-input"
                  placeholder="Enter password..."
                  required
                />
              </div>

              {error && (
                <div className="text-danger text-sm text-center p-2 border border-red-500 rounded">
                  ERROR: {error}
                </div>
              )}

              <button
                type="submit"
                disabled={isLoading}
                className="btn w-full"
              >
                {isLoading ? 'AUTHENTICATING...' : 'GRANT ACCESS'}
              </button>
            </form>

            <div className="mt-6 border-t border-green-500 pt-4">
              <p className="text-center text-sm text-muted mb-3">QUICK ACCESS (DEV MODE)</p>
              <div className="space-y-2">
                <button
                  type="button"
                  onClick={() => quickLogin('andrei', 'AndreI2024!')}
                  className="btn btn-secondary w-full text-xs"
                >
                  LOGIN AS ANDREI (SUPREME LEADER)
                </button>
                <button
                  type="button"
                  onClick={() => quickLogin('daemon_alpha', 'DaemonAlpha123!')}
                  className="btn btn-secondary w-full text-xs"
                >
                  LOGIN AS DAEMON (HACKER)
                </button>
                <button
                  type="button"
                  onClick={() => quickLogin('netadmin_alice', 'NetworkAlice123!')}
                  className="btn btn-secondary w-full text-xs"
                >
                  LOGIN AS NETWORK ADMIN (RESISTANCE)
                </button>
              </div>
            </div>
          </div>

          <div className="text-center text-xs text-muted">
            <p>SYSTEM VERSION 2.0.1 | CHAOS PROTOCOL ACTIVE</p>
            <p>CONNECTION ENCRYPTED • LOG MONITORING ENABLED</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
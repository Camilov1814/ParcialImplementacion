import React, { useState } from 'react';
import { apiService } from '../../services/api';
import { useAuth } from '../../contexts/AuthContext';

const SimpleReportTest: React.FC = () => {
  const { user } = useAuth();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState<string>('');

  const testMinimalReport = async () => {
    setIsLoading(true);
    setResult('');
    
    try {
      const minimalData = {
        title: title || 'Test Report',
        description: description || 'Test Description',
        type: 'anonymous'
      };
      
      console.log('Testing minimal report:', minimalData);
      
      const response = await apiService.createReport(minimalData);
      setResult(`✅ SUCCESS: ${JSON.stringify(response)}`);
    } catch (error: any) {
      console.error('Test failed:', error);
      setResult(`❌ ERROR: ${error.response?.data?.message || error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const testWithReporterType = async () => {
    setIsLoading(true);
    setResult('');
    
    try {
      const getReportType = (userRole: string) => {
        switch (userRole) {
          case 'network_admin': return 'resistance';
          case 'daemon': return 'capture';
          case 'andrei': return 'anonymous';
          default: return 'anonymous';
        }
      };

      const dataWithType = {
        title: title || 'Test Report with Type',
        description: description || 'Test Description with Type',
        type: getReportType(user?.role || 'anonymous')
      };
      
      console.log('Testing with type:', dataWithType);
      
      const response = await apiService.createReport(dataWithType);
      setResult(`✅ SUCCESS: ${JSON.stringify(response)}`);
    } catch (error: any) {
      console.error('Test failed:', error);
      setResult(`❌ ERROR: ${error.response?.data?.message || error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  const testWithStatus = async () => {
    setIsLoading(true);
    setResult('');
    
    try {
      const getReportType = (userRole: string) => {
        switch (userRole) {
          case 'network_admin': return 'resistance';
          case 'daemon': return 'capture';
          case 'andrei': return 'anonymous';
          default: return 'anonymous';
        }
      };

      const dataWithType = {
        title: title || 'Test Report Complete',
        description: description || 'Test Description Complete',
        type: getReportType(user?.role || 'anonymous')
      };
      
      console.log('Testing complete data structure:', dataWithType);
      
      const response = await apiService.createReport(dataWithType);
      setResult(`✅ SUCCESS: ${JSON.stringify(response)}`);
    } catch (error: any) {
      console.error('Test failed:', error);
      setResult(`❌ ERROR: ${error.response?.data?.message || error.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container mt-4">
      <div className="card">
        <h2 className="text-cyan mb-4">🧪 REPORT API TEST</h2>
        <div className="space-y-4">
          <div className="form-group">
            <label className="form-label">Test Title:</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="form-input"
              placeholder="Optional - will use default if empty"
            />
          </div>
          
          <div className="form-group">
            <label className="form-label">Test Description:</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="form-input form-textarea"
              placeholder="Optional - will use default if empty"
              rows={3}
            />
          </div>

          <div className="grid grid-3 gap-4">
            <button
              onClick={testMinimalReport}
              disabled={isLoading}
              className="btn"
            >
              Test Minimal
            </button>
            <button
              onClick={testWithReporterType}
              disabled={isLoading}
              className="btn btn-secondary"
            >
              Test + Type (Role-based)
            </button>
            <button
              onClick={testWithStatus}
              disabled={isLoading}
              className="btn btn-danger"
            >
              Test Complete
            </button>
          </div>

          {isLoading && (
            <div className="text-center text-cyan">Testing...</div>
          )}

          {result && (
            <div className="terminal-window">
              <div className="text-green-400 text-sm">
                {result}
              </div>
            </div>
          )}

          <div className="p-3 bg-yellow-900 bg-opacity-20 rounded border border-yellow-500">
            <p className="text-yellow-400 text-sm font-semibold">Debug Instructions:</p>
            <ul className="text-muted text-sm mt-2 space-y-1">
              <li>• Check browser console for detailed logs</li>
              <li>• Try each test button to see which data structure works</li>
              <li>• Current user role: <span className="text-primary">{user?.role}</span></li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SimpleReportTest;
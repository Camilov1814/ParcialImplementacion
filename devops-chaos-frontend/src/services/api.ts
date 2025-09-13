import axios from 'axios';
import type { AxiosInstance } from 'axios';
import type { LoginRequest, LoginResponse, User, DashboardStats, DaemonStats, LeaderboardEntry, Report, Punishment, ResistanceStats } from '../types';

class ApiService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: 'http://localhost:8080/api',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.api.interceptors.request.use((config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });
  }

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await this.api.post('/auth/login', credentials);
    // Backend returns { success: true, message: string, data: { token: string, user: User } }
    if (response.data.success && response.data.data) {
      return {
        token: response.data.data.token,
        user: response.data.data.user
      };
    }
    throw new Error(response.data.message || 'Login failed');
  }

  async getProfile(): Promise<User> {
    const response = await this.api.get('/auth/me');
    // Backend returns { success: true, message: string, data: User }
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to get profile');
  }

  async getAndreiDashboard(): Promise<DashboardStats> {
    const response = await this.api.get('/dashboard/andrei');
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to get dashboard data');
  }

  async getDaemonDashboard(): Promise<DaemonStats> {
    const response = await this.api.get('/dashboard/daemon');
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to get daemon dashboard');
  }

  async getResistanceData(): Promise<ResistanceStats> {
    const response = await this.api.get('/resistance');
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to get resistance data');
  }

  async getUsers(): Promise<User[]> {
    const response = await this.api.get('/users');
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to get users');
  }

  async getNetworkAdminsForCapture(): Promise<User[]> {
    try {
      // Try to get users with role filter for network admins
      const response = await this.api.get('/users?role=network_admin');
      if (response.data.success && response.data.data) {
        return response.data.data;
      }
      throw new Error(response.data.message || 'Failed to get network admins');
    } catch (error: any) {
      if (error.response?.status === 403) {
        // If access is denied, return an empty array or handle differently
        console.warn('Access denied to users endpoint, returning empty target list');
        return [];
      }
      throw error;
    }
  }

  async captureNetworkAdmin(id: number): Promise<void> {
    const response = await this.api.post(`/users/${id}/capture`);
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to capture network admin');
    }
  }

  async getReports(): Promise<Report[]> {
    const response = await this.api.get('/reports');
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to get reports');
  }

  async createReport(report: Partial<Report>): Promise<Report> {
    console.log('API: Sending report data:', report);
    const response = await this.api.post('/reports', report);
    console.log('API: Received response:', response.data);
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to create report');
  }

  async updateReportStatus(reportId: number, status: string): Promise<void> {
    const response = await this.api.put(`/reports/${reportId}/status`, { status });
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to update report status');
    }
  }

  async getPunishments(): Promise<Punishment[]> {
    const response = await this.api.get('/punishments');
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to get punishments');
  }

  async getPunishmentDetails(punishmentId: number): Promise<Punishment> {
    const response = await this.api.get(`/punishments/${punishmentId}`);
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to get punishment details');
  }

  async createPunishment(punishment: Partial<Punishment>): Promise<Punishment> {
    const response = await this.api.post('/punishments', punishment);
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to create punishment');
  }

  async updatePunishment(punishmentId: number, updates: { status?: string; description?: string; expires_at?: string }): Promise<void> {
    const response = await this.api.put(`/punishments/${punishmentId}`, updates);
    if (!response.data.success) {
      throw new Error(response.data.message || 'Failed to update punishment');
    }
  }

  async getLeaderboard(): Promise<LeaderboardEntry[]> {
    const response = await this.api.get('/statistics/leaderboard');
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.message || 'Failed to get leaderboard');
  }
}

export const apiService = new ApiService();
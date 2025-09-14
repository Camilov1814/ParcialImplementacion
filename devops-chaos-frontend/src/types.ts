// User types
export interface User {
  id: number;
  username: string;
  name?: string;
  email: string;
  role: string;
  status: string;
  created_at?: string;
}

export type UserRole = 'andrei' | 'daemon' | 'network_admin';

// Auth types
export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  success: boolean;
  message: string;
  token?: string;
  user?: User;
}

export interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}

// Stats types
export interface DashboardStats {
  total_users: number;
  total_daemons: number;
  total_network_admins?: number;
  captured_admins: number;
  punished_daemons?: number;
  pending_reports?: number;
  total_reports: number;
}

export interface DaemonStats {
  captures_count: number;
  reports_count: number;
  points: number;
  ranking: number;
  status: string;
}

// Leaderboard types
export interface LeaderboardEntry {
  id: number;
  username: string;
  name?: string;
  points: number;
  captures_count: number;
  captures: number; // Alias for backwards compatibility
  reports_count?: number;
  status?: string;
}

// Report types
export interface Report {
  id: number;
  title: string;
  description: string;
  type: string;
  status: string;
  author_id?: number;
  author?: string;
  created_at: string;
}

// Punishment types
export interface Punishment {
  id: number;
  target_id: number;
  target_name?: string;
  target?: { username: string };
  assigned_by?: number;
  assigner_name?: string;
  assigner?: { username: string };
  type: string;
  description?: string;
  status: string;
  created_at: string;
  expires_at?: string;
}

// Capture types
export interface Capture {
  id: number;
  target_name: string;
  capture_date: string;
  status: string;
  points: number;
  difficulty: string;
  method: string;
}

export interface CaptureDetail {
  id: number;
  daemon_name: string;
  target_name: string;
  capture_date: string;
  status: string;
  points: number;
  difficulty: string;
  method: string;
}

// API Response types
export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
}

export interface DashboardResponse {
  welcome_message: string;
  system_stats: DashboardStats;
  recent_reports: Report[];
  top_daemons: LeaderboardEntry[];
  recent_activity: any[];
  all_captures: CaptureDetail[];
}

export interface DaemonDashboardResponse {
  welcome_message: string;
  user_stats: DaemonStats;
  active_missions: any[];
  leaderboard: LeaderboardEntry[];
  recent_chaos: any[];
  active_punishments: Punishment[];
  recent_captures: Capture[];
}

// Resistance types
export interface ResistanceStats {
  survival_tips: number;
  emergency_contacts: number;
  safe_locations: number;
  resistance_members: number;
  totalAdmins: number;
  freeAdmins: number;
  capturedAdmins: number;
  anonymousReports: number;
}
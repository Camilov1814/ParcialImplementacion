import type { User } from './auth';

export interface DashboardStats {
  totalUsers: number;
  totalDaemons: number;
  capturedNetworkAdmins: number;
  totalReports: number;
}

export interface DaemonStats {
  captures: number;
  reports: number;
  points: number;
  ranking: number;
}

export interface LeaderboardEntry {
  id: number;
  username: string;
  name: string;
  points: number;
  captures: number;
}

export interface Report {
  id: number;
  title: string;
  description: string;
  type: string; // "resistance", "capture", "anonymous"
  status: string;
  createdAt?: string;
  author?: {
    id: number;
    username: string;
    name: string;
  };
}

export interface Punishment {
  id: number;
  target?: {
    id: number;
    username: string;
    name: string;
  };
  assigner?: {
    id: number;
    username: string;
    name: string;
  };
  target_name?: string;    // For list items
  assigner_name?: string;  // For list items
  type: string;
  description?: string;    // Only available in detailed view
  status: string;
  created_at: string;
  updated_at?: string;
  expires_at?: string;
}

export interface ResistanceStats {
  totalAdmins: number;
  freeAdmins: number;
  capturedAdmins: number;
  anonymousReports: number;
}
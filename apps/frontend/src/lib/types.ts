// Shared domain types mirroring the backend API contracts.

export type Role = 'USER' | 'ADMIN';
export type TaskStatus = 'TODO' | 'IN_PROGRESS' | 'COMPLETED';
export type TaskPriority = 'LOW' | 'MEDIUM' | 'HIGH';

export interface User {
  id: string;
  name: string;
  email: string;
  role: Role;
  createdAt: string;
  updatedAt: string;
}

export interface Attachment {
  id: string;
  taskId: string;
  userId: string;
  fileName: string;
  contentType: string;
  sizeBytes: number;
  createdAt: string;
}

export interface Task {
  id: string;
  userId: string;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  dueDate: string | null;
  createdAt: string;
  updatedAt: string;
  attachments?: Attachment[];
  owner?: User | null;
}

export interface ActivityLog {
  id: string;
  userId: string;
  taskId: string | null;
  action: string;
  metadata: string;
  createdAt: string;
  user?: User | null;
}

export interface PageMeta {
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export interface ApiEnvelope<T> {
  success: boolean;
  data: T;
  meta?: PageMeta;
  error?: { code: string; message: string; details?: unknown };
}

export interface AuthResponse {
  token: string;
  user: User;
}

export type SortOption =
  | 'created_desc'
  | 'created_asc'
  | 'due_date_asc'
  | 'due_date_desc'
  | 'priority_desc'
  | 'priority_asc'
  | 'title_asc'
  | 'title_desc';

export interface TaskFilters {
  search?: string;
  status?: TaskStatus | '';
  priority?: TaskPriority | '';
  sort?: SortOption;
  scope?: 'me' | 'all';
  page?: number;
  limit?: number;
}

export interface TaskInput {
  title: string;
  description?: string;
  status?: TaskStatus;
  priority?: TaskPriority;
  dueDate?: string | null;
}

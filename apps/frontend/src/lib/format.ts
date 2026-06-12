import { format, formatDistanceToNow, isValid, parseISO } from 'date-fns';
import type { TaskPriority, TaskStatus } from './types';

/** formatDate renders an ISO string as a short human date, or a dash if absent. */
export function formatDate(iso: string | null | undefined): string {
  if (!iso) return '—';
  const date = parseISO(iso);
  return isValid(date) ? format(date, 'MMM d, yyyy') : '—';
}

/** formatRelative renders an ISO string as e.g. "3 hours ago". */
export function formatRelative(iso: string | null | undefined): string {
  if (!iso) return '—';
  const date = parseISO(iso);
  return isValid(date) ? `${formatDistanceToNow(date)} ago` : '—';
}

/** toDateInputValue converts an ISO string to a yyyy-MM-dd value for <input type=date>. */
export function toDateInputValue(iso: string | null | undefined): string {
  if (!iso) return '';
  const date = parseISO(iso);
  return isValid(date) ? format(date, 'yyyy-MM-dd') : '';
}

/** formatBytes renders a byte count in human units. */
export function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export const STATUS_LABELS: Record<TaskStatus, string> = {
  TODO: 'To Do',
  IN_PROGRESS: 'In Progress',
  COMPLETED: 'Completed',
};

export const PRIORITY_LABELS: Record<TaskPriority, string> = {
  LOW: 'Low',
  MEDIUM: 'Medium',
  HIGH: 'High',
};

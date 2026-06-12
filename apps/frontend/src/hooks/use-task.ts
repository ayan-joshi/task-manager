'use client';

import {
  useMutation,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query';
import { toast } from 'sonner';
import { api, API_URL, extractError } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import type { ActivityLog, ApiEnvelope, Attachment, Task } from '@/lib/types';
import { taskKeys } from './use-tasks';

/** useTask fetches a single task with its attachments. */
export function useTask(id: string) {
  return useQuery({
    queryKey: taskKeys.detail(id),
    queryFn: async (): Promise<Task> => {
      const { data } = await api.get<ApiEnvelope<Task>>(`/tasks/${id}`);
      return data.data;
    },
    enabled: Boolean(id),
  });
}

/** useTaskActivity fetches the audit trail for a single task. */
export function useTaskActivity(id: string) {
  return useQuery({
    queryKey: ['tasks', 'activity', id],
    queryFn: async (): Promise<ActivityLog[]> => {
      const { data } = await api.get<ApiEnvelope<ActivityLog[]>>(`/tasks/${id}/activity`, {
        params: { page: 1, limit: 50 },
      });
      return data.data;
    },
    enabled: Boolean(id),
  });
}

/** useUploadAttachment uploads a file to a task via multipart/form-data. */
export function useUploadAttachment(taskId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (file: File) => {
      const form = new FormData();
      form.append('file', file);
      const { data } = await api.post<ApiEnvelope<Attachment>>(
        `/tasks/${taskId}/attachments`,
        form,
        { headers: { 'Content-Type': 'multipart/form-data' } },
      );
      return data.data;
    },
    onSuccess: () => {
      toast.success('Attachment uploaded');
      qc.invalidateQueries({ queryKey: taskKeys.detail(taskId) });
      qc.invalidateQueries({ queryKey: ['tasks', 'activity', taskId] });
    },
    onError: (error) => toast.error(extractError(error, 'Upload failed')),
  });
}

/** useDeleteAttachment removes an attachment from a task. */
export function useDeleteAttachment(taskId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (attachmentId: string) => {
      await api.delete(`/attachments/${attachmentId}`);
      return attachmentId;
    },
    onSuccess: () => {
      toast.success('Attachment removed');
      qc.invalidateQueries({ queryKey: taskKeys.detail(taskId) });
      qc.invalidateQueries({ queryKey: ['tasks', 'activity', taskId] });
    },
    onError: (error) => toast.error(extractError(error, 'Failed to remove attachment')),
  });
}

/**
 * downloadAttachment fetches a protected attachment with the auth header and
 * triggers a browser download. A plain anchor href cannot carry the bearer
 * token, so we fetch the bytes and create an object URL.
 */
export async function downloadAttachment(attachment: Attachment) {
  try {
    const token = useAuthStore.getState().token;
    const res = await fetch(`${API_URL}/attachments/${attachment.id}`, {
      headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    });
    if (!res.ok) throw new Error('download failed');
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = attachment.fileName;
    document.body.appendChild(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
  } catch {
    toast.error('Could not download file');
  }
}

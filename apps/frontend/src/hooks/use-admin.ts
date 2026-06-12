'use client';

import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { ActivityLog, ApiEnvelope, PageMeta, User } from '@/lib/types';

/** useAdminUsers fetches a paginated list of all users (admin only). */
export function useAdminUsers(page: number, limit = 20) {
  return useQuery({
    queryKey: ['admin', 'users', page, limit],
    queryFn: async () => {
      const { data } = await api.get<ApiEnvelope<User[]>>('/admin/users', {
        params: { page, limit },
      });
      return { users: data.data, meta: data.meta as PageMeta };
    },
  });
}

/** useActivityLogs fetches a paginated activity log (admin only). */
export function useActivityLogs(page: number, limit = 20) {
  return useQuery({
    queryKey: ['admin', 'activity-logs', page, limit],
    queryFn: async () => {
      const { data } = await api.get<ApiEnvelope<ActivityLog[]>>('/admin/activity-logs', {
        params: { page, limit },
      });
      return { logs: data.data, meta: data.meta as PageMeta };
    },
  });
}

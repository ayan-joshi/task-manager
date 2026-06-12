'use client';

import * as React from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { API_URL } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import { taskKeys } from './use-tasks';

/**
 * useTaskEvents subscribes to the backend SSE stream and invalidates task
 * queries whenever another session creates, updates, or deletes a task, so the
 * UI reflects server-side changes in real time without polling.
 *
 * EventSource cannot send Authorization headers, so the JWT is passed as a query
 * parameter and validated by the server.
 */
export function useTaskEvents() {
  const token = useAuthStore((s) => s.token);
  const qc = useQueryClient();

  React.useEffect(() => {
    if (!token) return;

    const source = new EventSource(`${API_URL}/events?token=${encodeURIComponent(token)}`);

    const invalidate = () => {
      void qc.invalidateQueries({ queryKey: taskKeys.all });
    };

    source.addEventListener('task.created', invalidate);
    source.addEventListener('task.updated', invalidate);
    source.addEventListener('task.deleted', invalidate);

    // The browser auto-reconnects on transient errors; nothing to do here.
    source.onerror = () => {};

    return () => source.close();
  }, [token, qc]);
}

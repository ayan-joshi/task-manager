'use client';

import { Activity } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { formatRelative } from '@/lib/format';
import type { ActivityLog } from '@/lib/types';

const ACTION_LABELS: Record<string, string> = {
  TASK_CREATED: 'Task created',
  TASK_UPDATED: 'Task updated',
  TASK_DELETED: 'Task deleted',
  TASK_STATUS_CHANGED: 'Status changed',
  ATTACHMENT_ADDED: 'Attachment added',
  ATTACHMENT_DELETED: 'Attachment removed',
};

function describe(log: ActivityLog): string {
  const base = ACTION_LABELS[log.action] ?? log.action;
  try {
    const meta = log.metadata ? (JSON.parse(log.metadata) as Record<string, unknown>) : null;
    if (meta?.status) return `${base} → ${String(meta.status)}`;
    if (meta?.fileName) return `${base}: ${String(meta.fileName)}`;
  } catch {
    // Ignore malformed metadata.
  }
  return base;
}

export function ActivityTimeline({
  logs,
  isLoading,
}: {
  logs: ActivityLog[];
  isLoading?: boolean;
}) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-base">
          <Activity className="h-4 w-4" />
          Activity
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="space-y-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-5 w-full" />
            ))}
          </div>
        ) : logs.length === 0 ? (
          <p className="text-sm text-muted-foreground">No activity recorded yet.</p>
        ) : (
          <ol className="relative space-y-4 border-l pl-4">
            {logs.map((log) => (
              <li key={log.id} className="relative">
                <span className="absolute -left-[21px] top-1.5 h-2 w-2 rounded-full bg-primary" />
                <p className="text-sm font-medium">{describe(log)}</p>
                <p className="text-xs text-muted-foreground">{formatRelative(log.createdAt)}</p>
              </li>
            ))}
          </ol>
        )}
      </CardContent>
    </Card>
  );
}

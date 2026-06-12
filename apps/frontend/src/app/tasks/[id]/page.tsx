'use client';

import * as React from 'react';
import Link from 'next/link';
import { useParams } from 'next/navigation';
import { ArrowLeft, CalendarDays, Clock, Loader2, Pencil } from 'lucide-react';
import { AuthGuard } from '@/components/auth-guard';
import { AppShell } from '@/components/app-shell';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { EmptyState } from '@/components/empty-state';
import { PriorityBadge, StatusBadge } from '@/components/tasks/task-badges';
import { AttachmentsPanel } from '@/components/tasks/attachments-panel';
import { ActivityTimeline } from '@/components/tasks/activity-timeline';
import { TaskFormDialog, type TaskFormSubmit } from '@/components/tasks/task-form-dialog';
import { useTask, useTaskActivity } from '@/hooks/use-task';
import { useUpdateTask } from '@/hooks/use-tasks';
import { formatDate, formatRelative } from '@/lib/format';
import { ClipboardX } from 'lucide-react';

function TaskDetailContent({ id }: { id: string }) {
  const { data: task, isLoading, isError } = useTask(id);
  const { data: activity, isLoading: activityLoading } = useTaskActivity(id);
  const updateTask = useUpdateTask({ page: 1, limit: 10, sort: 'created_desc' });
  const [editOpen, setEditOpen] = React.useState(false);

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError || !task) {
    return (
      <EmptyState
        icon={ClipboardX}
        title="Task not found"
        description="This task may have been deleted or you don't have access to it."
        action={
          <Button asChild variant="outline">
            <Link href="/dashboard">Back to tasks</Link>
          </Button>
        }
      />
    );
  }

  const handleSubmit = (values: TaskFormSubmit) => {
    updateTask.mutate({ id: task.id, input: values }, { onSuccess: () => setEditOpen(false) });
  };

  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <div className="flex items-center justify-between">
        <Button asChild variant="ghost" size="sm">
          <Link href="/dashboard">
            <ArrowLeft className="h-4 w-4" />
            Back
          </Link>
        </Button>
        <Button variant="outline" size="sm" onClick={() => setEditOpen(true)}>
          <Pencil className="h-4 w-4" />
          Edit
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-wrap items-center gap-2">
            <StatusBadge status={task.status} />
            <PriorityBadge priority={task.priority} />
          </div>
          <CardTitle className="text-2xl">{task.title}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <h3 className="mb-1 text-sm font-medium text-muted-foreground">Description</h3>
            {task.description ? (
              <p className="whitespace-pre-wrap text-sm">{task.description}</p>
            ) : (
              <p className="text-sm italic text-muted-foreground">No description provided.</p>
            )}
          </div>
          <div className="flex flex-wrap gap-6 text-sm">
            <div className="flex items-center gap-2 text-muted-foreground">
              <CalendarDays className="h-4 w-4" />
              Due: <span className="text-foreground">{formatDate(task.dueDate)}</span>
            </div>
            <div className="flex items-center gap-2 text-muted-foreground">
              <Clock className="h-4 w-4" />
              Created <span className="text-foreground">{formatRelative(task.createdAt)}</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-6 md:grid-cols-2">
        <AttachmentsPanel taskId={task.id} attachments={task.attachments ?? []} />
        <ActivityTimeline logs={activity ?? []} isLoading={activityLoading} />
      </div>

      <TaskFormDialog
        open={editOpen}
        onOpenChange={setEditOpen}
        task={task}
        onSubmit={handleSubmit}
        isPending={updateTask.isPending}
      />
    </div>
  );
}

export default function TaskDetailPage() {
  const params = useParams<{ id: string }>();
  return (
    <AuthGuard>
      <AppShell>
        <TaskDetailContent id={params.id} />
      </AppShell>
    </AuthGuard>
  );
}

'use client';

import * as React from 'react';
import { ClipboardList, LayoutGrid, List, Plus } from 'lucide-react';
import { AuthGuard } from '@/components/auth-guard';
import { AppShell } from '@/components/app-shell';
import { Button } from '@/components/ui/button';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { EmptyState } from '@/components/empty-state';
import { TaskFiltersBar } from '@/components/tasks/task-filters';
import { TaskTable } from '@/components/tasks/task-table';
import { TaskCardGrid } from '@/components/tasks/task-card-grid';
import { PaginationBar } from '@/components/tasks/pagination-bar';
import { TaskFormDialog, type TaskFormSubmit } from '@/components/tasks/task-form-dialog';
import { ConfirmDialog } from '@/components/confirm-dialog';
import {
  useCreateTask,
  useDeleteTask,
  useTasks,
  useUpdateTask,
} from '@/hooks/use-tasks';
import type { Task, TaskFilters } from '@/lib/types';

function DashboardContent() {
  const [filters, setFilters] = React.useState<TaskFilters>({
    page: 1,
    limit: 10,
    sort: 'created_desc',
    search: '',
    status: '',
    priority: '',
  });
  const [view, setView] = React.useState<'table' | 'card'>('table');
  const [formOpen, setFormOpen] = React.useState(false);
  const [editing, setEditing] = React.useState<Task | null>(null);
  const [deleting, setDeleting] = React.useState<Task | null>(null);

  const { data, isLoading, isError } = useTasks(filters);
  const createTask = useCreateTask(filters);
  const updateTask = useUpdateTask(filters);
  const deleteTask = useDeleteTask(filters);

  const patchFilters = (patch: Partial<TaskFilters>) =>
    setFilters((prev) => ({ ...prev, ...patch }));

  const openCreate = () => {
    setEditing(null);
    setFormOpen(true);
  };
  const openEdit = (task: Task) => {
    setEditing(task);
    setFormOpen(true);
  };

  const handleSubmit = (values: TaskFormSubmit) => {
    if (editing) {
      updateTask.mutate(
        { id: editing.id, input: values },
        { onSuccess: () => setFormOpen(false) },
      );
    } else {
      createTask.mutate(values, { onSuccess: () => setFormOpen(false) });
    }
  };

  const handleToggleComplete = (task: Task) => {
    updateTask.mutate({
      id: task.id,
      input: { status: task.status === 'COMPLETED' ? 'TODO' : 'COMPLETED' },
    });
  };

  const handleDelete = () => {
    if (!deleting) return;
    deleteTask.mutate(deleting.id, { onSuccess: () => setDeleting(null) });
  };

  const tasks = data?.tasks ?? [];
  const meta = data?.meta;
  const hasActiveFilters = Boolean(filters.search || filters.status || filters.priority);
  const showEmpty = !isLoading && !isError && tasks.length === 0;

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Tasks</h1>
          <p className="text-sm text-muted-foreground">
            Manage, filter, and track your work.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Tabs value={view} onValueChange={(v) => setView(v as 'table' | 'card')}>
            <TabsList>
              <TabsTrigger value="table" aria-label="Table view">
                <List className="h-4 w-4" />
              </TabsTrigger>
              <TabsTrigger value="card" aria-label="Card view">
                <LayoutGrid className="h-4 w-4" />
              </TabsTrigger>
            </TabsList>
          </Tabs>
          <Button onClick={openCreate}>
            <Plus className="h-4 w-4" />
            New task
          </Button>
        </div>
      </div>

      <TaskFiltersBar filters={filters} onChange={patchFilters} />

      {isError ? (
        <EmptyState
          icon={ClipboardList}
          title="Couldn't load tasks"
          description="There was a problem fetching your tasks. Please try again."
        />
      ) : showEmpty ? (
        <EmptyState
          icon={ClipboardList}
          title={hasActiveFilters ? 'No matching tasks' : 'No tasks yet'}
          description={
            hasActiveFilters
              ? 'Try adjusting your search or filters.'
              : 'Create your first task to get started.'
          }
          action={
            hasActiveFilters ? null : (
              <Button onClick={openCreate}>
                <Plus className="h-4 w-4" />
                New task
              </Button>
            )
          }
        />
      ) : view === 'table' ? (
        <TaskTable
          tasks={tasks}
          isLoading={isLoading}
          onEdit={openEdit}
          onDelete={setDeleting}
          onToggleComplete={handleToggleComplete}
        />
      ) : (
        <TaskCardGrid
          tasks={tasks}
          isLoading={isLoading}
          onEdit={openEdit}
          onDelete={setDeleting}
          onToggleComplete={handleToggleComplete}
        />
      )}

      {meta && !showEmpty && !isError ? (
        <PaginationBar meta={meta} onPageChange={(page) => patchFilters({ page })} />
      ) : null}

      <TaskFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        task={editing}
        onSubmit={handleSubmit}
        isPending={createTask.isPending || updateTask.isPending}
      />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(open) => !open && setDeleting(null)}
        title="Delete task"
        description={`Are you sure you want to delete "${deleting?.title ?? ''}"? This cannot be undone.`}
        confirmLabel="Delete"
        destructive
        isPending={deleteTask.isPending}
        onConfirm={handleDelete}
      />
    </div>
  );
}

export default function DashboardPage() {
  return (
    <AuthGuard>
      <AppShell>
        <DashboardContent />
      </AppShell>
    </AuthGuard>
  );
}

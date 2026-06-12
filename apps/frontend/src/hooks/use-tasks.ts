'use client';

import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from '@tanstack/react-query';
import { toast } from 'sonner';
import { api, extractError } from '@/lib/api';
import type { ApiEnvelope, PageMeta, Task, TaskFilters, TaskInput } from '@/lib/types';

interface TaskListResult {
  tasks: Task[];
  meta: PageMeta;
}

/** taskKeys centralizes query keys so mutations can target the right caches. */
export const taskKeys = {
  all: ['tasks'] as const,
  list: (filters: TaskFilters) => ['tasks', 'list', filters] as const,
  detail: (id: string) => ['tasks', 'detail', id] as const,
};

function buildParams(filters: TaskFilters): Record<string, string | number> {
  const params: Record<string, string | number> = {
    page: filters.page ?? 1,
    limit: filters.limit ?? 10,
    sort: filters.sort ?? 'created_desc',
  };
  if (filters.search) params.search = filters.search;
  if (filters.status) params.status = filters.status;
  if (filters.priority) params.priority = filters.priority;
  return params;
}

/** useTasks fetches a filtered, paginated page of tasks. */
export function useTasks(filters: TaskFilters) {
  return useQuery({
    queryKey: taskKeys.list(filters),
    queryFn: async (): Promise<TaskListResult> => {
      const { data } = await api.get<ApiEnvelope<Task[]>>('/tasks', {
        params: buildParams(filters),
      });
      return { tasks: data.data, meta: data.meta as PageMeta };
    },
    placeholderData: keepPreviousData,
  });
}

/** useCreateTask optimistically prepends the new task, rolling back on failure. */
export function useCreateTask(filters: TaskFilters) {
  const qc = useQueryClient();
  const key = taskKeys.list(filters);

  return useMutation({
    mutationFn: async (input: TaskInput) => {
      const { data } = await api.post<ApiEnvelope<Task>>('/tasks', input);
      return data.data;
    },
    onMutate: async (input) => {
      await qc.cancelQueries({ queryKey: key });
      const previous = qc.getQueryData<TaskListResult>(key);

      const optimistic: Task = {
        id: `optimistic-${Date.now()}`,
        userId: 'optimistic',
        title: input.title,
        description: input.description ?? '',
        status: input.status ?? 'TODO',
        priority: input.priority ?? 'MEDIUM',
        dueDate: input.dueDate ?? null,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };

      if (previous) {
        qc.setQueryData<TaskListResult>(key, {
          tasks: [optimistic, ...previous.tasks],
          meta: { ...previous.meta, total: previous.meta.total + 1 },
        });
      }
      return { previous };
    },
    onError: (error, _input, context) => {
      if (context?.previous) qc.setQueryData(key, context.previous);
      toast.error(extractError(error, 'Failed to create task'));
    },
    onSuccess: () => toast.success('Task created'),
    onSettled: () => qc.invalidateQueries({ queryKey: taskKeys.all }),
  });
}

/** useUpdateTask optimistically patches the cached task, rolling back on failure. */
export function useUpdateTask(filters: TaskFilters) {
  const qc = useQueryClient();
  const key = taskKeys.list(filters);

  return useMutation({
    mutationFn: async ({
      id,
      input,
    }: {
      id: string;
      input: Partial<TaskInput> & { clearDueDate?: boolean };
    }) => {
      const { data } = await api.put<ApiEnvelope<Task>>(`/tasks/${id}`, input);
      return data.data;
    },
    onMutate: async ({ id, input }) => {
      await qc.cancelQueries({ queryKey: key });
      const previous = qc.getQueryData<TaskListResult>(key);
      if (previous) {
        qc.setQueryData<TaskListResult>(key, {
          ...previous,
          tasks: previous.tasks.map((t) =>
            t.id === id ? { ...t, ...input, updatedAt: new Date().toISOString() } : t,
          ),
        });
      }
      return { previous };
    },
    onError: (error, _vars, context) => {
      if (context?.previous) qc.setQueryData(key, context.previous);
      toast.error(extractError(error, 'Failed to update task'));
    },
    onSuccess: (task) => {
      qc.setQueryData(taskKeys.detail(task.id), task);
      toast.success('Task updated');
    },
    onSettled: () => qc.invalidateQueries({ queryKey: taskKeys.all }),
  });
}

/** useDeleteTask optimistically removes the task, rolling back on failure. */
export function useDeleteTask(filters: TaskFilters) {
  const qc = useQueryClient();
  const key = taskKeys.list(filters);

  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/tasks/${id}`);
      return id;
    },
    onMutate: async (id) => {
      await qc.cancelQueries({ queryKey: key });
      const previous = qc.getQueryData<TaskListResult>(key);
      if (previous) {
        qc.setQueryData<TaskListResult>(key, {
          tasks: previous.tasks.filter((t) => t.id !== id),
          meta: { ...previous.meta, total: Math.max(0, previous.meta.total - 1) },
        });
      }
      return { previous };
    },
    onError: (error, _id, context) => {
      if (context?.previous) qc.setQueryData(key, context.previous);
      toast.error(extractError(error, 'Failed to delete task'));
    },
    onSuccess: () => toast.success('Task deleted'),
    onSettled: () => qc.invalidateQueries({ queryKey: taskKeys.all }),
  });
}

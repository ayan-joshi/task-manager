import * as React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) =>
    React.createElement('a', { href }, children),
}));

import { TaskTable } from '@/components/tasks/task-table';
import type { Task } from '@/lib/types';

function makeTask(overrides: Partial<Task>): Task {
  return {
    id: 'id',
    userId: 'user',
    title: 'Task',
    description: '',
    status: 'TODO',
    priority: 'MEDIUM',
    dueDate: null,
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
    ...overrides,
  };
}

describe('TaskTable', () => {
  const tasks: Task[] = [
    makeTask({ id: '1', title: 'Write quarterly report', status: 'TODO', priority: 'HIGH' }),
    makeTask({ id: '2', title: 'Review pull request', status: 'COMPLETED', priority: 'LOW' }),
  ];

  it('renders a row for each task with its title', () => {
    render(
      <TaskTable
        tasks={tasks}
        onEdit={vi.fn()}
        onDelete={vi.fn()}
        onToggleComplete={vi.fn()}
      />,
    );

    expect(screen.getByText('Write quarterly report')).toBeInTheDocument();
    expect(screen.getByText('Review pull request')).toBeInTheDocument();
    expect(screen.getAllByTestId('task-row')).toHaveLength(2);
  });

  it('renders skeleton placeholders while loading', () => {
    const { container } = render(
      <TaskTable
        tasks={[]}
        isLoading
        onEdit={vi.fn()}
        onDelete={vi.fn()}
        onToggleComplete={vi.fn()}
      />,
    );

    expect(screen.queryByTestId('task-row')).not.toBeInTheDocument();
    expect(container.querySelectorAll('.animate-pulse').length).toBeGreaterThan(0);
  });
});

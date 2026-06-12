import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

import { TaskFormDialog } from '@/components/tasks/task-form-dialog';

describe('TaskFormDialog (create task)', () => {
  it('blocks submission and shows an error when the title is empty', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    render(<TaskFormDialog open onOpenChange={vi.fn()} onSubmit={onSubmit} />);

    await user.click(screen.getByRole('button', { name: /create task/i }));

    expect(await screen.findByText('Title is required')).toBeInTheDocument();
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('submits the entered task with default status and priority', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    render(<TaskFormDialog open onOpenChange={vi.fn()} onSubmit={onSubmit} />);

    await user.type(screen.getByLabelText('Title'), 'My new task');
    await user.click(screen.getByRole('button', { name: /create task/i }));

    await waitFor(() => expect(onSubmit).toHaveBeenCalledTimes(1));
    const payload = onSubmit.mock.calls[0][0];
    expect(payload.title).toBe('My new task');
    expect(payload.status).toBe('TODO');
    expect(payload.priority).toBe('MEDIUM');
  });
});

import { Badge } from '@/components/ui/badge';
import { PRIORITY_LABELS, STATUS_LABELS } from '@/lib/format';
import type { TaskPriority, TaskStatus } from '@/lib/types';

const STATUS_VARIANT: Record<TaskStatus, React.ComponentProps<typeof Badge>['variant']> = {
  TODO: 'secondary',
  IN_PROGRESS: 'info',
  COMPLETED: 'success',
};

const PRIORITY_VARIANT: Record<TaskPriority, React.ComponentProps<typeof Badge>['variant']> = {
  LOW: 'outline',
  MEDIUM: 'warning',
  HIGH: 'destructive',
};

export function StatusBadge({ status }: { status: TaskStatus }) {
  return <Badge variant={STATUS_VARIANT[status]}>{STATUS_LABELS[status]}</Badge>;
}

export function PriorityBadge({ priority }: { priority: TaskPriority }) {
  return <Badge variant={PRIORITY_VARIANT[priority]}>{PRIORITY_LABELS[priority]}</Badge>;
}

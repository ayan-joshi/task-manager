'use client';

import Link from 'next/link';
import { CalendarDays, MoreHorizontal, Paperclip, Pencil, Trash2, CheckCircle2 } from 'lucide-react';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { PriorityBadge, StatusBadge } from './task-badges';
import { formatDate } from '@/lib/format';
import type { TaskTableProps } from './task-table';

export function TaskCardGrid({
  tasks,
  isLoading,
  onEdit,
  onDelete,
  onToggleComplete,
}: TaskTableProps) {
  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <Card key={i}>
            <CardHeader>
              <Skeleton className="h-5 w-3/4" />
            </CardHeader>
            <CardContent className="space-y-2">
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-2/3" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {tasks.map((task) => (
        <Card key={task.id} className="flex flex-col">
          <CardHeader className="flex-row items-start justify-between space-y-0 pb-2">
            <Link href={`/tasks/${task.id}`} className="font-semibold leading-tight hover:underline">
              {task.title}
            </Link>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" aria-label={`Actions for ${task.title}`}>
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => onToggleComplete(task)}>
                  <CheckCircle2 className="h-4 w-4" />
                  {task.status === 'COMPLETED' ? 'Mark as to do' : 'Mark complete'}
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => onEdit(task)}>
                  <Pencil className="h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem className="text-destructive" onClick={() => onDelete(task)}>
                  <Trash2 className="h-4 w-4" />
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </CardHeader>
          <CardContent className="flex-1">
            {task.description ? (
              <p className="line-clamp-3 text-sm text-muted-foreground">{task.description}</p>
            ) : (
              <p className="text-sm italic text-muted-foreground">No description</p>
            )}
          </CardContent>
          <CardFooter className="flex flex-wrap items-center gap-2 pt-0">
            <StatusBadge status={task.status} />
            <PriorityBadge priority={task.priority} />
            {task.dueDate ? (
              <span className="ml-auto flex items-center gap-1 text-xs text-muted-foreground">
                <CalendarDays className="h-3 w-3" />
                {formatDate(task.dueDate)}
              </span>
            ) : null}
            {task.attachments && task.attachments.length > 0 ? (
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <Paperclip className="h-3 w-3" />
                {task.attachments.length}
              </span>
            ) : null}
          </CardFooter>
        </Card>
      ))}
    </div>
  );
}

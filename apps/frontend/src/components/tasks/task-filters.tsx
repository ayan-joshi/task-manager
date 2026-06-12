'use client';

import { Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import type { SortOption, TaskFilters, TaskPriority, TaskStatus } from '@/lib/types';

const ALL = '__all__';

const SORT_OPTIONS: { value: SortOption; label: string }[] = [
  { value: 'created_desc', label: 'Newest first' },
  { value: 'created_asc', label: 'Oldest first' },
  { value: 'due_date_asc', label: 'Due date ↑' },
  { value: 'due_date_desc', label: 'Due date ↓' },
  { value: 'priority_desc', label: 'Priority: High → Low' },
  { value: 'priority_asc', label: 'Priority: Low → High' },
  { value: 'title_asc', label: 'Title A → Z' },
  { value: 'title_desc', label: 'Title Z → A' },
];

export function TaskFiltersBar({
  filters,
  onChange,
}: {
  filters: TaskFilters;
  onChange: (patch: Partial<TaskFilters>) => void;
}) {
  return (
    <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
      <div className="relative sm:col-span-2 lg:col-span-1">
        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search by title…"
          className="pl-8"
          aria-label="Search tasks"
          value={filters.search ?? ''}
          onChange={(e) => onChange({ search: e.target.value, page: 1 })}
        />
      </div>

      <Select
        value={filters.status || ALL}
        onValueChange={(v) =>
          onChange({ status: v === ALL ? '' : (v as TaskStatus), page: 1 })
        }
      >
        <SelectTrigger aria-label="Filter by status">
          <SelectValue placeholder="Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value={ALL}>All statuses</SelectItem>
          <SelectItem value="TODO">To Do</SelectItem>
          <SelectItem value="IN_PROGRESS">In Progress</SelectItem>
          <SelectItem value="COMPLETED">Completed</SelectItem>
        </SelectContent>
      </Select>

      <Select
        value={filters.priority || ALL}
        onValueChange={(v) =>
          onChange({ priority: v === ALL ? '' : (v as TaskPriority), page: 1 })
        }
      >
        <SelectTrigger aria-label="Filter by priority">
          <SelectValue placeholder="Priority" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value={ALL}>All priorities</SelectItem>
          <SelectItem value="LOW">Low</SelectItem>
          <SelectItem value="MEDIUM">Medium</SelectItem>
          <SelectItem value="HIGH">High</SelectItem>
        </SelectContent>
      </Select>

      <Select
        value={filters.sort ?? 'created_desc'}
        onValueChange={(v) => onChange({ sort: v as SortOption })}
      >
        <SelectTrigger aria-label="Sort tasks">
          <SelectValue placeholder="Sort" />
        </SelectTrigger>
        <SelectContent>
          {SORT_OPTIONS.map((opt) => (
            <SelectItem key={opt.value} value={opt.value}>
              {opt.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}

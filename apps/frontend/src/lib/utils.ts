import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * cn merges Tailwind class names, resolving conflicts so the last-specified
 * utility wins (e.g. `cn('p-2', 'p-4')` → `'p-4'`).
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

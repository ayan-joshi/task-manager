'use client';

import * as React from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/lib/auth-store';
import type { Role } from '@/lib/types';
import { Loader2 } from 'lucide-react';

/**
 * AuthGuard protects client routes. It waits for the persisted auth state to
 * rehydrate (avoiding a flash of redirect on refresh), then redirects
 * unauthenticated users to /login and users lacking `requireRole` to /dashboard.
 */
export function AuthGuard({
  children,
  requireRole,
}: {
  children: React.ReactNode;
  requireRole?: Role;
}) {
  const router = useRouter();
  const { token, user, hydrated } = useAuthStore();

  React.useEffect(() => {
    if (!hydrated) return;
    if (!token) {
      router.replace('/login');
      return;
    }
    if (requireRole && user?.role !== requireRole) {
      router.replace('/dashboard');
    }
  }, [hydrated, token, user, requireRole, router]);

  if (!hydrated || !token || (requireRole && user?.role !== requireRole)) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return <>{children}</>;
}

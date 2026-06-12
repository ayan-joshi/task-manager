'use client';

import * as React from 'react';
import { useRouter } from 'next/navigation';
import { Loader2 } from 'lucide-react';
import { useAuthStore } from '@/lib/auth-store';

/**
 * The index route immediately routes the visitor based on auth state, which is
 * only known on the client (the token lives in localStorage).
 */
export default function HomePage() {
  const router = useRouter();
  const { token, hydrated } = useAuthStore();

  React.useEffect(() => {
    if (!hydrated) return;
    router.replace(token ? '/dashboard' : '/login');
  }, [hydrated, token, router]);

  return (
    <div className="flex min-h-screen items-center justify-center">
      <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
    </div>
  );
}

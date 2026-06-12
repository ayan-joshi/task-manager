'use client';

import { useMutation } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';
import { api, extractError } from '@/lib/api';
import { useAuthStore } from '@/lib/auth-store';
import type { ApiEnvelope, AuthResponse } from '@/lib/types';
import type { LoginValues, SignupValues } from '@/lib/validations';

/** useLogin authenticates the user, persists the session, and redirects. */
export function useLogin() {
  const router = useRouter();
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: async (values: LoginValues) => {
      const { data } = await api.post<ApiEnvelope<AuthResponse>>('/auth/login', values);
      return data.data;
    },
    onSuccess: (auth) => {
      setAuth(auth.token, auth.user);
      toast.success(`Welcome back, ${auth.user.name}`);
      router.push('/dashboard');
    },
    onError: (error) => toast.error(extractError(error, 'Login failed')),
  });
}

/** useSignup registers a new account, persists the session, and redirects. */
export function useSignup() {
  const router = useRouter();
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: async (values: SignupValues) => {
      const { data } = await api.post<ApiEnvelope<AuthResponse>>('/auth/signup', values);
      return data.data;
    },
    onSuccess: (auth) => {
      setAuth(auth.token, auth.user);
      toast.success('Account created');
      router.push('/dashboard');
    },
    onError: (error) => toast.error(extractError(error, 'Signup failed')),
  });
}

/** useLogout clears the session and returns to the login screen. */
export function useLogout() {
  const router = useRouter();
  const clear = useAuthStore((s) => s.clear);
  return () => {
    clear();
    router.push('/login');
  };
}

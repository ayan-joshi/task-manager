import axios, { AxiosError } from 'axios';
import { useAuthStore } from './auth-store';
import type { ApiEnvelope } from './types';

export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080/api';

/**
 * api is a shared Axios instance. A request interceptor attaches the bearer
 * token from the auth store; a response interceptor clears auth on 401 so an
 * expired session redirects the user back to login.
 */
export const api = axios.create({
  baseURL: API_URL,
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiEnvelope<unknown>>) => {
    if (error.response?.status === 401) {
      // Avoid clobbering the login attempt itself.
      const url = error.config?.url ?? '';
      if (!url.includes('/auth/login') && !url.includes('/auth/signup')) {
        useAuthStore.getState().clear();
      }
    }
    return Promise.reject(error);
  },
);

/** extractError returns a human-readable message from an Axios error. */
export function extractError(error: unknown, fallback = 'Something went wrong'): string {
  if (axios.isAxiosError(error)) {
    const envelope = error.response?.data as ApiEnvelope<unknown> | undefined;
    if (envelope?.error?.message) {
      return envelope.error.message;
    }
    if (error.message) {
      return error.message;
    }
  }
  return fallback;
}

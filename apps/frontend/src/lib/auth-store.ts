import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from './types';

interface AuthState {
  token: string | null;
  user: User | null;
  hydrated: boolean;
  setAuth: (token: string, user: User) => void;
  setUser: (user: User) => void;
  clear: () => void;
  setHydrated: () => void;
}

/**
 * useAuthStore persists the auth token and user to localStorage so a page
 * refresh keeps the user logged in. `hydrated` flips to true once the persisted
 * state has been rehydrated on the client, which guards against rendering
 * protected UI before we know whether a token exists.
 */
export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      hydrated: false,
      setAuth: (token, user) => set({ token, user }),
      setUser: (user) => set({ user }),
      clear: () => set({ token: null, user: null }),
      setHydrated: () => set({ hydrated: true }),
    }),
    {
      name: 'task-manager-auth',
      partialize: (state) => ({ token: state.token, user: state.user }),
      onRehydrateStorage: () => (state) => state?.setHydrated(),
    },
  ),
);

'use client';

import { create } from 'zustand';
import { authAPI } from '@/lib/api';
import type { User } from '@/lib/api';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (username: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  fetchUser: () => Promise<void>;
  clearError: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: typeof window !== 'undefined' ? !!localStorage.getItem('auth_token') : false,
  isLoading: false,
  error: null,

  login: async (username: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      const data = await authAPI.login(username, password);
      set({ user: data.user, isAuthenticated: true, isLoading: false });
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      set({
        error: err.response?.data?.error || 'Login failed',
        isLoading: false,
      });
      throw error;
    }
  },

  register: async (username: string, email: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      const data = await authAPI.register(username, email, password);
      set({ user: data.user, isAuthenticated: true, isLoading: false });
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      set({
        error: err.response?.data?.error || 'Registration failed',
        isLoading: false,
      });
      throw error;
    }
  },

  logout: () => {
    authAPI.logout();
    set({ user: null, isAuthenticated: false });
  },

  fetchUser: async () => {
    if (typeof window === 'undefined' || !localStorage.getItem('auth_token')) {
      set({ isAuthenticated: false });
      return;
    }

    set({ isLoading: true });
    try {
      const user = await authAPI.getMe();
      set({ user, isAuthenticated: true, isLoading: false });
    } catch {
      set({ user: null, isAuthenticated: false, isLoading: false });
      authAPI.logout();
    }
  },

  clearError: () => set({ error: null }),
}));

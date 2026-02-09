import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import { useAuthStore } from '../authStore';
import * as api from '../../lib/api';

// Mock the API
vi.mock('../../lib/api', () => ({
  authAPI: {
    login: vi.fn(),
    register: vi.fn(),
    logout: vi.fn(),
    getMe: vi.fn(),
  },
}));

describe('AuthStore', () => {
  beforeEach(() => {
    // Clear store state
    useAuthStore.setState({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
    });

    // Clear localStorage
    localStorage.clear();

    // Reset mocks
    vi.clearAllMocks();
  });

  describe('login', () => {
    it('should successfully login a user', async () => {
      const mockUser = {
        user_id: 1,
        username: 'testuser',
        email: 'test@example.com',
        current_streak_days: 5,
        total_study_time_seconds: 3600,
      };

      vi.mocked(api.authAPI.login).mockResolvedValueOnce({
        user: mockUser,
        token: 'test-token',
      });

      const { result } = renderHook(() => useAuthStore());

      await act(async () => {
        await result.current.login('testuser', 'password123');
      });

      expect(result.current.user).toEqual(mockUser);
      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it('should handle login failure', async () => {
      const mockError = {
        response: {
          data: {
            error: 'Invalid credentials',
          },
        },
      };

      vi.mocked(api.authAPI.login).mockRejectedValueOnce(mockError);

      const { result } = renderHook(() => useAuthStore());

      await act(async () => {
        try {
          await result.current.login('testuser', 'wrongpassword');
        } catch (error) {
          // Expected to throw
        }
      });

      expect(result.current.user).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(result.current.error).toBe('Invalid credentials');
      expect(result.current.isLoading).toBe(false);
    });

    it('should set loading state during login', async () => {
      const mockUser = {
        user_id: 1,
        username: 'testuser',
        email: 'test@example.com',
        current_streak_days: 0,
        total_study_time_seconds: 0,
      };

      vi.mocked(api.authAPI.login).mockImplementationOnce(
        () => new Promise((resolve) => setTimeout(() => resolve({
          user: mockUser,
          token: 'test-token',
        }), 100))
      );

      const { result } = renderHook(() => useAuthStore());

      act(() => {
        result.current.login('testuser', 'password123');
      });

      expect(result.current.isLoading).toBe(true);

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });
    });
  });

  describe('register', () => {
    it('should successfully register a user', async () => {
      const mockUser = {
        user_id: 1,
        username: 'newuser',
        email: 'new@example.com',
        current_streak_days: 0,
        total_study_time_seconds: 0,
      };

      vi.mocked(api.authAPI.register).mockResolvedValueOnce({
        user: mockUser,
        token: 'test-token',
      });

      const { result } = renderHook(() => useAuthStore());

      await act(async () => {
        await result.current.register('newuser', 'new@example.com', 'password123');
      });

      expect(result.current.user).toEqual(mockUser);
      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.error).toBeNull();
    });

    it('should handle registration failure', async () => {
      const mockError = {
        response: {
          data: {
            error: 'Username already exists',
          },
        },
      };

      vi.mocked(api.authAPI.register).mockRejectedValueOnce(mockError);

      const { result } = renderHook(() => useAuthStore());

      await act(async () => {
        try {
          await result.current.register('existinguser', 'test@example.com', 'password123');
        } catch (error) {
          // Expected to throw
        }
      });

      expect(result.current.user).toBeNull();
      expect(result.current.error).toBe('Username already exists');
    });
  });

  describe('logout', () => {
    it('should clear user state on logout', () => {
      const { result } = renderHook(() => useAuthStore());

      // Set initial state
      act(() => {
        useAuthStore.setState({
          user: {
            user_id: 1,
            username: 'testuser',
            email: 'test@example.com',
            current_streak_days: 5,
            total_study_time_seconds: 3600,
          },
          isAuthenticated: true,
        });
      });

      // Logout
      act(() => {
        result.current.logout();
      });

      expect(result.current.user).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(api.authAPI.logout).toHaveBeenCalled();
    });
  });

  describe('fetchUser', () => {
    it('should fetch user data when token exists', async () => {
      const mockUser = {
        user_id: 1,
        username: 'testuser',
        email: 'test@example.com',
        current_streak_days: 5,
        total_study_time_seconds: 3600,
      };

      localStorage.setItem('auth_token', 'test-token');
      vi.mocked(api.authAPI.getMe).mockResolvedValueOnce(mockUser);

      const { result } = renderHook(() => useAuthStore());

      await act(async () => {
        await result.current.fetchUser();
      });

      expect(result.current.user).toEqual(mockUser);
      expect(result.current.isAuthenticated).toBe(true);
    });

    it('should clear auth state when token is invalid', async () => {
      localStorage.setItem('auth_token', 'invalid-token');
      vi.mocked(api.authAPI.getMe).mockRejectedValueOnce(new Error('Unauthorized'));

      const { result } = renderHook(() => useAuthStore());

      await act(async () => {
        await result.current.fetchUser();
      });

      expect(result.current.user).toBeNull();
      expect(result.current.isAuthenticated).toBe(false);
      expect(api.authAPI.logout).toHaveBeenCalled();
    });

    it('should not fetch when no token exists', async () => {
      const { result } = renderHook(() => useAuthStore());

      await act(async () => {
        await result.current.fetchUser();
      });

      expect(api.authAPI.getMe).not.toHaveBeenCalled();
      expect(result.current.isAuthenticated).toBe(false);
    });
  });

  describe('clearError', () => {
    it('should clear error state', () => {
      const { result } = renderHook(() => useAuthStore());

      // Set error
      act(() => {
        useAuthStore.setState({ error: 'Test error' });
      });

      expect(result.current.error).toBe('Test error');

      // Clear error
      act(() => {
        result.current.clearError();
      });

      expect(result.current.error).toBeNull();
    });
  });
});

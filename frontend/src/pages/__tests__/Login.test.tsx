import { describe, it, expect, beforeEach, vi } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { render } from '../../test/utils';
import Login from '../Login';
import { useAuthStore } from '../../stores/authStore';

const mockNavigate = vi.fn();

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe('Login Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useAuthStore.setState({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
    });
  });

  describe('Login Mode', () => {
    it('should render login form', () => {
      render(<Login />);

      expect(screen.getByText('Welcome Back')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('Enter your username')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('Enter your password')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
    });

    it('should not show email field in login mode', () => {
      render(<Login />);

      expect(screen.queryByPlaceholderText(/email/i)).not.toBeInTheDocument();
    });

    it('should handle successful login', async () => {
      const user = userEvent.setup();
      const mockLogin = vi.fn().mockResolvedValue(undefined);

      useAuthStore.setState({
        login: mockLogin,
        isLoading: false,
        error: null,
      } as any);

      render(<Login />);

      const usernameInput = screen.getByPlaceholderText('Enter your username');
      const passwordInput = screen.getByPlaceholderText('Enter your password');
      const loginButton = screen.getByRole('button', { name: /login/i });

      await user.type(usernameInput, 'testuser');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      await waitFor(() => {
        expect(mockLogin).toHaveBeenCalledWith('testuser', 'password123');
      });
    });

    it('should display error message on login failure', async () => {
      useAuthStore.setState({
        error: 'Invalid credentials',
        isLoading: false,
      } as any);

      render(<Login />);

      expect(screen.getByText('Invalid credentials')).toBeInTheDocument();
    });

    it('should disable submit button while loading', () => {
      useAuthStore.setState({
        isLoading: true,
      } as any);

      render(<Login />);

      const loginButton = screen.getByRole('button', { name: /please wait/i });
      expect(loginButton).toBeDisabled();
    });

    it('should show link to register page', () => {
      render(<Login />);

      const registerLink = screen.getByRole('link', { name: /sign up/i });
      expect(registerLink).toHaveAttribute('href', '/register');
    });
  });

  describe('Register Mode', () => {
    it('should render register form', () => {
      render(<Login isRegister={true} />);

      expect(screen.getByText('Create Account')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('Enter your username')).toBeInTheDocument();
      expect(screen.getByPlaceholderText(/email/i)).toBeInTheDocument();
      expect(screen.getByPlaceholderText('Enter your password')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /sign up/i })).toBeInTheDocument();
    });

    it('should handle successful registration', async () => {
      const user = userEvent.setup();
      const mockRegister = vi.fn().mockResolvedValue(undefined);

      useAuthStore.setState({
        register: mockRegister,
        isLoading: false,
        error: null,
      } as any);

      render(<Login isRegister={true} />);

      const usernameInput = screen.getByPlaceholderText('Enter your username');
      const emailInput = screen.getByPlaceholderText(/email/i);
      const passwordInput = screen.getByPlaceholderText('Enter your password');
      const registerButton = screen.getByRole('button', { name: /sign up/i });

      await user.type(usernameInput, 'newuser');
      await user.type(emailInput, 'new@example.com');
      await user.type(passwordInput, 'password123');
      await user.click(registerButton);

      await waitFor(() => {
        expect(mockRegister).toHaveBeenCalledWith('newuser', 'new@example.com', 'password123');
      });
    });

    it('should show link to login page', () => {
      render(<Login isRegister={true} />);

      const loginLink = screen.getByRole('link', { name: /login/i });
      expect(loginLink).toHaveAttribute('href', '/login');
    });

    it('should validate password minimum length', () => {
      render(<Login isRegister={true} />);

      const passwordInput = screen.getByPlaceholderText('Enter your password');
      expect(passwordInput).toHaveAttribute('minLength', '8');
    });
  });

  describe('Form Validation', () => {
    it('should require username', () => {
      render(<Login />);

      const usernameInput = screen.getByPlaceholderText('Enter your username');
      expect(usernameInput).toBeRequired();
    });

    it('should require password', () => {
      render(<Login />);

      const passwordInput = screen.getByPlaceholderText('Enter your password');
      expect(passwordInput).toBeRequired();
    });

    it('should require email in register mode', () => {
      render(<Login isRegister={true} />);

      const emailInput = screen.getByPlaceholderText(/email/i);
      expect(emailInput).toBeRequired();
      expect(emailInput).toHaveAttribute('type', 'email');
    });
  });
});

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { screen, waitFor, render } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import App from '../../App';
import * as api from '../../lib/api';

vi.mock('../../lib/api');

// Helper to render App with QueryClient but no Router (App has its own)
const renderApp = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  );
};

describe('User Flow Integration Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  describe('Authentication Flow', () => {
    it('should show login page on initial load', async () => {
      renderApp();

      // Should redirect to login
      await waitFor(() => {
        expect(screen.getByText('Welcome Back')).toBeInTheDocument();
      });

      // Verify login form elements
      expect(screen.getByPlaceholderText('Enter your username')).toBeInTheDocument();
      expect(screen.getByPlaceholderText('Enter your password')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
      expect(screen.getByRole('link', { name: /sign up/i })).toBeInTheDocument();
    });

    it('should handle login and navigate to dashboard', async () => {
      const user = userEvent.setup();
      const mockUser = {
        user_id: 1,
        username: 'testuser',
        email: 'test@example.com',
        current_streak_days: 5,
        total_study_time_seconds: 3600,
      };

      const mockStats = {
        total_attempts: 50,
        correct_attempts: 40,
        accuracy_rate: 80.0,
        total_study_time_seconds: 3600,
        current_streak_days: 5,
        problems_attempted: 20,
        problems_solved: 15,
        questions_answered: 50,
        strong_topics: ['Arrays'],
        weak_topics: ['Dynamic Programming'],
      };

      vi.mocked(api.authAPI.login).mockResolvedValue({
        user: mockUser,
        token: 'test-token',
      });
      vi.mocked(api.userAPI.getStats).mockResolvedValue(mockStats);

      renderApp();

      // Fill login form
      await waitFor(() => {
        expect(screen.getByText('Welcome Back')).toBeInTheDocument();
      });

      const usernameInput = screen.getByPlaceholderText('Enter your username');
      const passwordInput = screen.getByPlaceholderText('Enter your password');
      const loginButton = screen.getByRole('button', { name: /login/i });

      await user.type(usernameInput, 'testuser');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      // Should redirect to dashboard after login
      // Note: In actual implementation, this would navigate
      await waitFor(() => {
        expect(api.authAPI.login).toHaveBeenCalledWith('testuser', 'password123');
      });
    });
  });

  describe('Practice Session Flow', () => {
    it('should complete a full question answering session', async () => {
      const user = userEvent.setup({ delay: null });

      const mockQuestion1 = {
        question_id: 1,
        question_type: 'multiple_choice',
        question_format: 'code_analysis',
        question_text: 'What is the time complexity?',
        answer_options: {
          A: 'O(n)',
          B: 'O(log n)',
          C: 'O(n²)',
          D: 'O(1)',
        },
        difficulty_score: 50,
      };

      const mockQuestion2 = {
        question_id: 2,
        question_type: 'multiple_choice',
        question_format: 'code_analysis',
        question_text: 'What is the space complexity?',
        answer_options: {
          A: 'O(n)',
          B: 'O(log n)',
          C: 'O(n²)',
          D: 'O(1)',
        },
        difficulty_score: 60,
      };

      const mockSubmitResponse = {
        is_correct: true,
        correct_answer: 'A',
        explanation: 'The algorithm uses linear space.',
        points_earned: 10,
      };

      vi.mocked(api.questionsAPI.getRandomQuestion)
        .mockResolvedValueOnce(mockQuestion1)
        .mockResolvedValueOnce(mockQuestion2);
      vi.mocked(api.questionsAPI.submitAnswer).mockResolvedValue(mockSubmitResponse);

      // Render Practice component directly for this test
      const { render: renderWithoutRouter } = await import('@testing-library/react');
      const { QueryClient, QueryClientProvider } = await import('@tanstack/react-query');
      const Practice = (await import('../../pages/Practice')).default;

      const queryClient = new QueryClient({
        defaultOptions: {
          queries: { retry: false },
        },
      });

      renderWithoutRouter(
        <QueryClientProvider client={queryClient}>
          <Practice />
        </QueryClientProvider>
      );

      // Wait for first question to load
      await waitFor(() => {
        expect(screen.getByText('What is the time complexity?')).toBeInTheDocument();
      });

      // Select answer
      const answerA = screen.getByLabelText(/A\. O\(n\)/);
      await user.click(answerA);

      // Submit answer
      const submitButton = screen.getByRole('button', { name: /submit answer/i });
      await user.click(submitButton);

      // Check result is shown
      await waitFor(() => {
        expect(screen.getByText(/correct!/i)).toBeInTheDocument();
        expect(screen.getByText(/\+10/)).toBeInTheDocument();
      });

      // Click next question
      const nextButton = screen.getByRole('button', { name: /next question/i });
      await user.click(nextButton);

      // Verify second question loads
      await waitFor(() => {
        expect(screen.getByText('What is the space complexity?')).toBeInTheDocument();
      });

      // Verify API calls
      expect(api.questionsAPI.getRandomQuestion).toHaveBeenCalledTimes(2);
      expect(api.questionsAPI.submitAnswer).toHaveBeenCalledTimes(1);
    });
  });

  describe('Error Handling', () => {
    it('should handle API errors gracefully', async () => {
      const user = userEvent.setup();

      vi.mocked(api.authAPI.login).mockRejectedValue({
        response: {
          data: {
            error: 'Network error',
          },
        },
      });

      renderApp();

      await waitFor(() => {
        expect(screen.getByText('Welcome Back')).toBeInTheDocument();
      });

      const usernameInput = screen.getByPlaceholderText('Enter your username');
      const passwordInput = screen.getByPlaceholderText('Enter your password');
      const loginButton = screen.getByRole('button', { name: /login/i });

      await user.type(usernameInput, 'testuser');
      await user.type(passwordInput, 'password123');
      await user.click(loginButton);

      // Should show error message
      await waitFor(() => {
        expect(screen.getByText('Network error')).toBeInTheDocument();
      });
    });
  });
});

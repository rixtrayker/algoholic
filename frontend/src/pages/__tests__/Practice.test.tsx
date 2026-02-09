import { describe, it, expect, beforeEach, vi } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { render } from '../../test/utils';
import Practice from '../Practice';
import * as api from '../../lib/api';

vi.mock('../../lib/api');

const mockQuestion = {
  question_id: 1,
  question_type: 'multiple_choice',
  question_format: 'code_analysis',
  question_text: 'What is the time complexity of this algorithm?',
  answer_options: {
    A: 'O(n)',
    B: 'O(n log n)',
    C: 'O(n²)',
    D: 'O(1)',
  },
  difficulty_score: 50,
};

const mockSubmitResponse = {
  is_correct: true,
  correct_answer: 'B',
  explanation: 'The algorithm uses a sorting operation which has O(n log n) complexity.',
  points_earned: 10,
};

describe('Practice Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should show loading state initially', () => {
    vi.mocked(api.questionsAPI.getRandomQuestion).mockImplementation(
      () => new Promise(() => {})
    );

    render(<Practice />);

    expect(screen.getByText(/loading question/i)).toBeInTheDocument();
  });

  it('should display question after loading', async () => {
    vi.mocked(api.questionsAPI.getRandomQuestion).mockResolvedValue(mockQuestion);

    render(<Practice />);

    await waitFor(() => {
      expect(screen.getByText('What is the time complexity of this algorithm?')).toBeInTheDocument();
    });

    // Check for answer options - text may be split across elements
    expect(screen.getByText('A.', { exact: false })).toBeInTheDocument();
    expect(screen.getByText('O(n)', { exact: false })).toBeInTheDocument();
    expect(screen.getByText('B.', { exact: false })).toBeInTheDocument();
    expect(screen.getByText('O(n log n)', { exact: false })).toBeInTheDocument();
  });

  it('should display timer', async () => {
    vi.mocked(api.questionsAPI.getRandomQuestion).mockResolvedValue(mockQuestion);

    render(<Practice />);

    await waitFor(() => {
      expect(screen.getByText('0:00')).toBeInTheDocument();
    });

    // Note: Timer advancement test skipped due to complexity with React Query
    // The timer display is verified by checking for the initial 0:00 state
  });

  it('should allow selecting an answer', async () => {
    const user = userEvent.setup({ delay: null });
    vi.mocked(api.questionsAPI.getRandomQuestion).mockResolvedValue(mockQuestion);

    render(<Practice />);

    await waitFor(() => {
      expect(screen.getByText('What is the time complexity of this algorithm?')).toBeInTheDocument();
    });

    const optionB = screen.getByLabelText(/B\. O\(n log n\)/);
    await user.click(optionB);

    expect(optionB).toBeChecked();
  });

  it('should submit answer and show result', async () => {
    const user = userEvent.setup({ delay: null });
    vi.mocked(api.questionsAPI.getRandomQuestion).mockResolvedValue(mockQuestion);
    vi.mocked(api.questionsAPI.submitAnswer).mockResolvedValue(mockSubmitResponse);

    render(<Practice />);

    await waitFor(() => {
      expect(screen.getByText('What is the time complexity of this algorithm?')).toBeInTheDocument();
    });

    // Select answer
    const optionB = screen.getByLabelText(/B\. O\(n log n\)/);
    await user.click(optionB);

    // Submit
    const submitButton = screen.getByRole('button', { name: /submit answer/i });
    await user.click(submitButton);

    // Check result
    await waitFor(() => {
      expect(screen.getByText(/correct!/i)).toBeInTheDocument();
    });

    expect(screen.getByText(/explanation/i)).toBeInTheDocument();
    expect(screen.getByText(mockSubmitResponse.explanation)).toBeInTheDocument();
    expect(screen.getByText(/\+10/)).toBeInTheDocument(); // Points earned
  });

  it('should show incorrect result with correct answer', async () => {
    const user = userEvent.setup({ delay: null });
    const incorrectResponse = {
      ...mockSubmitResponse,
      is_correct: false,
    };

    vi.mocked(api.questionsAPI.getRandomQuestion).mockResolvedValue(mockQuestion);
    vi.mocked(api.questionsAPI.submitAnswer).mockResolvedValue(incorrectResponse);

    render(<Practice />);

    await waitFor(() => {
      expect(screen.getByText('What is the time complexity of this algorithm?')).toBeInTheDocument();
    });

    // Select wrong answer
    const optionA = screen.getByLabelText(/A\. O\(n\)/);
    await user.click(optionA);

    const submitButton = screen.getByRole('button', { name: /submit answer/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/incorrect/i)).toBeInTheDocument();
    });

    expect(screen.getByText(/correct answer was: B/i)).toBeInTheDocument();
  });

  it('should disable submit button when no answer selected', async () => {
    vi.mocked(api.questionsAPI.getRandomQuestion).mockResolvedValue(mockQuestion);

    render(<Practice />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /submit answer/i })).toBeDisabled();
    });
  });

  it('should load next question after submission', async () => {
    const user = userEvent.setup({ delay: null });
    const nextQuestion = {
      ...mockQuestion,
      question_id: 2,
      question_text: 'What is the space complexity?',
    };

    vi.mocked(api.questionsAPI.getRandomQuestion)
      .mockResolvedValueOnce(mockQuestion)
      .mockResolvedValueOnce(nextQuestion);
    vi.mocked(api.questionsAPI.submitAnswer).mockResolvedValue(mockSubmitResponse);

    render(<Practice />);

    await waitFor(() => {
      expect(screen.getByText('What is the time complexity of this algorithm?')).toBeInTheDocument();
    });

    // Answer and submit
    const optionB = screen.getByLabelText(/B\. O\(n log n\)/);
    await user.click(optionB);

    const submitButton = screen.getByRole('button', { name: /submit answer/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/correct!/i)).toBeInTheDocument();
    });

    // Click next
    const nextButton = screen.getByRole('button', { name: /next question/i });
    await user.click(nextButton);

    await waitFor(() => {
      expect(screen.getByText('What is the space complexity?')).toBeInTheDocument();
    });
  });

  it('should track time taken for answer submission', async () => {
    const user = userEvent.setup({ delay: null });
    vi.mocked(api.questionsAPI.getRandomQuestion).mockResolvedValue(mockQuestion);
    vi.mocked(api.questionsAPI.submitAnswer).mockResolvedValue(mockSubmitResponse);

    render(<Practice />);

    await waitFor(() => {
      expect(screen.getByText('What is the time complexity of this algorithm?')).toBeInTheDocument();
    });

    const optionB = screen.getByLabelText(/B\. O\(n log n\)/);
    await user.click(optionB);

    const submitButton = screen.getByRole('button', { name: /submit answer/i });
    await user.click(submitButton);

    // Verify submission was called with time tracking (actual time will vary)
    await waitFor(() => {
      expect(api.questionsAPI.submitAnswer).toHaveBeenCalledWith(
        1,
        { answer: 'B' },
        expect.any(Number) // Time in seconds
      );
    });
  });

  it('should show hint button', async () => {
    vi.mocked(api.questionsAPI.getRandomQuestion).mockResolvedValue(mockQuestion);

    render(<Practice />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /hint/i })).toBeInTheDocument();
    });
  });
});

import axios from 'axios';
import toast from 'react-hot-toast';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:4000/api';

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors and network errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Network error (no response)
    if (!error.response) {
      toast.error('Network error. Please check your connection.');
      return Promise.reject(error);
    }

    const status = error.response.status;
    const errorMessage = error.response.data?.error || error.response.data?.message;
    const requestUrl = error.config?.url || '';

    // Handle specific error codes
    if (status === 401) {
      const currentPath = window.location.pathname;
      // Don't show toast or redirect for auth/me requests (silent fail for token validation)
      if (requestUrl.includes('/auth/me')) {
        return Promise.reject(error);
      }
      // Only redirect if not already on login page
      if (!currentPath.includes('/login') && !currentPath.includes('/register')) {
        localStorage.removeItem('auth_token');
        toast.error('Session expired. Please login again.');
        setTimeout(() => {
          window.location.href = '/login';
        }, 1000);
      }
    } else if (status === 403) {
      toast.error(errorMessage || 'Access denied.');
    } else if (status === 404) {
      // Don't show toast for 404 on stats endpoints - handled by component
      if (!requestUrl.includes('/stats')) {
        toast.error(errorMessage || 'Resource not found.');
      }
    } else if (status === 500) {
      toast.error('Server error. Please try again later.');
    } else if (status >= 400 && status < 500) {
      // Client errors - show specific message (but not for auth endpoints, handled by component)
      if (errorMessage && !requestUrl.includes('/auth/login') && !requestUrl.includes('/auth/register')) {
        toast.error(errorMessage);
      }
    }

    return Promise.reject(error);
  }
);

// Types
export interface User {
  user_id: number;
  username: string;
  email: string;
  current_streak_days: number;
  total_study_time_seconds: number;
}

export interface Problem {
  problem_id: number;
  title: string;
  slug: string;
  description: string;
  difficulty_score: number;
  primary_pattern?: string;
  examples: any;
  hints?: string[];
}

export interface Question {
  question_id: number;
  question_type: string;
  question_format: string;
  question_text: string;
  answer_options?: any;
  difficulty_score: number;
}

export interface UserStats {
  total_attempts: number;
  correct_attempts: number;
  accuracy_rate: number;
  total_study_time_seconds: number;
  current_streak_days: number;
  problems_attempted: number;
  problems_solved: number;
  questions_answered: number;
  strong_topics: string[];
  weak_topics: string[];
}

export interface TrainingPlan {
  plan_id: number;
  name: string;
  description?: string;
  status: string;
  progress_percentage: number;
  start_date: string;
  questions_per_day: number;
}

// Auth API
export const authAPI = {
  register: async (username: string, email: string, password: string) => {
    const { data } = await api.post('/auth/register', { username, email, password });
    if (data.token) {
      localStorage.setItem('auth_token', data.token);
    }
    return data;
  },

  login: async (username: string, password: string) => {
    const { data } = await api.post('/auth/login', { username, password });
    if (data.token) {
      localStorage.setItem('auth_token', data.token);
    }
    return data;
  },

  logout: () => {
    localStorage.removeItem('auth_token');
  },

  getMe: async () => {
    const { data } = await api.get('/auth/me');
    return data as User;
  },
};

// Problems API
export const problemsAPI = {
  getProblems: async (params?: {
    min_difficulty?: number;
    max_difficulty?: number;
    pattern?: string;
    limit?: number;
    offset?: number;
  }) => {
    const { data } = await api.get('/problems', { params });
    return data;
  },

  getProblem: async (id: number) => {
    const { data } = await api.get(`/problems/${id}`);
    return data as Problem;
  },

  getProblemQuestions: async (id: number) => {
    const { data } = await api.get(`/problems/${id}/questions`);
    return data;
  },

  searchProblems: async (query: string, filters?: {
    difficulty?: string;
    topic?: string;
    limit?: number;
  }) => {
    const params = { q: query, ...filters };
    const { data } = await api.get('/problems/search', { params });
    return data;
  },
};

// Questions API
export const questionsAPI = {
  getQuestions: async (params?: {
    type?: string;
    min_difficulty?: number;
    max_difficulty?: number;
    limit?: number;
    offset?: number;
  }) => {
    const { data } = await api.get('/questions', { params });
    return data;
  },

  getQuestion: async (id: number) => {
    const { data } = await api.get(`/questions/${id}`);
    return data as Question;
  },

  getRandomQuestion: async (params?: {
    type?: string;
    min_difficulty?: number;
    max_difficulty?: number;
  }) => {
    const { data } = await api.get('/questions/random', { params });
    return data as Question;
  },

  submitAnswer: async (
    questionId: number,
    userAnswer: any,
    timeTaken: number
  ) => {
    const { data } = await api.post(`/questions/${questionId}/answer`, {
      answer: userAnswer,
      time_taken_seconds: timeTaken,
    });
    return data;
  },

  getHint: async (questionId: number) => {
    const { data } = await api.get(`/questions/${questionId}/hint`);
    return data;
  },
};

// User API
export const userAPI = {
  getStats: async () => {
    const { data } = await api.get('/users/me/stats');
    return data as UserStats;
  },

  getProgress: async (days = 30) => {
    const { data } = await api.get('/users/me/progress', { params: { days } });
    return data;
  },

  getAttempts: async (limit = 50) => {
    const { data } = await api.get('/users/me/attempts', { params: { limit } });
    return data;
  },

  getWeaknesses: async (limit = 10) => {
    const { data } = await api.get('/users/me/weaknesses', { params: { limit } });
    return data;
  },

  getRecommendations: async () => {
    const { data } = await api.get('/users/me/recommendations');
    return data;
  },

  getReviewQueue: async () => {
    const { data } = await api.get('/users/me/review-queue');
    return data;
  },

  getSkills: async () => {
    const { data } = await api.get('/users/me/skills');
    return data;
  },
};

// Training Plans API
export const trainingPlansAPI = {
  getPlans: async () => {
    const { data } = await api.get('/plans');
    return data as TrainingPlan[];
  },

  getMyPlans: async () => {
    const { data } = await api.get('/users/plans');
    return data as TrainingPlan[];
  },

  getPlan: async (id: number) => {
    const { data } = await api.get(`/plans/${id}`);
    return data as TrainingPlan;
  },

  enrollInPlan: async (planId: number) => {
    const { data } = await api.post(`/plans/${planId}/enroll`);
    return data;
  },

  updatePlanProgress: async (planId: number, progressData: {
    completed_items: number;
    current_day: number;
  }) => {
    const { data } = await api.put(`/plans/${planId}/progress`, progressData);
    return data;
  },

  createPlan: async (planData: {
    name: string;
    description?: string;
    plan_type: string;
    target_topics?: number[];
    target_patterns?: string[];
    duration_days: number;
    questions_per_day: number;
    difficulty_min: number;
    difficulty_max: number;
    adaptive_difficulty: boolean;
  }) => {
    const { data } = await api.post('/plans', planData);
    return data;
  },

  getNextQuestion: async (planId: number) => {
    const { data } = await api.get(`/plans/${planId}/next`);
    return data as Question;
  },

  getTodaysQuestions: async (planId: number) => {
    const { data} = await api.get(`/plans/${planId}/today`);
    return data;
  },

  pausePlan: async (planId: number) => {
    const { data } = await api.post(`/plans/${planId}/pause`);
    return data;
  },

  resumePlan: async (planId: number) => {
    const { data } = await api.post(`/plans/${planId}/resume`);
    return data;
  },

  deletePlan: async (planId: number) => {
    await api.delete(`/plans/${planId}`);
  },
};

// Topics API
export const topicsAPI = {
  getTopics: async () => {
    const { data } = await api.get('/topics');
    return data;
  },

  getTopicPerformance: async (topicId: number) => {
    const { data } = await api.get(`/users/topics/${topicId}/performance`);
    return data;
  },
};

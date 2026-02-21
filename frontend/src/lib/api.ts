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
      // Don't show toast or redirect for auth endpoints (handled by components)
      if (requestUrl.includes('/auth/')) {
        return Promise.reject(error);
      }
      // For other 401s (expired sessions during normal use), just clear token
      // The PrivateRoute will handle the redirect
      localStorage.removeItem('auth_token');
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
    timeTaken: number,
    hintsUsed?: number,
    confidenceLevel?: number,
    trainingPlanId?: number
  ) => {
    const { data } = await api.post(`/questions/${questionId}/answer`, {
      user_answer: userAnswer,
      time_taken_seconds: timeTaken,
      hints_used: hintsUsed || 0,
      confidence_level: confidenceLevel,
      training_plan_id: trainingPlanId,
    });
    return data;
  },

  getHint: async (questionId: number) => {
    const { data } = await api.get(`/questions/${questionId}/hint`);
    return data;
  },

  getAttempts: async (questionId: number) => {
    const { data } = await api.get(`/questions/${questionId}/attempts`);
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
    const { data } = await api.get('/training-plans');
    return data as { plans: TrainingPlan[]; count: number };
  },

  getPlan: async (id: number) => {
    const { data } = await api.get(`/training-plans/${id}`);
    return data as TrainingPlan;
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
    const { data } = await api.post('/training-plans', planData);
    return data;
  },

  getNextQuestion: async (planId: number) => {
    const { data } = await api.get(`/training-plans/${planId}/next`);
    return data as Question;
  },

  getPlanItems: async (planId: number) => {
    const { data } = await api.get(`/training-plans/${planId}/items`);
    return data;
  },

  getTodaysQuestions: async (planId: number) => {
    const { data } = await api.get(`/training-plans/${planId}/today`);
    return data;
  },

  completeItem: async (planId: number, itemId: number) => {
    const { data } = await api.post(`/training-plans/${planId}/items/${itemId}/complete`);
    return data;
  },

  pausePlan: async (planId: number) => {
    const { data } = await api.post(`/training-plans/${planId}/pause`);
    return data;
  },

  resumePlan: async (planId: number) => {
    const { data } = await api.post(`/training-plans/${planId}/resume`);
    return data;
  },

  deletePlan: async (planId: number) => {
    await api.delete(`/training-plans/${planId}`);
  },
};

// Topics API
export const topicsAPI = {
  getTopics: async () => {
    const { data } = await api.get('/topics');
    return data;
  },

  getTopic: async (id: number) => {
    const { data } = await api.get(`/topics/${id}`);
    return data;
  },

  getTopicPrerequisites: async (id: number) => {
    const { data } = await api.get(`/topics/${id}/prerequisites`);
    return data;
  },

  getTopicPerformance: async (userId: number, topicId: number) => {
    const { data } = await api.get(`/topics/${userId}/performance/${topicId}`);
    return data;
  },
};

// Lists API
export const listsAPI = {
  getLists: async () => {
    const { data } = await api.get('/lists');
    return data;
  },

  getList: async (id: number) => {
    const { data } = await api.get(`/lists/${id}`);
    return data;
  },

  createList: async (listData: {
    name: string;
    description?: string;
    is_public: boolean;
  }) => {
    const { data } = await api.post('/lists', listData);
    return data;
  },

  updateList: async (id: number, listData: {
    name?: string;
    description?: string;
    is_public?: boolean;
  }) => {
    const { data } = await api.put(`/lists/${id}`, listData);
    return data;
  },

  deleteList: async (id: number) => {
    await api.delete(`/lists/${id}`);
  },

  addProblem: async (listId: number, problemId: number) => {
    const { data } = await api.post(`/lists/${listId}/problems`, { problem_id: problemId });
    return data;
  },

  removeProblem: async (listId: number, problemId: number) => {
    await api.delete(`/lists/${listId}/problems/${problemId}`);
  },

  getListProblems: async (listId: number) => {
    const { data } = await api.get(`/lists/${listId}/problems`);
    return data;
  },
};

// Activity API
export const activityAPI = {
  getActivityChart: async (days = 365) => {
    const { data } = await api.get('/activity/chart', { params: { days } });
    return data;
  },

  getActivityStats: async () => {
    const { data } = await api.get('/activity/stats');
    return data;
  },

  getPracticeHistory: async (days = 30) => {
    const { data } = await api.get('/activity/history', { params: { days } });
    return data;
  },

  recordActivity: async (activityData: {
    problems_count: number;
    questions_count: number;
    study_time_seconds: number;
  }) => {
    await api.post('/activity/record', activityData);
  },
};

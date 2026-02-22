import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:4000/api';

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (typeof window !== 'undefined' && error.response?.status === 401) {
      const requestUrl = error.config?.url || '';
      if (!requestUrl.includes('/auth/login') && !requestUrl.includes('/auth/register')) {
        localStorage.removeItem('auth_token');
        if (!window.location.pathname.includes('/login') && !window.location.pathname.includes('/register')) {
          window.location.href = '/login';
        }
      }
    }
    return Promise.reject(error);
  }
);

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
  examples: Record<string, unknown>;
  hints?: string[];
}

export interface Question {
  question_id: number;
  question_type: string;
  question_format: string;
  question_text: string;
  answer_options?: Record<string, unknown>;
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
  average_difficulty?: number;
}

export interface TrainingPlan {
  plan_id: number;
  user_id: number;
  name: string;
  description?: string;
  status: string;
  progress_percentage: number;
  start_date: string;
  questions_per_day: number;
  plan_type?: string;
  target_topics?: number[];
  target_patterns?: string[];
  duration_days?: number;
}

export interface Topic {
  topic_id: number;
  name: string;
  slug: string;
  description?: string;
  difficulty_level?: number;
  parent_topic_id?: number;
}

export interface UserList {
  list_id: number;
  name: string;
  description?: string;
  is_public: boolean;
  total_items: number;
  completed: number;
  created_at: string;
}

export interface WeakTopic {
  topic_id: number;
  name: string;
  proficiency_level: number;
}

export interface Recommendation {
  type: string;
  topic?: Topic;
  reason: string;
  priority: string;
  action: string;
}

export interface ReviewQueueItem {
  user_id: number;
  topic_id: number;
  proficiency_level: number;
  next_review_at: string;
  topic?: Topic;
}

export const authAPI = {
  register: async (username: string, email: string, password: string) => {
    const { data } = await api.post('/auth/register', { username, email, password });
    if (data.token && typeof window !== 'undefined') {
      localStorage.setItem('auth_token', data.token);
    }
    return data;
  },

  login: async (username: string, password: string) => {
    const { data } = await api.post('/auth/login', { username, password });
    if (data.token && typeof window !== 'undefined') {
      localStorage.setItem('auth_token', data.token);
    }
    return data;
  },

  logout: () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
    }
  },

  getMe: async () => {
    const { data } = await api.get('/auth/me');
    return data as User;
  },

  changePassword: async (oldPassword: string, newPassword: string) => {
    const { data } = await api.post('/auth/change-password', {
      old_password: oldPassword,
      new_password: newPassword,
    });
    return data;
  },
};

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

  getProblemBySlug: async (slug: string) => {
    const { data } = await api.get(`/problems/slug/${slug}`);
    return data as Problem;
  },

  getProblemTopics: async (id: number) => {
    const { data } = await api.get(`/problems/${id}/topics`);
    return data;
  },

  getSimilarProblems: async (id: number, limit = 5) => {
    const { data } = await api.get(`/problems/${id}/similar`, { params: { limit } });
    return data;
  },

  searchProblems: async (query: string, limit = 20) => {
    const { data } = await api.get('/search/problems', { params: { q: query, limit } });
    return data;
  },
};

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
    userAnswer: Record<string, unknown>,
    timeTaken: number,
    hintsUsed = 0,
    confidenceLevel?: number,
    trainingPlanId?: number
  ) => {
    const { data } = await api.post(`/questions/${questionId}/answer`, {
      user_answer: userAnswer,
      time_taken_seconds: timeTaken,
      hints_used: hintsUsed,
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

  getQuestionsByProblem: async (problemId: number) => {
    const { data } = await api.get(`/problems/${problemId}/questions`);
    return data;
  },
};

export const userAPI = {
  getStats: async () => {
    const { data } = await api.get('/users/me/stats');
    return data as UserStats;
  },

  getWeaknesses: async (limit = 10) => {
    const { data } = await api.get('/users/me/weaknesses', { params: { limit } });
    return data as { weak_topics: WeakTopic[]; count: number };
  },

  getRecommendations: async () => {
    const { data } = await api.get('/users/me/recommendations');
    return data as { recommendations: Recommendation[]; count: number };
  },

  getReviewQueue: async () => {
    const { data } = await api.get('/users/me/review-queue');
    return data as { review_queue: ReviewQueueItem[]; count: number };
  },

  getSkills: async () => {
    const { data } = await api.get('/users/me/skills');
    return data;
  },

  getSkillByTopic: async (topicId: number) => {
    const { data } = await api.get(`/users/me/skills/${topicId}`);
    return data;
  },

  getPreferences: async () => {
    const { data } = await api.get('/users/me/preferences');
    return data;
  },

  updatePreferences: async (preferences: Record<string, unknown>) => {
    const { data } = await api.put('/users/me/preferences', preferences);
    return data;
  },

  getRecentAttempts: async (limit = 20) => {
    const { data } = await api.get('/users/me/attempts', { params: { limit } });
    return data;
  },
};

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

export const topicsAPI = {
  getTopics: async () => {
    const { data } = await api.get('/topics');
    return data as { topics: Topic[] };
  },

  getTopic: async (id: number) => {
    const { data } = await api.get(`/topics/${id}`);
    return data as Topic;
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

export const listsAPI = {
  getLists: async () => {
    const { data } = await api.get('/lists');
    return data as UserList[];
  },

  getList: async (id: number) => {
    const { data } = await api.get(`/lists/${id}`);
    return data as UserList;
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

export const searchAPI = {
  searchProblems: async (query: string, limit = 10) => {
    const { data } = await api.get('/search/problems', { params: { q: query, limit } });
    return data;
  },

  searchQuestions: async (query: string, limit = 10) => {
    const { data } = await api.get('/search/questions', { params: { q: query, limit } });
    return data;
  },
};

export const graphAPI = {
  getLearningPath: async (from: number, to: number) => {
    const { data } = await api.get('/graph/learning-path', { params: { from, to } });
    return data;
  },
};

export const intelligenceAPI = {
  getStatus: async () => {
    const { data } = await api.get('/intelligence/status');
    return data;
  },
};

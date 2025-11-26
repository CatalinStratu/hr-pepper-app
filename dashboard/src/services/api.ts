import axios from 'axios';
import type { Question, Interview, Analytics, InterviewDetail } from '../types';

const API_BASE = import.meta.env.VITE_API_URL || '/api';

const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Questions API
export const questionsApi = {
  getAll: async (): Promise<Question[]> => {
    const response = await api.get('/questions');
    return response.data.questions || [];
  },

  getById: async (id: string): Promise<Question> => {
    const response = await api.get(`/questions/${id}`);
    return response.data;
  },

  create: async (question: Partial<Question>): Promise<Question> => {
    const response = await api.post('/questions', question);
    return response.data;
  },

  update: async (id: string, question: Partial<Question>): Promise<void> => {
    await api.put(`/questions/${id}`, question);
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/questions/${id}`);
  },
};

// Interviews API
export const interviewsApi = {
  getAll: async (status?: string): Promise<Interview[]> => {
    const params = status ? { status } : {};
    const response = await api.get('/interviews', { params });
    return response.data.interviews || [];
  },

  getById: async (id: string): Promise<InterviewDetail> => {
    const response = await api.get(`/interviews/${id}`);
    return response.data;
  },

  create: async (data: { candidate_name: string; candidate_email?: string; position: string }): Promise<Interview> => {
    const response = await api.post('/interviews', data);
    return response.data.interview;
  },

  update: async (id: string, data: { status?: string; notes?: string; score?: number }): Promise<void> => {
    await api.put(`/interviews/${id}`, data);
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/interviews/${id}`);
  },

  getReport: async (id: string): Promise<any> => {
    const response = await api.get(`/interviews/${id}/report`);
    return response.data;
  },
};

// Analytics API
export const analyticsApi = {
  get: async (): Promise<Analytics> => {
    const response = await api.get('/analytics');
    return response.data;
  },
};

export default api;

export interface Question {
  id: string;
  text: string;
  category: string;
  difficulty: string;
  time_limit: number;
  order: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Interview {
  id: string;
  candidate_name: string;
  candidate_email: string;
  position: string;
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled';
  started_at: string;
  completed_at?: string;
  total_score?: number;
  notes: string;
  created_at: string;
}

export interface Response {
  id: string;
  interview_id: string;
  question_id: string;
  question_text?: string;
  response_text: string;
  duration: number;
  sentiment_score?: number;
  confidence_score?: number;
  keywords_found?: string[];
  analysis?: string;
  created_at: string;
}

export interface Analytics {
  total_interviews: number;
  completed_interviews: number;
  average_score: number;
  average_duration: number;
  question_stats: QuestionStat[];
  recent_interviews: InterviewSummary[];
  score_distribution: Record<string, number>;
}

export interface QuestionStat {
  question_id: string;
  question_text: string;
  times_asked: number;
  average_score: number;
  average_duration: number;
}

export interface InterviewSummary {
  id: string;
  candidate_name: string;
  position: string;
  status: string;
  total_score?: number;
  completed_at?: string;
}

export interface InterviewDetail {
  interview: Interview;
  responses: Response[];
}

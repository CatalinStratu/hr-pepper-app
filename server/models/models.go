package models

import "time"

// Question represents an interview question
type Question struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	Category    string    `json:"category"`
	Difficulty  string    `json:"difficulty"` // easy, medium, hard
	TimeLimit   int       `json:"time_limit"` // seconds
	Order       int       `json:"order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Interview represents an interview session
type Interview struct {
	ID            string    `json:"id"`
	CandidateName string    `json:"candidate_name"`
	CandidateEmail string   `json:"candidate_email"`
	Position      string    `json:"position"`
	Status        string    `json:"status"` // pending, in_progress, completed, cancelled
	StartedAt     time.Time `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	TotalScore    *float64  `json:"total_score,omitempty"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

// Response represents a candidate's response to a question
type Response struct {
	ID              string    `json:"id"`
	InterviewID     string    `json:"interview_id"`
	QuestionID      string    `json:"question_id"`
	QuestionText    string    `json:"question_text,omitempty"`
	ResponseText    string    `json:"response_text"`
	ResponseAudio   string    `json:"response_audio,omitempty"` // base64 encoded audio
	Duration        int       `json:"duration"` // seconds
	SentimentScore  *float64  `json:"sentiment_score,omitempty"`
	ConfidenceScore *float64  `json:"confidence_score,omitempty"`
	KeywordsFound   []string  `json:"keywords_found,omitempty"`
	Analysis        string    `json:"analysis,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// Analytics represents interview analytics
type Analytics struct {
	TotalInterviews     int                    `json:"total_interviews"`
	CompletedInterviews int                    `json:"completed_interviews"`
	AverageScore        float64                `json:"average_score"`
	AverageDuration     float64                `json:"average_duration"`
	QuestionStats       []QuestionStat         `json:"question_stats"`
	RecentInterviews    []InterviewSummary     `json:"recent_interviews"`
	ScoreDistribution   map[string]int         `json:"score_distribution"`
}

// QuestionStat represents statistics for a question
type QuestionStat struct {
	QuestionID       string  `json:"question_id"`
	QuestionText     string  `json:"question_text"`
	TimesAsked       int     `json:"times_asked"`
	AverageScore     float64 `json:"average_score"`
	AverageDuration  float64 `json:"average_duration"`
}

// InterviewSummary represents a brief interview summary
type InterviewSummary struct {
	ID            string    `json:"id"`
	CandidateName string    `json:"candidate_name"`
	Position      string    `json:"position"`
	Status        string    `json:"status"`
	TotalScore    *float64  `json:"total_score,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}

// CreateQuestionRequest represents request to create a question
type CreateQuestionRequest struct {
	Text       string `json:"text" binding:"required"`
	Category   string `json:"category" binding:"required"`
	Difficulty string `json:"difficulty"`
	TimeLimit  int    `json:"time_limit"`
	Order      int    `json:"order"`
}

// CreateInterviewRequest represents request to start an interview
type CreateInterviewRequest struct {
	CandidateName  string `json:"candidate_name" binding:"required"`
	CandidateEmail string `json:"candidate_email"`
	Position       string `json:"position" binding:"required"`
}

// SubmitResponseRequest represents request to submit a response
type SubmitResponseRequest struct {
	InterviewID   string `json:"interview_id" binding:"required"`
	QuestionID    string `json:"question_id" binding:"required"`
	ResponseText  string `json:"response_text" binding:"required"`
	ResponseAudio string `json:"response_audio,omitempty"`
	Duration      int    `json:"duration"`
}

// UpdateInterviewRequest represents request to update interview status
type UpdateInterviewRequest struct {
	Status string   `json:"status"`
	Notes  string   `json:"notes"`
	Score  *float64 `json:"score,omitempty"`
}

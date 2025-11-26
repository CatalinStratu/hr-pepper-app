package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"hr-pepper-server/database"
	"hr-pepper-server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetInterviews returns all interviews
func GetInterviews(c *gin.Context) {
	status := c.Query("status")

	var rows *sql.Rows
	var err error

	if status != "" {
		rows, err = database.DB.Query(`
			SELECT id, candidate_name, candidate_email, position, status, started_at, completed_at, total_score, notes, created_at
			FROM interviews
			WHERE status = ?
			ORDER BY created_at DESC
		`, status)
	} else {
		rows, err = database.DB.Query(`
			SELECT id, candidate_name, candidate_email, position, status, started_at, completed_at, total_score, notes, created_at
			FROM interviews
			ORDER BY created_at DESC
		`)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interviews"})
		return
	}
	defer rows.Close()

	var interviews []models.Interview
	for rows.Next() {
		var i models.Interview
		var startedAt, completedAt sql.NullTime
		var totalScore sql.NullFloat64
		var candidateEmail, notes sql.NullString

		if err := rows.Scan(&i.ID, &i.CandidateName, &candidateEmail, &i.Position, &i.Status, &startedAt, &completedAt, &totalScore, &notes, &i.CreatedAt); err != nil {
			continue
		}

		if startedAt.Valid {
			i.StartedAt = startedAt.Time
		}
		if completedAt.Valid {
			i.CompletedAt = &completedAt.Time
		}
		if totalScore.Valid {
			i.TotalScore = &totalScore.Float64
		}
		if candidateEmail.Valid {
			i.CandidateEmail = candidateEmail.String
		}
		if notes.Valid {
			i.Notes = notes.String
		}

		interviews = append(interviews, i)
	}

	c.JSON(http.StatusOK, gin.H{"interviews": interviews})
}

// GetInterview returns a single interview with its responses
func GetInterview(c *gin.Context) {
	id := c.Param("id")

	var i models.Interview
	var startedAt, completedAt sql.NullTime
	var totalScore sql.NullFloat64
	var candidateEmail, notes sql.NullString

	err := database.DB.QueryRow(`
		SELECT id, candidate_name, candidate_email, position, status, started_at, completed_at, total_score, notes, created_at
		FROM interviews WHERE id = ?
	`, id).Scan(&i.ID, &i.CandidateName, &candidateEmail, &i.Position, &i.Status, &startedAt, &completedAt, &totalScore, &notes, &i.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
		return
	}

	if startedAt.Valid {
		i.StartedAt = startedAt.Time
	}
	if completedAt.Valid {
		i.CompletedAt = &completedAt.Time
	}
	if totalScore.Valid {
		i.TotalScore = &totalScore.Float64
	}
	if candidateEmail.Valid {
		i.CandidateEmail = candidateEmail.String
	}
	if notes.Valid {
		i.Notes = notes.String
	}

	// Fetch responses for this interview
	rows, err := database.DB.Query(`
		SELECT r.id, r.interview_id, r.question_id, q.text, r.response_text, r.duration, r.sentiment_score, r.confidence_score, r.analysis, r.created_at
		FROM responses r
		LEFT JOIN questions q ON r.question_id = q.id
		WHERE r.interview_id = ?
		ORDER BY r.created_at ASC
	`, id)

	var responses []models.Response
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var r models.Response
			var questionText sql.NullString
			var sentimentScore, confidenceScore sql.NullFloat64
			var analysis sql.NullString

			if err := rows.Scan(&r.ID, &r.InterviewID, &r.QuestionID, &questionText, &r.ResponseText, &r.Duration, &sentimentScore, &confidenceScore, &analysis, &r.CreatedAt); err != nil {
				continue
			}

			if questionText.Valid {
				r.QuestionText = questionText.String
			}
			if sentimentScore.Valid {
				r.SentimentScore = &sentimentScore.Float64
			}
			if confidenceScore.Valid {
				r.ConfidenceScore = &confidenceScore.Float64
			}
			if analysis.Valid {
				r.Analysis = analysis.String
			}

			responses = append(responses, r)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"interview": i,
		"responses": responses,
	})
}

// CreateInterview starts a new interview session
func CreateInterview(c *gin.Context) {
	var req models.CreateInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	now := time.Now()

	_, err := database.DB.Exec(`
		INSERT INTO interviews (id, candidate_name, candidate_email, position, status, started_at, created_at)
		VALUES (?, ?, ?, ?, 'in_progress', ?, ?)
	`, id, req.CandidateName, req.CandidateEmail, req.Position, now, now)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interview"})
		return
	}

	// Fetch questions for the interview
	rows, err := database.DB.Query(`
		SELECT id, text, category, difficulty, time_limit, question_order
		FROM questions
		WHERE is_active = 1
		ORDER BY question_order ASC
	`)

	var questions []models.Question
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var q models.Question
			if err := rows.Scan(&q.ID, &q.Text, &q.Category, &q.Difficulty, &q.TimeLimit, &q.Order); err != nil {
				continue
			}
			questions = append(questions, q)
		}
	}

	interview := models.Interview{
		ID:             id,
		CandidateName:  req.CandidateName,
		CandidateEmail: req.CandidateEmail,
		Position:       req.Position,
		Status:         "in_progress",
		StartedAt:      now,
		CreatedAt:      now,
	}

	c.JSON(http.StatusCreated, gin.H{
		"interview": interview,
		"questions": questions,
	})
}

// UpdateInterview updates an interview status
func UpdateInterview(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var completedAt interface{}
	if req.Status == "completed" {
		now := time.Now()
		completedAt = now
	} else {
		completedAt = nil
	}

	result, err := database.DB.Exec(`
		UPDATE interviews SET status = ?, notes = ?, total_score = ?, completed_at = ?
		WHERE id = ?
	`, req.Status, req.Notes, req.Score, completedAt, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update interview"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interview updated successfully"})
}

// DeleteInterview deletes an interview
func DeleteInterview(c *gin.Context) {
	id := c.Param("id")

	// Delete responses first
	database.DB.Exec("DELETE FROM responses WHERE interview_id = ?", id)

	result, err := database.DB.Exec("DELETE FROM interviews WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete interview"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interview deleted successfully"})
}

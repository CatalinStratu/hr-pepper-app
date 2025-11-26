package handlers

import (
	"database/sql"
	"net/http"

	"hr-pepper-server/database"
	"hr-pepper-server/models"

	"github.com/gin-gonic/gin"
)

// GetAnalytics returns interview analytics and statistics
func GetAnalytics(c *gin.Context) {
	analytics := models.Analytics{
		ScoreDistribution: make(map[string]int),
	}

	// Total interviews
	database.DB.QueryRow("SELECT COUNT(*) FROM interviews").Scan(&analytics.TotalInterviews)

	// Completed interviews
	database.DB.QueryRow("SELECT COUNT(*) FROM interviews WHERE status = 'completed'").Scan(&analytics.CompletedInterviews)

	// Average score
	var avgScore sql.NullFloat64
	database.DB.QueryRow("SELECT AVG(total_score) FROM interviews WHERE total_score IS NOT NULL").Scan(&avgScore)
	if avgScore.Valid {
		analytics.AverageScore = avgScore.Float64
	}

	// Average duration (sum of all response durations per interview)
	var avgDuration sql.NullFloat64
	database.DB.QueryRow(`
		SELECT AVG(total_duration) FROM (
			SELECT interview_id, SUM(duration) as total_duration
			FROM responses
			GROUP BY interview_id
		)
	`).Scan(&avgDuration)
	if avgDuration.Valid {
		analytics.AverageDuration = avgDuration.Float64
	}

	// Score distribution
	rows, err := database.DB.Query(`
		SELECT
			CASE
				WHEN total_score >= 0.9 THEN 'excellent'
				WHEN total_score >= 0.7 THEN 'good'
				WHEN total_score >= 0.5 THEN 'average'
				WHEN total_score >= 0.3 THEN 'below_average'
				ELSE 'poor'
			END as score_range,
			COUNT(*) as count
		FROM interviews
		WHERE total_score IS NOT NULL
		GROUP BY score_range
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var scoreRange string
			var count int
			if err := rows.Scan(&scoreRange, &count); err == nil {
				analytics.ScoreDistribution[scoreRange] = count
			}
		}
	}

	// Question statistics
	qRows, err := database.DB.Query(`
		SELECT
			q.id,
			q.text,
			COUNT(r.id) as times_asked,
			COALESCE(AVG(r.sentiment_score), 0) as avg_score,
			COALESCE(AVG(r.duration), 0) as avg_duration
		FROM questions q
		LEFT JOIN responses r ON q.id = r.question_id
		WHERE q.is_active = 1
		GROUP BY q.id, q.text
		ORDER BY times_asked DESC
		LIMIT 10
	`)
	if err == nil {
		defer qRows.Close()
		for qRows.Next() {
			var stat models.QuestionStat
			if err := qRows.Scan(&stat.QuestionID, &stat.QuestionText, &stat.TimesAsked, &stat.AverageScore, &stat.AverageDuration); err == nil {
				analytics.QuestionStats = append(analytics.QuestionStats, stat)
			}
		}
	}

	// Recent interviews
	iRows, err := database.DB.Query(`
		SELECT id, candidate_name, position, status, total_score, completed_at
		FROM interviews
		ORDER BY created_at DESC
		LIMIT 10
	`)
	if err == nil {
		defer iRows.Close()
		for iRows.Next() {
			var summary models.InterviewSummary
			var totalScore sql.NullFloat64
			var completedAt sql.NullTime
			if err := iRows.Scan(&summary.ID, &summary.CandidateName, &summary.Position, &summary.Status, &totalScore, &completedAt); err == nil {
				if totalScore.Valid {
					summary.TotalScore = &totalScore.Float64
				}
				if completedAt.Valid {
					summary.CompletedAt = &completedAt.Time
				}
				analytics.RecentInterviews = append(analytics.RecentInterviews, summary)
			}
		}
	}

	c.JSON(http.StatusOK, analytics)
}

// GetInterviewReport generates a detailed report for a specific interview
func GetInterviewReport(c *gin.Context) {
	id := c.Param("id")

	// Fetch interview details
	var candidateName, position, status string
	var totalScore sql.NullFloat64
	err := database.DB.QueryRow(`
		SELECT candidate_name, position, status, total_score
		FROM interviews WHERE id = ?
	`, id).Scan(&candidateName, &position, &status, &totalScore)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
		return
	}

	// Fetch all responses with questions
	rows, err := database.DB.Query(`
		SELECT
			q.text as question,
			q.category,
			r.response_text,
			r.duration,
			r.sentiment_score,
			r.confidence_score,
			r.analysis
		FROM responses r
		JOIN questions q ON r.question_id = q.id
		WHERE r.interview_id = ?
		ORDER BY r.created_at ASC
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch responses"})
		return
	}
	defer rows.Close()

	var responses []gin.H
	var totalSentiment, totalConfidence float64
	var responseCount int

	for rows.Next() {
		var question, category, responseText, analysis string
		var duration int
		var sentimentScore, confidenceScore sql.NullFloat64

		if err := rows.Scan(&question, &category, &responseText, &duration, &sentimentScore, &confidenceScore, &analysis); err != nil {
			continue
		}

		resp := gin.H{
			"question":     question,
			"category":     category,
			"response":     responseText,
			"duration":     duration,
			"analysis":     analysis,
		}

		if sentimentScore.Valid {
			resp["sentiment_score"] = sentimentScore.Float64
			totalSentiment += sentimentScore.Float64
		}
		if confidenceScore.Valid {
			resp["confidence_score"] = confidenceScore.Float64
			totalConfidence += confidenceScore.Float64
		}

		responses = append(responses, resp)
		responseCount++
	}

	// Calculate averages
	avgSentiment := 0.0
	avgConfidence := 0.0
	if responseCount > 0 {
		avgSentiment = totalSentiment / float64(responseCount)
		avgConfidence = totalConfidence / float64(responseCount)
	}

	report := gin.H{
		"candidate_name":       candidateName,
		"position":             position,
		"status":               status,
		"total_score":          totalScore.Float64,
		"responses":            responses,
		"response_count":       responseCount,
		"average_sentiment":    avgSentiment,
		"average_confidence":   avgConfidence,
	}

	c.JSON(http.StatusOK, report)
}

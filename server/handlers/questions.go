package handlers

import (
	"net/http"
	"time"

	"hr-pepper-server/database"
	"hr-pepper-server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetQuestions returns all active interview questions
func GetQuestions(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, text, category, difficulty, time_limit, question_order, is_active, created_at, updated_at
		FROM questions
		WHERE is_active = 1
		ORDER BY question_order ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		if err := rows.Scan(&q.ID, &q.Text, &q.Category, &q.Difficulty, &q.TimeLimit, &q.Order, &q.IsActive, &q.CreatedAt, &q.UpdatedAt); err != nil {
			continue
		}
		questions = append(questions, q)
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

// GetQuestion returns a single question by ID
func GetQuestion(c *gin.Context) {
	id := c.Param("id")

	var q models.Question
	err := database.DB.QueryRow(`
		SELECT id, text, category, difficulty, time_limit, question_order, is_active, created_at, updated_at
		FROM questions WHERE id = ?
	`, id).Scan(&q.ID, &q.Text, &q.Category, &q.Difficulty, &q.TimeLimit, &q.Order, &q.IsActive, &q.CreatedAt, &q.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, q)
}

// CreateQuestion creates a new interview question
func CreateQuestion(c *gin.Context) {
	var req models.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	now := time.Now()

	if req.Difficulty == "" {
		req.Difficulty = "medium"
	}
	if req.TimeLimit == 0 {
		req.TimeLimit = 120
	}

	_, err := database.DB.Exec(`
		INSERT INTO questions (id, text, category, difficulty, time_limit, question_order, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, 1, ?, ?)
	`, id, req.Text, req.Category, req.Difficulty, req.TimeLimit, req.Order, now, now)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	question := models.Question{
		ID:         id,
		Text:       req.Text,
		Category:   req.Category,
		Difficulty: req.Difficulty,
		TimeLimit:  req.TimeLimit,
		Order:      req.Order,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	c.JSON(http.StatusCreated, question)
}

// UpdateQuestion updates an existing question
func UpdateQuestion(c *gin.Context) {
	id := c.Param("id")

	var req models.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	result, err := database.DB.Exec(`
		UPDATE questions SET text = ?, category = ?, difficulty = ?, time_limit = ?, question_order = ?, updated_at = ?
		WHERE id = ?
	`, req.Text, req.Category, req.Difficulty, req.TimeLimit, req.Order, now, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question updated successfully"})
}

// DeleteQuestion soft-deletes a question
func DeleteQuestion(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("UPDATE questions SET is_active = 0 WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

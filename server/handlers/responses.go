package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"
	"time"

	"hr-pepper-server/database"
	"hr-pepper-server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubmitResponse handles candidate response submission and analysis
func SubmitResponse(c *gin.Context) {
	var req models.SubmitResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify interview exists and is in progress
	var status string
	err := database.DB.QueryRow("SELECT status FROM interviews WHERE id = ?", req.InterviewID).Scan(&status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
		return
	}
	if status != "in_progress" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Interview is not in progress"})
		return
	}

	// Verify question exists
	var questionExists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM questions WHERE id = ?)", req.QuestionID).Scan(&questionExists)
	if err != nil || !questionExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// Analyze the response
	analysis := analyzeResponse(req.ResponseText)

	id := uuid.New().String()
	now := time.Now()

	keywordsJSON, _ := json.Marshal(analysis.Keywords)

	_, err = database.DB.Exec(`
		INSERT INTO responses (id, interview_id, question_id, response_text, response_audio, duration, sentiment_score, confidence_score, keywords_found, analysis, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, req.InterviewID, req.QuestionID, req.ResponseText, req.ResponseAudio, req.Duration, analysis.SentimentScore, analysis.ConfidenceScore, string(keywordsJSON), analysis.Analysis, now)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save response"})
		return
	}

	response := models.Response{
		ID:              id,
		InterviewID:     req.InterviewID,
		QuestionID:      req.QuestionID,
		ResponseText:    req.ResponseText,
		Duration:        req.Duration,
		SentimentScore:  &analysis.SentimentScore,
		ConfidenceScore: &analysis.ConfidenceScore,
		KeywordsFound:   analysis.Keywords,
		Analysis:        analysis.Analysis,
		CreatedAt:       now,
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": response,
		"analysis": gin.H{
			"sentiment_score":  analysis.SentimentScore,
			"confidence_score": analysis.ConfidenceScore,
			"keywords":         analysis.Keywords,
			"feedback":         analysis.Analysis,
		},
	})
}

// ResponseAnalysis represents the analysis result
type ResponseAnalysis struct {
	SentimentScore  float64
	ConfidenceScore float64
	Keywords        []string
	Analysis        string
}

// analyzeResponse performs basic text analysis on the response
func analyzeResponse(text string) ResponseAnalysis {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	wordCount := len(words)

	// Sentiment analysis based on keyword matching
	positiveWords := []string{
		"achieved", "accomplished", "successful", "improved", "developed", "created",
		"led", "managed", "increased", "exceeded", "innovative", "passionate",
		"dedicated", "collaborative", "team", "growth", "opportunity", "excited",
		"skills", "experience", "learn", "contribute", "challenge", "solution",
	}
	negativeWords := []string{
		"failed", "couldn't", "didn't", "never", "problem", "difficult",
		"hate", "boring", "quit", "fired", "conflict", "mistake",
	}

	positiveCount := 0
	negativeCount := 0
	foundKeywords := []string{}

	for _, word := range words {
		word = strings.Trim(word, ".,!?;:")
		for _, pos := range positiveWords {
			if word == pos {
				positiveCount++
				if !contains(foundKeywords, word) {
					foundKeywords = append(foundKeywords, word)
				}
			}
		}
		for _, neg := range negativeWords {
			if word == neg {
				negativeCount++
			}
		}
	}

	// Calculate sentiment score (0-1)
	sentimentScore := 0.5
	if positiveCount+negativeCount > 0 {
		sentimentScore = float64(positiveCount) / float64(positiveCount+negativeCount)
	}
	sentimentScore = math.Round(sentimentScore*100) / 100

	// Calculate confidence score based on response length and structure
	confidenceScore := calculateConfidenceScore(text, wordCount)

	// Generate analysis feedback
	analysis := generateAnalysisFeedback(wordCount, sentimentScore, confidenceScore, foundKeywords)

	return ResponseAnalysis{
		SentimentScore:  sentimentScore,
		ConfidenceScore: confidenceScore,
		Keywords:        foundKeywords,
		Analysis:        analysis,
	}
}

func calculateConfidenceScore(text string, wordCount int) float64 {
	score := 0.0

	// Length score (optimal: 50-200 words)
	if wordCount >= 50 && wordCount <= 200 {
		score += 0.4
	} else if wordCount >= 30 && wordCount <= 300 {
		score += 0.25
	} else if wordCount >= 10 {
		score += 0.1
	}

	// Structure indicators
	structureIndicators := []string{
		"first", "second", "then", "after", "finally", "because",
		"for example", "specifically", "in addition", "as a result",
	}
	for _, indicator := range structureIndicators {
		if strings.Contains(text, indicator) {
			score += 0.05
		}
	}

	// Specific examples (numbers, percentages)
	if strings.ContainsAny(text, "0123456789%") {
		score += 0.1
	}

	// Capped at 1.0
	if score > 1.0 {
		score = 1.0
	}

	return math.Round(score*100) / 100
}

func generateAnalysisFeedback(wordCount int, sentiment, confidence float64, keywords []string) string {
	var feedback []string

	// Length feedback
	if wordCount < 30 {
		feedback = append(feedback, "Response is quite brief. Consider providing more detail.")
	} else if wordCount > 300 {
		feedback = append(feedback, "Response is lengthy. Consider being more concise.")
	} else {
		feedback = append(feedback, "Response length is appropriate.")
	}

	// Sentiment feedback
	if sentiment >= 0.7 {
		feedback = append(feedback, "Positive tone with good use of achievement-oriented language.")
	} else if sentiment >= 0.4 {
		feedback = append(feedback, "Neutral tone. Could benefit from more positive framing.")
	} else {
		feedback = append(feedback, "Consider reframing with more positive language.")
	}

	// Keywords feedback
	if len(keywords) >= 3 {
		feedback = append(feedback, "Good use of professional keywords: "+strings.Join(keywords[:min(5, len(keywords))], ", ")+".")
	} else if len(keywords) > 0 {
		feedback = append(feedback, "Limited professional keywords detected.")
	}

	// Confidence feedback
	if confidence >= 0.6 {
		feedback = append(feedback, "Well-structured response with clear examples.")
	} else if confidence >= 0.3 {
		feedback = append(feedback, "Consider adding more specific examples or structure.")
	} else {
		feedback = append(feedback, "Would benefit from concrete examples and better structure.")
	}

	return strings.Join(feedback, " ")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

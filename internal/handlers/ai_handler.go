package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"irontrack-backend/internal/database"
	"irontrack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

func getAIClient(ctx context.Context) (*genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key not configured")
	}
	return genai.NewClient(ctx, option.WithAPIKey(apiKey))
}

type GeneratePlanRequest struct {
	models.UserProfile
	Language string `json:"language"`
}

func GenerateWorkoutPlan(c *gin.Context) {
	var req GeneratePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	client, err := getAIClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	model.ResponseMIMEType = "application/json"
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are an expert fitness coach. Create structured, safe, and effective workout plans tailored to the user's biometrics and goals. Output JSON matching the schema: {name, description, targetGoal, exercises: [{name, defaultSets (int), defaultReps (int), muscleGroup, instructions}]}. IMPORTANT: defaultSets and defaultReps must be strictly integers, not strings or ranges."))

	languageInstruction := ""
	if req.Language != "" {
		languageInstruction = fmt.Sprintf("Please generate the plan content (name, description, instructions) in %s.", req.Language)
	}

	prompt := fmt.Sprintf(`Create a workout plan for a user with the following profile:
      - Gender: %s
      - Age: %s
      - Height: %s
      - Weight: %s
      - Goal: %s
      - Target Duration: %s
      - Experience Level: %s
      
      %s
      Provide a comprehensive single-day example routine that fits this duration and physical profile.`,
		req.Gender, req.Age, req.Height, req.Weight, req.MainGoal, req.WorkoutDuration, req.ExperienceLevel, languageInstruction)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate content: " + err.Error()})
		return
	}

	var resultText string
	if len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				resultText += string(txt)
			}
		}
	}

	// Validate and parse into WorkoutPlan model
	var plan models.WorkoutPlan
	if err := json.Unmarshal([]byte(resultText), &plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON from AI: " + err.Error()})
		return
	}

	// Populate metadata fields
	userID := c.GetString("userID")
	plan.ID = uuid.New().String()
	plan.UserID = userID
	plan.IsAiGenerated = true
	plan.CreatedAt = time.Now()

	logAIRequest(userID, "generate_plan")

	c.JSON(http.StatusOK, plan)
}

type ReportRequest struct {
	Range     string `json:"range"`
	Count     int    `json:"count"`
	Hours     int    `json:"hours"`
	Minutes   int    `json:"minutes"`
	TopMuscle string `json:"topMuscle"`
}

func GenerateProgressReport(c *gin.Context) {
	var stats ReportRequest
	if err := c.ShouldBindJSON(&stats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	client, err := getAIClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are an encouraging data-driven fitness coach."))

	prompt := fmt.Sprintf(`
        Analyze the following workout statistics for the user over the selected period:
        - Time Period: %s
        - Total Workouts: %d
        - Total Duration: %d hours and %d minutes
        - Most Trained Muscle Group: %s
        
        Provide a brief, encouraging summary of their performance and 3 specific, actionable suggestions for improvement or balance (e.g., if they only train chest, suggest back/legs). Keep the tone motivational but professional. Limit to 150 words.
      `, stats.Range, stats.Count, stats.Hours, stats.Minutes, stats.TopMuscle)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate content: " + err.Error()})
		return
	}

	var resultText string
	if len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				resultText += string(txt)
			}
		}
	}

	userID := c.GetString("userID")
	logAIRequest(userID, "generate_report")

	c.JSON(http.StatusOK, gin.H{"response": resultText})
}

func logAIRequest(userID, typ string) {
	if userID == "" {
		return
	}
	_ = database.DB.Create(&models.AIRequestLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      typ,
		CreatedAt: time.Now(),
	}).Error
}

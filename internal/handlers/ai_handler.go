package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"irontrack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func getAIClient(ctx context.Context) (*genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key not configured")
	}
	return genai.NewClient(ctx, option.WithAPIKey(apiKey))
}

func GenerateWorkoutPlan(c *gin.Context) {
	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
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
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are an expert fitness coach. Create structured, safe, and effective workout plans tailored to the user's biometrics and goals. Output JSON matching the schema: {name, description, targetGoal, exercises: [{name, defaultSets, defaultReps, muscleGroup, instructions}]}"))

	prompt := fmt.Sprintf(`Create a workout plan for a user with the following profile:
      - Gender: %s
      - Age: %s
      - Height: %s
      - Weight: %s
      - Goal: %s
      - Target Duration: %s
      - Experience Level: %s
      
      Provide a comprehensive single-day example routine that fits this duration and physical profile.`,
		profile.Gender, profile.Age, profile.Height, profile.Weight, profile.MainGoal, profile.WorkoutDuration, profile.ExperienceLevel)

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

	// Validate JSON
	var js map[string]interface{}
	if err := json.Unmarshal([]byte(resultText), &js); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON from AI"})
		return
	}

	c.JSON(http.StatusOK, js)
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

	model := client.GenerativeModel("gemini-1.5-flash-001")
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

	c.JSON(http.StatusOK, gin.H{"response": resultText})
}

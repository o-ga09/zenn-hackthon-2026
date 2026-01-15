package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/labstack/echo"
	"github.com/o-ga09/zenn-hackthon-2026/internal/domain"
	"github.com/o-ga09/zenn-hackthon-2026/pkg/config"
)

type IAgentServer interface {
	CreateVLog(echo.Context) error
}

type AgentServer struct {
	storage    domain.IImageStorage
	recipeFlow *core.Flow[*RecipeInput, *Recipe, struct{}]
}

// Define input schema
type RecipeInput struct {
	Ingredient          string `json:"ingredient" jsonschema:"description=Main ingredient or cuisine type"`
	DietaryRestrictions string `json:"dietaryRestrictions,omitempty" jsonschema:"description=Any dietary restrictions"`
}

// Define output schema
type Recipe struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	PrepTime     string   `json:"prepTime"`
	CookTime     string   `json:"cookTime"`
	Servings     int      `json:"servings"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Tips         []string `json:"tips,omitempty"`
}

func NewAgentServer(ctx context.Context, storage domain.IImageStorage) *AgentServer {
	g := config.GetGenkitCtx(ctx)
	// Define a recipe generator flow
	recipeGeneratorFlow := genkit.DefineFlow(g, "recipeGeneratorFlow", func(ctx context.Context, input *RecipeInput) (*Recipe, error) {
		// Create a prompt based on the input
		dietaryRestrictions := input.DietaryRestrictions
		if dietaryRestrictions == "" {
			dietaryRestrictions = "none"
		}

		prompt := fmt.Sprintf(`Create a recipe with the following requirements:
            Main ingredient: %s
            Dietary restrictions: %s`, input.Ingredient, dietaryRestrictions)

		// Generate structured recipe data using the same schema
		recipe, _, err := genkit.GenerateData[Recipe](ctx, g,
			ai.WithPrompt(prompt),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate recipe: %w", err)
		}

		return recipe, nil
	})
	return &AgentServer{
		recipeFlow: recipeGeneratorFlow,
		storage:    storage,
	}
}

func (s *AgentServer) CreateVLog(c echo.Context) error {
	ctx := c.Request().Context()
	input := &RecipeInput{
		Ingredient:          "chicken",
		DietaryRestrictions: "gluten-free",
	}
	output, err := s.recipeFlow.Run(ctx, input)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, output)
}

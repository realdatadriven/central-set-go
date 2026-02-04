

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/firebase/genkit/go"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/googlegenai" // swap to ollama if preferred
)

var (
	g         *genkit.Genkit
	modelName = "googleai/gemini-1.5-flash" // or "ollama/llama3.1" etc.
	systemPrompt = `...paste the full improved prompt from above...`
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Messages []Message `json:"messages"`
}

type Response struct {
	Content string `json:"content"`
}

func initGenkit() {
	ctx := context.Background()
	var err error
	g, err = genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}), // or ollama.Ollama{ServerAddress: "http://localhost:11434"}
		genkit.WithLogLevel("debug"), // optional
	)
	if err != nil {
		log.Fatalf("Failed to init Genkit: %v", err)
	}
	log.Println("Genkit initialized with model:", modelName)
}

func etlxAssistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Build Genkit messages: prepend system if needed
	msgs := []ai.Message{}
	hasSystem := false
	for _, m := range req.Messages {
		if m.Role == "system" {
			hasSystem = true
		}
		role := ai.RoleUser
		if m.Role == "assistant" {
			role = ai.RoleModel
		}
		msgs = append(msgs, ai.Message{
			Role:    role,
			Content: []ai.Part{ai.NewTextPart(m.Content)},
		})
	}

	if !hasSystem {
		msgs = append([]ai.Message{
			{
				Role:    ai.RoleSystem,
				Content: []ai.Part{ai.NewTextPart(systemPrompt)},
			},
		}, msgs...)
	}

	// Generate
	resp, err := genkit.Generate(ctx, g,
		ai.WithModelName(modelName),
		ai.WithMessages(msgs...),
		// Optional: ai.WithTemperature(0.2), ai.WithMaxTokens(4000)
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Response{Content: resp.Text()})
}
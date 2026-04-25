package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/firebase/genkit/go/plugins/anthropic"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/ollama"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Provider         string    `json:"provider"`
	Model            string    `json:"model"`
	BaseURL          string    `json:"base_url,omitempty"`
	SystemPrompt     string    `json:"system_prompt,omitempty"`
	SystemPromptFile string    `json:"system_prompt_file,omitempty"`
	Messages         []Message `json:"messages"`
	Temperature      *float32  `json:"temperature,omitempty"`
	MaxTokens        *int      `json:"max_tokens,omitempty"`
}

type Response struct {
	Content string `json:"content"`
}

func initGenkit(ctx context.Context, req Request) (*genkit.Genkit, error) {
	switch req.Provider {

	case "googleai":
		return genkit.Init(ctx,
			genkit.WithPlugins(&googlegenai.GoogleAI{}),
		), nil

	case "ollama":
		base := req.BaseURL
		if base == "" {
			base = "http://localhost:11434"
		}
		return genkit.Init(ctx,
			genkit.WithPlugins(&ollama.Ollama{
				ServerAddress: base,
			}),
		), nil

	case "openai":
		return genkit.Init(ctx,
			genkit.WithPlugins(&openai.OpenAI{
				APIKey: os.Getenv("OPENAI_API_KEY"),
			}),
		), nil

	case "anthropic":
		return genkit.Init(ctx,
			genkit.WithPlugins(&anthropic.Anthropic{
				APIKey: os.Getenv("ANTHROPIC_API_KEY"),
			}),
		), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

func buildSystemPrompt(req Request) string {
	// Default fallback (important safeguard)
	defaultPrompt := "You are a helpful AI assistant."
	prompt := ""
	if req.SystemPrompt != "" {
		prompt = req.SystemPrompt
	}
	if req.SystemPromptFile != "" {
		if data, err := os.ReadFile(req.SystemPromptFile); err == nil {
			prompt += "\n" + string(data)
		}
	}
	if prompt == "" {
		prompt = defaultPrompt
	}
	return prompt
}

func aiAssistHandler(w http.ResponseWriter, r *http.Request) {
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

	g, err := initGenkit(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert messages
	msgs := []*ai.Message{}
	hasSystem := false

	for _, m := range req.Messages {
		if m.Role == "system" {
			hasSystem = true
		}

		role := ai.RoleUser
		if m.Role == "assistant" {
			role = ai.RoleModel
		}

		msgs = append(msgs, &ai.Message{
			Role:    role,
			Content: []*ai.Part{ai.NewTextPart(m.Content)},
		})
	}

	// System prompt
	systemPrompt := buildSystemPrompt(req)

	if !hasSystem {
		msgs = append([]*ai.Message{
			{
				Role:    ai.RoleSystem,
				Content: []*ai.Part{ai.NewTextPart(systemPrompt)},
			},
		}, msgs...)
	}

	fullModel := fmt.Sprintf("%s/%s", req.Provider, req.Model)

	// Build options dynamically
	opts := []ai.GenerateOption{
		ai.WithModelName(fullModel),
		ai.WithMessages(msgs...),
	}

	if req.Temperature != nil {
		// opts = append(opts, ai.WithTemperature(*req.Temperature))
	}

	if req.MaxTokens != nil {
		// opts = append(opts, ai.WithMaxTokens(*req.MaxTokens))
	}

	resp, err := genkit.Generate(ctx, g, opts...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Response{
		Content: resp.Text(),
	})
}

/*
{
  "provider": "googleai",
  "model": "gemini-1.5-flash",
  "messages": [
    {"role": "user", "content": "Hello"}
  ]
}
{
  "provider": "ollama",
  "model": "llama3.1",
  "base_url": "http://localhost:11434",
  "messages": [...]
}
{
  "provider": "googleai",
  "model": "gemini-1.5-flash",
  "system_prompt": "You are a ETLX assistant.",
  "system_prompt_file": "etlxllm.txt",
  "messages": [
    	{"role": "user", "content": "Analyze this dataset"}
  ]
}
func main() {
	http.HandleFunc("/etlx-assist", etlxAssistHandler)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}*/

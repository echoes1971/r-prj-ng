package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var ollamaURL string
var ollamaModel string
var ollamaAppName string

func OllamaInit(appName, url, model string) error {
	ollamaURL = url
	ollamaModel = model
	ollamaAppName = appName

	if ollamaURL != "" && ollamaModel != "" {
		go UpdateOllamaDefaultPageResponse("en")
		log.Printf("Ollama initialized with URL: %s and Model: %s\n", ollamaURL, ollamaModel)
	} else {
		log.Println("Ollama not configured - using fallback responses")
		lastDefaultPageResponse = "<h2>Welcome! 👋</h2><p>Please log in to continue using the application.</p>"
		// lastDefaultPageResponse = "\u003ch2\u003eWelcome to the Galaxy 🚀\u003c/h2\u003e\n\n\u003cp\u003e Warning: Abandon all hope, ye who enter here... or at least try our web app! 😂\u003c/p\u003e\n\n\u003cb\u003eGalactic Hitcher's Guide\u003c/b\u003e\u003cbr\u003e\n\u003cbr\u003e\n\n1. First, buckle up and strap yourself in, because things are about to get weird. 🚀\n2. Our web app is powered by a combination of magic dust and cutting-edge technology (just kidding, it's just JavaScript). ✨\n3. Be prepared for epic battles with bugs and occasional crashes into the mothership (aka error 404). 💥\n4. But don't worry, our team of expert space rangers will be here to guide you through the galaxy and fix any problems that come your way. 🚀\n\nSo, are you ready to embark on this intergalactic adventure? Let's get started! 🔴"

	}

	return nil
}

func OllamaHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Prompt string `json:"prompt"`
	}
	type Response struct {
		Response string `json:"response"`
		Error    string `json:"error,omitempty"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondSimpleError(w, ErrInvalidRequest, "Invalid request body", http.StatusBadRequest)
		return
	}

	if ollamaURL == "" || ollamaModel == "" {
		res := Response{Error: "Ollama service not configured"}
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}

	respText, err := CallOllama(req.Prompt)
	if err != nil {
		res := Response{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		return
	}

	res := Response{Response: respText}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

var lastDefaultPageResponse string

func DefaultPageOllamaHandler(w http.ResponseWriter, r *http.Request) {
	// Get lang query parameter
	langParam := r.URL.Query().Get("lang")
	if langParam == "" {
		langParam = "en"
	}
	log.Print("Requested language tag: ", langParam)

	// Just use English for now
	langParam = "en"
	go UpdateOllamaDefaultPageResponse(langParam)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"response": lastDefaultPageResponse,
	})
}

var prompts = []string{
	"Provide a short funny star trek based welcome message for the default page of ",
	"Provide a short funny star wars based welcome message for the default page of ",
	"Provide a short funny hitcher's guide to galaxy based welcome message for the default page of ",
	"Provide a short funny positive welcome message for the default page of ",
}

func UpdateOllamaDefaultPageResponse(languageTag string) {
	if ollamaURL == "" || ollamaModel == "" {
		log.Print("Ollama not configured, skipping default page response update.")
		return
	}

	// ollamaModel = modelList[0]
	prompt := prompts[2]
	prompt += " a web application"
	// prompt += " " + ollamaAppName + " web application"
	switch languageTag {
	case "it", "it-IT":
		prompt += prompt + " in Italian."
	case "en", "en-US":
		prompt += prompt + " in English."
	case "fr", "fr-FR":
		prompt += prompt + " in French."
	case "de", "de-DE":
		prompt += prompt + " in German."
	}
	prompt += " Use emojis."
	prompt += " Use only <h2>, <p>, <b>, <br> tags. No links <a>. Output html code only."

	// prompt = "Créez un court message de bienvenue humoristique, "
	// prompt += "inspiré de animal crossing, "
	// prompt += "pour la page d'accueil d'une application web. "
	// prompt += "Rédigez ce message en français. "
	// prompt += "Utilisez des émojis. "
	// prompt += "Utilisez uniquement les tags <h2>, <p>, <b> et <br>. Pas de liens <a>. Générez uniquement du code HTML."

	log.Print("Prompt: ", prompt)

	respText, err := CallOllama(prompt)
	if err != nil {
		log.Printf("Error getting default page response from Ollama: %v\n", err)
		respText = "Welcome to our application!"
	}

	lastDefaultPageResponse = respText
}

// CallOllama sends a prompt to the Ollama API and returns the response
func CallOllama(prompt string) (string, error) {
	if ollamaURL == "" || ollamaModel == "" {
		return "", fmt.Errorf("Ollama not configured")
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": ollamaModel,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream": false,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s", ollamaURL), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf("Ollama response: %s\n", string(body))

	// Parse the response
	// {"model":"llama3.2:latest","created_at":"2025-11-13T13:17:45.487900903Z","message":{"role":"assistant","content":"I'm an artificial intelligence model known as Llama. Llama stands for \"Large Language Model Meta AI.\""},"done":true,"done_reason":"stop","total_duration":1593273648,"load_duration":103413131,"prompt_eval_count":29,"prompt_eval_duration":132183696,"eval_count":23,"eval_duration":1338392262}

	var ollamaResponseSingle struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.Unmarshal(body, &ollamaResponseSingle); err == nil && ollamaResponseSingle.Message.Content != "" {
		content := ollamaResponseSingle.Message.Content
		content = strings.ReplaceAll(content, "```html", "")
		content = strings.ReplaceAll(content, "```", "")
		content = strings.ReplaceAll(content, "[app name]", ollamaAppName)
		content = strings.ReplaceAll(content, "[App Name]", ollamaAppName)
		return content, nil
	}

	return "", fmt.Errorf("invalid response from Ollama")

}

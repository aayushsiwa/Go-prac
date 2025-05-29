// morse/api/handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

// Morse code mappings
var textToMorse = map[rune]string{
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".", 'F': "..-.",
	'G': "--.", 'H': "....", 'I': "..", 'J': ".---", 'K': "-.-", 'L': ".-..",
	'M': "--", 'N': "-.", 'O': "---", 'P': ".--.", 'Q': "--.-", 'R': ".-.",
	'S': "...", 'T': "-", 'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-",
	'Y': "-.--", 'Z': "--..",
	'0': "-----", '1': ".----", '2': "..---", '3': "...--", '4': "....-",
	'5': ".....", '6': "-....", '7': "--...", '8': "---..", '9': "----.",
	' ': "/", '.': ".-.-.-", ',': "--..--", '?': "..--..", '\'': ".----.",
	'!': "-.-.--", '/': "-..-.", '(': "-.--.", ')': "-.--.-", '&': ".-...",
	':': "---...", ';': "-.-.-.", '=': "-...-", '+': ".-.-.", '-': "-....-",
	'_': "..--.-", '"': ".-..-.", '$': "...-..-", '@': ".--.-.",
}

var morseToText map[string]rune

func init() {
	// Create reverse mapping
	morseToText = make(map[string]rune)
	for char, morse := range textToMorse {
		morseToText[morse] = char
	}
}

type ConvertRequest struct {
	Text string `json:"text"`
	Type string `json:"type"` // "to_morse" or "to_text"
}

type ConvertResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func textToMorseCode(text string) string {
	text = strings.ToUpper(text)
	var result []string

	for _, char := range text {
		if morse, exists := textToMorse[char]; exists {
			result = append(result, morse)
		} else if char == ' ' {
			result = append(result, "/")
		}
	}

	return strings.Join(result, " ")
}

func morseToTextCode(morse string) string {
	// Clean up the morse code input
	morse = strings.TrimSpace(morse)

	// Split by multiple spaces or forward slashes for word separation
	words := regexp.MustCompile(`\s*/\s*|\s{2,}`).Split(morse, -1)

	var resultWords []string

	for _, word := range words {
		if word == "" {
			continue
		}

		// Split each word by single spaces for character separation
		chars := strings.Fields(word)
		var resultChars []rune

		for _, char := range chars {
			if char == "/" {
				resultChars = append(resultChars, ' ')
			} else if text, exists := morseToText[char]; exists {
				resultChars = append(resultChars, text)
			}
		}

		if len(resultChars) > 0 {
			resultWords = append(resultWords, string(resultChars))
		}
	}

	return strings.Join(resultWords, " ")
}

// Handler is the main entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		response := ConvertResponse{Error: "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ConvertResponse{Error: "Invalid JSON"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var result string

	switch req.Type {
	case "to_morse":
		result = textToMorseCode(req.Text)
	case "to_text":
		result = morseToTextCode(req.Text)
	default:
		response := ConvertResponse{Error: "Invalid conversion type"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ConvertResponse{Result: result}
	json.NewEncoder(w).Encode(response)
}

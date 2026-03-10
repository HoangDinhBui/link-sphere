package response

import (
	"encoding/json"
	"net/http"
)

// Response is the standard API response format.
type Response struct {
	CodeStatus int         `json:"code_status"`
	Message    string      `json:"message"`
	Result     bool        `json:"result"`
	Errors     interface{} `json:"errors"`
	Data       interface{} `json:"data"`
}

// JSON sends a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		CodeStatus: status,
		Message:    http.StatusText(status),
		Result:     status >= 200 && status < 300,
		Errors:     map[string]interface{}{},
		Data:       data,
	})
}

// Error sends an error JSON response.
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		CodeStatus: status,
		Message:    message,
		Result:     false,
		Errors:     map[string]string{"error": message},
		Data:       nil,
	})
}

// Success sends a success JSON response with a message.
func Success(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		CodeStatus: http.StatusOK,
		Message:    message,
		Result:     true,
		Errors:     map[string]interface{}{},
		Data:       data,
	})
}

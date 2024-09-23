package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UserIDKey is the key used to store the user ID in the context.
type ContextKey string

const UserIDKey ContextKey = "userID"

func Encode[T any](w http.ResponseWriter, status int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

type Validator interface {
	Valid() (problems map[string]string)
}

func DecodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}

	if problems := v.Valid(); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}

	return v, nil, nil
}

type ErrorList struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func HandleInvalidRequest(w http.ResponseWriter, problems map[string]string) {
	var errors []Error

	if len(problems) > 0 {
		errors = make([]Error, 0, len(problems))
		for field, problem := range problems {
			err := Error{
				Message: fmt.Sprintf("invalid %s", field),
				Details: problem,
			}
			errors = append(errors, err)
		}
	} else {
		errors = make([]Error, 0, 1)
		errors = append(errors, Error{
			Message: "invalid input",
			Details: "failed to parse request",
		})
	}

	HandleError(w, http.StatusBadRequest, errors...)
}

func HandleError(w http.ResponseWriter, status int, err ...Error) {
	errList := ErrorList{Errors: err}
	if err := Encode(w, status, errList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

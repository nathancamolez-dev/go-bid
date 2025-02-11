package jsonutils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nathancamolez-dev/go-bid/internal/validator"
)

func EncodeJson[T any](w http.ResponseWriter, r *http.Request, statusCode int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("failed to encode json %w", err)

	}
	return nil
}

func DecodeValidJson[T validator.Validator](r *http.Request) (T, map[string]string, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, nil, fmt.Errorf("decode json %w", err)
	}
	if problems := data.Valid(r.Context()); len(problems) > 0 {
		return data, problems, fmt.Errorf("invalid %T : %d problems", data, len(problems))
	}
	return data, nil, nil
}

func DecodeJson[T any](r *http.Request) (T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("failed to decode json %w", err)
	}
	return data, nil
}

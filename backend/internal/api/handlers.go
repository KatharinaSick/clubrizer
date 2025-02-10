package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"log"
	"net/http"
)

// type targetFunc[Out any] func(context.Context) (Out, error)
type targetFuncWithInput[In any, Out any] func(context.Context, In) (*Out, error)

/*
TODO implement those functions

func handle[Out any](f targetFunc[Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call target function
		out, err := f(r.Context())
		if err != nil {
			// TODO return correct status code based on error
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeResponse(w, out)
	})
}

func handleWithParams[In any, Out any](f targetFuncWithInput[In, Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in In

		// Retrieve & validate payload
		// TODO get params
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		validate := validator.New()
		if err := validate.Struct(in); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		// Call target function
		out, err := f(r.Context(), in)
		if err != nil {
			// TODO return correct status code based on error
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeResponse(w, out)
	})
}
*/

func handleWithBody[In any, Out any](f targetFuncWithInput[In, Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in In

		// Retrieve & validate payload
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse json: %s", err), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		if err := validate.Struct(in); err != nil {
			http.Error(w, fmt.Sprintf("failed to parse json: %s", err), http.StatusBadRequest)
			return
		}

		// Call target function
		out, err := f(r.Context(), in)
		if err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		writeResponse(w, out)
	})
}

func writeResponse[Out any](w http.ResponseWriter, out Out) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(out)
	if err != nil {
		// TODO use proper logger
		log.Printf("failed to encode created note: %v", err)
		return
	}
}

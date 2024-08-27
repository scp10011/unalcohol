package unalcohol

import "net/http"

func DefaultStatusBadRequest(r *http.Request, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusBadRequest)
	return nil
}

func DefaultStatusMethodNotAllowed(r *http.Request, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}

func DefaultStatusInternalServerError(r *http.Request, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusInternalServerError)
	return nil
}

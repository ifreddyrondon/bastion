package utils

// Auxiliary structure that contains the necessary data to give a response
// with the error context.
type responseError struct {
	Message string `json:"message"`
	Errors  string `json:"error"`
	Status  int    `json:"status"`
}

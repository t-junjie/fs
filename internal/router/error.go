package router

// APIErrorResponse is an API error that is marshaled into a JSON response.
type APIErrorResponse struct {
	ErrorMsg string `json:"error"`
}

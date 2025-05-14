package gin_handler

type HTTPError struct {
	Code    string `json:"code" example:"EXAMPLE_ERROR_CODE"`
	Message string `json:"message" example:"example error message"`
}

package common

type SwaggerErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

type SwaggerValidationErrorResponse struct {
	Status string            `json:"status" example:"error"`
	Errors map[string]string `json:"errors"`
}

type SwaggerCreatedResponse struct {
	Status  int    `json:"status" example:"201"`
	Message string `json:"message" example:"Transaction Created"`
}

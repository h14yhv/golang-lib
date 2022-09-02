package rest

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail"`
}

package models

type Response struct {
	Message string `json:"message"`
}

func CustomResponse(newMessage string) Response {
	return Response{Message: newMessage}
}

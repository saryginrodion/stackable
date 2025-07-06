package stackable

type HttpError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e HttpError) Error() string {
	return e.Message
}

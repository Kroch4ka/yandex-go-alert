package util

type StatusErr struct {
	Status  int
	Message string
}

func (se StatusErr) Error() string {
	return se.Message
}

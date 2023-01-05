package utility

type RequestManager interface {
	Request(id string, url string, uqlPayload string) ([]byte, error)
}

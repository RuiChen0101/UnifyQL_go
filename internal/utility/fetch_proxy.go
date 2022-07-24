package utility

type FetchProxy interface {
	Request(id string, url string, uqlPayload string) ([]byte, error)
}

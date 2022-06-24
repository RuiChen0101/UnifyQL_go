package utility

type FetchProxy interface {
	Request(url string, uqlPayload string) ([]byte, error)
}

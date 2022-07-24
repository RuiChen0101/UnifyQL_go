package fake

import "errors"

type FetchRequest struct {
	Id         string
	Url        string
	UqlPayload string
}

type FakeFetchProxy struct {
	expectRes []string
	requests  []FetchRequest
}

func NewFakeFetchProxy(expectRes []string) *FakeFetchProxy {
	return &FakeFetchProxy{
		expectRes: expectRes,
	}
}

func (ffp *FakeFetchProxy) GetRecord(index int) FetchRequest {
	return ffp.requests[index]
}

func (ffp *FakeFetchProxy) Request(id string, url string, uqlPayload string) ([]byte, error) {
	if len(ffp.expectRes) <= 0 {
		return nil, errors.New("Too many request")
	}
	req := ffp.expectRes[0]
	ffp.expectRes = ffp.expectRes[1:]
	ffp.requests = append(ffp.requests, FetchRequest{
		Id:         id,
		Url:        url,
		UqlPayload: uqlPayload,
	})
	return []byte(req), nil
}

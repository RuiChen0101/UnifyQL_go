package fake

type FakeIdGenerator struct{}

func (fig *FakeIdGenerator) NanoId8() string {
	return "12345678"
}

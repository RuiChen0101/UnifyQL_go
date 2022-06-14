package utility

import gonanoid "github.com/matoous/go-nanoid/v2"

type IdGenerator interface {
	NanoId8() string
}

type DefaultIdGenerator struct{}

func (dig *DefaultIdGenerator) NanoId8() string {
	id, _ := gonanoid.Generate("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 8)
	return id
}

package utils

import (
	"github.com/satori/go.uuid"
	"strings"
)

func GenerateUUID() string {
    t, _ := uuid.NewV1()
    sig := t.String()
    sig = strings.Replace(sig, "-", "", -1)
    return sig
}

package registry_test

import (
	"testing"

	"github.com/currybab/tokgo"
	"github.com/currybab/tokgo/registry"
)

func TestDefaultEncodingRegistry(t *testing.T) {
	reg := registry.NewDefaultEncodingRegistry()
	enc, _ := reg.GetEncodingByType(tokgo.CL100K_BASE)
	if tt := enc.Decode(enc.EncodeToIntArray("hello world")); tt != "hello world" {
		panic(tt)
	}
}

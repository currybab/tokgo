package registry_test

import (
	"testing"

	mod "github.com/currybab/tokgo/mod"
	tokgo "github.com/currybab/tokgo/registry"
)

func TestDefaultEncodingRegistry(t *testing.T) {
	reg := tokgo.NewDefaultEncodingRegistry()
	enc, _ := reg.GetEncodingByType(mod.CL100K_BASE)
	if tt := enc.Decode(enc.EncodeToIntArray("hello world")); tt != "hello world" {
		panic(tt)
	}
}

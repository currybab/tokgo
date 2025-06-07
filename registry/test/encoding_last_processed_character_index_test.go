package registry_test

import (
	"testing"

	"github.com/currybab/tokgo/mod"
	tokgo "github.com/currybab/tokgo/registry"
	"github.com/stretchr/testify/assert"
)

var ENCODING, _ = tokgo.NewDefaultEncodingRegistry().GetEncodingByType(mod.CL100K_BASE)

func TestNullInput(t *testing.T) {
	encodingResult := ENCODING.Encode("", 10)
	assert.Equal(t, encodingResult.GetLastProcessedCharacterIndex(), -1)
}

func TestEmptyInput(t *testing.T) {
	encodingResult := ENCODING.Encode("", 10)
	assert.Equal(t, encodingResult.GetLastProcessedCharacterIndex(), -1)
}

func TestShortInput(t *testing.T) {
	encodingResult := ENCODING.Encode("Hello World!", 10)
	assert.Equal(t, encodingResult.GetLastProcessedCharacterIndex(), 11)
}

func TestLongInput(t *testing.T) {
	encodingResult := ENCODING.Encode("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce condimentum enim ac tellus malesuada, a consectetur nibh efficitur. 🚀🚀🚀", 10)
	assert.Equal(t, encodingResult.GetLastProcessedCharacterIndex(), 55)
}

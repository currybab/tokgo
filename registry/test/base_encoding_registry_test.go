package registry_test

import (
	"testing"

	regexp "github.com/dlclark/regexp2"

	mod "github.com/currybab/tokgo/mod"
	"github.com/stretchr/testify/assert"
)

var registry mod.EncodingRegistry

type DummyEncoding struct {
}

func (d *DummyEncoding) EncodeToIntArray(text string) []int {
	return nil
}

func (d *DummyEncoding) Encode(text string, maxTokens int) *mod.EncodingResult {
	return nil
}

func (d *DummyEncoding) EncodeOrdinaryToIntArray(text string) []int {
	return nil
}

func (d *DummyEncoding) EncodeOrdinary(text string, maxTokens int) *mod.EncodingResult {
	return nil
}

func (d *DummyEncoding) CountTokens(text string) int {
	return 0
}

func (d *DummyEncoding) CountTokensOrdinary(text string) int {
	return 0
}

func (d *DummyEncoding) Decode(tokens []int) string {
	return ""
}

func (d *DummyEncoding) DecodeBytes(tokens []int) []byte {
	return nil
}

func (d *DummyEncoding) GetName() string {
	return "dummy"
}

func getEncodingReturnsCorrectEncoding(t *testing.T) {
	for _, encodingType := range mod.EncodingTypeValues() {
		encoding, err := registry.GetEncodingByType(encodingType)
		assert.Nil(t, err)
		assert.NotNil(t, encoding)
		assert.Equal(t, encodingType.GetName(), encoding.GetName())
	}
}

func getEncodingByNameReturnsCorrectEncoding(t *testing.T) {
	for _, encodingType := range mod.EncodingTypeValues() {
		encoding, err := registry.GetEncoding(encodingType.GetName())
		assert.Nil(t, err)
		assert.NotNil(t, encoding)
		assert.Equal(t, encodingType.GetName(), encoding.GetName())
	}
}

func getEncodingForModelReturnsCorrectEncoding(t *testing.T) {
	for _, modelType := range mod.ModelTypeValues() {
		encoding, err := registry.GetEncodingForModelType(modelType)
		assert.Nil(t, err)
		assert.NotNil(t, encoding)
		assert.Equal(t, modelType.GetEncodingType().GetName(), encoding.GetName())
	}
}

func getEncodingForModelByNameReturnsCorrectEncoding(t *testing.T) {
	for _, modelType := range mod.ModelTypeValues() {
		encoding, err := registry.GetEncodingForModel(modelType.GetName())
		assert.Nil(t, err)
		assert.NotNil(t, encoding)
		assert.Equal(t, modelType.GetEncodingType().GetName(), encoding.GetName())
	}
}

func getEncodingForModelByPrefixReturnsCorrectEncodingForGpt432k(t *testing.T) {
	encoding, err := registry.GetEncodingForModel("gpt-4-32k-0314")
	assert.Nil(t, err)
	assert.NotNil(t, encoding)
	assert.Equal(t, mod.GPT_4.GetEncodingType().GetName(), encoding.GetName())
}

func getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4(t *testing.T) {
	encoding, err := registry.GetEncodingForModel("gpt-4-0314")
	assert.Nil(t, err)
	assert.NotNil(t, encoding)
	assert.Equal(t, mod.GPT_4.GetEncodingType().GetName(), encoding.GetName())
}

func getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4o(t *testing.T) {
	encoding, err := registry.GetEncodingForModel("gpt-4o-123")
	assert.Nil(t, err)
	assert.NotNil(t, encoding)
	assert.Equal(t, mod.GPT_4O.GetEncodingType().GetName(), encoding.GetName())
}

func getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4oMini(t *testing.T) {
	encoding, err := registry.GetEncodingForModel("gpt-4o-mini-123")
	assert.Nil(t, err)
	assert.NotNil(t, encoding)
	assert.Equal(t, mod.GPT_4O_MINI.GetEncodingType().GetName(), encoding.GetName())
}

func getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4Turbo(t *testing.T) {
	encoding, err := registry.GetEncodingForModel("gpt-4-turbo-123")
	assert.Nil(t, err)
	assert.NotNil(t, encoding)
	assert.Equal(t, mod.GPT_4_TURBO.GetEncodingType().GetName(), encoding.GetName())
}

func getEncodingForModelByPrefixReturnsCorrectEncodingForGpt3Turbo(t *testing.T) {
	encoding, err := registry.GetEncodingForModel("gpt-3.5-turbo-0301")
	assert.Nil(t, err)
	assert.NotNil(t, encoding)
	assert.Equal(t, mod.GPT_3_5_TURBO.GetEncodingType().GetName(), encoding.GetName())
}

func getEncodingForModelByPrefixReturnsCorrectEncodingForGpt3Turbo16k(t *testing.T) {
	encoding, err := registry.GetEncodingForModel("gpt-3.5-turbo-16k-0613")
	assert.Nil(t, err)
	assert.NotNil(t, encoding)
	assert.Equal(t, mod.GPT_3_5_TURBO_16K.GetEncodingType().GetName(), encoding.GetName())
}

func canRegisterCustomEncoding(t *testing.T) {
	encoding := &DummyEncoding{}

	registry.RegisterCustomEncoding(encoding)

	retrievedEncoding, err := registry.GetEncoding(encoding.GetName())
	assert.Nil(t, err)
	assert.NotNil(t, retrievedEncoding)
	assert.Equal(t, encoding, retrievedEncoding)
}

func canRegisterCustomGptBpe(t *testing.T) {
	params := mod.NewGptBytePairEncodingParams(
		"custom",
		regexp.MustCompile("test", regexp.None),
		map[string]int{},
		map[string]int{},
	)

	customRegistry, err := registry.RegisterGptBytePairEncoding(params)
	registry = customRegistry
	assert.Nil(t, err)
	assert.NotNil(t, registry)

	retrievedEncoding, err := registry.GetEncoding(params.GetName())
	assert.Nil(t, err)
	assert.NotNil(t, retrievedEncoding)
	assert.Equal(t, params.GetName(), retrievedEncoding.GetName())
}

func throwsIfCustomEncodingIsAlreadyRegistered(t *testing.T) {
	encoding := &DummyEncoding{}

	customRegistry, err := registry.RegisterCustomEncoding(encoding)
	registry = customRegistry
	assert.Nil(t, err)
	assert.NotNil(t, registry)

	_, err = registry.RegisterCustomEncoding(encoding)
	assert.NotNil(t, err)
}

func getEncodingReturnsEmptyOptionalForNonExistingEncodingName(t *testing.T) {
	result, _ := registry.GetEncoding("nonexistent")
	assert.Nil(t, result)
}

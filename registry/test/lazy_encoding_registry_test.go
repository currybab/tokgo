package registry_test

import (
	"testing"

	tokgo "github.com/currybab/tokgo/registry"
)

func setupLazyEncodingRegistry() {
	registry = tokgo.NewLazyEncodingRegistry()
}

func TestLazyGetEncodingReturnsCorrectEncoding(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingReturnsCorrectEncoding(t)
}

func TestLazyGetEncodingByNameReturnsCorrectEncoding(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingByNameReturnsCorrectEncoding(t)
}

func TestLazyGetEncodingForModelReturnsCorrectEncoding(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingForModelReturnsCorrectEncoding(t)
}

func TestLazyGetEncodingForModelByNameReturnsCorrectEncoding(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingForModelByNameReturnsCorrectEncoding(t)
}

func TestLazyGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt432k(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt432k(t)
}

func TestLazyGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt4(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4(t)
}

func TestLazyGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt4o(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4o(t)
}

func TestLazyGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt4oMini(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4oMini(t)
}

func TestLazyGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt4Turbo(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4Turbo(t)
}

func TestLazyGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt3Turbo(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt3Turbo(t)
}

func TestLazyGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt3Turbo16k(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt3Turbo16k(t)
}

func TestLazyCanRegisterCustomEncoding(t *testing.T) {
	setupLazyEncodingRegistry()
	canRegisterCustomEncoding(t)
}

func TestLazyCanRegisterCustomGptBpe(t *testing.T) {
	setupLazyEncodingRegistry()
	canRegisterCustomGptBpe(t)
}

func TestLazyThrowsIfCustomEncodingIsAlreadyRegistered(t *testing.T) {
	setupLazyEncodingRegistry()
	throwsIfCustomEncodingIsAlreadyRegistered(t)
}

func TestLazyGetEncodingReturnsEmptyOptionalForNonExistingEncodingName(t *testing.T) {
	setupLazyEncodingRegistry()
	getEncodingReturnsEmptyOptionalForNonExistingEncodingName(t)
}

package registry_test

import (
	"testing"

	tokgo "github.com/currybab/tokgo/registry"
)

func setupDefaultEncodingRegistry() {
	registry = tokgo.NewDefaultEncodingRegistry()
}

func TestGetEncodingReturnsCorrectEncoding(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingReturnsCorrectEncoding(t)
}

func TestGetEncodingByNameReturnsCorrectEncoding(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingByNameReturnsCorrectEncoding(t)
}

func TestGetEncodingForModelReturnsCorrectEncoding(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingForModelReturnsCorrectEncoding(t)
}

func TestGetEncodingForModelByNameReturnsCorrectEncoding(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingForModelByNameReturnsCorrectEncoding(t)
}

func TestGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt432k(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt432k(t)
}

func TestGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt4(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4(t)
}

func TestGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt4o(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4o(t)
}

func TestGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt4oMini(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4oMini(t)
}

func TestGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt4Turbo(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt4Turbo(t)
}

func TestGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt3Turbo(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt3Turbo(t)
}

func TestGetEncodingForModelByPrefixReturnsCorrectEncodingForGpt3Turbo16k(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingForModelByPrefixReturnsCorrectEncodingForGpt3Turbo16k(t)
}

func TestCanRegisterCustomEncoding(t *testing.T) {
	setupDefaultEncodingRegistry()
	canRegisterCustomEncoding(t)
}

func TestCanRegisterCustomGptBpe(t *testing.T) {
	setupDefaultEncodingRegistry()
	canRegisterCustomGptBpe(t)
}

func TestThrowsIfCustomEncodingIsAlreadyRegistered(t *testing.T) {
	setupDefaultEncodingRegistry()
	throwsIfCustomEncodingIsAlreadyRegistered(t)
}

func TestGetEncodingReturnsEmptyOptionalForNonExistingEncodingName(t *testing.T) {
	setupDefaultEncodingRegistry()
	getEncodingReturnsEmptyOptionalForNonExistingEncodingName(t)
}

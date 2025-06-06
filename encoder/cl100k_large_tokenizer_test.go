package encoder_test

import (
	"os"
	"testing"

	"github.com/currybab/tokgo"
	"github.com/currybab/tokgo/encoding"
)

// Environment Setup
func setup(t *testing.T) (string, bool) {
	originalThresholdValue, originalWasSet := os.LookupEnv(tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY)
	setEnvErr := os.Setenv(tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, "0")
	if setEnvErr != nil {
		t.Fatalf("Failed to set environment variable '%s' to '0': %v", tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, setEnvErr)
	}
	ENCODING = encoding.Cl100kBase()
	return originalThresholdValue, originalWasSet
}

func teardown(t *testing.T, originalThresholdValue string, originalWasSet bool) {
	if !originalWasSet {
		unsetEnvErr := os.Unsetenv(tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY)
		if unsetEnvErr != nil {
			t.Logf("Error trying to unset environment variable '%s': %v", tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, unsetEnvErr)
		}
	} else {
		restoreEnvErr := os.Setenv(tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, originalThresholdValue)
		if restoreEnvErr != nil {
			t.Logf("Failed to restore environment variable '%s' to '%s': %v", tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, originalThresholdValue, restoreEnvErr)
		}
	}

	ENCODING = encoding.Cl100kBase()
}

func TestLargeCl100kMeasureEncodingSpeeds(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	measureEncodingSpeeds(t)
}

func TestLargeCl100kEdgeCaseRoundTrips(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	cl100kEdgeCaseRoundTrips(t)
}

func TestLargeCl100kEncodeRoundTripWithRandomStrings(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	cl100kEncodeRoundTripWithRandomStrings(t)
}

func TestLargeCl100kEncodeOrdinaryRoundTripWithRandomStrings(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	cl100kEncodeOrdinaryRoundTripWithRandomStrings(t)
}

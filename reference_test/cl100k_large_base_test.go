package referencetest

import (
	"os"
	"strings"
	"testing"

	"github.com/currybab/tokgo/encoding"
	tokmod "github.com/currybab/tokgo/mod"
	"github.com/stretchr/testify/assert"
)

var CL100K_LARGE_BASE_ENCODING tokmod.Encoding

func setup(t *testing.T) (string, bool) {
	originalThresholdValue, originalWasSet := os.LookupEnv(tokmod.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY)
	setEnvErr := os.Setenv(tokmod.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, "0")
	if setEnvErr != nil {
		t.Fatalf("Failed to set environment variable '%s' to '0': %v", tokmod.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, setEnvErr)
	}
	CL100K_LARGE_BASE_ENCODING = encoding.Cl100kBase()
	return originalThresholdValue, originalWasSet
}

func teardown(t *testing.T, originalThresholdValue string, originalWasSet bool) {
	if !originalWasSet {
		unsetEnvErr := os.Unsetenv(tokmod.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY)
		if unsetEnvErr != nil {
			t.Logf("Error trying to unset environment variable '%s': %v", tokmod.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, unsetEnvErr)
		}
	} else {
		restoreEnvErr := os.Setenv(tokmod.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, originalThresholdValue)
		if restoreEnvErr != nil {
			t.Logf("Failed to restore environment variable '%s' to '%s': %v", tokmod.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, originalThresholdValue, restoreEnvErr)
		}
	}
}

func TestCL100kLargeBaseEncodesCorrectly(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, output string, _ string) {
		expected := parseEncodingString(output)
		actual := CL100K_LARGE_BASE_ENCODING.EncodeOrdinaryToIntArray(input)
		assert.Equal(t, expected, actual)
	})
}

func TestCL100kLargeBaseEncodesStable(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := CL100K_LARGE_BASE_ENCODING.Decode(CL100K_LARGE_BASE_ENCODING.EncodeToIntArray(input))
		assert.Equal(t, input, actual)
	})
}

func TestCL100kLargeBaseEncodesCorrectlyWithMaxTokensSet(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, output string, outputMaxTokens10 string) {
		expected := parseEncodingString(output)
		expectedWithMaxTokens := parseEncodingString(outputMaxTokens10)
		encodingResult := CL100K_LARGE_BASE_ENCODING.Encode(input, 10)

		assert.Equal(t, expectedWithMaxTokens, encodingResult.GetTokens())
		assert.Equal(t, len(expected) > len(expectedWithMaxTokens), encodingResult.IsTruncated())
	})
}

func TestCL100kLargeBaseEncodesStableWithMaxTokensSet(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := CL100K_LARGE_BASE_ENCODING.Decode(CL100K_LARGE_BASE_ENCODING.Encode(input, 10).GetTokens())

		assert.True(t, strings.HasPrefix(input, actual))
	})
}

func TestCL100kLargeBaseEncodeOrdinaryEncodesCorrectly(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, output string, _ string) {
		expected := parseEncodingString(output)
		actual := CL100K_LARGE_BASE_ENCODING.EncodeOrdinaryToIntArray(input)
		assert.Equal(t, expected, actual)
	})
}

func TestCL100kLargeBaseEncodeOrdinaryEncodesCorrectlyWithMaxTokensSet(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, output string, outputMaxTokens10 string) {
		expected := parseEncodingString(output)
		expectedWithMaxTokens := parseEncodingString(outputMaxTokens10)
		encodingResult := CL100K_LARGE_BASE_ENCODING.EncodeOrdinary(input, 10)

		assert.Equal(t, expectedWithMaxTokens, encodingResult.GetTokens())
		assert.Equal(t, len(expected) > len(expectedWithMaxTokens), encodingResult.IsTruncated())
	})
}

func TestCL100kLargeBaseEncodeOrdinaryEncodesStable(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := CL100K_LARGE_BASE_ENCODING.Decode(CL100K_LARGE_BASE_ENCODING.EncodeOrdinaryToIntArray(input))
		assert.Equal(t, input, actual)
	})
}

func TestCL100kLargeBaseEncodeOrdinaryEncodesStableWithMaxTokensSet(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := CL100K_LARGE_BASE_ENCODING.Decode(CL100K_LARGE_BASE_ENCODING.Encode(input, 10).GetTokens())
		assert.True(t, strings.HasPrefix(input, actual))
	})
}

func TestCL100kLargeBaseEncodeOrdinaryEncodesSpecialTokensCorrectly(t *testing.T) {
	originalThresholdValue, originalWasSet := setup(t)
	defer teardown(t, originalThresholdValue, originalWasSet)
	input := "Hello<|endoftext|>, <|fim_prefix|> <|fim_middle|> world <|fim_suffix|> ! <|endofprompt|>"
	actual := CL100K_LARGE_BASE_ENCODING.Decode(CL100K_LARGE_BASE_ENCODING.EncodeOrdinaryToIntArray(input))

	assert.Equal(t, input, actual)
}

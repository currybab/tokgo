package referencetest

import (
	"strings"
	"testing"

	tokmod "github.com/currybab/tokgo/mod"
	tokgo "github.com/currybab/tokgo/registry"
	"github.com/stretchr/testify/assert"
)

var O200K_BASE_ENCODING, _ = tokgo.NewDefaultEncodingRegistry().GetEncodingByType(tokmod.O200K_BASE)

func TestO200kBaseEncodesCorrectly(t *testing.T) {
	WrapTest(t, "../resources/test/o200k_base_encodings.csv", func(input string, output string, _ string) {
		expected := parseEncodingString(output)
		actual := O200K_BASE_ENCODING.EncodeOrdinaryToIntArray(input)
		assert.Equal(t, expected, actual)
	})
}

func TestO200kBaseEncodesStable(t *testing.T) {
	WrapTest(t, "../resources/test/o200k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := O200K_BASE_ENCODING.Decode(O200K_BASE_ENCODING.EncodeToIntArray(input))
		assert.Equal(t, input, actual)
	})
}

func TestO200kBaseEncodesCorrectlyWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/o200k_base_encodings.csv", func(input string, output string, outputMaxTokens10 string) {
		expected := parseEncodingString(output)
		expectedWithMaxTokens := parseEncodingString(outputMaxTokens10)
		encodingResult := O200K_BASE_ENCODING.Encode(input, 10)

		assert.Equal(t, expectedWithMaxTokens, encodingResult.GetTokens())
		assert.Equal(t, len(expected) > len(expectedWithMaxTokens), encodingResult.IsTruncated())
	})
}

func TestO200kBaseEncodesStableWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/o200k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := O200K_BASE_ENCODING.Decode(O200K_BASE_ENCODING.Encode(input, 10).GetTokens())

		assert.True(t, strings.HasPrefix(input, actual))
	})
}

func TestO200kBaseEncodeOrdinaryEncodesCorrectly(t *testing.T) {
	WrapTest(t, "../resources/test/o200k_base_encodings.csv", func(input string, output string, _ string) {
		expected := parseEncodingString(output)
		actual := O200K_BASE_ENCODING.EncodeOrdinaryToIntArray(input)
		assert.Equal(t, expected, actual)
	})
}

func TestO200kBaseEncodeOrdinaryEncodesCorrectlyWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/o200k_base_encodings.csv", func(input string, output string, outputMaxTokens10 string) {
		expected := parseEncodingString(output)
		expectedWithMaxTokens := parseEncodingString(outputMaxTokens10)
		encodingResult := O200K_BASE_ENCODING.EncodeOrdinary(input, 10)

		assert.Equal(t, expectedWithMaxTokens, encodingResult.GetTokens())
		assert.Equal(t, len(expected) > len(expectedWithMaxTokens), encodingResult.IsTruncated())
	})
}

func TestO200kBaseEncodeOrdinaryEncodesStable(t *testing.T) {
	WrapTest(t, "../resources/test/o200k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := O200K_BASE_ENCODING.Decode(O200K_BASE_ENCODING.EncodeOrdinaryToIntArray(input))
		assert.Equal(t, input, actual)
	})
}

func TestO200kBaseEncodeOrdinaryEncodesStableWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/o200k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := O200K_BASE_ENCODING.Decode(O200K_BASE_ENCODING.Encode(input, 10).GetTokens())
		assert.True(t, strings.HasPrefix(input, actual))
	})
}

func TestO200kBaseEncodeOrdinaryEncodesSpecialTokensCorrectly(t *testing.T) {
	input := "Hello<|endoftext|>, <|fim_prefix|> <|fim_middle|> world <|fim_suffix|> ! <|endofprompt|>"
	actual := O200K_BASE_ENCODING.Decode(O200K_BASE_ENCODING.EncodeOrdinaryToIntArray(input))

	assert.Equal(t, input, actual)
}

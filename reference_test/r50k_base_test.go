package referencetest

import (
	"strings"
	"testing"

	tokmod "github.com/currybab/tokgo/mod"
	tokgo "github.com/currybab/tokgo/registry"
	"github.com/stretchr/testify/assert"
)

var R50K_BASE_ENCODING, _ = tokgo.NewDefaultEncodingRegistry().GetEncodingByType(tokmod.R50K_BASE)

func TestR50kBaseEncodesCorrectly(t *testing.T) {
	WrapTest(t, "../resources/test/r50k_base_encodings.csv", func(input string, output string, _ string) {
		expected := parseEncodingString(output)
		actual := R50K_BASE_ENCODING.EncodeOrdinaryToIntArray(input)
		assert.Equal(t, expected, actual)
	})
}

func TestR50kBaseBaseEncodesStable(t *testing.T) {
	WrapTest(t, "../resources/test/r50k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := R50K_BASE_ENCODING.Decode(R50K_BASE_ENCODING.EncodeToIntArray(input))
		assert.Equal(t, input, actual)
	})
}

func TestR50kBaseBaseEncodesCorrectlyWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/r50k_base_encodings.csv", func(input string, output string, outputMaxTokens10 string) {
		expected := parseEncodingString(output)
		expectedWithMaxTokens := parseEncodingString(outputMaxTokens10)
		encodingResult := R50K_BASE_ENCODING.Encode(input, 10)

		assert.Equal(t, expectedWithMaxTokens, encodingResult.GetTokens())
		assert.Equal(t, len(expected) > len(expectedWithMaxTokens), encodingResult.IsTruncated())
	})
}

func TestR50kBaseBaseEncodesStableWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/r50k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := R50K_BASE_ENCODING.Decode(R50K_BASE_ENCODING.Encode(input, 10).GetTokens())

		assert.True(t, strings.HasPrefix(input, actual))
	})
}

func TestR50kBaseBaseEncodeOrdinaryEncodesCorrectly(t *testing.T) {
	WrapTest(t, "../resources/test/r50k_base_encodings.csv", func(input string, output string, _ string) {
		expected := parseEncodingString(output)
		actual := R50K_BASE_ENCODING.EncodeOrdinaryToIntArray(input)
		assert.Equal(t, expected, actual)
	})
}

func TestR50kBaseBaseEncodeOrdinaryEncodesCorrectlyWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/r50k_base_encodings.csv", func(input string, output string, outputMaxTokens10 string) {
		expected := parseEncodingString(output)
		expectedWithMaxTokens := parseEncodingString(outputMaxTokens10)
		encodingResult := R50K_BASE_ENCODING.EncodeOrdinary(input, 10)

		assert.Equal(t, expectedWithMaxTokens, encodingResult.GetTokens())
		assert.Equal(t, len(expected) > len(expectedWithMaxTokens), encodingResult.IsTruncated())
	})
}

func TestR50kBaseBaseEncodeOrdinaryEncodesStable(t *testing.T) {
	WrapTest(t, "../resources/test/r50k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := R50K_BASE_ENCODING.Decode(R50K_BASE_ENCODING.EncodeOrdinaryToIntArray(input))
		assert.Equal(t, input, actual)
	})
}

func TestR50kBaseBaseEncodeOrdinaryEncodesStableWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/r50k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := R50K_BASE_ENCODING.Decode(R50K_BASE_ENCODING.Encode(input, 10).GetTokens())
		assert.True(t, strings.HasPrefix(input, actual))
	})
}

func TestR50kBaseBaseEncodeOrdinaryEncodesSpecialTokensCorrectly(t *testing.T) {
	input := "Hello<|endoftext|>, <|fim_prefix|> <|fim_middle|> world <|fim_suffix|> ! <|endofprompt|>"
	actual := R50K_BASE_ENCODING.Decode(R50K_BASE_ENCODING.EncodeOrdinaryToIntArray(input))

	assert.Equal(t, input, actual)
}

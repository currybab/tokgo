package referencetest

import (
	"strings"
	"testing"

	tokmod "github.com/currybab/tokgo/mod"
	tokgo "github.com/currybab/tokgo/registry"
	"github.com/stretchr/testify/assert"
)

var P50K_EDIT_ENCODING, _ = tokgo.NewDefaultEncodingRegistry().GetEncodingByType(tokmod.P50K_EDIT)

func TestP50kEditEncodesCorrectly(t *testing.T) {
	WrapTest(t, "../resources/test/p50k_edit_encodings.csv", func(input string, output string, _ string) {
		expected := parseEncodingString(output)
		actual := P50K_EDIT_ENCODING.EncodeOrdinaryToIntArray(input)
		assert.Equal(t, expected, actual)
	})
}

func TestP50kEditEncodesStable(t *testing.T) {
	WrapTest(t, "../resources/test/p50k_edit_encodings.csv", func(input string, _ string, _ string) {
		actual := P50K_EDIT_ENCODING.Decode(P50K_EDIT_ENCODING.EncodeToIntArray(input))
		assert.Equal(t, input, actual)
	})
}

func TestP50kEditEncodesCorrectlyWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/p50k_edit_encodings.csv", func(input string, output string, outputMaxTokens10 string) {
		expected := parseEncodingString(output)
		expectedWithMaxTokens := parseEncodingString(outputMaxTokens10)
		encodingResult := P50K_EDIT_ENCODING.Encode(input, 10)

		assert.Equal(t, expectedWithMaxTokens, encodingResult.GetTokens())
		assert.Equal(t, len(expected) > len(expectedWithMaxTokens), encodingResult.IsTruncated())
	})
}

func TestP50kEditEncodesStableWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/p50k_edit_encodings.csv", func(input string, _ string, _ string) {
		actual := P50K_EDIT_ENCODING.Decode(P50K_EDIT_ENCODING.Encode(input, 10).GetTokens())

		assert.True(t, strings.HasPrefix(input, actual))
	})
}

func TestP50kEditEncodeOrdinaryEncodesCorrectly(t *testing.T) {
	WrapTest(t, "../resources/test/p50k_edit_encodings.csv", func(input string, output string, _ string) {
		expected := parseEncodingString(output)
		actual := P50K_EDIT_ENCODING.EncodeOrdinaryToIntArray(input)
		assert.Equal(t, expected, actual)
	})
}

func TestP50kEditEncodeOrdinaryEncodesCorrectlyWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/p50k_edit_encodings.csv", func(input string, output string, outputMaxTokens10 string) {
		expected := parseEncodingString(output)
		expectedWithMaxTokens := parseEncodingString(outputMaxTokens10)
		encodingResult := P50K_EDIT_ENCODING.EncodeOrdinary(input, 10)

		assert.Equal(t, expectedWithMaxTokens, encodingResult.GetTokens())
		assert.Equal(t, len(expected) > len(expectedWithMaxTokens), encodingResult.IsTruncated())
	})
}

func TestP50kEditEncodeOrdinaryEncodesStable(t *testing.T) {
	WrapTest(t, "../resources/test/p50k_edit_encodings.csv", func(input string, _ string, _ string) {
		actual := P50K_EDIT_ENCODING.Decode(P50K_EDIT_ENCODING.EncodeOrdinaryToIntArray(input))
		assert.Equal(t, input, actual)
	})
}

func TestP50kEditEncodeOrdinaryEncodesStableWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/p50k_edit_encodings.csv", func(input string, _ string, _ string) {
		actual := P50K_EDIT_ENCODING.Decode(P50K_EDIT_ENCODING.Encode(input, 10).GetTokens())
		assert.True(t, strings.HasPrefix(input, actual))
	})
}

func TestP50kEditEncodeOrdinaryEncodesSpecialTokensCorrectly(t *testing.T) {
	input := "Hello<|endoftext|>, <|fim_prefix|> <|fim_middle|> world <|fim_suffix|> ! <|endofprompt|>"
	actual := P50K_EDIT_ENCODING.Decode(P50K_EDIT_ENCODING.EncodeOrdinaryToIntArray(input))

	assert.Equal(t, input, actual)
}

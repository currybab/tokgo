package referencetest

import (
	"strings"
	"testing"

	tokmod "github.com/currybab/tokgo/mod"
	tokgo "github.com/currybab/tokgo/registry"
	"github.com/stretchr/testify/assert"
)

var CL100K_BASE_ENCODING, _ = tokgo.NewDefaultEncodingRegistry().GetEncodingByType(tokmod.CL100K_BASE)

func TestCL100kBaseEncodesCorrectly(t *testing.T) {
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, output string, _ string) {
		expected := parseEncodingString(output)
		actual := CL100K_BASE_ENCODING.EncodeOrdinaryToIntArray(input)
		assert.Equal(t, expected, actual)
	})
}

func TestCL100kBaseEncodesStable(t *testing.T) {
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := CL100K_BASE_ENCODING.Decode(CL100K_BASE_ENCODING.EncodeToIntArray(input))
		assert.Equal(t, input, actual)
	})
}

func TestCL100kBaseEncodesCorrectlyWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, output string, outputMaxTokens10 string) {
		expected := parseEncodingString(output)
		expectedWithMaxTokens := parseEncodingString(outputMaxTokens10)
		encodingResult := CL100K_BASE_ENCODING.Encode(input, 10)

		assert.Equal(t, expectedWithMaxTokens, encodingResult.GetTokens())
		assert.Equal(t, len(expected) > len(expectedWithMaxTokens), encodingResult.IsTruncated())
	})
}

func TestCL100kBaseEncodesStableWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := CL100K_BASE_ENCODING.Decode(CL100K_BASE_ENCODING.Encode(input, 10).GetTokens())

		assert.True(t, strings.HasPrefix(input, actual))
	})
}

func TestCL100kBaseEncodeOrdinaryEncodesCorrectly(t *testing.T) {
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, output string, _ string) {
		expected := parseEncodingString(output)
		actual := CL100K_BASE_ENCODING.EncodeOrdinaryToIntArray(input)
		assert.Equal(t, expected, actual)
	})
}

func TestCL100kBaseEncodeOrdinaryEncodesCorrectlyWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, output string, outputMaxTokens10 string) {
		expected := parseEncodingString(output)
		expectedWithMaxTokens := parseEncodingString(outputMaxTokens10)
		encodingResult := CL100K_BASE_ENCODING.EncodeOrdinary(input, 10)

		assert.Equal(t, expectedWithMaxTokens, encodingResult.GetTokens())
		assert.Equal(t, len(expected) > len(expectedWithMaxTokens), encodingResult.IsTruncated())
	})
}

func TestCL100kBaseEncodeOrdinaryEncodesStable(t *testing.T) {
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := CL100K_BASE_ENCODING.Decode(CL100K_BASE_ENCODING.EncodeOrdinaryToIntArray(input))
		assert.Equal(t, input, actual)
	})
}

func TestCL100kBaseEncodeOrdinaryEncodesStableWithMaxTokensSet(t *testing.T) {
	WrapTest(t, "../resources/test/cl100k_base_encodings.csv", func(input string, _ string, _ string) {
		actual := CL100K_BASE_ENCODING.Decode(CL100K_BASE_ENCODING.Encode(input, 10).GetTokens())
		assert.True(t, strings.HasPrefix(input, actual))
	})
}

func TestCL100kBaseEncodeOrdinaryEncodesSpecialTokensCorrectly(t *testing.T) {
	input := "Hello<|endoftext|>, <|fim_prefix|> <|fim_middle|> world <|fim_suffix|> ! <|endofprompt|>"
	actual := CL100K_BASE_ENCODING.Decode(CL100K_BASE_ENCODING.EncodeOrdinaryToIntArray(input))

	assert.Equal(t, input, actual)
}

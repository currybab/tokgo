package tokgo

var VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY = "VERY_LARGE_TOKENIZER_BYTE_THRESHOLD"

type Encoding interface {
	Encode(text string) *EncodingResult
	EncodeOrdinaryToIntArray(text string) []int
	EncodeOrdinaryToEncodingResult(text string, maxTokens int) *EncodingResult
	CountTokens(text string) int
	CountTokensOrdinary(text string) int
	Decode(tokens []int) string
	DecodeBytes(tokens []int) []byte
	GetName() string
}

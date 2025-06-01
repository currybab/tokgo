package tokgo

var VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY = "VERY_LARGE_TOKENIZER_BYTE_THRESHOLD"

type Encoding interface {
	EncodeToIntArray(text string) []int
	Encode(text string, maxTokens int) *EncodingResult
	EncodeOrdinaryToIntArray(text string) []int
	EncodeOrdinary(text string, maxTokens int) *EncodingResult
	CountTokens(text string) int
	CountTokensOrdinary(text string) int
	Decode(tokens []int) string
	DecodeBytes(tokens []int) []byte
	GetName() string
}

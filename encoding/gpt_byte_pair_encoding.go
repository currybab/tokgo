package encoding

import (
	"fmt"
	"math"
	"strings"

	regexp "github.com/dlclark/regexp2"

	"github.com/currybab/tokgo"
	"github.com/currybab/tokgo/encoder"
	"github.com/currybab/tokgo/parser"
)

type internalResult struct {
	tokens                      []int
	truncated                   bool
	tokenCount                  int
	lastProcessedCharacterIndex int // -1 == text was null or string was empty()
}

func newInternalResult(tokens []int, tokenCount int, truncated bool, lastProcessedCharacterIndex int) *internalResult {
	if tokenCount < 0 {
		tokenCount = len(tokens)
	}
	return &internalResult{
		tokens:                      tokens,
		truncated:                   truncated,
		tokenCount:                  tokenCount,
		lastProcessedCharacterIndex: lastProcessedCharacterIndex,
	}
}

func (i *internalResult) ToEncodingResult() *tokgo.EncodingResult {
	if len(i.tokens) != i.tokenCount {
		panic(fmt.Sprintf("Token count does not match token list size (tokenCount=%v, tokens size=%v)", i.tokenCount, len(i.tokens)))
	}

	return tokgo.NewEncodingResult(i.tokens, i.truncated, i.lastProcessedCharacterIndex)
}

func (i *internalResult) ToTokenCount() int {
	return i.tokenCount
}

type GptBytePairEncoding struct {
	Encoder        *encoder.TokenEncoder
	name           string
	pattern        *regexp.Regexp
	specialEncoder *encoder.SpecialEncoder
}

func NewGptBytePairEncoding(params *tokgo.GptBytePairEncodingParams) *GptBytePairEncoding {
	return &GptBytePairEncoding{
		name:           params.GetName(),
		pattern:        params.GetPattern(),
		Encoder:        encoder.NewTokenEncoder(params.GetEncoder()),
		specialEncoder: encoder.NewSpecialEncoder(params.GetSpecialTokensEncoder()),
	}
}

func (e *GptBytePairEncoding) encodeInternal(text string, maxTokenCount int, keepEncodings bool) *internalResult {
	if text == "" {
		return newInternalResult([]int{}, -1, false, -1)
	}

	e.specialEncoder.CheckForSpecialTokens(text)

	return e.encodeOrdinaryInternal(text, maxTokenCount, keepEncodings)
}

func (e *GptBytePairEncoding) encodeOrdinaryInternal(text string, maxTokenCount int, keepEncodings bool) *internalResult {
	if text == "" {
		return newInternalResult([]int{}, -1, false, -1)
	}

	out := make([]int, 0)
	tokenCount := e.encodeOrdinaryInternalToInt(text, maxTokenCount, keepEncodings, &out)

	if keepEncodings && maxTokenCount != math.MaxInt {
		// Make sure we didn't break the multibyte character
		for tokensToRemove := 0; tokensToRemove <= len(out); tokensToRemove++ {
			size := len(out) - tokensToRemove
			tokens := make([]int, size)
			for i := 0; i < size; i++ {
				tokens[i] = out[i]
			}
			decoded := e.Decode(tokens)
			if strings.HasPrefix(text, decoded) {
				// If decoded text is equal to the head of the original text, we can safely return the tokens
				return newInternalResult(tokens, -1, len(text) > len(decoded), len(decoded)-1)
			}
		}
	}

	return newInternalResult(out, tokenCount, false, len(text)-1)
}

func (e *GptBytePairEncoding) encodeOrdinaryInternalToInt(text string, maxTokenCount int, keepEncodings bool, out *[]int) int {
	if e.pattern == nil {
		// if cl100k
		tokenCount := []int{0}
		ranks := make([]int, 0)
		parser.Split(text, func(utf8BytesList []byte) bool {
			tokenCount[0] += e.Encoder.AddTokensAndGetCount(maxTokenCount, keepEncodings, utf8BytesList, out, &ranks)
			return tokenCount[0] >= maxTokenCount
		})
		return tokenCount[0]
	}

	tokenCount := 0
	ranks := make([]int, 0, 10)

	match, _ := e.pattern.FindStringMatch(text)
	for tokenCount < maxTokenCount && match != nil {
		bytes := match.Group.String()
		tokenCount += e.Encoder.AddTokensAndGetCount(maxTokenCount, keepEncodings, []byte(bytes), out, &ranks)
		match, _ = e.pattern.FindNextMatch(match)
	}
	return tokenCount
}

func (e *GptBytePairEncoding) EncodeToIntArray(text string) []int {
	return e.Encode(text, math.MaxInt).GetTokens()
}

func (e *GptBytePairEncoding) Encode(text string, maxTokens int) *tokgo.EncodingResult {
	return e.encodeInternal(text, maxTokens, true).ToEncodingResult()
}

func (e *GptBytePairEncoding) EncodeOrdinaryToIntArray(text string) []int {
	return e.EncodeOrdinary(text, math.MaxInt).GetTokens()
}

func (e *GptBytePairEncoding) EncodeOrdinary(text string, maxTokens int) *tokgo.EncodingResult {
	return e.encodeOrdinaryInternal(text, maxTokens, true).ToEncodingResult()
}

func (e *GptBytePairEncoding) CountTokens(text string) int {
	return e.encodeInternal(text, math.MaxInt, false).ToTokenCount()
}

func (e *GptBytePairEncoding) CountTokensOrdinary(text string) int {
	return e.encodeOrdinaryInternal(text, math.MaxInt, false).ToTokenCount()
}

func (e *GptBytePairEncoding) Decode(tokens []int) string {
	return string(e.DecodeBytes(tokens))
}

func (e *GptBytePairEncoding) DecodeBytes(tokens []int) []byte {
	out := make([]byte, 0, 10*len(tokens))
	for i := 0; i < len(tokens); i++ {
		decodedToken := e.Encoder.DecodeToken(tokens[i], e.specialEncoder)
		out = append(out, decodedToken...)
	}
	return out
}

func (e *GptBytePairEncoding) GetName() string {
	return e.name
}

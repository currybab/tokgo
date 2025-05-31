package tokgo

import (
	"fmt"
)

// The result of encoding operation.
type EncodingResult struct {
	tokens                      []int
	truncated                   bool
	lastProcessedCharacterIndex int
}

func NewEncodingResult(tokens []int, truncated bool, lastProcessedCharacterIndex ...int) *EncodingResult {
	lastIndex := -1
	if len(lastProcessedCharacterIndex) > 0 {
		lastIndex = lastProcessedCharacterIndex[0]
	}
	return &EncodingResult{
		tokens:                      tokens,
		truncated:                   truncated,
		lastProcessedCharacterIndex: lastIndex,
	}
}

func (er *EncodingResult) GetTokens() []int {
	return er.tokens
}

func (er *EncodingResult) IsTruncated() bool {
	return er.truncated
}

func (er *EncodingResult) GetLastProcessedCharacterIndex() int {
	return er.lastProcessedCharacterIndex
}

func (er *EncodingResult) ToString() string {
	return fmt.Sprintf("EncodingResult{tokens=%v, truncated=%v, lastProcessedCharacterIndex=%v}", er.tokens, er.truncated, er.lastProcessedCharacterIndex)
}

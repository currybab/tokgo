package tokgo

import (
	regexp "github.com/dlclark/regexp2"
)

type GptBytePairEncodingParams struct {
	name                 string
	pattern              *regexp.Regexp
	encoder              map[string]int
	specialTokensEncoder map[string]int
}

func NewGptBytePairEncodingParams(
	name string,
	pattern *regexp.Regexp,
	encoder map[string]int,
	specialTokensEncoder map[string]int,
) *GptBytePairEncodingParams {
	return &GptBytePairEncodingParams{
		name:                 name,
		pattern:              pattern,
		encoder:              encoder,
		specialTokensEncoder: specialTokensEncoder,
	}
}

func (g *GptBytePairEncodingParams) GetName() string {
	return g.name
}

func (g *GptBytePairEncodingParams) GetPattern() *regexp.Regexp {
	return g.pattern
}

func (g *GptBytePairEncodingParams) GetEncoder() map[string]int {
	return g.encoder
}

func (g *GptBytePairEncodingParams) GetSpecialTokensEncoder() map[string]int {
	return g.specialTokensEncoder
}

package encoder

import "strings"

const (
	SPECIAL_START = "<|"
	SPECIAL_END   = "|>"
)

type SpecialEncoder struct {
	encodedToDecoded map[int]string
}

func NewSpecialEncoder(encoder map[string]int) *SpecialEncoder {
	encodedToDecoded := make(map[int]string, len(encoder))
	for key, value := range encoder {
		if !strings.Contains(key, SPECIAL_START) || !strings.Contains(key, SPECIAL_END) {
			panic("Special tokens must contain " + SPECIAL_START + " and " + SPECIAL_END + " (but was " + key + ")")
		}

		encodedToDecoded[value] = key
	}

	return &SpecialEncoder{
		encodedToDecoded: encodedToDecoded,
	}
}

func (s *SpecialEncoder) DecodeIfPresent(encodedToken int) []byte {
	result, ok := s.encodedToDecoded[encodedToken]
	if ok {
		return []byte(result)
	}
	return nil
}

func (s *SpecialEncoder) CheckForSpecialTokens(text string) {
	if strings.Contains(text, SPECIAL_START) && strings.Contains(text, SPECIAL_END) {
		for _, specialToken := range s.encodedToDecoded {
			if strings.Contains(text, specialToken) {
				panic("Encoding special tokens is not supported.")
			}
		}
	}
}

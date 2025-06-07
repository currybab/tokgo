package encoder

import (
	"math"
	"os"
	"strconv"

	"github.com/currybab/tokgo/mod"
	"github.com/emirpasic/gods/v2/maps/treemap"
)

const (
	MAX_RANK   int = math.MaxInt32 - 1
	dummy_rank int = math.MaxInt32
)

type TokenEncoder struct {
	encoders                            []map[string]int
	decoder                             map[int][]byte
	VERY_LARGE_TOKENIZER_BYTE_THRESHOLD int
}

func NewTokenEncoder(encoder map[string]int) *TokenEncoder {
	if len(encoder) > 0 {
		thresholdKey := os.Getenv(mod.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY)
		if thresholdKey == "" {
			thresholdKey = "500"
		}
		VERY_LARGE_TOKENIZER_BYTE_THRESHOLD, _ := strconv.Atoi(thresholdKey)
		tempEncoders := treemap.New[int, map[string]int]()
		for k, v := range encoder {
			length := len(k)
			dict, ok := tempEncoders.Get(length)
			if !ok {
				dict = make(map[string]int)
				tempEncoders.Put(length, dict)
			}
			dict[k] = v
		}
		//noinspection unchecked
		keys := tempEncoders.Keys()
		encoders := make([]map[string]int, keys[len(keys)-1]+1)
		for _, k := range keys {
			v, _ := tempEncoders.Get(k)
			encoders[k] = v
		}

		decoder := make(map[int][]byte, len(encoder))
		for k, v := range encoder {
			decoder[v] = []byte(k)
		}
		return &TokenEncoder{
			encoders:                            encoders,
			decoder:                             decoder,
			VERY_LARGE_TOKENIZER_BYTE_THRESHOLD: VERY_LARGE_TOKENIZER_BYTE_THRESHOLD,
		}
	} else {
		//noinspection unchecked
		return &TokenEncoder{
			encoders: []map[string]int{},
			decoder:  map[int][]byte{},
		}
	}
}

func getMinRankIndex(ranks []int) int {
	minRankIndex := -1
	minRank := MAX_RANK

	i := 0
	length := len(ranks) - 3
	for ; i < length-2; i += 4 { // Unrolled loop
		{
			r := ranks[i]
			if r < minRank {
				minRankIndex = i
				minRank = r
			}
		}
		{
			r := ranks[i+1]
			if r < minRank {
				minRankIndex = i + 1
				minRank = r
			}
		}
		{
			r := ranks[i+2]
			if r < minRank {
				minRankIndex = i + 2
				minRank = r
			}
		}
		{
			r := ranks[i+3]
			if r < minRank {
				minRankIndex = i + 3
				minRank = r
			}
		}
	}

	for ; i <= length; i++ {
		r := ranks[i]
		if r < minRank {
			minRankIndex = i
			minRank = r
		}
	}

	return minRankIndex
}

func getNextIndex(ranks []int, nextIndex int) int {
	for nextIndex < len(ranks) && ranks[nextIndex] == dummy_rank {
		nextIndex++
	}
	return nextIndex
}

func getPreviousIndex(ranks []int, previousIndex int) int {
	for previousIndex >= 0 && ranks[previousIndex] == dummy_rank {
		previousIndex--
	}
	return previousIndex
}

func (t *TokenEncoder) AddTokensAndGetCount(maxTokenCount int, keepEncodings bool, byteArray []byte, out *[]int, ranks *[]int) int {
	match := byteArray
	encoded := t.encode(match)
	if encoded != MAX_RANK {
		if keepEncodings {
			*out = append(*out, encoded)
		}
		return 1
	} else {
		if len(match) < t.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD {
			return t.calculateTokensSmall(maxTokenCount, keepEncodings, out, ranks, match)
		} else {
			return CalculateTokensLarge(t, maxTokenCount, keepEncodings, out, match)
		}
	}
}

func (t *TokenEncoder) calculateTokensSmall(maxTokenCount int, keepEncodings bool, out *[]int, ranks *[]int, match []byte) int {
	length := len(match)
	if length <= 1 {
		panic("Already filtered out")
	}
	*ranks = make([]int, 0, length+1)

	minRankIndex := -1
	for i, minRank := 0, MAX_RANK; i < length+1; i++ {
		encoded := t.Encode(match, i, i+2)
		if encoded != MAX_RANK {
			if encoded < minRank {
				minRankIndex = i
				minRank = encoded
			}
		}
		*ranks = append(*ranks, encoded)
	}
	tokenCount := t.MergeBytesAndGetTokenCount(match, length, *ranks, minRankIndex)
	if keepEncodings {
		for start, end := 0, 1; end < len(*ranks) && len(*out) < maxTokenCount; end++ {
			if (*ranks)[end] != dummy_rank {
				token := t.Encode(match, start, end)
				if token != MAX_RANK {
					*out = append(*out, token)
					start = end
				}
			}
		}
	}
	return tokenCount
}

func (t *TokenEncoder) MergeBytesAndGetTokenCount(piece []byte, length int, ranks []int, minRankIndex int) int {
	if getMinRankIndex(ranks) != minRankIndex {
		panic("getMinRankIndex(ranks) != minRankIndex")
	}
	for minRankIndex >= 0 {
		previousIndex := getPreviousIndex(ranks, minRankIndex-1)
		nextIndex := getNextIndex(ranks, minRankIndex+1)
		nextNextIndex := getNextIndex(ranks, nextIndex+1)
		nextNextNextIndex := getNextIndex(ranks, nextNextIndex+1)

		if previousIndex >= 0 {
			if ranks[previousIndex] == dummy_rank {
				panic("ranks[previousIndex] == dummy_rank")
			}
			newRank := t.Encode(piece, previousIndex, nextNextIndex)
			ranks[previousIndex] = newRank
		}
		if ranks[minRankIndex] == dummy_rank {
			panic("ranks[minRankIndex] == dummy_rank")
		}
		newRank := t.Encode(piece, minRankIndex, nextNextNextIndex)
		ranks[minRankIndex] = newRank

		ranks[nextIndex] = dummy_rank

		length--
		if length < 3 {
			break // single tokens were already filtered out, let's skip a minimum calculation
		} else {
			minRankIndex = getMinRankIndex(ranks)
		}
	}
	if getMinRankIndex(ranks) >= 0 {
		panic("getMinRankIndex(ranks) >= 0")
	}
	return length
}

func (t *TokenEncoder) encode(payload []byte) int {
	if len(payload) < len(t.encoders) {
		encoder := t.encoders[len(payload)]
		if len(encoder) > 0 {
			result, ok := encoder[string(payload)]
			if ok {
				return result
			}
		}
	}
	return MAX_RANK
}

func (t *TokenEncoder) Encode(piece []byte, start int, end int) int {
	if end > len(piece) || end-start == len(piece) {
		return MAX_RANK
	} else {
		return t.encode(piece[start:end])
	}
}

func (t *TokenEncoder) DecodeToken(token int, specialEncodeer *SpecialEncoder) []byte {
	if decodeToken, ok := t.decoder[token]; ok {
		return decodeToken
	}
	return specialEncodeer.DecodeIfPresent(token)
}

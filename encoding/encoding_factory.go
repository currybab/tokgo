package encoding

import (
	"bufio"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	regexp "github.com/dlclark/regexp2"

	"github.com/currybab/tokgo/mod"
	"github.com/currybab/tokgo/parser"
)

// Special token constants
const (
	ENDOFTEXT   = "<|endoftext|>"
	FIM_PREFIX  = "<|fim_prefix|>"
	FIM_MIDDLE  = "<|fim_middle|>"
	FIM_SUFFIX  = "<|fim_suffix|>"
	ENDOFPROMPT = "<|endofprompt|>"
)

// Special token maps
var (
	SPECIAL_TOKENS_X50K_BASE = map[string]int{
		ENDOFTEXT: 50256,
	}
	SPECIAL_TOKENS_P50K_EDIT = map[string]int{
		ENDOFTEXT:  50256,
		FIM_PREFIX: 50281,
		FIM_MIDDLE: 50282,
		FIM_SUFFIX: 50283,
	}
	SPECIAL_TOKENS_CL100K_BASE = map[string]int{
		ENDOFTEXT:   100257,
		FIM_PREFIX:  100258,
		FIM_MIDDLE:  100259,
		FIM_SUFFIX:  100260,
		ENDOFPROMPT: 100276,
	}
	SPECIAL_TOKENS_O200K_BASE = map[string]int{
		ENDOFTEXT:   199999,
		ENDOFPROMPT: 200018,
	}
)

func R50kBase() mod.Encoding {
	return from50kParameters(
		"r50k_base",
		"r50k_base.tiktoken",
		SPECIAL_TOKENS_X50K_BASE,
	)
}

func P50kBase() mod.Encoding {
	return from50kParameters(
		"p50k_base",
		"p50k_base.tiktoken",
		SPECIAL_TOKENS_X50K_BASE,
	)
}

func P50kEdit() mod.Encoding {
	return from50kParameters(
		"p50k_edit",
		"p50k_base.tiktoken",
		SPECIAL_TOKENS_P50K_EDIT,
	)
}

func Cl100kBase() mod.Encoding {
	mergeableRanks, err := LoadMergeableRanks("cl100k_base.tiktoken")
	if err != nil {
		panic(err)
	}
	// regex, err := regexp.Compile("'(?:[sdmt]|ll|ve|re)|[^\r\n\\p{L}\\p{N}]?+\\p{L}+|\\p{N}{1,3}| ?[^\\s\\p{L}\\p{N}]++[\r\n]*|\\s*[\r\n]|\\s+(?!\\S)|\\s+", regexp.None)
	// if err != nil {
	// 	panic(err)
	// }
	params := mod.NewGptBytePairEncodingParams(
		"cl100k_base",
		nil,
		mergeableRanks,
		SPECIAL_TOKENS_CL100K_BASE,
	)
	return NewCl100kGptBytePairEncoding(params)
}

func O200kBase() mod.Encoding {
	mergeableRanks, err := LoadMergeableRanks("o200k_base.tiktoken")
	if err != nil {
		panic(err)
	}
	patterns := []string{
		`[^\r\n\p{L}\p{N}]?[\p{Lu}\p{Lt}\p{Lm}\p{Lo}\p{M}]*[\p{Ll}\p{Lm}\p{Lo}\p{M}]+(?i:'s|'t|'re|'ve|'m|'ll|'d)?`,
		`[^\r\n\p{L}\p{N}]?[\p{Lu}\p{Lt}\p{Lm}\p{Lo}\p{M}]+[\p{Ll}\p{Lm}\p{Lo}\p{M}]*(?i:'s|'t|'re|'ve|'m|'ll|'d)?`,
		`\p{N}{1,3}`,
		` ?[^\s\p{L}\p{N}]+[\r\n/]*`,
		`\s*[\r\n]+`,
		`\s+(?!\S)`,
		`\s+`,
	}
	regex, err := regexp.Compile(strings.Join(patterns, "|"), regexp.None)
	if err != nil {
		panic(err)
	}
	params := mod.NewGptBytePairEncodingParams(
		"o200k_base",
		regex,
		mergeableRanks,
		SPECIAL_TOKENS_O200K_BASE,
	)
	return FromParameters(params)
}

func from50kParameters(name, fileName string, specialTokens map[string]int) mod.Encoding {
	regex, err := regexp.Compile(`'(?:[sdmt]|ll|ve|re)| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`, regexp.None)
	if err != nil {
		panic(err)
	}
	mergeableRanks, err := LoadMergeableRanks(fileName)
	if err != nil {
		panic(err)
	}
	params := mod.NewGptBytePairEncodingParams(
		name,
		regex,
		mergeableRanks,
		specialTokens,
	)
	return FromParameters(params)
}

// getResourcePath returns the path to a resource file relative to this source file
func getResourcePath(fileName string) string {
	_, currentFilePath, _, _ := runtime.Caller(1)
	// Go up two directories: from factory to project root
	baseDir := filepath.Dir(filepath.Dir(currentFilePath))
	return filepath.Join(baseDir, "resources", fileName)
}

func LoadMergeableRanks(fileName string) (map[string]int, error) {
	file, err := os.Open(getResourcePath(fileName))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := make(map[string]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid line: " + line)
		}
		tokenBytes, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			return nil, err
		}
		rank, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		result[string(tokenBytes)] = rank
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

type Cl100kGptBytePairEncoding struct {
	*GptBytePairEncoding
}

func NewCl100kGptBytePairEncoding(params *mod.GptBytePairEncodingParams) mod.Encoding {
	return &Cl100kGptBytePairEncoding{
		GptBytePairEncoding: NewGptBytePairEncoding(params),
	}
}

func (e *Cl100kGptBytePairEncoding) encodeOrdinaryInternalToInt(text string, maxTokenCount int, keepEncodings bool, out *[]int) int {
	tokenCount := []int{0}
	ranks := make([]int, 0)
	parser.Split(text, func(utf8BytesList []byte) bool {
		tokenCount[0] += e.GptBytePairEncoding.Encoder.AddTokensAndGetCount(maxTokenCount, keepEncodings, utf8BytesList, out, &ranks)
		return tokenCount[0] >= maxTokenCount
	})
	return tokenCount[0]
}

func FromParameters(params *mod.GptBytePairEncodingParams) mod.Encoding {
	return NewGptBytePairEncoding(params)
}

package factory

import (
	"bufio"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/nerdface-ai/tokgo"
	"github.com/nerdface-ai/tokgo/encoding"
	"github.com/nerdface-ai/tokgo/parser"
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

func R50kBase() tokgo.Encoding {
	return from50kParameters(
		"r50k_base",
		"r50k_base.tiktoken",
		SPECIAL_TOKENS_X50K_BASE,
	)
}

func P50kBase() tokgo.Encoding {
	return from50kParameters(
		"p50k_base",
		"p50k_base.tiktoken",
		SPECIAL_TOKENS_X50K_BASE,
	)
}

func P50kEdit() tokgo.Encoding {
	return from50kParameters(
		"p50k_edit",
		"p50k_base.tiktoken",
		SPECIAL_TOKENS_P50K_EDIT,
	)
}

func Cl100kBase() tokgo.Encoding {
	mergeableRanks, err := loadMergeableRanks("cl100k_base.tiktoken")
	if err != nil {
		panic(err)
	}
	regex := compileRegex("'(?:[sdmt]|ll|ve|re)|[^\r\n\\p{L}\\p{N}]?+\\p{L}+|\\p{N}{1,3}| ?[^\\s\\p{L}\\p{N}]++[\r\n]*|\\s*[\r\n]|\\s+(?!\\S)|\\s+", false)
	params := tokgo.NewGptBytePairEncodingParams(
		"cl100k_base",
		regex,
		mergeableRanks,
		SPECIAL_TOKENS_CL100K_BASE,
	)
	return NewCl100kGptBytePairEncoding(params)
}

func O200kBase() tokgo.Encoding {
	mergeableRanks, err := loadMergeableRanks("o200k_base.tiktoken")
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
	regex := compileRegex(strings.Join(patterns, "|"), false)
	params := tokgo.NewGptBytePairEncodingParams(
		"o200k_base",
		regex,
		mergeableRanks,
		SPECIAL_TOKENS_O200K_BASE,
	)
	return FromParameters(params)
}

func from50kParameters(name, fileName string, specialTokens map[string]int) tokgo.Encoding {
	regex := compileRegex(`'(?:[sdmt]|ll|ve|re)| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+`, false)
	mergeableRanks, err := loadMergeableRanks(fileName)
	if err != nil {
		panic(err)
	}
	params := tokgo.NewGptBytePairEncodingParams(
		name,
		regex,
		mergeableRanks,
		specialTokens,
	)
	return FromParameters(params)
}

func compileRegex(pattern string, caseInsensitive bool) *regexp.Regexp {
	flags := ""
	if caseInsensitive {
		flags = "(?i)"
	}
	re, err := regexp.Compile(flags + pattern)
	if err != nil {
		panic(err)
	}
	return re
}

// getResourcePath returns the path to a resource file relative to this source file
func getResourcePath(fileName string) string {
	_, currentFilePath, _, _ := runtime.Caller(1)
	// Go up two directories: from factory to project root
	baseDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFilePath)))
	return filepath.Join(baseDir, "resources", fileName)
}

func loadMergeableRanks(fileName string) (map[string]int, error) {
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

type cl100kGptBytePairEncoding struct {
	*encoding.GptBytePairEncoding
}

func NewCl100kGptBytePairEncoding(params *tokgo.GptBytePairEncodingParams) tokgo.Encoding {
	return &cl100kGptBytePairEncoding{
		GptBytePairEncoding: encoding.NewGptBytePairEncoding(params),
	}
}

func (e *cl100kGptBytePairEncoding) EncodeOrdinaryInternal(text string, maxTokenCount int, keepEncodings bool, out *[]int) int {
	tokenCount := []int{0}
	ranks := make([]int, 0)
	parser.Split(text, func(utf8BytesList []byte) bool {
		tokenCount[0] += e.GptBytePairEncoding.Encoder.AddTokensAndGetCount(maxTokenCount, keepEncodings, utf8BytesList, out, &ranks)
		return tokenCount[0] >= maxTokenCount
	})
	return tokenCount[0]
}

func FromParameters(params *tokgo.GptBytePairEncodingParams) tokgo.Encoding {
	return NewCl100kGptBytePairEncoding(params)
}

package parser_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/currybab/tokgo/parser"
	"github.com/dlclark/regexp2"
	"github.com/stretchr/testify/assert"
)

const (
	MIN_CODE_POINT = 0
	MAX_CODE_POINT = 0x10FFFF
)

func fetchUnicodeData() (map[int]string, error) {
	url := "https://www.unicode.org/Public/UCD/latest/ucd/UnicodeData.txt"

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	unicodeMap := make(map[int]string)
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) < 2 {
			return nil, errors.New("malformed UnicodeData.txt line")
		}
		codePoint, err := strconv.ParseInt(parts[0], 16, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid code point: %v", err)
		}
		name := parts[1]
		unicodeMap[int(codePoint)] = name
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return unicodeMap, nil
}

func TestToUtf8BytesOnFetchedUnicodeData(t *testing.T) {
	unicodeData, err := fetchUnicodeData()
	if err != nil {
		t.Fatalf("Failed to fetch unicode data: %v", err)
	}

	var wg sync.WaitGroup
	for codePoint, name := range unicodeData {
		wg.Add(1)
		go func(cp int, n string) {
			defer wg.Done()
			expectedChar := string(rune(cp))
			expectedBytes := []byte(expectedChar)

			actualBytes := make([]byte, 0, utf8.UTFMax)
			actualBytes = utf8.AppendRune(actualBytes, rune(cp))

			if !bytes.Equal(expectedBytes, actualBytes) {
				t.Errorf("Expected `%v` (string: %q, char: %q) but was `%v` (char: %q). CodePoint: 0x%X, Name: %q",
					expectedBytes,
					expectedChar,
					string(rune(cp)),
					actualBytes,
					string(actualBytes),
					cp,
					n,
				)
			}
		}(codePoint, name)
	}
	wg.Wait()
}

func TestIsShortContraction(t *testing.T) {
	pattern := regexp2.MustCompile("^(?:'s|'t|'re|'ve|'m|'ll|'d|'ſ)$", regexp2.IgnoreCase)
	// jtokkit pattern doesn't have 'ſ so we add it here

	for cp := MIN_CODE_POINT; cp <= MAX_CODE_POINT; cp++ {
		asString := "'" + string(rune(cp))

		matchesRegex, err := pattern.MatchString(asString)
		if err != nil {
			t.Errorf("Failed to match regex: %v", err)
			return
		}

		actual := parser.IsShortContraction(cp)
		// t.Logf("%s -> %v", asString, actual)
		if matchesRegex != actual {
			t.Errorf("Mismatch at code point: `%s` (%s) %v", asString, string(rune(cp)), matchesRegex)
		}
	}
}

func TestIsLongContraction(t *testing.T) {
	if os.Getenv("NOT_SKIP_TEST") != "1" {
		t.Skip("Skipping long contraction test")
	}
	pattern := regexp2.MustCompile("^(?:'s|'t|'re|'ve|'m|'ll|'d|'ſ)$", regexp2.IgnoreCase)
	// jtokkit pattern doesn't have 'ſ so we add it here

	for cp1 := MIN_CODE_POINT; cp1 <= MAX_CODE_POINT; cp1++ {
		for cp2 := MIN_CODE_POINT; cp2 <= MAX_CODE_POINT; cp2++ {
			if cp2 == 10 { // different with jtokkit
				continue
			}
			asString := "'" + string(rune(cp1)) + string(rune(cp2))

			matchesRegex, err := pattern.MatchString(asString)
			if err != nil {
				t.Errorf("Failed to match regex: %v", err)
				return
			}

			actual := parser.IsLongContraction(cp1, cp2)
			// t.Logf("%s -> %v", asString, actual)
			if matchesRegex != actual {
				t.Errorf("Mismatch at code point: `%s` (%s) %d %d", asString, string(rune(cp1))+string(rune(cp2)), cp1, cp2)
			}
		}
	}
}

func TestIsNumeric(t *testing.T) {
	count := 0
	assert.Equal(t, parser.IsNumeric(-1), false)
	pattern := regexp2.MustCompile("^\\p{N}$", regexp2.IgnoreCase)
	for cp := MIN_CODE_POINT; cp <= MAX_CODE_POINT; cp++ {
		charAsString := string(rune(cp))
		matchesRegex, err := pattern.MatchString(charAsString)
		if err != nil {
			t.Errorf("Failed to match regex: %v", err)
			return
		}
		actual := parser.IsNumeric(cp)
		if matchesRegex {
			count++
		}

		assert.Equal(t, matchesRegex, actual)
	}
	t.Log(count)
}

func TestIsLetter(t *testing.T) {
	count := 0
	assert.Equal(t, parser.IsLetter(-1), false)
	pattern := regexp2.MustCompile("^\\p{L}$", regexp2.IgnoreCase)
	for cp := MIN_CODE_POINT; cp <= MAX_CODE_POINT; cp++ {
		charAsString := string(rune(cp))
		matchesRegex, err := pattern.MatchString(charAsString)
		if err != nil {
			t.Errorf("Failed to match regex: %v", err)
			return
		}
		actual := parser.IsLetter(cp)
		if matchesRegex {
			count++
		}

		assert.Equal(t, matchesRegex, actual)
	}
	t.Log(count)
}

func TestIsUnicodeWhitespace(t *testing.T) {
	count := 0
	assert.Equal(t, parser.IsWhitespace(-1), false)
	pattern := regexp2.MustCompile("^\\s$", regexp2.IgnoreCase)
	for cp := MIN_CODE_POINT; cp <= MAX_CODE_POINT; cp++ {
		charAsString := string(rune(cp))
		matchesRegex, err := pattern.MatchString(charAsString)
		if err != nil {
			t.Errorf("Failed to match regex: %v", err)
			return
		}
		actual := parser.IsWhitespace(cp)
		if matchesRegex {
			count++
		}

		assert.Equal(t, matchesRegex, actual)
	}
	t.Log(count)
}

func TestIsLetterOrNumeric(t *testing.T) {
	count := 0
	assert.Equal(t, parser.IsLetterOrNumeric(-1), false)
	pattern := regexp2.MustCompile("^[\\p{L}\\p{N}]$", regexp2.IgnoreCase)
	for cp := MIN_CODE_POINT; cp <= MAX_CODE_POINT; cp++ {
		charAsString := string(rune(cp))
		matchesRegex, err := pattern.MatchString(charAsString)
		if err != nil {
			t.Errorf("Failed to match regex: %v", err)
			return
		}
		actual := parser.IsLetterOrNumeric(cp)
		if matchesRegex {
			count++
		}

		assert.Equal(t, matchesRegex, actual)
	}
	t.Log(count)
}

func TestIsNotWhitespaceOrLetterOrNumeric(t *testing.T) {
	count := 0
	assert.Equal(t, parser.IsNotWhitespaceOrLetterOrNumeric(-1), false)
	pattern := regexp2.MustCompile("^[^\\s\\p{L}\\p{N}]$", regexp2.IgnoreCase)
	for cp := MIN_CODE_POINT; cp <= MAX_CODE_POINT; cp++ {
		charAsString := string(rune(cp))
		matchesRegex, err := pattern.MatchString(charAsString)
		if err != nil {
			t.Errorf("Failed to match regex: %v", err)
			return
		}
		actual := parser.IsNotWhitespaceOrLetterOrNumeric(cp)
		if matchesRegex {
			count++
		}

		assert.Equal(t, matchesRegex, actual)
	}
	t.Log(count)
}

func TestIsNotNewlineOrLetterOrNumeric(t *testing.T) {
	count := 0
	assert.Equal(t, parser.IsNotNewlineOrLetterOrNumeric(-1), false)
	pattern := regexp2.MustCompile("^[^\r\n\\p{L}\\p{N}]$", regexp2.IgnoreCase)
	for cp := MIN_CODE_POINT; cp <= MAX_CODE_POINT; cp++ {
		charAsString := string(rune(cp))
		matchesRegex, err := pattern.MatchString(charAsString)
		if err != nil {
			t.Errorf("Failed to match regex: %v", err)
			return
		}
		actual := parser.IsNotNewlineOrLetterOrNumeric(cp)
		if matchesRegex {
			count++
		}

		assert.Equal(t, matchesRegex, actual)
	}
	t.Log(count)
}

func TestIsNewline(t *testing.T) {
	count := 0
	assert.Equal(t, parser.IsNewline(-1), false)
	pattern := regexp2.MustCompile("^[\r\n]$", regexp2.IgnoreCase)
	for cp := MIN_CODE_POINT; cp <= MAX_CODE_POINT; cp++ {
		charAsString := string(rune(cp))
		matchesRegex, err := pattern.MatchString(charAsString)
		if err != nil {
			t.Errorf("Failed to match regex: %v", err)
			return
		}
		actual := parser.IsNewline(cp)
		if matchesRegex {
			count++
		}

		assert.Equal(t, matchesRegex, actual)
	}
	t.Log(count)
}

package encoder_test

import (
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/currybab/tokgo"
	"github.com/currybab/tokgo/encoding"
	"github.com/currybab/tokgo/parser"
	"github.com/emirpasic/gods/v2/maps/treemap"
	"github.com/stretchr/testify/assert"
)

const (
	MIN_CODE_POINT = 0
	MAX_CODE_POINT = 0x10FFFF
)
const PUNCTUATION = "'\".,?!:()"

var LETTERS = generateUnicodeCategoryString(parser.IsLetter)
var NUMBERS = generateUnicodeCategoryString(parser.IsNumeric)
var WHITESPACES = generateUnicodeCategoryString(parser.IsWhitespace)
var NEWLINES = "\n\r"
var NOT_NEWLINE_OR_LETTER_OR_NUMERIC = generateUnicodeCategoryString(parser.IsNotNewlineOrLetterOrNumeric)
var NOT_WHITESPACE_OR_LETTER_OR_NUMERIC = generateUnicodeCategoryString(parser.IsNotWhitespaceOrLetterOrNumeric)
var SPECIAL = []string{"'s", "'t", "'re", "'ve", "'m", "'ll", "'d", "'ſ", "'x", "🤚🏾", "😩", "　", "½"}
var ENCODING = encoding.Cl100kBase()

func generateUnicodeCategoryString(predicate func(ch int) bool) string {
	var sb strings.Builder
	for r := MIN_CODE_POINT; r <= MAX_CODE_POINT; r++ {
		if utf8.ValidRune(rune(r)) && unicode.IsPrint(rune(r)) && predicate(r) {
			sb.WriteRune(rune(r))
		}
	}
	return sb.String()
}

func normalizeStringForTesting(testString string) string {
	s := strings.ReplaceAll(testString, "\r", "\\r")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, " ", "␣")
	return s
}

func getEncoding() tokgo.Encoding {
	// fmt.Println("Using encoding:", os.Getenv(tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY))
	return ENCODING
}

func getAllTokens() []string {
	tokens, _ := encoding.LoadMergeableRanks("cl100k_base.tiktoken")
	var result []string
	for token := range tokens {
		result = append(result, token)
	}
	return result
}

func generateRandomUtf8String(singleTokenStrings []string) string {
	var testString string
	for {
		length := rand.Intn(9) + 1
		var sb strings.Builder
		for i := 0; i < length; i++ {
			category := rand.Intn(20)
			charRunes := getRandomCharFromCategory(category, singleTokenStrings)
			charSegment := string(charRunes)
			if rand.Float64() >= 0.5 {
				if rand.Float64() >= 0.5 {
					charSegment = strings.ToUpper(charSegment)
				} else {
					charSegment = strings.ToLower(charSegment)
				}
			}
			sb.WriteString(charSegment)
		}
		testString = sb.String()
		if utf8.ValidString(testString) {
			break
		}
	}
	return testString
}

func getRandomCharFromCategory(category int, singleTokenStrings []string) []rune {
	switch category {
	case 0:
		return []rune{' '}
	case 1:
		return []rune{' ', ' '} // two spaces
	case 2, 3, 4: // letters
		charA := 'a'
		if rand.Intn(2) == 0 {
			charA = 'A'
		}
		return []rune{charA + rune(rand.Intn(26))}
	case 5: // punctuation
		if len(PUNCTUATION) == 0 {
			return []rune{'?'}
		}
		runes := []rune(PUNCTUATION)
		return []rune{runes[rand.Intn(len(runes))]}
	case 6, 7: // newlines
		if len(NEWLINES) == 0 {
			return []rune{'\n'}
		}
		runes := []rune(NEWLINES)
		return []rune{runes[rand.Intn(len(runes))]}
	case 8: // numbers
		if len(NUMBERS) == 0 {
			return []rune{'0' + rune(rand.Intn(10))}
		}
		runes := []rune(NUMBERS)
		return []rune{runes[rand.Intn(len(runes))]}
	case 9: // whitespaces
		if len(WHITESPACES) == 0 {
			return []rune{' '}
		}
		runes := []rune(WHITESPACES)
		return []rune{runes[rand.Intn(len(runes))]}
	case 10, 11: // letters (broader set from generatedUnicodeCategoryStringGo)
		if len(LETTERS) == 0 {
			return getRandomCharFromCategory(2, singleTokenStrings) // fallback to basic letters
		}
		runes := []rune(LETTERS)
		return []rune{runes[rand.Intn(len(runes))]}
	case 12, 13: // NOT_NEWLINE_OR_LETTER_OR_NUMERIC
		if len(NOT_NEWLINE_OR_LETTER_OR_NUMERIC) == 0 {
			return []rune{'*'}
		}
		runes := []rune(NOT_NEWLINE_OR_LETTER_OR_NUMERIC)
		return []rune{runes[rand.Intn(len(runes))]}
	case 14: // NOT_WHITESPACE_OR_LETTER_OR_NUMERIC
		if len(NOT_WHITESPACE_OR_LETTER_OR_NUMERIC) == 0 {
			return []rune{'%'}
		}
		runes := []rune(NOT_WHITESPACE_OR_LETTER_OR_NUMERIC)
		return []rune{runes[rand.Intn(len(runes))]}
	case 15, 16: // emojis (simple range)
		return []rune{0x1F600 + rune(rand.Intn(0x50))}
	case 17: // special tokens from the list
		if len(SPECIAL) == 0 {
			return []rune{'S'} // Fallback for empty SPECIAL list
		}
		return []rune(SPECIAL[rand.Intn(len(SPECIAL))])
	case 18: // single token strings from the broader list
		if len(singleTokenStrings) == 0 {
			return []rune{'T'} // Fallback for empty singleTokenStrings list
		}
		return []rune(singleTokenStrings[rand.Intn(len(singleTokenStrings))])
	case 19: // any defined unicode char (simplified)
		for {
			r := rune(rand.Intn(MAX_CODE_POINT + 1))
			if utf8.ValidRune(r) && unicode.IsPrint(r) { // Ensure it's printable and valid
				return []rune{r}
			}
		}
	default:
		panic("Invalid category")
	}
}

func measureEncodingSpeeds(t *testing.T) {
	if os.Getenv("NOT_SKIP_TEST") != "1" {
		t.Skip("Skipping performance test as it was disabled in the original Java version and can be lengthy.")
	}

	var input strings.Builder
	measurements := treemap.New[int, int64]()
	iterations := 20
	// Math.max(i + 1, i * 1.01) equivalent
	maxFloat := func(a, b float64) float64 {
		if a > b {
			return a
		}
		return b
	}

	for i := 1.0; i < 1000; i = maxFloat(i+1, i*1.01) {
		for float64(input.Len()) < i {
			input.WriteString("a")
		}
		inputString := input.String()

		// Warmup
		for j := 0; j < 10*iterations; j++ {
			warmup := getEncoding().EncodeToIntArray(inputString)
			assert.NotEmpty(t, warmup, "Warmup encoding should not be empty")
		}

		startTime := time.Now().UnixNano()
		for j := 0; j < iterations; j++ {
			encodingResult := getEncoding().EncodeToIntArray(inputString)
			assert.NotEmpty(t, encodingResult, "Encoding result should not be empty")
		}
		endTime := time.Now().UnixNano()
		measurements.Put(int(i), (endTime-startTime)/int64(iterations))
	}

	// To print in sorted order of keys, like TreeMap
	it := measurements.Iterator()
	for it.Begin(); it.Next(); {
		t.Logf("%d\t%d\n", it.Key(), it.Value())
	}
}

func cl100kEdgeCaseRoundTrips(t *testing.T) {
	testStrings := []string{
		"\n",
		" ",
		"a : b",
		"  a",
		"\n \n ",
		"\n \n",
		"\n ",
		"\n \n!",
		"\n \n   ",
		"\n  !",
		"\n A",
		"  \n\r  \r\n  \r \n  A\nA \n A",
		",\n\n",
		" ***\n\n\n\n",

		"   !",
		"   A",
		"   0",
		"   *",

		"   \n!",
		"   \nA",
		"   \n0",
		"   \n*",

		"   \n !",
		"   \n A",
		"   \n 0",
		"   \n *",

		"Many words map to one token, but some don't: indivisible.\n\nUnicode characters like emojis may be split into many tokens containing the underlying bytes: 🤚🏾\n\nSequences of characters commonly found next to each other may be grouped together: 1234567890",
		"I paid $123,456 to 9876543210 people!",
		"Mixed script: 你好 world! 🌍",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"Unicode snowman: ☃️",
		"I'm:  0\n",
		"We'll meet at 3 o'clock.",
		"Hello, world! It's a beautiful day...",
		"In 2023, I'll be 25 years old.",
		"Hello \n\n World  !",
		" It's 2:30pm;\n\n\n\nlet's eat, sleep , and code!",
		"'Thank God, here it is.' But when we took up the trunk...",
		"What in the world are you doing???!!!",
		"user@example.com",
		"this is a 'quoted' word",
		"　　a",
		"'ſ",
		"'ſ𣶸𣄬Ƙ淚",
		"😩\n",
		"03½",
		"* ע", // Note: Java string has a non-breaking space (U+00A0) here
		"مرحبا بالعالم! كيف حالك؟ 😎",
		"\u0000\U000147E1 a\u0000b-\u0000\u0000\u0000 \u0000", // Java string literals with Unicode escapes
		"🌍 a",
		"(𥧙h",
		",   𰉄", // Note: Java string has EM SPACE (U+2003)
		"  󵨐)",  // Note: Java string has THIN SPACE (U+2009) and FIGURE SPACE (U+2007)
		"ﮀ\n ",
		"😐𪶫X",
		"෫𞅄",
		"𬕹\n  ",
		" 😈b\n\U000212AE'ſ\U00023ED8\U0002312CƘ淚",
		"𗭾  󻥹\n\U0002D6F0蛇",
		"こんにちは世界",
	}

	for _, testString := range testStrings {
		t.Logf("Validating `%s`\n", normalizeStringForTesting(testString))

		actualTokens := getEncoding().EncodeToIntArray(testString)
		decoded := getEncoding().Decode(actualTokens)
		assert.Equal(t, testString, string(decoded), "Round trip failed for: %s", normalizeStringForTesting(testString))
	}
}

func cl100kEncodeRoundTripWithRandomStrings(t *testing.T) {
	singleTokenStrings := getAllTokens()
	var wg sync.WaitGroup
	for i := 0; i < 100_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testString := generateRandomUtf8String(singleTokenStrings)

			maxTokenCount := rand.Intn(2*len(testString)) + 1
			actualTokens := getEncoding().EncodeToIntArray(testString)
			assert.Equal(t, len(actualTokens), getEncoding().CountTokens(testString))

			decodedTokens := getEncoding().Decode(actualTokens)
			assert.Equal(t, testString, decodedTokens)

			actualTrimmedTokens := getEncoding().Encode(testString, maxTokenCount).GetTokens()
			decodedTrimmedTokens := getEncoding().Decode(actualTrimmedTokens)
			assert.True(t, strings.HasPrefix(testString, decodedTrimmedTokens))
		}()
	}
	wg.Wait()
}

func cl100kEncodeOrdinaryRoundTripWithRandomStrings(t *testing.T) {
	singleTokenStrings := getAllTokens()
	var wg sync.WaitGroup
	for i := 0; i < 100_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testString := generateRandomUtf8String(singleTokenStrings)

			maxTokenCount := rand.Intn(2*len(testString)) + 1
			actualTokens := getEncoding().EncodeOrdinaryToIntArray(testString)
			assert.Equal(t, len(actualTokens), getEncoding().CountTokensOrdinary(testString))

			decodedTokens := getEncoding().Decode(actualTokens)
			assert.Equal(t, testString, decodedTokens)

			actualTrimmedTokens := getEncoding().EncodeOrdinary(testString, maxTokenCount).GetTokens()
			decodedTrimmedTokens := getEncoding().Decode(actualTrimmedTokens)
			assert.True(t, strings.HasPrefix(testString, decodedTrimmedTokens))
		}()
	}
	wg.Wait()
}

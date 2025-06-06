package encoder_test

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/currybab/tokgo"
	"github.com/currybab/tokgo/encoding"
	"github.com/emirpasic/gods/v2/maps/treemap"

	"github.com/stretchr/testify/assert"
)

// Environment variable key we want to manipulate for these tests.

func TestMain(m *testing.M) {
	// --- Setup: Store original environment variable and set the new one ---
	originalThresholdValue := os.Getenv(tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY)
	err := os.Setenv(tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, "0")
	if err != nil {
		log.Fatalf("Failed to set environment variable '%s' for large tokenizer test: %v", tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, err)
	}

	// Re-initialize encoding with the new environment variable setting.
	LARGE_ENCODING = encoding.Cl100kBase()
	if LARGE_ENCODING == nil {
		log.Fatal("Failed to re-initialize LARGE_ENCODING after setting environment variable")
	}

	// Run tests in this file
	exitCode := m.Run()

	// --- Teardown: Restore original environment variable ---
	if originalThresholdValue == "" {
		os.Unsetenv(tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY)
	} else {
		err = os.Setenv(tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, originalThresholdValue)
		if err != nil {
			log.Printf("Failed to restore environment variable '%s': %v", tokgo.VERY_LARGE_TOKENIZER_BYTE_THRESHOLD_KEY, err)
		}
	}

	os.Exit(exitCode)
}

var LARGE_ENCODING = encoding.Cl100kBase()

func getLargeEncoding() tokgo.Encoding {
	return LARGE_ENCODING
}

func TestLargeCl100kMeasureEncodingSpeeds(t *testing.T) {
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
			warmup := getLargeEncoding().EncodeToIntArray(inputString)
			assert.NotEmpty(t, warmup, "Warmup encoding should not be empty")
		}

		startTime := time.Now().UnixNano()
		for j := 0; j < iterations; j++ {
			encodingResult := getLargeEncoding().EncodeToIntArray(inputString)
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

func TestLargeCl100kEdgeCaseRoundTrips(t *testing.T) {
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

		actualTokens := getLargeEncoding().EncodeToIntArray(testString)
		decoded := getLargeEncoding().Decode(actualTokens)
		assert.Equal(t, testString, string(decoded), "Round trip failed for: %s", normalizeStringForTesting(testString))
	}
}

func TestLargeCl100kEncodeRoundTripWithRandomStrings(t *testing.T) {
	singleTokenStrings := getAllTokens()
	var wg sync.WaitGroup
	for i := 0; i < 100_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testString := generateRandomUtf8String(singleTokenStrings)

			maxTokenCount := rand.Intn(2*len(testString)) + 1
			actualTokens := getLargeEncoding().EncodeToIntArray(testString)
			assert.Equal(t, len(actualTokens), getLargeEncoding().CountTokens(testString))

			decodedTokens := getLargeEncoding().Decode(actualTokens)
			assert.Equal(t, testString, decodedTokens)

			actualTrimmedTokens := getLargeEncoding().Encode(testString, maxTokenCount).GetTokens()
			decodedTrimmedTokens := getLargeEncoding().Decode(actualTrimmedTokens)
			assert.True(t, strings.HasPrefix(testString, decodedTrimmedTokens))
		}()
	}
	wg.Wait()
}

func TestLargeCl100kEncodeOrdinaryRoundTripWithRandomStrings(t *testing.T) {
	singleTokenStrings := getAllTokens()
	var wg sync.WaitGroup
	for i := 0; i < 100_000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testString := generateRandomUtf8String(singleTokenStrings)

			maxTokenCount := rand.Intn(2*len(testString)) + 1
			actualTokens := getLargeEncoding().EncodeOrdinaryToIntArray(testString)
			assert.Equal(t, len(actualTokens), getLargeEncoding().CountTokensOrdinary(testString))

			decodedTokens := getLargeEncoding().Decode(actualTokens)
			assert.Equal(t, testString, decodedTokens)

			actualTrimmedTokens := getLargeEncoding().EncodeOrdinary(testString, maxTokenCount).GetTokens()
			decodedTrimmedTokens := getLargeEncoding().Decode(actualTrimmedTokens)
			assert.True(t, strings.HasPrefix(testString, decodedTrimmedTokens))
		}()
	}
	wg.Wait()
}

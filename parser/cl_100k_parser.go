package parser

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	SDTM               = "sdtmSDTMſ"
	SIMPLE_WHITESPACES = "\t\n\u000B\u000C\r"
)

var REMAINING_WHITESPACES = []int{
	0x1680, 0x2000, 0x2001, 0x2002, 0x2003, 0x2004, 0x2005, 0x2006, 0x2007, 0x2008, 0x2009, 0x200A, 0x2028, 0x2029, 0x202F, 0x205F, 0x3000,
}

// ByteArrayList is a simple dynamic array of bytes
type ByteArrayList struct {
	data []byte
}

// Clear empties the byte array
func (b *ByteArrayList) Clear() {
	b.data = b.data[:0]
}

// Add adds a byte to the array
func (b *ByteArrayList) Add(value byte) {
	b.data = append(b.data, value)
}

// Get the underlying byte array
func (b *ByteArrayList) GetData() []byte {
	return b.data
}

// FragmentConsumer is a function that processes ByteArrayList and returns a boolean
type FragmentConsumer func([]byte) bool

// Split tokenizes the input string into UTF-8 fragments
func Split(input string, fragmentConsumer FragmentConsumer) {
	if !IsValidUTF8(input) {
		panic("Input is not UTF-8: " + input)
	}

	utf8Bytes := []byte{}
	finished := false
	inputRunes := []rune(input)

	for endIndex := 0; endIndex < len(inputRunes) && !finished; {
		startIndex := endIndex
		c0 := int(inputRunes[startIndex])
		_ = utf8.RuneLen(rune(c0)) // Not used but kept for potential future use
		nextIndex := startIndex + 1

		var c1 int = -1
		if nextIndex < len(inputRunes) {
			c1 = int(inputRunes[nextIndex])
		}

		if c0 == '\'' && c1 > 0 {
			if IsShortContraction(c1) {
				// 1) `\'[sdtm]` - contractions, such as the suffixes of `he\'s`, `I\'d`, `\'tis`, `I\'m`
				endIndex += 2
				finished = fragmentConsumer(AddUtf8Bytes(inputRunes, startIndex, endIndex, utf8Bytes))
				continue
			} else if startIndex+2 < len(inputRunes) && IsLongContraction(c1, int(inputRunes[startIndex+2])) {
				// 1) `\'(?:ll|ve|re)` - contractions, such as the suffixes of `you\'ll`, `we\'ve`, `they\'re`
				endIndex += 3
				finished = fragmentConsumer(AddUtf8Bytes(inputRunes, startIndex, endIndex, utf8Bytes))
				continue
			}
		}

		// Character length not needed in Go implementation
		_ = c1 // Ensure c1 is used

		if (IsNotNewlineOrLetterOrNumeric(c0) && IsLetter(c1)) || IsLetter(c0) {
			// 2) `[^\r\n\p{L}\p{N}]?+\p{L}+` - words such as ` of`, `th`, `It`, ` not`
			endIndex += 1
			if IsLetter(c1) {
				endIndex += 1
				for endIndex < len(inputRunes) {
					c0 = int(inputRunes[endIndex])
					if !IsLetter(c0) {
						break
					}
					endIndex += 1
				}
			}
			finished = fragmentConsumer(AddUtf8Bytes(inputRunes, startIndex, endIndex, utf8Bytes))
		} else if IsNumeric(c0) {
			// 3) `\p{N}{1,3}` - numbers, such as `4`, `235` or `3½`
			endIndex += 1
			if IsNumeric(c1) {
				endIndex += 1
				if endIndex < len(inputRunes) {
					c0 = int(inputRunes[endIndex])
					if IsNumeric(c0) {
						endIndex += 1
					}
				}
			}
			finished = fragmentConsumer(AddUtf8Bytes(inputRunes, startIndex, endIndex, utf8Bytes))
		} else if IsNotWhitespaceOrLetterOrNumeric(c0) || ((c0 == ' ') && IsNotWhitespaceOrLetterOrNumeric(c1)) {
			// 4) ` ?[^\s\p{L}\p{N}]++[\r\n]*` - punctuation, such as `,`, ` .`, `"`
			endIndex += 1
			if endIndex < len(inputRunes) && IsNotWhitespaceOrLetterOrNumeric(c1) {
				endIndex += 1
				for endIndex < len(inputRunes) {
					c0 = int(inputRunes[endIndex])
					if !IsNotWhitespaceOrLetterOrNumeric(c0) {
						break
					}
					endIndex += 1
				}
			}
			for endIndex < len(inputRunes) && IsNewline(int(inputRunes[endIndex])) {
				endIndex += 1
			}
			finished = fragmentConsumer(AddUtf8Bytes(inputRunes, startIndex, endIndex, utf8Bytes))
		} else {
			// 5) `\s*[\r\n]+` - line endings such as `\r\n    \r\n`
			// 6) `\s+(?!\S)` - whitespaces such as `               ` or ` `
			// 7) `\s+` - unmatched remaining spaces, such as ` `
			if !IsWhitespace(c0) {
				panic("Invalid character")
			}

			lastNewLineIndex := -1
			if IsNewline(c0) {
				lastNewLineIndex = endIndex
			}
			endIndex += 1

			if IsWhitespace(c1) {
				if IsNewline(c1) {
					lastNewLineIndex = endIndex
				}
				endIndex += 1
				for endIndex < len(inputRunes) {
					c0 = int(inputRunes[endIndex])
					if !IsWhitespace(c0) {
						break
					}
					if IsNewline(c0) {
						lastNewLineIndex = endIndex
					}
					endIndex += 1
				}
			}

			if lastNewLineIndex > -1 {
				finalEndIndex := endIndex
				endIndex = lastNewLineIndex + 1
				if endIndex < finalEndIndex {
					if startIndex >= endIndex {
						panic("startIndex must be less than endIndex")
					}
					finished = fragmentConsumer(AddUtf8Bytes(inputRunes, startIndex, endIndex, utf8Bytes))
					startIndex = endIndex
					endIndex = finalEndIndex
				}
			}

			if !finished {
				if lastNewLineIndex+1 < endIndex && !IsWhitespace(c0) {
					endIndex -= 1
				}
				if startIndex < endIndex {
					finished = fragmentConsumer(AddUtf8Bytes(inputRunes, startIndex, endIndex, utf8Bytes))
				}
			}
		}
	}
}

// IsShortContraction checks if a character is a short contraction
func IsShortContraction(ch int) bool {
	return strings.ContainsRune(SDTM, rune(ch))
}

// IsLongContraction checks if a character pair forms a long contraction
func IsLongContraction(ch1, ch2 int) bool {
	if (ch1 == 'l' && ch2 == 'l') ||
		(ch1 == 'v' && ch2 == 'e') ||
		(ch1 == 'r' && ch2 == 'e') {
		return true
	} else {
		lch1 := unicode.ToUpper(rune(ch1))
		lch2 := unicode.ToUpper(rune(ch2))
		return (lch1 == 'L' && lch2 == 'L') ||
			(lch1 == 'V' && lch2 == 'E') ||
			(lch1 == 'R' && lch2 == 'E')
	}
}

// IsValidUTF8 checks if a string is valid UTF-8
func IsValidUTF8(input string) bool {
	return utf8.ValidString(input)
}

// IsLetter checks if a code point is a letter
func IsLetter(ch int) bool {
	if ch < 0xaa {
		return (ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z')
	} else if ch <= 0x323af {
		return unicode.IsLetter(rune(ch))
	}
	return false
}

// IsNumeric checks if a code point is a number
func IsNumeric(ch int) bool {
	if ch < 0xb2 {
		return ch >= '0' && ch <= '9'
	} else if ch <= 0x1fbf9 {
		runeChar := rune(ch)
		return unicode.IsNumber(runeChar)
	}
	return false
}

// IsLetterOrNumeric checks if a code point is a letter or number
func IsLetterOrNumeric(ch int) bool {
	if ch < 0xaa {
		return (ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9')
	} else if ch <= 0x323af {
		runeChar := rune(ch)
		return unicode.IsLetter(runeChar) || unicode.IsNumber(runeChar)
	}
	return false
}

// IsWhitespace checks if a code point is whitespace
func IsWhitespace(ch int) bool {
	if ch <= '\r' {
		return strings.ContainsRune(SIMPLE_WHITESPACES, rune(ch))
	} else if ch < 0x85 {
		return ch == ' '
	} else {
		return ch == 0x85 ||
			ch == 0xA0 ||
			(ch >= 0x1680 && ch <= 0x3000 && binarySearch(REMAINING_WHITESPACES, ch) >= 0)
	}
}

// IsNewline checks if a code point is a newline
func IsNewline(ch int) bool {
	return ch == '\r' || ch == '\n'
}

// IsNotWhitespaceOrLetterOrNumeric checks if a code point is not whitespace, letter, or numeric
func IsNotWhitespaceOrLetterOrNumeric(ch int) bool {
	if ch < '0' {
		return ch >= 0 && ch != ' ' && (ch > '\r' || ch < '\t')
	} else {
		return !IsLetterOrNumeric(ch) && !IsWhitespace(ch)
	}
}

// IsNotNewlineOrLetterOrNumeric checks if a code point is not newline, letter, or numeric
func IsNotNewlineOrLetterOrNumeric(ch int) bool {
	if ch < '0' {
		return ch >= 0 && (ch == ' ' || !IsNewline(ch))
	} else {
		return !IsLetterOrNumeric(ch)
	}
}

// AddUtf8Bytes converts a rune slice to UTF-8 bytes
func AddUtf8Bytes(input []rune, start, end int, dst []byte) []byte {
	dst = dst[:0]
	for i := start; i < end; i++ {
		cp := int(input[i])
		if cp < 0x80 {
			dst = append(dst, byte(cp))
		} else if cp < 0x800 {
			dst = append(dst, byte(0xc0|(cp>>0x6)))
			dst = append(dst, byte(0x80|(cp&0x3f)))
		} else if cp < 0x10000 {
			dst = append(dst, byte(0xe0|(cp>>0xc)))
			dst = append(dst, byte(0x80|((cp>>0x6)&0x3f)))
			dst = append(dst, byte(0x80|(cp&0x3f)))
		} else {
			if cp >= 0x110000 {
				panic(fmt.Sprintf("Invalid code point: %d", cp))
			}
			dst = append(dst, byte(0xf0|(cp>>0x12)))
			dst = append(dst, byte(0x80|((cp>>0xc)&0x3f)))
			dst = append(dst, byte(0x80|((cp>>0x6)&0x3f)))
			dst = append(dst, byte(0x80|(cp&0x3f)))
		}
	}
	return dst
}

// binarySearch performs a binary search on a sorted int slice
func binarySearch(array []int, key int) int {
	index := sort.SearchInts(array, key)
	if index < len(array) && array[index] == key {
		return index
	}
	return -1
}

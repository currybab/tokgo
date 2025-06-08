package referencetest

import (
	"encoding/csv"
	"io" // For io.EOF
	"os"
	"strconv"
	"strings"
	"testing"
)

func parseEncodingString(encodingString string) []int {
	list := strings.Split(encodingString[1:len(encodingString)-1], ",")
	var result = make([]int, len(list))
	for i, s := range list {
		result[i], _ = strconv.Atoi(strings.TrimSpace(s))
	}
	return result
}

// CsvIterator helps to iterate over records in a CSV file.
type CsvIterator struct {
	reader  *csv.Reader
	file    *os.File
	current []string
	lastErr error
}

// NewCsvIterator creates a new CsvIterator for the given file path.
// If skipHeader is true, it will attempt to read and discard the first line of the CSV.
// The caller is responsible for calling Close() on the iterator when done.
func NewCsvIterator(filePath string, skipHeader bool) (*CsvIterator, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	// Optional: Configure the reader here if needed, for example:
	// reader.Comma = ';'
	// reader.FieldsPerRecord = -1 // To allow variable number of fields

	if skipHeader {
		if _, err := reader.Read(); err != nil {
			// If there's an error reading the header (e.g., empty file, or just a header leading to EOF),
			// close the file and return the error.
			file.Close() // Attempt to close the file
			return nil, err
		}
	}

	return &CsvIterator{
		reader: reader,
		file:   file,
	}, nil
}

// Next advances the iterator to the next record.
// It returns true if a record was successfully read, false otherwise (EOF or error).
// After Next returns false, the Err method should be checked to distinguish
// between EOF and other errors.
func (it *CsvIterator) Next() bool {
	record, err := it.reader.Read()
	if err != nil {
		if err == io.EOF { // io.EOF is a clean end of file, not an error to be reported by Err()
			it.lastErr = nil
		} else {
			it.lastErr = err
		}
		it.current = nil // Clear current record on EOF or error
		return false
	}
	it.current = record
	it.lastErr = nil
	return true
}

// Record returns the current record. It should only be called after a successful call to Next.
// The returned slice should not be modified by the caller if it's to be reused by the iterator.
func (it *CsvIterator) Record() []string {
	return it.current
}

// Err returns the first non-EOF error that was encountered by the iterator.
func (it *CsvIterator) Err() error {
	return it.lastErr
}

// Close closes the underlying file. It should be called when done with the iterator
// to release system resources.
func (it *CsvIterator) Close() error {
	if it.file != nil {
		return it.file.Close()
	}
	return nil
}

func WrapTest(t *testing.T, filename string, testFunc func(string, string, string)) {
	it, err := NewCsvIterator(filename, true) // Assuming the CSV has a header to skip, true)
	if err != nil {
		t.Fatal(err)
	}
	defer it.Close()

	for it.Next() {
		record := it.Record()
		input := record[0]
		output := record[1]
		outputMaxTokens10 := record[2]
		testFunc(input, output, outputMaxTokens10)
	}
}

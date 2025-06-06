package encoder_test

import (
	"testing"
)

func TestMeasureEncodingSpeeds(t *testing.T) {
	measureEncodingSpeeds(t)
}

func TestCl100kEdgeCaseRoundTrips(t *testing.T) {
	cl100kEdgeCaseRoundTrips(t)
}

func TestCl100kEncodeRoundTripWithRandomStrings(t *testing.T) {
	cl100kEncodeRoundTripWithRandomStrings(t)
}

func TestCl100kEncodeOrdinaryRoundTripWithRandomStrings(t *testing.T) {
	cl100kEncodeOrdinaryRoundTripWithRandomStrings(t)
}

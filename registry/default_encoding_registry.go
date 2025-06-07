package tokgo

import "github.com/currybab/tokgo/mod"

type DefaultEncodingRegistry struct {
	*AbstractEncodingRegistry
}

func (r *DefaultEncodingRegistry) initializeDefaultEncodings() {
	for _, encodingType := range mod.EncodingTypeValues() {
		r.AddEncoding(encodingType)
	}
}

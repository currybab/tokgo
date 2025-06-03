package registry

import "github.com/currybab/tokgo"

type DefaultEncodingRegistry struct {
	*AbstractEncodingRegistry
}

func (r *DefaultEncodingRegistry) initializeDefaultEncodings() {
	for _, encodingType := range tokgo.EncodingTypeValues() {
		r.AddEncoding(encodingType)
	}
}

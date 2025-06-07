package tokgo

import "github.com/currybab/tokgo/mod"

func NewDefaultEncodingRegistry() mod.EncodingRegistry {
	reg := &DefaultEncodingRegistry{
		AbstractEncodingRegistry: &AbstractEncodingRegistry{},
	}
	reg.initializeDefaultEncodings()
	return reg
}

func NewLazyEncodingRegistry() mod.EncodingRegistry {
	return &LazyEncodingRegistry{
		AbstractEncodingRegistry: &AbstractEncodingRegistry{},
	}
}

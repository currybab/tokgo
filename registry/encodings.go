package registry

import "github.com/currybab/tokgo"

func NewDefaultEncodingRegistry() tokgo.EncodingRegistry {
	reg := &DefaultEncodingRegistry{
		AbstractEncodingRegistry: &AbstractEncodingRegistry{},
	}
	reg.initializeDefaultEncodings()
	return reg
}

func NewLazyEncodingRegistry() tokgo.EncodingRegistry {
	return &LazyEncodingRegistry{
		AbstractEncodingRegistry: &AbstractEncodingRegistry{},
	}
}

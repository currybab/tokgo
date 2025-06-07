package tokgo

import "github.com/currybab/tokgo/mod"

type LazyEncodingRegistry struct {
	*AbstractEncodingRegistry
}

func (r *LazyEncodingRegistry) GetEncoding(encodingName string) (mod.Encoding, error) {
	encodingType, exists := mod.EncodingTypeFromName(encodingName)
	if exists {
		err := r.AddEncoding(encodingType)
		if err != nil {
			return nil, err
		}
	}
	return r.AbstractEncodingRegistry.GetEncoding(encodingName)
}

func (r *LazyEncodingRegistry) GetEncodingByType(encodingType mod.EncodingType) (mod.Encoding, error) {
	err := r.AddEncoding(encodingType)
	if err != nil {
		return nil, err
	}
	return r.AbstractEncodingRegistry.GetEncodingByType(encodingType)
}

func (r *LazyEncodingRegistry) GetEncodingForModelType(modelType mod.ModelType) (mod.Encoding, error) {
	err := r.AddEncoding(modelType.GetEncodingType())
	if err != nil {
		return nil, err
	}
	return r.AbstractEncodingRegistry.GetEncodingForModelType(modelType)
}

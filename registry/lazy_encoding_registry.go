package registry

import "github.com/nerdface-ai/tokgo"

type LazyEncodingRegistry struct {
	*AbstractEncodingRegistry
}

func (r *LazyEncodingRegistry) GetEncoding(encodingName string) (tokgo.Encoding, error) {
	encodingType, exists := tokgo.EncodingTypeFromName(encodingName)
	if exists {
		err := r.AddEncoding(encodingType)
		if err != nil {
			return nil, err
		}
	}
	return r.AbstractEncodingRegistry.GetEncoding(encodingName)
}

func (r *LazyEncodingRegistry) GetEncodingByType(encodingType tokgo.EncodingType) (tokgo.Encoding, error) {
	err := r.AddEncoding(encodingType)
	if err != nil {
		return nil, err
	}
	return r.AbstractEncodingRegistry.GetEncodingByType(encodingType)
}

func (r *LazyEncodingRegistry) GetEncodingForModelType(modelType tokgo.ModelType) (tokgo.Encoding, error) {
	err := r.AddEncoding(modelType.GetEncodingType())
	if err != nil {
		return nil, err
	}
	return r.AbstractEncodingRegistry.GetEncodingForModelType(modelType)
}

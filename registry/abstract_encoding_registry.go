package tokgo

import (
	"fmt"
	"sync"

	"github.com/currybab/tokgo/encoding"
	"github.com/currybab/tokgo/mod"
)

type AbstractEncodingRegistry struct {
	encodings sync.Map // map[string]mod.Encoding
}

func (a *AbstractEncodingRegistry) GetEncoding(encodingName string) (mod.Encoding, error) {
	encoding, exists := a.encodings.Load(encodingName)
	if !exists {
		return nil, fmt.Errorf("encoding %s not found", encodingName)
	}
	return encoding.(mod.Encoding), nil
}

func (a *AbstractEncodingRegistry) GetEncodingByType(encodingType mod.EncodingType) (mod.Encoding, error) {
	encoding, exists := a.encodings.Load(encodingType.GetName())
	if !exists {
		return nil, fmt.Errorf("encoding %s not found", encodingType.GetName())
	}
	return encoding.(mod.Encoding), nil
}

func (a *AbstractEncodingRegistry) GetEncodingForModel(modelName string) (mod.Encoding, error) {
	modelType, exists := mod.ModelTypeFromName(modelName)
	if exists {
		return a.GetEncodingForModelType(*modelType)
	}
	return nil, fmt.Errorf("model %s not found", modelName)
}

func (a *AbstractEncodingRegistry) GetEncodingForModelType(modelType mod.ModelType) (mod.Encoding, error) {
	encoding, exists := a.encodings.Load(modelType.GetEncodingType().GetName())
	if !exists {
		return nil, fmt.Errorf("encoding %s not found", modelType.GetEncodingType().GetName())
	}
	return encoding.(mod.Encoding), nil
}

func (a *AbstractEncodingRegistry) RegisterGptBytePairEncoding(parameters *mod.GptBytePairEncodingParams) (mod.EncodingRegistry, error) {
	return a.RegisterCustomEncoding(encoding.FromParameters(parameters))
}

func (a *AbstractEncodingRegistry) RegisterCustomEncoding(encoding mod.Encoding) (mod.EncodingRegistry, error) {
	encodingName := encoding.GetName()
	_, exists := a.encodings.Load(encodingName)
	if exists {
		return nil, fmt.Errorf("encoding %s already registered", encodingName)
	}
	a.encodings.Store(encodingName, encoding)
	return a, nil
}

func (a *AbstractEncodingRegistry) AddEncoding(encodingType mod.EncodingType) error {
	switch encodingType {
	case mod.R50K_BASE:
		a.encodings.Store(encodingType.GetName(), encoding.R50kBase())
	case mod.P50K_BASE:
		a.encodings.Store(encodingType.GetName(), encoding.P50kBase())
	case mod.P50K_EDIT:
		a.encodings.Store(encodingType.GetName(), encoding.P50kEdit())
	case mod.CL100K_BASE:
		a.encodings.Store(encodingType.GetName(), encoding.Cl100kBase())
	case mod.O200K_BASE:
		a.encodings.Store(encodingType.GetName(), encoding.O200kBase())
	default:
		return fmt.Errorf("unknown encoding type %s", encodingType.GetName())
	}
	return nil
}

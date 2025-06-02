package registry

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nerdface-ai/tokgo"
	"github.com/nerdface-ai/tokgo/encoding"
)

type AbstractEncodingRegistry struct {
	encodings sync.Map // map[string]tokgo.Encoding
}

func (a *AbstractEncodingRegistry) GetEncoding(encodingName string) (tokgo.Encoding, error) {
	encoding, exists := a.encodings.Load(encodingName)
	if !exists {
		return nil, fmt.Errorf("encoding %s not found", encodingName)
	}
	return encoding.(tokgo.Encoding), nil
}

func (a *AbstractEncodingRegistry) GetEncodingByType(encodingType tokgo.EncodingType) (tokgo.Encoding, error) {
	encoding, exists := a.encodings.Load(encodingType.GetName())
	if !exists {
		return nil, fmt.Errorf("encoding %s not found", encodingType.GetName())
	}
	return encoding.(tokgo.Encoding), nil
}

func (a *AbstractEncodingRegistry) GetEncodingForModel(modelName string) (tokgo.Encoding, error) {
	modelType, exists := tokgo.ModelTypeFromName(modelName)
	if exists {
		return a.GetEncodingForModelType(*modelType)
	}

	if strings.HasPrefix(modelName, tokgo.GPT_4O.GetName()) {
		return a.GetEncodingForModelType(tokgo.GPT_4O)
	}

	if strings.HasPrefix(modelName, tokgo.GPT_4_32K.GetName()) {
		return a.GetEncodingForModelType(tokgo.GPT_4_32K)
	}

	if strings.HasPrefix(modelName, tokgo.GPT_4.GetName()) {
		return a.GetEncodingForModelType(tokgo.GPT_4)
	}

	if strings.HasPrefix(modelName, tokgo.GPT_3_5_TURBO_16K.GetName()) {
		return a.GetEncodingForModelType(tokgo.GPT_3_5_TURBO_16K)
	}

	if strings.HasPrefix(modelName, tokgo.GPT_3_5_TURBO.GetName()) {
		return a.GetEncodingForModelType(tokgo.GPT_3_5_TURBO)
	}

	return nil, fmt.Errorf("model %s not found", modelName)
}

func (a *AbstractEncodingRegistry) GetEncodingForModelType(modelType tokgo.ModelType) (tokgo.Encoding, error) {
	encoding, exists := a.encodings.Load(modelType.GetEncodingType().GetName())
	if !exists {
		return nil, fmt.Errorf("encoding %s not found", modelType.GetEncodingType().GetName())
	}
	return encoding.(tokgo.Encoding), nil
}

func (a *AbstractEncodingRegistry) RegisterGptBytePairEncoding(parameters *tokgo.GptBytePairEncodingParams) (tokgo.EncodingRegistry, error) {
	return a.RegisterCustomEncoding(encoding.FromParameters(parameters))
}

func (a *AbstractEncodingRegistry) RegisterCustomEncoding(encoding tokgo.Encoding) (tokgo.EncodingRegistry, error) {
	encodingName := encoding.GetName()
	_, exists := a.encodings.Load(encodingName)
	if exists {
		return nil, fmt.Errorf("encoding %s already registered", encodingName)
	}
	a.encodings.Store(encodingName, encoding)
	return a, nil
}

func (a *AbstractEncodingRegistry) AddEncoding(encodingType tokgo.EncodingType) error {
	switch encodingType {
	case tokgo.R50K_BASE:
		a.encodings.Store(encodingType.GetName(), encoding.R50kBase())
	case tokgo.P50K_BASE:
		a.encodings.Store(encodingType.GetName(), encoding.P50kBase())
	case tokgo.P50K_EDIT:
		a.encodings.Store(encodingType.GetName(), encoding.P50kEdit())
	case tokgo.CL100K_BASE:
		a.encodings.Store(encodingType.GetName(), encoding.Cl100kBase())
	case tokgo.O200K_BASE:
		a.encodings.Store(encodingType.GetName(), encoding.O200kBase())
	default:
		return fmt.Errorf("unknown encoding type %s", encodingType.GetName())
	}
	return nil
}

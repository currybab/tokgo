package mod

type EncodingRegistry interface {
	GetEncoding(encodingName string) (Encoding, error)
	GetEncodingByType(encodingType EncodingType) (Encoding, error)
	GetEncodingForModel(modelName string) (Encoding, error)
	GetEncodingForModelType(modelType ModelType) (Encoding, error)
	RegisterGptBytePairEncoding(parameters *GptBytePairEncodingParams) (EncodingRegistry, error)
	RegisterCustomEncoding(encoding Encoding) (EncodingRegistry, error)
}

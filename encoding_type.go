package tokgo

type EncodingType string

const (
	R50K_BASE   EncodingType = "r50k_base"
	P50K_BASE   EncodingType = "p50k_base"
	P50K_EDIT   EncodingType = "p50k_edit"
	CL100K_BASE EncodingType = "cl100k_base"
	O200K_BASE  EncodingType = "o200k_base"
)

var encodingTypeMap = map[string]EncodingType{
	"r50k_base":   R50K_BASE,
	"p50k_base":   P50K_BASE,
	"p50k_edit":   P50K_EDIT,
	"cl100k_base": CL100K_BASE,
	"o200k_base":  O200K_BASE,
}

func (e EncodingType) GetName() string {
	return string(e)
}

func EncodingTypeFromName(name string) (EncodingType, bool) {
	val, exists := encodingTypeMap[name]
	return val, exists
}

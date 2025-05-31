package tokgo

type ModelType struct {
	name             string
	encodingType     EncodingType
	maxContextLength int
}

func (m *ModelType) GetName() string {
	return m.name
}

func (m *ModelType) GetEncodingType() EncodingType {
	return m.encodingType
}

func (m *ModelType) GetMaxContextLength() int {
	return m.maxContextLength
}

var (
	// Chat models
	GPT_4             = ModelType{"gpt-4", CL100K_BASE, 8192}
	GPT_4O            = ModelType{"gpt-4o", O200K_BASE, 128000}
	GPT_4O_MINI       = ModelType{"gpt-4o-mini", O200K_BASE, 128000}
	GPT_4_32K         = ModelType{"gpt-4-32k", CL100K_BASE, 32768}
	GPT_4_TURBO       = ModelType{"gpt-4-turbo", CL100K_BASE, 128000}
	GPT_3_5_TURBO     = ModelType{"gpt-3.5-turbo", CL100K_BASE, 16385}
	GPT_3_5_TURBO_16K = ModelType{"gpt-3.5-turbo-16k", CL100K_BASE, 16385}

	// Text models
	TEXT_DAVINCI_003 = ModelType{"text-davinci-003", P50K_BASE, 4097}
	TEXT_DAVINCI_002 = ModelType{"text-davinci-002", P50K_BASE, 4097}
	TEXT_DAVINCI_001 = ModelType{"text-davinci-001", R50K_BASE, 2049}
	TEXT_CURIE_001   = ModelType{"text-curie-001", R50K_BASE, 2049}
	TEXT_BABBAGE_001 = ModelType{"text-babbage-001", R50K_BASE, 2049}
	TEXT_ADA_001     = ModelType{"text-ada-001", R50K_BASE, 2049}
	DAVINCI          = ModelType{"davinci", R50K_BASE, 2049}
	CURIE            = ModelType{"curie", R50K_BASE, 2049}
	BABBAGE          = ModelType{"babbage", R50K_BASE, 2049}
	ADA              = ModelType{"ada", R50K_BASE, 2049}

	// Code models
	CODE_DAVINCI_002 = ModelType{"code-davinci-002", P50K_BASE, 8001}
	CODE_DAVINCI_001 = ModelType{"code-davinci-001", P50K_BASE, 8001}
	CODE_CUSHMAN_002 = ModelType{"code-cushman-002", P50K_BASE, 2048}
	CODE_CUSHMAN_001 = ModelType{"code-cushman-001", P50K_BASE, 2048}
	DAVINCI_CODEX    = ModelType{"davinci-codex", P50K_BASE, 4096}
	CUSHMAN_CODEX    = ModelType{"cushman-codex", P50K_BASE, 2048}

	// Edit models
	TEXT_DAVINCI_EDIT_001 = ModelType{"text-davinci-edit-001", P50K_EDIT, 3000}
	CODE_DAVINCI_EDIT_001 = ModelType{"code-davinci-edit-001", P50K_EDIT, 3000}

	// Embeddings
	TEXT_EMBEDDING_ADA_002 = ModelType{"text-embedding-ada-002", CL100K_BASE, 8191}
	TEXT_EMBEDDING_3_SMALL = ModelType{"text-embedding-3-small", CL100K_BASE, 8191}
	TEXT_EMBEDDING_3_LARGE = ModelType{"text-embedding-3-large", CL100K_BASE, 8191}

	// Old embeddings
	TEXT_SIMILARITY_DAVINCI_001  = ModelType{"text-similarity-davinci-001", R50K_BASE, 2046}
	TEXT_SIMILARITY_CURIE_001    = ModelType{"text-similarity-curie-001", R50K_BASE, 2046}
	TEXT_SIMILARITY_BABBAGE_001  = ModelType{"text-similarity-babbage-001", R50K_BASE, 2046}
	TEXT_SIMILARITY_ADA_001      = ModelType{"text-similarity-ada-001", R50K_BASE, 2046}
	TEXT_SEARCH_DAVINCI_DOC_001  = ModelType{"text-search-davinci-doc-001", R50K_BASE, 2046}
	TEXT_SEARCH_CURIE_DOC_001    = ModelType{"text-search-curie-doc-001", R50K_BASE, 2046}
	TEXT_SEARCH_BABBAGE_DOC_001  = ModelType{"text-search-babbage-doc-001", R50K_BASE, 2046}
	TEXT_SEARCH_ADA_DOC_001      = ModelType{"text-search-ada-doc-001", R50K_BASE, 2046}
	CODE_SEARCH_BABBAGE_CODE_001 = ModelType{"code-search-babbage-code-001", R50K_BASE, 2046}
	CODE_SEARCH_ADA_CODE_001     = ModelType{"code-search-ada-code-001", R50K_BASE, 2046}
)

var nameToModelType = map[string]*ModelType{
	"gpt-4":                        &GPT_4,
	"gpt-4o":                       &GPT_4O,
	"gpt-4o-mini":                  &GPT_4O_MINI,
	"gpt-4-32k":                    &GPT_4_32K,
	"gpt-4-turbo":                  &GPT_4_TURBO,
	"gpt-3.5-turbo":                &GPT_3_5_TURBO,
	"gpt-3.5-turbo-16k":            &GPT_3_5_TURBO_16K,
	"text-davinci-003":             &TEXT_DAVINCI_003,
	"text-davinci-002":             &TEXT_DAVINCI_002,
	"text-davinci-001":             &TEXT_DAVINCI_001,
	"text-curie-001":               &TEXT_CURIE_001,
	"text-babbage-001":             &TEXT_BABBAGE_001,
	"text-ada-001":                 &TEXT_ADA_001,
	"davinci":                      &DAVINCI,
	"curie":                        &CURIE,
	"babbage":                      &BABBAGE,
	"ada":                          &ADA,
	"code-davinci-002":             &CODE_DAVINCI_002,
	"code-davinci-001":             &CODE_DAVINCI_001,
	"code-cushman-002":             &CODE_CUSHMAN_002,
	"code-cushman-001":             &CODE_CUSHMAN_001,
	"davinci-codex":                &DAVINCI_CODEX,
	"cushman-codex":                &CUSHMAN_CODEX,
	"text-davinci-edit-001":        &TEXT_DAVINCI_EDIT_001,
	"code-davinci-edit-001":        &CODE_DAVINCI_EDIT_001,
	"text-embedding-ada-002":       &TEXT_EMBEDDING_ADA_002,
	"text-embedding-3-small":       &TEXT_EMBEDDING_3_SMALL,
	"text-embedding-3-large":       &TEXT_EMBEDDING_3_LARGE,
	"text-similarity-davinci-001":  &TEXT_SIMILARITY_DAVINCI_001,
	"text-similarity-curie-001":    &TEXT_SIMILARITY_CURIE_001,
	"text-similarity-babbage-001":  &TEXT_SIMILARITY_BABBAGE_001,
	"text-similarity-ada-001":      &TEXT_SIMILARITY_ADA_001,
	"text-search-davinci-doc-001":  &TEXT_SEARCH_DAVINCI_DOC_001,
	"text-search-curie-doc-001":    &TEXT_SEARCH_CURIE_DOC_001,
	"text-search-babbage-doc-001":  &TEXT_SEARCH_BABBAGE_DOC_001,
	"text-search-ada-doc-001":      &TEXT_SEARCH_ADA_DOC_001,
	"code-search-babbage-code-001": &CODE_SEARCH_BABBAGE_CODE_001,
	"code-search-ada-code-001":     &CODE_SEARCH_ADA_CODE_001,
}

func ModelTypeFromName(name string) (*ModelType, bool) {
	model, exists := nameToModelType[name]
	return model, exists
}

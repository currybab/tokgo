package examples

import (
	"testing"

	tokmod "github.com/currybab/tokgo/mod"
	tokgo "github.com/currybab/tokgo/registry"
)

func TestWorking(t *testing.T) {
	registry := tokgo.NewDefaultEncodingRegistry()
	enc, _ := registry.GetEncodingByType(tokmod.CL100K_BASE)
	if tt := enc.Decode(enc.EncodeToIntArray("hello world")); tt != "hello world" {
		panic(tt)
	}

	if tt := enc.Decode(enc.EncodeToIntArray("안녕하세요.")); tt != "안녕하세요." {
		panic(tt)
	}

	// Or get the tokenizer corresponding to a specific OpenAI model
	enc, _ = registry.GetEncodingForModelType(tokmod.TEXT_EMBEDDING_ADA_002)
	if tt := enc.Decode(enc.Encode("hello world", 10).GetTokens()); tt != "hello world" {
		panic(tt)
	}

	enc, _ = registry.GetEncodingForModelType(tokmod.GPT_4O)
	if tt := enc.Decode(enc.Encode("hello world", 10).GetTokens()); tt != "hello world" {
		panic(tt)
	}

	if tt := enc.Decode(enc.EncodeToIntArray("안녕하세요.")); tt != "안녕하세요." {
		panic(tt)
	}
	t.Logf("%v", enc.Decode(enc.EncodeToIntArray("안녕하세요.")))
}

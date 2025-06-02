package main

import (
	"github.com/nerdface-ai/tokgo"
	"github.com/nerdface-ai/tokgo/registry"
)

func main() {
	reg := registry.NewDefaultEncodingRegistry()
	enc, _ := reg.GetEncodingByType(tokgo.CL100K_BASE)
	if enc.Decode(enc.EncodeToIntArray("hello world")) != "hello world" {
		panic("hello world")
	}

	// Or get the tokenizer corresponding to a specific OpenAI model
	enc, _ = reg.GetEncodingForModelType(tokgo.TEXT_EMBEDDING_ADA_002)
	if enc.Decode(enc.Encode("hello world", 10).GetTokens()) != "hello world" {
		panic("hello world")
	}
}

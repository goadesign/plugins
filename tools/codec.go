package tools

type (
	// JSONCodec is a generic interface for marshaling and unmarshaling JSON values
	// for tool payloads and results. Generated code in service toolsets uses JSONCodec
	// to encode Go structs into canonical JSON (for sending tool payloads to LLMs or APIs)
	// and decode JSON responses into the appropriate Go struct types. For each tool,
	// generated code provides strongly-typed JSONCodec values (e.g., JSONCodec[*MyPayload])
	// to ensure compile-time type safety when serializing and deserializing tool inputs
	// and outputs.
	JSONCodec[T any] struct {
		ToJSON   func(v T) ([]byte, error)
		FromJSON func(data []byte) (T, error)
	}
)

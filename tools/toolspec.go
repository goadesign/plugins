package tools

type (
	// ToolSpec captures metadata for a registered tool, including payload
	// and result specs.
	ToolSpec struct {
		// Name of the tool.
		Name string
		// Name of the tool set that the tool is defined in.
		Set string
		// Name of the service that the set is defined in.
		Service string
		// Name of the method that the tool is derived from, if any.
		Method *string
		// Payload type specification.
		Payload TypeSpec
		// Result type specification.
		Result TypeSpec
	}

	// TypeSpec describes the name and schema for a tool payload or result.
	TypeSpec struct {
		// Name of the type.
		Name string
		// JSON schema for the type.
		Schema []byte
		// Codec is a JSON codec for the type.
		Codec JSONCodec[any]
	}
)

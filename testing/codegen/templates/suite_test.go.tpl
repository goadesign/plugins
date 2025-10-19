{{- printf "Run%sHarness exercises the generated harness against your service implementation." .Service.Name | comment }}
{{- printf "Call this helper from your test, passing your service implementation." | comment }}
func Run{{ .Service.StructName }}Harness(t *testing.T, svc {{ .Service.PkgName }}.Service) {
	t.Helper()
{{- if .UseCtx }}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
{{- end }}

	h := {{ .TestPkg }}.NewHarness(t, svc)
	defer h.Close()
	
	td := {{ .TestPkg }}.NewTestData()

{{- range .NonStream }}
	{{- $m := . }}
	t.Run("{{ $m.Method.Name }}", func(t *testing.T) {
		{{ if $m.Method.ResultRef }}result, {{ end }}err := h.Client.{{ $m.Method.VarName }}({{ if $.UseCtx }}ctx{{ else }}context.Background(){{ end }}{{- if $m.Method.PayloadEx }}, td.Valid{{ goify $m.Method.Name true }}Payload(){{ end }})
		if err != nil {
			t.Errorf("{{ $m.Method.Name }} failed: %v", err)
		}
	{{- if $m.Method.ResultRef }}
		if result == nil {
			t.Error("{{ $m.Method.Name }} returned nil result")
		}
	{{- end }}
	})
{{- end }}

{{- range .Stream }}
	{{- $m := . }}
	t.Run("{{ $m.Method.Name }}_Stream", func(t *testing.T) {
	{{- if $m.Method.PayloadEx }}
		stream, err := h.Client.{{ $m.Method.VarName }}({{ if $.UseCtx }}ctx{{ else }}context.Background(){{ end }}, td.Valid{{ goify $m.Method.Name true }}Payload())
	{{- else }}
		stream, err := h.Client.{{ $m.Method.VarName }}({{ if $.UseCtx }}ctx{{ else }}context.Background(){{ end }})
	{{- end }}
		if err != nil {
			t.Errorf("Failed to create {{ $m.Method.Name }} stream: %v", err)
		}
		if stream == nil {
			t.Fatal("{{ $m.Method.Name }} returned nil stream")
		}
		
	{{- if eq $m.Method.StreamKind 3 }}
		// Server stream - receive at least one message
		{{- if $m.Method.StreamingResultRef }}
		msg, err := stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("{{ $m.Method.Name }} recv failed: %v", err)
		}
		if err != io.EOF && msg == nil {
			t.Error("{{ $m.Method.Name }} recv returned nil message without EOF")
		}
		{{- else }}
		_, err = stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("{{ $m.Method.Name }} recv failed: %v", err)
		}
		{{- end }}
	{{- else if eq $m.Method.StreamKind 2 }}
		// Client stream - send test data and close
		{{- if $m.Method.StreamingPayloadRef }}
		// Stream has typed payloads, send multiple
		for i := 0; i < 3; i++ {
			payload := td.Valid{{ goify $m.Method.Name true }}Payload()
			if err := stream.Send(payload); err != nil {
				t.Errorf("{{ $m.Method.Name }} send failed: %v", err)
				break
			}
		}
		{{- end }}
		result, err := stream.CloseAndRecv()
		if err != nil {
			t.Errorf("{{ $m.Method.Name }} close and recv failed: %v", err)
		}
		{{- if $m.Method.ResultRef }}
		if result == nil {
			t.Error("{{ $m.Method.Name }} returned nil result")
		}
		{{- end }}
	{{- else if eq $m.Method.StreamKind 4 }}
		// Bidirectional stream - send and receive
		{{- if $m.Method.StreamingPayloadRef }}
		// Send a test message
		payload := td.Valid{{ goify $m.Method.Name true }}Payload()
		if err := stream.Send(payload); err != nil {
			t.Errorf("{{ $m.Method.Name }} send failed: %v", err)
		}
		{{- end }}
		
		// Try to receive a response
		{{- if $m.Method.StreamingResultRef }}
		msg, err := stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("{{ $m.Method.Name }} recv failed: %v", err)
		}
		if err != io.EOF && msg == nil {
			t.Error("{{ $m.Method.Name }} recv returned nil message without EOF")
		}
		{{- else }}
		_, err = stream.Recv()
		if err != nil && err != io.EOF {
			t.Errorf("{{ $m.Method.Name }} recv failed: %v", err)
		}
		{{- end }}
		
		// Close the stream
		if err := stream.Close(); err != nil {
			t.Errorf("{{ $m.Method.Name }} close failed: %v", err)
		}
	{{- end }}
	})
{{- end }}
}

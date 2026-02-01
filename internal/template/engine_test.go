package template

import (
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("expected engine, got nil")
	}
	if engine.templates == nil {
		t.Fatal("expected templates map, got nil")
	}
}

func TestAddTemplate_Valid(t *testing.T) {
	engine := NewEngine()
	err := engine.AddTemplate("test", "Hello {{ .name }}")
	if err != nil {
		t.Errorf("failed to add template: %v", err)
	}
}

func TestAddTemplate_Invalid(t *testing.T) {
	engine := NewEngine()
	err := engine.AddTemplate("test", "Hello {{ .name")
	if err == nil {
		t.Error("expected error for invalid template, got nil")
	}
}

func TestRender_Success(t *testing.T) {
	engine := NewEngine()
	_ = engine.AddTemplate("greeting", "Hello {{ .name }}!")

	data := map[string]interface{}{
		"name": "World",
	}

	result, err := engine.Render("greeting", data)
	if err != nil {
		t.Fatalf("failed to render: %v", err)
	}

	expected := "Hello World!"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestRender_MissingTemplate(t *testing.T) {
	engine := NewEngine()
	_, err := engine.Render("nonexistent", nil)
	if err == nil {
		t.Error("expected error for missing template, got nil")
	}
}

func TestRender_ComplexTemplate(t *testing.T) {
	engine := NewEngine()
	_ = engine.AddTemplate("cert", "{{ .tlsCrt }}")
	_ = engine.AddTemplate("key", "{{ .tlsKey }}")

	data := map[string]interface{}{
		"tlsCrt": "CERT_DATA",
		"tlsKey": "KEY_DATA",
	}

	cert, err := engine.Render("cert", data)
	if err != nil {
		t.Fatalf("failed to render cert: %v", err)
	}
	if cert != "CERT_DATA" {
		t.Errorf("expected 'CERT_DATA', got '%s'", cert)
	}

	key, err := engine.Render("key", data)
	if err != nil {
		t.Fatalf("failed to render key: %v", err)
	}
	if key != "KEY_DATA" {
		t.Errorf("expected 'KEY_DATA', got '%s'", key)
	}
}

func TestRenderAll_Success(t *testing.T) {
	engine := NewEngine()
	_ = engine.AddTemplate("username", "{{ .username }}")
	_ = engine.AddTemplate("password", "{{ .password }}")

	data := map[string]interface{}{
		"username": "admin",
		"password": "secret",
	}

	results, err := engine.RenderAll(data)
	if err != nil {
		t.Fatalf("failed to render all: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if results["username"] != "admin" {
		t.Errorf("expected username 'admin', got '%s'", results["username"])
	}

	if results["password"] != "secret" {
		t.Errorf("expected password 'secret', got '%s'", results["password"])
	}
}

func TestRender_MissingField(t *testing.T) {
	engine := NewEngine()
	_ = engine.AddTemplate("test", "{{ .missing }}")

	data := map[string]interface{}{
		"other": "value",
	}

	result, err := engine.Render("test", data)
	if err != nil {
		t.Fatalf("failed to render: %v", err)
	}

	if result != "<no value>" {
		t.Errorf("expected '<no value>', got '%s'", result)
	}
}

func TestRender_MultilineTemplate(t *testing.T) {
	engine := NewEngine()
	tmpl := `-----BEGIN CERTIFICATE-----
{{ .cert }}
-----END CERTIFICATE-----`
	_ = engine.AddTemplate("cert", tmpl)

	data := map[string]interface{}{
		"cert": "MIICljCCAX4CCQCKz8Zr8vJKZDANBg",
	}

	result, err := engine.Render("cert", data)
	if err != nil {
		t.Fatalf("failed to render: %v", err)
	}

	expected := `-----BEGIN CERTIFICATE-----
MIICljCCAX4CCQCKz8Zr8vJKZDANBg
-----END CERTIFICATE-----`
	if result != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

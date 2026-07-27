package kufar

import "testing"

func TestParseSearchURL_ValidMinskRent(t *testing.T) {
	raw := "https://re.kufar.by/l/minsk/snyat/kvartiru?cur=USD&prc=r%3A0%2C500"
	params, err := ParseSearchURL(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if params["cat"] != "1010" {
		t.Fatalf("cat=%q", params["cat"])
	}
	if params["typ"] != "let" {
		t.Fatalf("typ=%q", params["typ"])
	}
	if params["gtsy"] == "" {
		t.Fatal("gtsy empty")
	}
}

func TestValidateSearchURL_RejectsNonKufar(t *testing.T) {
	err := ValidateSearchURL("https://example.com/search")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateSearchURL_RejectsIncomplete(t *testing.T) {
	err := ValidateSearchURL("https://re.kufar.by/l/minsk")
	if err == nil {
		t.Fatal("expected error for incomplete search url")
	}
}

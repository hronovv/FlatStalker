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
	if params["prc"] != "r:0,500" {
		t.Fatalf("prc=%q", params["prc"])
	}
}

func TestParseSearchURL_SkipsTrackingParams(t *testing.T) {
	raw := "https://re.kufar.by/l/minsk/snyat/kvartiru?cur=USD&_gl=1abc&r_track=xyz&prc=r%3A0%2C500"
	params, err := ParseSearchURL(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := params["_gl"]; ok {
		t.Fatal("expected _gl to be skipped")
	}
	if _, ok := params["r_track"]; ok {
		t.Fatal("expected r_track to be skipped")
	}
	if params["prc"] != "r:0,500" {
		t.Fatalf("prc=%q", params["prc"])
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

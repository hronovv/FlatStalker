package kufar

import (
	"strings"
	"testing"
)

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

func TestParseSearchURL_MinskMicrodistrict(t *testing.T) {
	raw := "https://re.kufar.by/l/minsk-kuncevshchina-mkrn/snyat/kvartiru?cur=USD"
	params, err := ParseSearchURL(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if params["gtsy"] != gtsyMinsk {
		t.Fatalf("gtsy=%q", params["gtsy"])
	}
	if params["rgn"] != "7" {
		t.Fatalf("rgn=%q", params["rgn"])
	}
	if params["red"] != "v.or:130" {
		t.Fatalf("red=%q", params["red"])
	}
	if params["mee"] != "" {
		t.Fatalf("metro should not be set together with microdistrict, mee=%q", params["mee"])
	}
	if params["cur"] != "USD" {
		t.Fatalf("cur=%q", params["cur"])
	}
}

func TestParseSearchURL_MinskMetro(t *testing.T) {
	raw := "https://re.kufar.by/l/minsk/snyat/kvartiru/metro-kuncevshchina"
	params, err := ParseSearchURL(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if params["gtsy"] != gtsyMinsk {
		t.Fatalf("gtsy=%q", params["gtsy"])
	}
	if params["mee"] != "v.or:9" {
		t.Fatalf("mee=%q", params["mee"])
	}
	if params["red"] != "" {
		t.Fatalf("microdistrict should not be set together with metro, red=%q", params["red"])
	}
}

func TestParseSearchURL_BrestCity(t *testing.T) {
	raw := "https://re.kufar.by/l/brest/snyat/kvartiru?cur=USD"
	params, err := ParseSearchURL(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(params["gtsy"], "locality-brest") {
		t.Fatalf("gtsy=%q", params["gtsy"])
	}
	if params["rgn"] != "1" {
		t.Fatalf("rgn=%q", params["rgn"])
	}
	if params["ar"] != "1" {
		t.Fatalf("ar=%q", params["ar"])
	}
}

func TestParseSearchURL_RoomsPath(t *testing.T) {
	raw := "https://re.kufar.by/l/minsk/snyat/kvartiru/1k"
	params, err := ParseSearchURL(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if params["rms"] != "v.or:1" {
		t.Fatalf("rms=%q", params["rms"])
	}
}

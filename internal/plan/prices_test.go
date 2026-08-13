package plan

import "testing"

func TestPriceCatalog(t *testing.T) {
	cat := PriceCatalog()
	if cat.Currency != "BYN" {
		t.Fatalf("currency: %s", cat.Currency)
	}
	if cat.Plus["30"] != "9.90" || cat.Pro["30"] != "19.90" {
		t.Fatalf("30-day shop prices: plus=%s pro=%s", cat.Plus["30"], cat.Pro["30"])
	}
	if cat.Plus["1"] != "0.70" || cat.Pro["365"] != "144.00" {
		t.Fatalf("ends: plus1=%s pro365=%s", cat.Plus["1"], cat.Pro["365"])
	}
	if len(cat.PeriodDays) != len(Periods) {
		t.Fatalf("periods: %v", cat.PeriodDays)
	}
}

func TestPriceKopUnknown(t *testing.T) {
	if _, ok := PriceKop(Plus, 2); ok {
		t.Fatal("2 days should be missing")
	}
	if kop, ok := PriceKop(Free, 0); !ok || kop != 0 {
		t.Fatalf("free: kop=%d ok=%v", kop, ok)
	}
}

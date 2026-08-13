package plan

import (
	"fmt"
	"strconv"
)

const Currency = "BYN"

// Periods are the paid durations, in days.
var Periods = []int{1, 3, 5, 7, 15, 30, 90, 180, 365}

// Kop is the price in kopecks (1 BYN = 100). 30-day PLUS is the anchor;
// shorter periods cost more per day, longer ones less. PRO is 2× PLUS
// except the 30-day shop prices 9.90 / 19.90.
var plusKop = map[int]int{
	1:   70,
	3:   160,
	5:   230,
	7:   290,
	15:  550,
	30:  990,
	90:  2520,
	180: 4450,
	365: 7200,
}

var proKop = map[int]int{
	1:   140,
	3:   320,
	5:   460,
	7:   580,
	15:  1100,
	30:  1990,
	90:  5040,
	180: 8900,
	365: 14400,
}

func PriceKop(name string, days int) (int, bool) {
	switch Normalize(name) {
	case Plus:
		kop, ok := plusKop[days]
		return kop, ok
	case Pro:
		kop, ok := proKop[days]
		return kop, ok
	default:
		return 0, days == 0
	}
}

func FormatBYN(kop int) string {
	return fmt.Sprintf("%d.%02d", kop/100, kop%100)
}

// Catalog is the public price list for Mini App / invoices.
type Catalog struct {
	Currency   string            `json:"currency"`
	PeriodDays []int             `json:"period_days"`
	Plus       map[string]string `json:"plus"`
	Pro        map[string]string `json:"pro"`
}

func PriceCatalog() Catalog {
	plus := make(map[string]string, len(Periods))
	pro := make(map[string]string, len(Periods))
	for _, days := range Periods {
		key := strconv.Itoa(days)
		plus[key] = FormatBYN(plusKop[days])
		pro[key] = FormatBYN(proKop[days])
	}
	return Catalog{
		Currency:   Currency,
		PeriodDays: append([]int(nil), Periods...),
		Plus:       plus,
		Pro:        pro,
	}
}

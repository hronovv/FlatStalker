package kufar

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Ad is a single Kufar listing from the search API.
type Ad struct {
	AdID      int64
	URL       string
	Subject   string
	PriceBYN  *int64
	Rooms     string
	Size      *int
	Address   string
	Metro     []string
	BodyShort string
	ListTime  string
}

type apiParam struct {
	P  string          `json:"p"`
	V  json.RawMessage `json:"v"`
	VL json.RawMessage `json:"vl"`
}

type apiAd struct {
	AdID              int64           `json:"ad_id"`
	AdLink            string          `json:"ad_link"`
	Subject           string          `json:"subject"`
	PriceBYN          json.RawMessage `json:"price_byn"`
	BodyShort         string          `json:"body_short"`
	ListTime          string          `json:"list_time"`
	AdParameters      []apiParam      `json:"ad_parameters"`
	AccountParameters []apiParam      `json:"account_parameters"`
}

func AdFromAPI(raw json.RawMessage) (Ad, error) {
	var ad apiAd
	if err := json.Unmarshal(raw, &ad); err != nil {
		return Ad{}, err
	}

	params := indexParams(ad.AdParameters)
	account := indexParams(ad.AccountParameters)

	metro := parseStringList(params["metro"].VL)
	if len(metro) == 0 {
		metro = parseStringList(params["metro"].V)
	}

	var size *int
	if sizeValue := stringValue(params["size"].V); sizeValue != "" {
		if parsed, err := strconv.Atoi(sizeValue); err == nil {
			size = &parsed
		}
	}

	priceBYN := parseOptionalInt64(ad.PriceBYN)

	listingURL := ad.AdLink
	if listingURL == "" {
		listingURL = fmt.Sprintf("https://re.kufar.by/vi/%d", ad.AdID)
	}

	subject := strings.TrimSpace(ad.Subject)
	if subject == "" {
		subject = "Без названия"
	}

	return Ad{
		AdID:      ad.AdID,
		URL:       listingURL,
		Subject:   subject,
		PriceBYN:  priceBYN,
		Rooms:     stringValue(params["rooms"].VL),
		Size:      size,
		Address:   stringValue(account["address"].V),
		Metro:     metro,
		BodyShort: ad.BodyShort,
		ListTime:  ad.ListTime,
	}, nil
}

type indexedParam struct {
	V  json.RawMessage
	VL json.RawMessage
}

func indexParams(params []apiParam) map[string]indexedParam {
	result := make(map[string]indexedParam, len(params))
	for _, param := range params {
		result[param.P] = indexedParam{V: param.V, VL: param.VL}
	}
	return result
}

func parseStringList(raw json.RawMessage) []string {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var items []string
	if err := json.Unmarshal(raw, &items); err == nil {
		return items
	}

	var single string
	if err := json.Unmarshal(raw, &single); err == nil && single != "" {
		return []string{single}
	}

	return nil
}

func stringValue(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}

	var value string
	if err := json.Unmarshal(raw, &value); err == nil {
		return value
	}

	return strings.Trim(string(raw), `"`)
}

func parseOptionalInt64(raw json.RawMessage) *int64 {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var value int64
	if err := json.Unmarshal(raw, &value); err == nil {
		return &value
	}

	asString := stringValue(raw)
	if asString == "" {
		return nil
	}

	parsed, err := strconv.ParseInt(asString, 10, 64)
	if err != nil {
		return nil
	}

	return &parsed
}

func (a Ad) FormatPrice() string {
	if a.PriceBYN == nil {
		return "цена не указана"
	}

	rub := float64(*a.PriceBYN) / 100
	if rub == float64(int64(rub)) {
		return fmt.Sprintf("%d р./мес.", int64(rub))
	}

	return fmt.Sprintf("%.2f р./мес.", rub)
}

func (a Ad) FormatListTime() string {
	if a.ListTime == "" {
		return ""
	}

	parsed, err := time.Parse(time.RFC3339, a.ListTime)
	if err != nil {
		parsed, err = time.Parse("2006-01-02T15:04:05Z", a.ListTime)
		if err != nil {
			return a.ListTime
		}
	}

	location, err := time.LoadLocation("Europe/Minsk")
	if err != nil {
		location = time.UTC
	}

	return parsed.In(location).Format("02.01.2006, 15:04")
}

func (a Ad) FormatMessage() string {
	lines := []string{
		"🏠 " + a.Subject,
		"💰 " + a.FormatPrice(),
	}

	if published := a.FormatListTime(); published != "" {
		lines = append(lines, "🕒 "+published)
	}

	if a.Rooms != "" {
		details := []string{a.Rooms + " комн."}
		if a.Size != nil {
			details = append(details, fmt.Sprintf("%d м²", *a.Size))
		}
		lines = append(lines, "📐 "+strings.Join(details, ", "))
	}

	if a.Address != "" {
		lines = append(lines, "📍 "+a.Address)
	}

	if len(a.Metro) > 0 {
		lines = append(lines, "🚇 "+strings.Join(a.Metro, ", "))
	}

	if a.BodyShort != "" {
		short := strings.TrimSpace(a.BodyShort)
		if len([]rune(short)) > 200 {
			short = string([]rune(short)[:197]) + "..."
		}
		lines = append(lines, "", short)
	}

	lines = append(lines, "", a.URL)
	return strings.Join(lines, "\n")
}

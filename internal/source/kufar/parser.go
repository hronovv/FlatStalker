package kufar

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

const APIBase = "https://api.kufar.by/search-api/v2/search/rendered-paginated"

var pathDefaults = map[string]map[string]string{
	"minsk": {
		"gtsy": "country-belarus~province-minsk~locality-minsk",
		"rgn":  "7",
	},
	"snyat":           {"typ": "let"},
	"kvartiru":        {"cat": "1010"},
	"bez-posrednikov": {"cmp": "0"},
}

// ParseSearchURL turns a Kufar search page URL into API query params.
func ParseSearchURL(rawURL string) (map[string]string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("некорректная ссылка")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("ссылка должна начинаться с https://")
	}
	host := strings.ToLower(parsed.Host)
	if !strings.Contains(host, "kufar.by") {
		return nil, fmt.Errorf("нужна ссылка на поиск kufar.by")
	}

	params := make(map[string]string)

	for segment := range strings.SplitSeq(strings.Trim(parsed.Path, "/"), "/") {
		if strings.HasPrefix(segment, "l") {
			continue
		}
		if defaults, ok := pathDefaults[segment]; ok {
			for key, value := range defaults {
				params[key] = value
			}
		}
	}

	query := parsed.Query()
	for key, values := range query {
		if len(values) == 0 {
			continue
		}
		decoded, err := url.QueryUnescape(values[len(values)-1])
		if err != nil {
			decoded = values[len(values)-1]
		}
		params[key] = decoded
	}

	if _, ok := params["size"]; !ok {
		params["size"] = "30"
	}
	if _, ok := params["sort"]; !ok {
		params["sort"] = "lst.d"
	}
	if _, ok := params["cur"]; !ok {
		params["cur"] = "BYR"
	}

	required := []string{"cat", "typ", "gtsy"}
	var missing []string
	for _, name := range required {
		if _, ok := params[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf(
			"это не ссылка на поиск аренды Kufar (не хватает: %s)",
			strings.Join(missing, ", "),
		)
	}

	return params, nil
}

// ValidateSearchURL checks that the URL can be used as a watch target.
func ValidateSearchURL(rawURL string) error {
	_, err := ParseSearchURL(strings.TrimSpace(rawURL))
	return err
}

func BuildAPIURL(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", key, params[key]))
	}

	return APIBase + "?" + strings.Join(parts, "&")
}

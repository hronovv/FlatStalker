package kufar

import (
	"fmt"
	"maps"
	"net/url"
	"sort"
	"strings"
	"unicode"
)

const APIBase = "https://api.kufar.by/search-api/v2/search/rendered-paginated"

func isAPIQueryParam(key string) bool {
	return !strings.HasPrefix(key, "_") && !strings.HasPrefix(key, "r_")
}

func isListingPrefix(segment string) bool {
	if segment == "l" {
		return true
	}
	if !strings.HasPrefix(segment, "l") || len(segment) < 2 {
		return false
	}
	for _, r := range segment[1:] {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
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
	unknownGeo := false

	for segment := range strings.SplitSeq(strings.Trim(parsed.Path, "/"), "/") {
		segment = strings.ToLower(strings.TrimSpace(segment))
		if segment == "" || isListingPrefix(segment) {
			continue
		}
		if extra, ok := pathTokens[segment]; ok {
			maps.Copy(params, extra)
			continue
		}
		if slug, ok := strings.CutPrefix(segment, "metro-"); ok {
			if mee, found := metroIDs[slug]; found {
				params["mee"] = mee
			} else {
				unknownGeo = true
			}
			continue
		}
		if extra, ok := placeParams[segment]; ok {
			maps.Copy(params, extra)
			continue
		}
		unknownGeo = true
	}

	query := parsed.Query()
	for key, values := range query {
		if !isAPIQueryParam(key) || len(values) == 0 {
			continue
		}
		decoded, err := url.QueryUnescape(values[len(values)-1])
		if err != nil {
			decoded = values[len(values)-1]
		}
		params[key] = decoded
	}

	if params["cat"] != "" && params["typ"] != "" && (params["gtsy"] == "" || unknownGeo) {
		_ = resolveMissingGeo(rawURL, params)
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

	if params["gtsy"] == gtsyBelarus && !pathHasSegment(parsed.Path, "belarus") {
		return nil, fmt.Errorf("не удалось определить город по ссылке Kufar")
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

func pathHasSegment(path, want string) bool {
	for segment := range strings.SplitSeq(strings.Trim(path, "/"), "/") {
		if strings.EqualFold(segment, want) {
			return true
		}
	}
	return false
}

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

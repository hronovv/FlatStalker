package kufar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	nextDataRe    = regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)
	resolveClient = &http.Client{Timeout: 12 * time.Second}
	resolvedGeo   sync.Map // path -> map[string]string
)

var geoQueryKeys = []string{"gtsy", "rgn", "ar", "red", "mee"}

func resolveMissingGeo(rawURL string, params map[string]string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	path := strings.ToLower(parsed.EscapedPath())
	if path == "" {
		path = strings.ToLower(parsed.Path)
	}
	if cached, ok := resolvedGeo.Load(path); ok {
		mergeMissing(params, cached.(map[string]string))
		return nil
	}

	geo, err := fetchPageGeo(rawURL)
	if err != nil {
		return err
	}
	resolvedGeo.Store(path, geo)
	mergeMissing(params, geo)
	return nil
}

func fetchPageGeo(rawURL string) (map[string]string, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "text/html")

	resp, err := resolveClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть ссылку Kufar: %w", err)
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		resp.Body.Close()
		time.Sleep(2 * time.Second)
		resp, err = resolveClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("не удалось открыть ссылку Kufar: %w", err)
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("kufar вернул %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать страницу Kufar: %w", err)
	}

	match := nextDataRe.FindSubmatch(body)
	if match == nil {
		return nil, fmt.Errorf("не удалось разобрать гео из ссылки Kufar")
	}

	var payload struct {
		Query map[string]any `json:"query"`
	}
	if err := json.Unmarshal(match[1], &payload); err != nil {
		return nil, fmt.Errorf("не удалось разобрать гео из ссылки Kufar")
	}

	geo := make(map[string]string, len(geoQueryKeys))
	for _, key := range geoQueryKeys {
		if v := stringifyQueryValue(payload.Query[key]); v != "" {
			geo[key] = v
		}
	}
	if geo["gtsy"] == "" {
		return nil, fmt.Errorf("не удалось разобрать гео из ссылки Kufar")
	}
	return geo, nil
}

func stringifyQueryValue(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(t)
	case float64:
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return strings.TrimSpace(fmt.Sprintf("%v", t))
	case json.Number:
		return t.String()
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return ""
		}
		return strings.Trim(string(b), `"`)
	}
}

func mergeMissing(dst, src map[string]string) {
	for k, v := range src {
		if dst[k] == "" && v != "" {
			dst[k] = v
		}
	}
}

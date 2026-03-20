package quarksearch

import (
	"html"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	quarkURLPattern = regexp.MustCompile(`https?://pan\.(?:quark|qoark)\.cn/[^\s<>"']+`)
	passwordPattern = regexp.MustCompile(`(?i)(?:提取码|密码|pwd|pass(?:word)?)\s*[：:\-]?\s*([a-z0-9]{4,8})`)
	tagPattern      = regexp.MustCompile(`<[^>]+>`)
	spacePattern    = regexp.MustCompile(`\s+`)
)

func normalizeList(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimSpace(strings.ToLower(value))
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func resolveChannels(requested, configured []string) []string {
	requested = normalizeList(requested)
	if len(requested) > 0 {
		return requested
	}
	configured = normalizeList(configured)
	if len(configured) > 0 {
		return configured
	}
	return []string{"tgsearchers3"}
}

func resolvePluginNames(requested, configured []string, supported map[string]Plugin) []string {
	configured = normalizeList(configured)
	if len(configured) == 0 {
		return nil
	}

	allowed := make(map[string]struct{}, len(configured))
	for _, name := range configured {
		allowed[name] = struct{}{}
	}

	requested = normalizeList(requested)
	if len(requested) == 0 {
		return filterSupportedPluginNames(configured, supported)
	}

	result := make([]string, 0, len(requested))
	seen := make(map[string]struct{}, len(requested))
	for _, name := range requested {
		if _, exists := allowed[name]; !exists {
			continue
		}
		if _, exists := supported[name]; !exists {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, name)
	}
	return result
}

func filterSupportedPluginNames(names []string, supported map[string]Plugin) []string {
	result := make([]string, 0, len(names))
	for _, name := range names {
		if _, exists := supported[name]; exists {
			result = append(result, name)
		}
	}
	return result
}

func cleanHTMLText(raw string) string {
	text := tagPattern.ReplaceAllString(raw, " ")
	text = html.UnescapeString(text)
	return cleanText(text)
}

func cleanText(raw string) string {
	raw = strings.ReplaceAll(raw, "\u00a0", " ")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	lines := strings.Split(raw, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(spacePattern.ReplaceAllString(line, " "))
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}

func extractTitle(raw string) string {
	text := cleanText(raw)
	if text == "" {
		return ""
	}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return truncateText(line, 120)
	}
	return truncateText(text, 120)
}

func truncateText(text string, limit int) string {
	if limit <= 0 || len(text) <= limit {
		return text
	}
	return strings.TrimSpace(text[:limit]) + "..."
}

func extractPassword(text string) string {
	matches := passwordPattern.FindStringSubmatch(strings.ToLower(text))
	if len(matches) < 2 {
		return ""
	}
	return strings.TrimSpace(matches[1])
}

func extractPasswordFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	for _, key := range []string{"pwd", "password", "passcode"} {
		value := strings.TrimSpace(u.Query().Get(key))
		if value != "" {
			return value
		}
	}
	return ""
}

func normalizeQuarkURL(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	value = strings.TrimRight(value, ".,;)]}>'\"")
	value = strings.Replace(value, "pan.qoark.cn", "pan.quark.cn", 1)
	decoded, err := url.QueryUnescape(value)
	if err == nil && strings.Contains(decoded, "pan.quark.cn") {
		value = decoded
	}
	return value
}

func extractQuarkLinks(text string, hrefs []string) []Link {
	result := make([]Link, 0, 4)
	seen := make(map[string]struct{}, 4)
	password := extractPassword(text)

	appendURL := func(raw string) {
		raw = normalizeQuarkURL(raw)
		if raw == "" {
			return
		}
		if !quarkURLPattern.MatchString(raw) {
			return
		}
		if _, exists := seen[raw]; exists {
			return
		}
		seen[raw] = struct{}{}
		linkPassword := extractPasswordFromURL(raw)
		if linkPassword == "" {
			linkPassword = password
		}
		result = append(result, Link{
			Type:     "quark",
			URL:      raw,
			Password: linkPassword,
		})
	}

	for _, match := range quarkURLPattern.FindAllString(text, -1) {
		appendURL(match)
	}
	for _, href := range hrefs {
		appendURL(href)
	}
	return result
}

func containsAllKeywords(text, keyword string) bool {
	if strings.TrimSpace(keyword) == "" {
		return true
	}
	text = strings.ToLower(strings.TrimSpace(text))
	for _, part := range strings.Fields(strings.ToLower(keyword)) {
		if !strings.Contains(text, part) {
			return false
		}
	}
	return true
}

func applyFilterConfig(results []SearchResult, filter *FilterConfig) []SearchResult {
	if filter == nil {
		return results
	}
	includes := normalizeList(filter.Include)
	excludes := normalizeList(filter.Exclude)
	filtered := make([]SearchResult, 0, len(results))

	for _, result := range results {
		haystack := strings.ToLower(strings.TrimSpace(result.Title + "\n" + result.Content))
		if haystack == "" {
			haystack = strings.ToLower(strings.TrimSpace(result.Title))
		}

		excluded := false
		for _, keyword := range excludes {
			if strings.Contains(haystack, keyword) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		if len(includes) > 0 {
			matched := false
			for _, keyword := range includes {
				if strings.Contains(haystack, keyword) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		filtered = append(filtered, result)
	}

	return filtered
}

func filterResultsByKeyword(results []SearchResult, keyword string) []SearchResult {
	filtered := make([]SearchResult, 0, len(results))
	for _, result := range results {
		if containsAllKeywords(result.Title, keyword) || containsAllKeywords(result.Content, keyword) {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func sortResults(results []SearchResult, keyword string) {
	sort.SliceStable(results, func(i, j int) bool {
		leftScore := resultScore(results[i], keyword)
		rightScore := resultScore(results[j], keyword)
		if leftScore != rightScore {
			return leftScore > rightScore
		}
		if !results[i].Datetime.Equal(results[j].Datetime) {
			return results[i].Datetime.After(results[j].Datetime)
		}
		return results[i].Title < results[j].Title
	})
}

func resultScore(result SearchResult, keyword string) int {
	score := 0
	if containsAllKeywords(result.Title, keyword) {
		score += 2
	}
	if containsAllKeywords(result.Content, keyword) {
		score++
	}
	if !result.Datetime.IsZero() {
		score++
	}
	return score
}

func mergeResultsByType(results []SearchResult) map[string][]MergedLink {
	merged := make([]MergedLink, 0, len(results))
	seen := make(map[string]struct{}, len(results))

	for _, result := range results {
		note := strings.TrimSpace(result.Title)
		if note == "" {
			note = extractTitle(result.Content)
		}

		for _, link := range result.Links {
			if strings.TrimSpace(link.URL) == "" || strings.ToLower(strings.TrimSpace(link.Type)) != "quark" {
				continue
			}
			key := normalizeQuarkURL(link.URL)
			if key == "" {
				continue
			}
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			merged = append(merged, MergedLink{
				URL:      link.URL,
				Password: link.Password,
				Note:     note,
				Datetime: result.Datetime,
				Source:   resultSource(result),
				Images:   append([]string(nil), result.Images...),
			})
		}
	}

	sort.SliceStable(merged, func(i, j int) bool {
		if !merged[i].Datetime.Equal(merged[j].Datetime) {
			return merged[i].Datetime.After(merged[j].Datetime)
		}
		return merged[i].Note < merged[j].Note
	})

	if len(merged) == 0 {
		return map[string][]MergedLink{}
	}
	return map[string][]MergedLink{"quark": merged}
}

func resultSource(result SearchResult) string {
	if strings.TrimSpace(result.Channel) != "" {
		return "tg:" + strings.TrimSpace(result.Channel)
	}
	if idx := strings.Index(result.UniqueID, "-"); idx > 0 {
		return "plugin:" + result.UniqueID[:idx]
	}
	return ""
}

func filterResponseByType(response SearchResponse, resultType string) SearchResponse {
	switch strings.ToLower(strings.TrimSpace(resultType)) {
	case "all":
		return response
	case "results":
		return SearchResponse{
			Total:   len(response.Results),
			Results: response.Results,
		}
	default:
		total := 0
		if links, exists := response.MergedByType["quark"]; exists {
			total = len(links)
		}
		return SearchResponse{
			Total:        total,
			MergedByType: response.MergedByType,
		}
	}
}

func parseTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

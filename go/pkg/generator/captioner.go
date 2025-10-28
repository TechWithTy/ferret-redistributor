package generator

import (
    "regexp"
    "strings"
)

var nonWord = regexp.MustCompile(`[^a-z0-9]+`)

// CaptionFor builds a platform-specific caption from variant metadata and topic.
// Returns caption and a space-separated hashtag string.
func CaptionFor(platform, topic string, v Variant) (string, string) {
    tags := MakeTags(topic)
    title := strings.TrimSpace(v.Title)
    cta := strings.TrimSpace(v.CTA)
    switch strings.ToLower(platform) {
    case "linkedin":
        caption := joinNonEmpty([]string{title, cta}, "\n\n")
        return caption, tags
    case "twitter":
        base := joinNonEmpty([]string{title, tags}, " \n")
        if len(base) > 260 { base = base[:260] }
        return base, ""
    case "instagram":
        caption := joinNonEmpty([]string{title, tags, cta}, "\n\n")
        return caption, tags
    default:
        caption := joinNonEmpty([]string{title, cta}, " \n")
        return caption, tags
    }
}

// MakeTags builds up to 3 hashtags from topic words.
func MakeTags(topic string) string {
    t := strings.ToLower(topic)
    t = nonWord.ReplaceAllString(t, " ")
    fields := strings.Fields(t)
    out := make([]string, 0, 3)
    for _, w := range fields {
        if w == "" { continue }
        out = append(out, "#"+strings.TrimSpace(w))
        if len(out) >= 3 { break }
    }
    return strings.Join(out, " ")
}

func joinNonEmpty(parts []string, sep string) string {
    out := make([]string, 0, len(parts))
    for _, p := range parts {
        if strings.TrimSpace(p) != "" { out = append(out, p) }
    }
    return strings.Join(out, sep)
}


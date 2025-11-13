package helpers

import (
	"html"
	"net/url"
	"regexp"
	"strings"
)

var xssPattern = []*regexp.Regexp {
	regexp.MustCompile(`(?i)<\s*script.*?>.*?<\s*/\s*script\s*>`),
	regexp.MustCompile(`(?i)on\w+\s*=`),
	regexp.MustCompile(`(?i)javascript:`),
}

var sqlInjectionPatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)\bOR\b\s+\d+\s*=\s*\d+`),          
    regexp.MustCompile(`(?i)\bAND\b\s+\d+\s*=\s*\d+`),         
    regexp.MustCompile(`(?i)UNION\s+SELECT`),                  
    regexp.MustCompile(`(?i)--`),                              
    regexp.MustCompile(`(?i);\s*DROP\s+TABLE`),
}

func CheckPatterns(input string) bool {
	normInput := normalizeInput(input)

	switch {
	case detectXSS(normInput) == true:
		return true
	case detectSQLInjection(normInput) == true:
		return true
	default: 
		return false
	}
}

func detectXSS(input string) bool{
	for _,re := range xssPattern{
		if re.MatchString(input){
			return true
		}
	}
	return false
}

func detectSQLInjection(input string) bool{
	for _,re := range sqlInjectionPatterns {
		if re.MatchString(input){
			return true
		}
	}
	return false
}

func normalizeInput(input string) string{
	if u, err := url.QueryUnescape(input); err == nil {
		input = u 
	}
	input = html.UnescapeString(input)
	input = strings.ReplaceAll(input, "\x00", "")
	input = strings.TrimSpace(input)
	return input
}

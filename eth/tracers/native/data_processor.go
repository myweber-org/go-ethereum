
package data_processor

import (
	"regexp"
	"strings"
	"unicode"
)

type Processor struct {
	stripSpaces   bool
	removeSpecial bool
	toLowercase   bool
}

func NewProcessor(stripSpaces, removeSpecial, toLowercase bool) *Processor {
	return &Processor{
		stripSpaces:   stripSpaces,
		removeSpecial: removeSpecial,
		toLowercase:   toLowercase,
	}
}

func (p *Processor) CleanString(input string) string {
	result := input

	if p.stripSpaces {
		result = strings.TrimSpace(result)
		result = strings.Join(strings.Fields(result), " ")
	}

	if p.removeSpecial {
		reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
		result = reg.ReplaceAllString(result, "")
	}

	if p.toLowercase {
		result = strings.ToLower(result)
	}

	return result
}

func (p *Processor) NormalizeWhitespace(input string) string {
	var builder strings.Builder
	prevSpace := false

	for _, r := range input {
		if unicode.IsSpace(r) {
			if !prevSpace {
				builder.WriteRune(' ')
				prevSpace = true
			}
		} else {
			builder.WriteRune(r)
			prevSpace = false
		}
	}

	return strings.TrimSpace(builder.String())
}

func (p *Processor) ProcessBatch(inputs []string) []string {
	results := make([]string, len(inputs))
	for i, input := range inputs {
		results[i] = p.CleanString(input)
	}
	return results
}
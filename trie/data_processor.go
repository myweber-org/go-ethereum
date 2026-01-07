
package data_processor

import (
	"regexp"
	"strings"
)

type Processor struct {
	allowedPattern *regexp.Regexp
}

func NewProcessor(allowedPattern string) (*Processor, error) {
	compiled, err := regexp.Compile(allowedPattern)
	if err != nil {
		return nil, err
	}
	return &Processor{allowedPattern: compiled}, nil
}

func (p *Processor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	if p.allowedPattern == nil {
		return trimmed
	}
	return p.allowedPattern.FindString(trimmed)
}

func (p *Processor) ValidateInput(input string) bool {
	if input == "" {
		return false
	}
	if p.allowedPattern != nil && !p.allowedPattern.MatchString(input) {
		return false
	}
	return true
}

func (p *Processor) Process(input string) (string, bool) {
	cleaned := p.CleanInput(input)
	valid := p.ValidateInput(cleaned)
	return cleaned, valid
}

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
	if p.allowedPattern != nil {
		return p.allowedPattern.FindString(trimmed)
	}
	return trimmed
}

func (p *Processor) Validate(input string) bool {
	if input == "" {
		return false
	}
	if p.allowedPattern != nil {
		return p.allowedPattern.MatchString(input)
	}
	return true
}

func (p *Processor) Process(input string) (string, bool) {
	cleaned := p.CleanInput(input)
	valid := p.Validate(cleaned)
	return cleaned, valid
}
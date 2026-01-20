
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TransformToUpper(input string) string {
	return strings.ToUpper(input)
}

func PrettyPrintJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func main() {
	email := "test@example.com"
	fmt.Printf("Email %s valid: %v\n", email, ValidateEmail(email))

	str := "hello world"
	fmt.Printf("Original: %s, Transformed: %s\n", str, TransformToUpper(str))

	sample := map[string]interface{}{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
	}
	pretty, _ := PrettyPrintJSON(sample)
	fmt.Println("Pretty JSON:")
	fmt.Println(pretty)
}package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return normalized
}

func (dp *DataProcessor) NormalizeCase(input string, toUpper bool) string {
	cleaned := dp.CleanString(input)
	if toUpper {
		return strings.ToUpper(cleaned)
	}
	return strings.ToLower(cleaned)
}

func (dp *DataProcessor) ExtractAlphanumeric(input string) string {
	alnumRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleaned := dp.CleanString(input)
	return alnumRegex.ReplaceAllString(cleaned, "")
}

func main() {
	processor := NewDataProcessor()
	
	sample := "  Hello   World!  This  is  a  test.  "
	
	cleaned := processor.CleanString(sample)
	println("Cleaned:", cleaned)
	
	upper := processor.NormalizeCase(sample, true)
	println("Uppercase:", upper)
	
	alnum := processor.ExtractAlphanumeric(sample)
	println("Alphanumeric only:", alnum)
}
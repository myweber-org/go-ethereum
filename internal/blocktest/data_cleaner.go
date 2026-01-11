package data

func DeduplicateStrings(slice []string) []string {
    seen := make(map[string]struct{})
    result := make([]string, 0, len(slice))
    
    for _, item := range slice {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
            result = append(result, item)
        }
    }
    
    return result
}

package data

func FilterAndTransform(numbers []int, predicate func(int) bool, transform func(int) int) []int {
    var result []int
    for _, num := range numbers {
        if predicate(num) {
            result = append(result, transform(num))
        }
    }
    return result
}
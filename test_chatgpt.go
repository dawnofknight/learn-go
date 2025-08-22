package main

import "fmt"

// This function has some issues and could be improved
func calculateSum(numbers []int) int {
    var result int
    for i := 0; i < len(numbers); i++ {
        result = result + numbers[i]
    }
    return result
}

// This function is inefficient
func findMax(arr []int) int {
    max := arr[0]
    for i := 0; i < len(arr); i++ {
        if arr[i] > max {
            max = arr[i]
        }
    }
    return max
}

func main() {
    nums := []int{1, 2, 3, 4, 5}
    sum := calculateSum(nums)
    max := findMax(nums)
    
    fmt.Printf("Sum: %d, Max: %d\n", sum, max)
}

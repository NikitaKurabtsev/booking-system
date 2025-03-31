// // Для задачи упрощения пути в файловой системе нужно обработать строки таким образом,
// // чтобы они корректно отображали относительные и абсолютные пути,
// // включая обработку специальных компонентов пути: `..` (переход на уровень выше) и `.` (текущая директория).
//
// // originalPath := "/foo/../test/../test/../foo/server/http/"
//
// // Примеры использования
// // /foo/../test/../test/../foo/server/http/ -->должно выводить /foo/server/http
// package main

// func simplifyPath(path string) string {
// 	var stack []string
// 	pathSequence := strings.Split(path, "/")

// 	for _, p := range pathSequence {
// 		if p == ".." {
// 			if len(stack) > 0 {
// 				stack = stack[:len(stack)-1]
// 			}
// 		} else if p != "" && p != "." {
// 			stack = append(stack, p)
// 		}
// 	}

//		result := "/" + strings.Join(stack, "/")
//		return result
//	}
// func simplifyPath(path string) string {
// 	var stack []string
// 	pathSequence := strings.Split(path, "/")

// 	for _, p := range pathSequence {
// 		if p == ".." {
// 			if len(stack) > 0 {
// 				stack = stack[:len(stack)-1]
// 			}
// 		} else if p != "" && p != "." {
// 			stack = append(stack, p)
// 		}
// 	}

// 	result := "/" + strings.Join(stack, "/")

// 	return result
// }

// "/foo/../test/../test/../foo/server/http/../http/"

// func simplifyPath(path string) string {
// 	var stack []string
// 	pathSequence := strings.Split(path, "/")

// 	for _, p := range pathSequence {
// 		switch p {
// 		case "..":
// 			if len(stack) > 0 {
// 				stack = stack[:len(stack)-1]
// 			}
// 		case ".", "":
// 			continue
// 		default:
// 			stack = append(stack, p)
// 		}
// 	}

// 	result := "/" + strings.Join(stack, "/")

//		return result
//	}
// func simplifyPath(path string) string {
// 	var stack []string
// 	pathSequence := strings.Split(path, "/")

// 	for _, p := range pathSequence {
// 		switch p {
// 		case "", ".":
// 			continue
// 		case "..":
// 			if len(stack) > 0 {
// 				stack = stack[:len(stack)-1]
// 			}
// 		default:
// 			stack = append(stack, p)
// 		}

// 	}

// 	result := "/" + strings.Join(stack, "/")

// 	return result
// }

// func simplyPath(path string) string {
// 	var stack []string
// 	pathSequence := strings.Split(path, "/")

// 	for _, p := range pathSequence {
// 		if p != "" && p != "." {
// 			stack = append(stack, p)
// 		} else if p == ".." {
// 			if len(stack) > 0 {
// 				stack = stack[:len(stack)-1]
// 			}
// 		}
// 	}

// 	result := "/" + strings.Join(stack, "/")

//		return result
//	}
package main

import (
	"fmt"
	"strings"
)

func simplifyPath(path string) string {
	var stack []string
	pathSequence := strings.Split(path, "/")

	for _, p := range pathSequence {
		if p == ".." {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		} else if p != "" && p != "." {
			stack = append(stack, p)
		}
	}

	result := "/" + strings.Join(stack, "/")

	return result
}

func main() {
	originalPath := "/foo/../test/../test/../foo/server/http/"

	og := "/foo/../test/../test/../foo/server/http/../http/"

	finalPath := simplifyPath(originalPath)
	fl := simplifyPath(og)

	fmt.Println(originalPath) // Output: /foo/../test/../test/../foo/server/http/
	fmt.Println(finalPath)    // Output: /foo/server/http

	fmt.Println(og) // Output: /foo/../test/../test/../foo/server/http/../http/
	fmt.Println(fl) // Output: /foo/server/http
}

// Дана последовательность в виде строки, нужно произвести сжатие символов
// in: "AAAABBBCCDDDAAA"
// out: "A4B3C2D3A3"

// package main
//
// import (
//
//	"fmt"
//	"strconv"
//	"strings"
//
// )
//

//	func compress(inputStr string) string {
//		var result strings.Builder
//		counter := 1
//		current := inputStr[0]
//
//		for i := 1; i < len(inputStr); i++ {
//			if inputStr[i] == current {
//				counter++
//			} else {
//				result.WriteRune(rune(current))
//				result.WriteString(strconv.Itoa(counter))
//
//				counter = 1
//				current = inputStr[i]
//			}
//		}
//
//		result.WriteRune(rune(current))
//		result.WriteString(strconv.Itoa(counter))
//
//		return result.String()
//	}

// func compress(inputStr string) string {
// 	var result strings.Builder
// 	current := inputStr[0]
// 	counter := 1

// 	for i := 1; i < len(inputStr); i++ {
// 		if inputStr[i] == current {
// 			counter++
// 		} else {
// 			result.WriteRune(rune(current))
// 			result.WriteString(strconv.Itoa(counter))

// 			current = inputStr[i]
// 			counter = 1
// 		}
// 	}

// 	result.WriteRune(rune(current))
// 	result.WriteString(strconv.Itoa(counter))

// 	return result.String()
// }

// func compress(inputString string) string {
// 	var result strings.Builder
// 	current := inputString[0]
// 	counter := 1

// 	for i := 1; i < len(inputString); i++ {
// 		if inputString[i] == current {
// 			counter++
// 		} else {
// 			result.WriteString(string(current))
// 			result.WriteString(string(counter))

// 			current = inputString[i]
// 			counter = 1
// 		}
// 	}

// 	result.WriteString(string(current))
// 	result.WriteString(string(counter))

//		return result.String()
//	}

//	func convertToString(s any) string {
//		switch v := s.(type) {
//		case byte:
//			return string(v)
//		case int:
//			return strconv.Itoa(v)
//		default:
//			return "unknown type"
//		}
//	}
//
//	func compress(inputStr string) string {
//		var outputStr strings.Builder
//		current := inputStr[0]
//		counter := 1
//
//		for i := 1; i < len(inputStr); i++ {
//			if inputStr[i] == current {
//				counter++
//			} else {
//				outputStr.WriteString(string(current))
//				outputStr.WriteString(strconv.Itoa(counter))
//
//				current = inputStr[i]
//				counter = 1
//			}
//		}
//		outputStr.WriteString(string(current))
//		outputStr.WriteString(strconv.Itoa(counter))
//
//		return outputStr.String()
//	}
// package main

// import (
// 	"fmt"
// 	"strconv"
// 	"strings"
// )

// func compress(inputStr string) string {
// 	var outputStr strings.Builder
// 	current := inputStr[0]
// 	counter := 1

// 	for i := 1; i < len(inputStr); i++ {
// 		if inputStr[i] == current {
// 			counter++
// 		} else {
// 			outputStr.WriteString(string(current))
// 			outputStr.WriteString(strconv.Itoa(counter))

// 			current = inputStr[i]
// 			counter = 1
// 		}
// 	}

// 	outputStr.WriteString(string(current))
// 	outputStr.WriteString(strconv.Itoa(counter))

// 	return outputStr.String()
// }

// func main() {
// 	inputStr := "AAAABBBCCDDDAAA"
// 	compressedStr := compress(inputStr)
// 	fmt.Println(compressedStr) // Вывод: "A4B3C2D3A3"
// }

// Дан массив целых чисел `nums`, вернуть значение `true`,
// если какое-либо значение встречается в массиве хотя бы дважды,
// и вернуть значение `false`, если каждый элемент различен.
//
//
//Example 1:
//
//Input: nums = [1,2,3,1]
//Output: true
//
//Example 2:
//
//Input: nums = [1,2,3,4]
//Output: false
//
//Example 3:
//
//Input: nums = [1,1,1,3,3,4,3,2,4,2]
//Output: true

// package main

// import "fmt"

// func isNotUnique(nums []int) bool {
// 	set := make(map[int]struct{})

// 	for _, num := range nums {
// 		if _, ok := set[num]; ok {
// 			return true
// 		}
// 		set[num] = struct{}{}
// 	}

// 	return false
// }

// func doItRecursive(n int) {
// 	if n > 0 {
// 		fmt.Println("...")
// 		doItRecursive(n - 1)
// 	} else {
// 		fmt.Println("OK 200")
// 	}

// }

// func main() {
// 	exmpl1 := []int{1, 2, 3, 1}
// 	exmpl2 := []int{1, 2, 3, 4}
// 	exmpl3 := []int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}

// 	fmt.Println(isNotUnique(exmpl1), "must be TRUE")
// 	fmt.Println(isNotUnique(exmpl2), "must be FALSE")
// 	fmt.Println(isNotUnique(exmpl3), "must be TRUE")

// 	doItRecursive(10)

// }

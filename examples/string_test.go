package examples

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// 1. Always prefer strconv.FormatInt(num) over fmt.Sprintf("%d", num)
// 2. strings.Builder can be used with fmt.Fprintf for appending formatted strings
// 3. builtin.go defines functions like append, make, delete, copy, new, close, panic, print, println, len, cap and various types.
// 	  These are called predeclared identifiers. Types: bool, byte, ..., uintptr, Constants: true/false/iota, Zero Value: nil, Functions: above ones
// 4. byte is uint8 and rune is int32 only.
// 5. error is an interface which defines only one function Error()
// 6. The size of an array is part of its type, whereas slices can have a dynamic size because they are wrappers around arrays.
//	  e.g. var a [10]int and is not equivalent to var a [11]int
// 7. Capacity tells you how large your subset can grow before it will no longer fit in the array that is backing the slice.
//    combining slices with the append function gives us a type that is very similar to arrays, but is capable of growing over time to handle more elements.

var result string

// BenchmarkStringConcatenation/strings.Builder-12         	 3433200	       365.4 ns/op	     470 B/op	       0 allocs/op
// BenchmarkStringConcatenation/fmt.Sprintf-12             	 2138028	       563.7 ns/op	     176 B/op	       6 allocs/op
// BenchmarkStringConcatenation/strings.Join-12            	11226514	       109.7 ns/op	      80 B/op	       1 allocs/op
// BenchmarkStringConcatenation/with-plus-12               	   10000	      342055 ns/op	 2144945 B/op	       5 allocs/op
func BenchmarkStringConcatenation(b *testing.B) {

	var input []string = []string{"Mayuresh", "Vivek", "Ajit", "Nilesh", "Sandeep"}

	b.Run("strings.Builder", func(b *testing.B) {
		b.ReportAllocs()
		var builder strings.Builder
		for i := 0; i < b.N; i++ {
			for _, in := range input {
				builder.WriteString(in)
				builder.WriteString("_delimiter_")
			}
		}
		result = builder.String()
		// fmt.Println("String Builder: " + builder.String())
	})

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ReportAllocs()
		var output string
		for i := 0; i < b.N; i++ {
			output = fmt.Sprintf("%s_delimiter_%s_delimiter_%s_delimiter_%s_delimiter_%s_delimiter_", input[0], input[1], input[2], input[3], input[4])
		}
		result = output
		// fmt.Println("fmt.Sprintf: " + output)
	})

	b.Run("strings.Join", func(b *testing.B) {
		b.ReportAllocs()
		var output string
		for i := 0; i < b.N; i++ {
			output = strings.Join(input, "_delimiter_")
		}
		result = output
		// fmt.Println("strings.Join: " + output)
	})

	b.Run("with-plus", func(b *testing.B) {
		b.ReportAllocs()
		var output string
		for i := 0; i < b.N; i++ {
			for _, in := range input {
				output += in + "_delimiter_"
			}
		}
		result = output
		// fmt.Println("with-plus: " + output)
	})
}

// BenchmarkIntegerString/strconv.Itoa-12         	25905498	        45.75 ns/op	       7 B/op	       0 allocs/op
// BenchmarkIntegerString/fmt.Sprintf-12          	 9348594	       135.9 ns/op	      16 B/op	       1 allocs/op
func BenchmarkIntegerString(b *testing.B) {

	b.Run("strconv.Itoa", func(b *testing.B) {
		b.ReportAllocs()
		var output string
		for i := 0; i < b.N; i++ {
			output = strconv.Itoa(i)
		}
		result = output
	})

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ReportAllocs()
		var output string
		for i := 0; i < b.N; i++ {
			output = fmt.Sprintf("%d", i)
		}
		result = output
	})
}

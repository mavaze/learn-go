package examples

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func Test23(t *testing.T) {
	a := make([]int, 5)
	printSlice("a", a) // a len=5 cap=5 [0 0 0 0 0]

	b := make([]int, 0, 5)
	printSlice("b", b) // b len=0 cap=5 []

	c := b[:2]
	printSlice("c", c) // c len=2 cap=5 [0 0]

	d := c[2:5]
	printSlice("d", d) // d len=3 cap=3 [0 0 0]

	var e [5]int // e len=5 cap=5 [0 0 0 0 0]
	fmt.Printf("e len=%d cap=%d %v\n", len(e), cap(e), e)

	var f [5]int = [5]int{1, 1, 1} // f len=5 cap=5 [1 1 1 0 0]
	fmt.Printf("f len=%d cap=%d %v\n", len(f), cap(f), 5)
}

func printSlice(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n", s, len(x), cap(x), x)
}

// func TestArraySlice(t *testing.T) {
// 	arr := [...]int{1, 2, 3, 4, 5}
// 	var sl [3]int = arr[2:4]
// 	t.Log(arr)
// 	t.Log(sl)
// }

func BenchmarkSliceAppend(b *testing.B) {
	a := make([]int, 0, b.N)
	for i := 0; i < b.N; i++ {
		a = append(a, i)
	}
}

func BenchmarkSliceSet(b *testing.B) {
	a := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		a[i] = i
	}
}

var p1 []int
var p2 []int

func BenchmarkSliceAppend1(b *testing.B) {
	b.Run("append", func(b *testing.B) {
		b.ReportAllocs()
		a := make([]int, 0, b.N)
		for i := 0; i < b.N; i++ {
			a = append(a, i)
		}
		p1 = a
	})

	b.Run("index", func(b *testing.B) {
		b.ReportAllocs()
		a := make([]int, b.N)
		for i := 0; i < b.N; i++ {
			a[i] = i
		}
		p2 = a
	})
}

func BenchmarkSliceConversion(b *testing.B) {
	numbers := make([]int, 100)
	for i := range numbers {
		numbers[i] = i
	}
	b.Run("no-capacity", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			numbersToStringsBad(numbers)
		}
	})
	b.Run("with-capacity", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			numbersToStringsBetter(numbers)
		}
	})
}

func numbersToStringsBad(numbers []int) []string {
	vals := []string{}
	for _, n := range numbers {
		vals = append(vals, strconv.Itoa(n))
	}
	return vals
}

func numbersToStringsBetter(numbers []int) []string {
	vals := make([]string, 0, len(numbers))
	for _, n := range numbers {
		vals = append(vals, strconv.Itoa(n))
	}
	return vals
}

type employee struct {
	name        string
	age         *int
	salary      int
	departments []string // structs containing slices cannot be compared using == operator, however arrays can
}

func TestComparison(t *testing.T) {
	age1 := new(int)
	age2 := new(int)
	*age1 = 30
	*age2 = 31
	emp1 := employee{name: "Sam", age: age1, salary: 2000, departments: []string{"CS"}}
	emp2 := employee{name: "Sam", age: age2, salary: 2000, departments: []string{"CS"}}
	if reflect.DeepEqual(emp1, emp2) {
		// if emp1 == emp2 { // compilation error if employee contains any slice, but arrays can be just fine
		fmt.Println("emp1 and emp2 are equal")
	} else {
		fmt.Println("emp1 and emp2 are not equal")
	}
}

func TestCopySlices(t *testing.T) {
	a1 := []int{1, 2, 3, 4}
	a2 := []int{8, 9, 10, 11, 12}
	a3 := []int{30, 31, 32}

	n := copy(a1, a2) // copies elements of same size of dest a1, i.e. a1's length doesn't grow to 5 to match with that of a2
	fmt.Printf("a1: %v, a2: %v, a3: %v, n: %d", a1, a2, a3, n) // a1: [8 9 10 11]
	fmt.Println()

	n = copy(a2[1:], a3)
	fmt.Printf("a1: %v, a2: %v, a3: %v, n: %d", a1, a2, a3, n) // a2: [8 30 31 32 12], 
	fmt.Println()
}

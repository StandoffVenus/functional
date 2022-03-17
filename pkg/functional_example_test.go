package functional_test

import (
	"fmt"
	"strconv"
	"strings"

	functional "github.com/standoffvenus/functional/pkg"
)

func ExampleReduce() {
	summuation := func(ints ...int) int {
		return functional.Reduce(ints, func(accum int, cur int) int {
			return accum + cur
		})
	}

	total := summuation(0, 1, 2, 3, 4, 5)

	fmt.Println("Total:", total)
	// Output: Total: 15
}

func ExampleMap() {
	type Person struct {
		FirstName string
		LastName  string
	}

	getFullNames := func(people ...Person) []string {
		return functional.Map(people, func(p Person) string {
			return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
		})
	}

	people := []Person{
		{
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}

	names := getFullNames(people...)

	fmt.Println("People:", strings.Join(names, ", "))
	// Output: People: John Doe, Jane Doe
}

func ExampleFilter() {
	type Person struct {
		Age int
	}

	adultFilter := func(p Person) bool {
		return p.Age > 18
	}

	people := []Person{
		{
			Age: 14,
		},
		{
			Age: 61,
		},
		{
			Age: 22,
		},
	}

	adults := functional.Filter(people, adultFilter)

	fmt.Println("Adults:", adults)
	// Output: Adults: [{61} {22}]
}

func ExampleCompose() {
	stringToInt := func(s string) int {
		i, _ := strconv.ParseInt(s, 10, 64)
		return int(i)
	}
	intToPtr := func(i int) *int {
		return &i
	}

	// Parses string to int, then returns pointer to int
	composer := functional.Compose(intToPtr, stringToInt)

	fmt.Println("Integer:", *composer("42"))
	// Output: Integer: 42
}

func ExampleChain() {
	type Number = int
	addOne := func(i Number) Number { return i + 1 }
	square := func(i Number) Number { return i * i }
	halve := func(i Number) Number { return i / 2 }

	// Recall that the chained function will invoke
	// the provided functions in reverse order:
	// addOne() -> square() -> halve()
	chain := functional.Chain(halve, square, addOne)

	// addOne(5) -> 6
	// square(6) -> 36
	// halve(36) -> 18
	fmt.Println("Result:", chain(5))
	// Output: Result: 18
}

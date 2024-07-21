package main

import (
	"fmt"
	"log"
	"math"
)

// 2.1 Variable shadowing
// usually happens with var and when you have a function that returns an err
// solve using temporary variables
// solve using "var err error"
// can be solved with linters?
func shadow(cond bool) {
	var client string
	condFunc := func() (string, error) {
		return "cond", nil
	}
	nonCondFunc := func() (string, error) {
		return "nonCond", nil
	}
	if cond {
		client, _ := condFunc()
		log.Println(client)
	} else {
		client, _ := nonCondFunc()
		log.Println(client)
	}
	log.Println(client)
}

func solveShadow(cond bool) {
	var client string
	var err error
	condFunc := func() (string, error) {
		return "cond", nil
	}
	nonCondFunc := func() (string, error) {
		return "nonCond", nil
	}
	if cond {
		client, err = condFunc()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		log.Println(client)
	} else {
		client, err = nonCondFunc()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		log.Println(client)
	}
	log.Println(client)
}

// 2.2 Unnecessary nested code
// allign happy path to the left and error handling to the right
// use early returns

// 2.3 Misusing init functions
// init functions are executed in fancy ways
// using (import _ "package") will only execute the init function of the package (if you don't use any other functions)
// they have bad error handling, so you can only panic (maybe a package wants a fallback or retry)
// they will run when testing
// (bad) example with database: alters a global var (better use a function)
// (good) example with http: setting static routes (handles) or just setting static configuration
func init() {
	log.Println("Initializing...")
}

// 2.4 Overusing getters and setters
// allow new functionality to be added later (field validation, logging, etc.) or mutex
// hide internal representation (encapsulation)
// don't abuse them, they are not always needed
// use them like Balance() and SetBalance() (not GetBalance() and SetBalance())

// 2.5 Interface pollution
// they specify behavior of an object
// (good) example with io.Reader, io.Writer
// creating functions that take interfaces as arguments eases unit testing
// "The bigger the interface, the weaker the abstraction" - Rob Pike
// we shouldn’t design with interfaces but wait for a concrete need
// "Don't design with interfaces, discover them" - Rob Pike
//
//	When to use interfaces:
//		1. common behavior (sort.Interface as an example)
//		2. decoupling (if you rely on abstraction, you can change without changing the code)
// 			2.1. liskov substitution principle (LSP)
// 			2.2. helps with testing (example with customer store)
//		3. restricting behavior
//			3.1. example with configGetter (creating a new interface for a subset of methods)
//

// 2.6 Interface on the producer side (defined in the same package as the concrete implementation)
// The "implements" is implicit, so the consumer can define its own abstraction
// It's not up to the producer to force an abstraction, but the consumer to define it
// The client can define the most accurate interface for its needs (Interface Segregation Principle from SOLID)
// If we build producer-side interfaces, make them as minimal as possible

// 2.7 Returning interfaces
// "Be conservative in what you do, be liberal in what you accept from others." - TCP
// Returning structs instead of interfaces
// Accepting interfaces if possible

// Producer-side code
// package springer or package scraper
type SpringerScraper struct {
	// ...
}

func (s *SpringerScraper) WithISBN(isbn string) (string, error) {
	return "ISBN", nil
}

func (s *SpringerScraper) WithURL(url string) (string, error) {
	return "URL", nil
}

func (s *SpringerScraper) WithTitle(title string) (string, error) {
	return "Title", nil
}

// Consumer-side code that handles only the ISBN and URL
// package aggregator
type Scraper interface {
	WithISBN(isbn string) (string, error)
	WithURL(url string) (string, error)
}

func scrape(s Scraper) {
	byISBN, _ := s.WithISBN("123")
	byURL, _ := s.WithURL("https://example.com")
	log.Println(byISBN, byURL)
}

func runScrape() {
	s := &SpringerScraper{}
	scrape(s)
}

// 2.8 any says nothing
// ...

// 3.2.3 Detecting integer overflows during addition
func add(x, y int) (int, error) {
	if x > math.MaxInt-y {
		panic("addition overflow!")
	}

	return x + y, nil
}

func mul(x, y int) int {
	if x == 0 || y == 0 {
		return 0
	}

	res := x * y
	if x == 1 || y == 1 {
		return res
	}
	if x == math.MinInt || y == math.MinInt {
		panic("mul overflow 1!")
	}
	if res/x != y {
		panic("mul overflow 2!")
	}

	return x * y
}

// 3.3 Not understanding floating points
// IEEE 754 deep dive undestanding
// maybe left to do some manual conversions like -1.00001
// enough to understand the 16 bit floats with:
// 1 sign bit
// 9 mantissa
// 6 exponent

// 3.4 slice len and capacity
func sliceMagic() {
	s1 := make([]int, 3, 5)
	s1[1] = 1
	s1[2] = 2
	log.Println("s1", s1, len(s1), cap(s1))

	s2 := s1[2:4]
	log.Println("s2", s2, len(s2), cap(s2))

	s2 = append(s2, 3, 4)
	log.Println("s2", s2, len(s2), cap(s2))
}

// 3.5 inefficient init of slices
// using 0 and capacity -> usage of append
func initSlice1() {
	s := make([]int, 0, 100)
	for i := 0; i <= 100; i++ {
		s = append(s, i)
	}
}

// using len -> usage of direct assignment
// a bit more efficient
func initSlice2() {
	s := make([]int, 100)
	for i := 0; i < len(s); i++ {
		s[i] = i
	}
}

func customLog(i int, s []int) {
	fmt.Printf("%d: empty=%t\tnil=%t\n", i, len(s) == 0, s == nil)
}

// 3.6 nil vs empty slices
func nilEmptySlices() {
	var s []int
	customLog(1, s)

	s = []int(nil)
	customLog(2, s)

	s = []int{}
	customLog(3, s)

	s = make([]int, 0)
	customLog(4, s)
}

// prefer returning empty slices rather than
// allocated slices, as they are safe to use with append()
func returnNilSlices() []string {
	var s []string

	if false {
		s = append(s, "abc")
	}
	if false {
		s = append(s, "xyz")
	}

	return s
}

// if you don't want unintended behaviour, use "full slice expression"
func sliceGoodPractices() {
	s := make([]int, 10)
	bad := s[:2]
	bad = append(bad, 2)

	good := s[:2:2]
	good = append(good, 3)

	log.Println("bad: ", bad)
	log.Println("good: ", good)
	log.Println("s: ", s)
}

// if taking small slices, prefer using copy() - this will avoid leaks
func sliceLeaks() {
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = i
		log.Printf("s[%d]: %p", i, &s[i])
	}

	// copies from beginning
	good := make([]int, 2)
	copy(good, s)
	log.Printf("s: %v, %p", s, &s)
	log.Printf("good: %v, %p", good, &good)

	// best if you want from inside
	another := make([]int, 2)
	copy(another, s[5:7])
	log.Printf("another: %v, %p", another, &another)

	// more complex, but not necessary
	this := make([]int, 2)
	for i := 5; i < 7; i++ {
		this[i-5] = s[i]
	}
	log.Printf("this: %v, %p", this, &this)
}

type foo struct {
	bar []byte
}

func otherLeaks() []foo {
	foos := make([]foo, 100)
	for i := 0; i <= len(foos); i++ {
		foos[i] = foo{bar: make([]byte, 1000)}
	}

	// either do this if the slice is small enough
	keep := make([]foo, 2)
	copy(keep, foos[:2])

	// or if slice is big enough, set the unwanted data to nil
	for i := 2; i <= len(foos); i++ {
		foos[i].bar = nil
	}

	return keep
}

func main() {
	// solveShadow(true)
	// runScrape()
	// mul(math.MinInt, 2)
	// sliceMagic()
	// nilEmptySlices()
	// sliceGoodPractices()
	// sliceLeaks()

}

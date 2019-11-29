package db

import (
	"fmt"
	"testing"
)

func init() {
	fmt.Println()
	fmt.Println("----++> RESULT <++----")
	fmt.Println()
}

func Test_WRONG_Slices(t *testing.T) {
	var a [5]int
	var b *[5]int
	a = [5]int{1, 2, 3, 4, 5}
	b = &a
	fmt.Println("a, b:", a, b)
	a = [5]int{3, 4, 5, 6, 7}
	fmt.Println("a, b:", a, b)
}

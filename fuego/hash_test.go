package cache

import (
	"fmt"
	"testing"
)

func Test_hashing(t *testing.T) {
	val1 := 1
	val2 := "1"
	val3 := struct {
		Val int
	}{Val: 1}

	h1 := ApplyHash(val1)
	h2 := ApplyHash(val2)
	h3 := ApplyHash(val3)

	if h1 == h2 {
		t.Fail()
	}

	if h1 == h3 {
		t.Fail()
	}

	if h2 == h3 {
		t.Fail()
	}
	fmt.Println(h1, h2, h3, m)
}

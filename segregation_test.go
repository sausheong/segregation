package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestSplit(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	d := split("2:3:4")
	if d[0] != 2 || d[1] != 3 || d[2] != 4 {
		t.Failed()
	}
}

func TestCalc(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	f := calc("2:2")

	d := []int{0, 0}
	for i := 0; i < 1000; i++ {
		p := rand.Float64()
		// t.Log(p)
		d[f(p)]++
	}
	t.Log(d)
}

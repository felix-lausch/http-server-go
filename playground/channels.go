package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	var walk func(t *tree.Tree)

	walk = func(t *tree.Tree) {
		if t == nil {
			return
		}

		walk(t.Left)
		ch <- t.Value
		walk(t.Right)
	}

	walk(t)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		x, xOk := <-ch1
		y, yOk := <-ch2

		if !xOk && !yOk {
			break
		} else if xOk != yOk {
			return false
		} else if x != y {
			return false
		}
	}

	return true
}

func channels() {
	// TestWalk()

	TestSame(1, 1)
	TestSame(1, 2)
}

func TestSame(x, y int) {
	t1 := tree.New(x)
	t2 := tree.New(y)

	fmt.Println(Same(t1, t2))
}

func TestWalk() {
	ch := make(chan int)

	go Walk(tree.New(1), ch)

	for val := range ch {
		fmt.Printf("%v,", val)
	}

	fmt.Println()
}

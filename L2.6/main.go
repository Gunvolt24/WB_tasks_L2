package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	fmt.Println("main (before): ", s, len(s), cap(s))
	fmt.Printf("main (before): %p->%p\n", &s, s)
	modifySlice(s)
	fmt.Println("main (after): ", s, len(s), cap(s))
	fmt.Printf("main (after): %p->%p\n", &s, s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	fmt.Println("  modifySlice (i): ", i, len(i), cap(i))
	fmt.Printf("  modifySlice (i): %p->%p\n", &i, i)
	i[0] = "3"
	fmt.Println("  modifySlice (i[0] = 3): ", i, len(i), cap(i))
	fmt.Printf("  modifySlice (i[0] = 3): %p->%p\n", &i, i)
	i = append(i, "4")
	fmt.Println("  modifySlice (append 4): ", i, len(i), cap(i))
	fmt.Printf("  modifySlice (append 4): %p->%p\n", &i, i)
	i[1] = "5"
	fmt.Println("  modifySlice: i[1] = 5", i, len(i), cap(i))
	fmt.Printf("  modifySlice i[1] = 5: %p->%p\n", &i, i)
	i = append(i, "6")
	fmt.Println("  modifySlice (append 6): ", i, len(i), cap(i))
	fmt.Printf("  modifySlice (append 6): %p->%p\n", &i, i)
}

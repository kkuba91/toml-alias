package main

// #cgo LDFLAGS: -L. -laliaslib
// #include "aliaslib.h"
import "C"

func main() {
	C.processAll()
}

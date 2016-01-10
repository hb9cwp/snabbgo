// from
//  https://golang.org/src/syscall/mmap_unix_test.go?m=text

// see also
//  [go-nuts] How to use mmap to get a pointer to a memory?
//  http://grokbase.com/t/gg/golang-nuts/12atgwwczh/go-nuts-how-to-use-mmap-to-get-a-pointer-to-a-memory
//
// Q:
// A: syscall.Mmap() returns a slice, why not just use the slice?
//    You could use a int64 as an index and pretend it is a pointer.
//    d[p] instead of just p, p is technically pointing to your value (just
//    inside of d...)
//    You also don't need unsafe for arithmetic on a int64 index
//    If you are doing something very strange that needs a pointer it might be a
//    better idea to use cgo to call mmap yourself.
//
// Q: ... the first argument need to be a file descriptor...
//    But I have no file to open.. I do not use mmap to map a file into the
//    memory, I just want use mmap to allocate some memory?
// A: Pass -1 as the file descriptor. Note that mmap is not portable but
//    that will work on most systems. See "man mmap" for details.

// +build darwin dragonfly freebsd linux netbsd openbsd

//package syscall_test
package main

import (
	"fmt"
	"reflect"
	"syscall"
	"testing"
	"unsafe"
)

//Mmap(t *testing.T) {
//func main(t *testing.T) {
func main() {
	const n = 1e1
	size := int(unsafe.Sizeof(0)) * n
	var t testing.T
	/*
		 man mmap
		 void *mmap(void *addr, size_t length, int prot, int flags,
		            int fd, off_t offset);

		 https://golang.org/src/syscall/syscall_linux.go?h=Mmap#L915
		 func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
				return mapper.Mmap(fd, offset, length, prot, flags)
		 }

		 from http://stackoverflow.com/questions/9203526/mapping-an-array-to-a-file-via-mmap-in-go
		 map_file, err := os.Create("/tmp/test.dat")
	*/
	//mmap, err := syscall.Mmap(-1, 0, syscall.Getpagesize(), syscall.PROT_NONE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
	mmap, err := syscall.Mmap(-1, 0, n*size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
	//mmap, err := syscall.Mmap(-1, 0, n * size, syscall.PROT_NONE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
	//mmap, err := syscall.Mmap(-1, 0, n * size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		t.Fatalf("Mmap: %v", err)
	}
	fmt.Printf("syscall.Mmap(): OK\n")

	fmt.Println(reflect.TypeOf(mmap)) // mmap is a slice []byte (or equivalent []uint8)
	fmt.Println(mmap[0])
	//mmap[0] = 'r'
	mmap[0] = byte(99) // byte is an alias for uint8 and is equivalent to uint8 in all ways
	//mmap[0] = uint8(0)
	mmap[1] = byte(101)
	fmt.Println(mmap[0])
	fmt.Println(mmap[1])

	/*
		map_array := (*[n]int)(unsafe.Pointer(&mmap[0]))
		for i := 0; i < n; i++ {
			map_array[i] = i * i
		}
		fmt.Println(*map_array)
	*/

	if err := syscall.Munmap(mmap); err != nil {
		t.Fatalf("Munmap: %v", err)
	}
	fmt.Printf("syscall.Munmap(): OK\n")
}

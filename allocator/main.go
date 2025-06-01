package main

import (
	"fmt"
	"sort"
	"unsafe"
)

func main() {
	memory := []byte{0xAA, 0x00, 0xBB, 0x00, 0x00, 0xCC}

	pointers := []unsafe.Pointer{
		unsafe.Pointer(&memory[0]),
		unsafe.Pointer(&memory[2]),
		unsafe.Pointer(&memory[5]),
	}

	fmt.Println("До дефрагментации:")
	for i, b := range memory {
		fmt.Printf("memory[%v] = 0x%v\n", i, b)
	}

	Defragment(memory, pointers)

	fmt.Println("После дефрагментации:")
	for i, b := range memory {
		fmt.Printf("memory[%v] = 0x%v\n", i, b)
	}
}

type ptrInfo struct {
	index    int
	position uintptr
	ptr      unsafe.Pointer
}

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	baseAddr := uintptr(unsafe.Pointer(&memory[0]))

	ptrInfos := make([]ptrInfo, len(pointers))

	for i, ptr := range pointers {
		ptrInfos[i] = ptrInfo{
			index:    i,
			position: uintptr(ptr) - baseAddr,
			ptr:      ptr,
		}
	}

	sort.Slice(ptrInfos, func(i, j int) bool {
		return ptrInfos[i].position < ptrInfos[j].position
	})

	writePos := uintptr(0)
	for _, info := range ptrInfos {
		if info.position != writePos {
			memory[writePos] = memory[info.position]
		}

		pointers[info.index] = unsafe.Pointer(baseAddr + writePos)
		writePos++
	}

	for i := writePos; i < uintptr(len(memory)); i++ {
		memory[i] = 0x00
	}
}

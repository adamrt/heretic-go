// This file contains a way to read binary data from the FFT ISO.
// It should be expanded to also read the FFT bin file.
package fft

import (
	"encoding/binary"
	"log"
	"os"
)

const sectorSize int64 = 2048

func NewISOReader(filename string) ISOReader {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("open iso: %v", err)
	}
	return ISOReader{f}
}

type ISOReader struct {
	file *os.File
}

func (r ISOReader) Close() {
	r.file.Close()
}

// seekSector will seek to the specified sector of the iso file.
func (r ISOReader) seekSector(sector int64) {
	to := sector * sectorSize
	_, err := r.file.Seek(to, 0)
	if err != nil {
		log.Fatalf("seek to sector: %v", err)
	}
}

// seekPointer will seek to the specified sector, plus a little more, of the iso
// file. This is useful when using MeshFileHeader intra-file pointers.
func (r ISOReader) seekPointer(sector int64, ptr int64) {
	to := sector*sectorSize + ptr
	_, err := r.file.Seek(to, 0)
	if err != nil {
		log.Fatalf("seek to pointer: %v", err)
	}
}

func (r ISOReader) uint8() uint8 {
	size := 1
	data := make([]byte, size)
	n, err := r.file.Read(data)
	if err != nil || n != size {
		log.Fatal(err)
	}
	return data[0]
}

func (r ISOReader) uint16() uint16 {
	size := 2
	data := make([]byte, size)
	n, err := r.file.Read(data)
	if err != nil || n != size {
		log.Fatal(err)
	}
	return binary.LittleEndian.Uint16(data)
}

func (r ISOReader) uint32() uint32 {
	size := 4
	data := make([]byte, size)
	n, err := r.file.Read(data)
	if err != nil || n != size {
		log.Fatal(err)
	}
	return binary.LittleEndian.Uint32(data)
}

func (r ISOReader) int8() int8   { return int8(r.uint8()) }
func (r ISOReader) int16() int16 { return int16(r.uint16()) }
func (r ISOReader) int32() int32 { return int32(r.uint32()) }

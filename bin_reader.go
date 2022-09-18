package main

import (
	"log"
	"os"
	"strconv"
)

const SectorSize = 2048

type Map struct {
	textures []Texture
}

func NewBinReader(filename string) BinReader {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return BinReader{
		iso: f,
	}
}

type BinReader struct {
	iso *os.File
}

func (r BinReader) Map(num int) Map {
	gnsSector := GNSSectors[num]
	offset := gnsSector * SectorSize

	_, err := r.iso.Seek(offset, 0)
	if err != nil {
		log.Fatal(err)
	}

	records := r.parseGNSRecords()

	textures := []Texture{}
	// primaryMeshRecord := GNSRecord{}

	for _, record := range records {
		if record.Type() == ResourceTypeTexture {
			_, err := r.iso.Seek(record.Offset(), 0)
			if err != nil {
				log.Fatal(err)
			}
			textureData := make([]byte, record.Len())
			n, err := r.iso.Read(textureData)
			if err != nil || int64(n) != record.Len() {
				log.Fatalf("want %d got %d: %v", record.Len(), n, err)
			}
			textures = append(textures, NewTextureFFT(textureData))
		}
	}

	for index, texture := range textures {
		texture.WritePPM(strconv.Itoa(index) + ".ppm")
	}
	return Map{
		textures: textures,
	}
}

func (r BinReader) parseGNSRecords() []GNSRecord {
	records := []GNSRecord{}
	for {
		record := make(GNSRecord, 20)
		n, err := r.iso.Read(record)
		if err != nil {
			log.Fatal(err)
		} else if n != 20 {
			log.Fatalf("want 20 got %d", n)
		}

		if record.Type() == ResourceTypeIgnore {
			break
		}
		records = append(records, record)
	}
	return records
}

func (r BinReader) Close() {
	r.iso.Close()
}

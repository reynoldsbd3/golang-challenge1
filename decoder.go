package drum

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	p, err := getPattern(file)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getPattern(in io.Reader) (*Pattern, error) {

	var length uint8
	var magicNumberRef = [13]byte{
		0x53, 0x50, 0x4c, 0x49, // SPLI
		0x43, 0x45, 0x00, 0x00, // CE\0\0
		0x00, 0x00, 0x00, 0x00, // \0\0\0\0
		0x00, // \0
	}
	var magicNumber [13]byte
	p := &Pattern{}
	rawVersion := make([]byte, 32)

	err := binary.Read(in, binary.LittleEndian, &magicNumber)
	if err != nil {
		return nil, err
	}

	if magicNumber != magicNumberRef {
		return nil, errors.New("Invalid magic number")
	}

	err = binary.Read(in, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	err = binary.Read(in, binary.LittleEndian, &rawVersion)
	if err != nil {
		return nil, err
	}

	p.Version = string(rawVersion)

	err = binary.Read(in, binary.LittleEndian, &p.Tempo)
	if err != nil {
		return nil, err
	}

	return p, nil
}

type Pattern struct {
	Version string
	Tempo   float32
	Tracks  []Track
}

func (this Pattern) String() string {
	str := fmt.Sprintf("Saved with HW version: %s\n", this.Version)
	str += fmt.Sprintf("Tempo: %g\n", this.Tempo)
	return str
}

type Track struct {
	Id    uint32
	Label string
	Steps [16]byte
}

func (this Track) String() string {
	str := fmt.Sprintf("(%d) %s\t|", this.Id, this.Label)
	for i, step := range this.Steps {
		if step == 0 {
			str += "-"
		} else {
			str += "X"
		}

		if i%4 == 3 {
			str += "|"
		}
	}
	return str
}

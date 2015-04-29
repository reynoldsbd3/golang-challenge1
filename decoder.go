package drum

import (
	"encoding/binary"
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

	p := &Pattern{}

	h, err := getPatternHeader(in)
	if err != nil {
		return nil, err
	}

	p.Header = *h

	binary.Read(in, binary.LittleEndian, p.Tempo)

	return p, nil
}

func getPatternHeader(in io.Reader) (*PatternHeader, error) {

	h := &PatternHeader{}

	err := binary.Read(in, binary.BigEndian, h)
	if err != nil {
		return nil, err
	}

	return h, nil
}

type PatternHeader struct {
	Magic   [6]byte
	_       [4]byte
	Size    uint32
	Version [32]byte
}

type Pattern struct {
	Header PatternHeader
	Tempo  float32
	Tracks []Track
}

type TrackHeader struct {
	Id        uint32
	LabelSize uint8
	Label     []byte
}

type Track struct {
	Header TrackHeader
	Steps  [16]byte
}

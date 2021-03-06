package drum

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
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

	var length uint64
	magicNumber := make([]byte, 6)
	p := &Pattern{}
	rawVersion := make([]byte, 32)

	err := binary.Read(in, binary.LittleEndian, &magicNumber)
	if err != nil {
		return nil, err
	}

	if string(magicNumber) != "SPLICE" {
		return nil, fmt.Errorf("invalid magic number")
	}

	err = binary.Read(in, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	err = binary.Read(in, binary.LittleEndian, &rawVersion)
	if err != nil {
		return nil, err
	}

	p.Version = string(rawVersion[:strings.Index(string(rawVersion), "\x00")])
	length -= 32

	err = binary.Read(in, binary.LittleEndian, &p.Tempo)
	if err != nil {
		return nil, err
	}

	length -= 4

	var t *Track
	var trackLength uint64
	for length > 0 {

		t, trackLength, err = getTrack(in)
		if err != nil {
			return nil, err
		}

		p.Tracks = append(p.Tracks, *t)
		length -= trackLength
	}

	return p, nil
}

func getTrack(in io.Reader) (*Track, uint64, error) {

	var labelLength uint8
	t := &Track{}

	err := binary.Read(in, binary.LittleEndian, &t.ID)
	if err != nil {
		return nil, 0, err
	}

	err = binary.Read(in, binary.LittleEndian, &labelLength)
	if err != nil {
		return nil, 4, err
	}

	rawLabel := make([]byte, labelLength)

	err = binary.Read(in, binary.LittleEndian, &rawLabel)
	if err != nil {
		return nil, 5, err
	}

	t.Label = string(rawLabel)

	err = binary.Read(in, binary.LittleEndian, &t.Steps)
	if err != nil {
		return nil, uint64(5 + labelLength), err
	}

	return t, uint64(21 + labelLength), nil
}

// Pattern represents a pattern to be reconstructed and played by the drum machine.
type Pattern struct {
	Version string
	Tempo   float32
	Tracks  []Track
}

func (pattern Pattern) String() string {
	str := fmt.Sprintf("Saved with HW Version: %s\n", pattern.Version)
	str += fmt.Sprintf("Tempo: %g\n", pattern.Tempo)
	for _, track := range pattern.Tracks {
		str += track.String()
	}
	return str
}

// Track represents one of the tracks comprising a pattern.
type Track struct {
	ID    uint32
	Label string
	Steps [16]byte
}

func (track Track) String() string {
	str := fmt.Sprintf("(%d) %s\t|", track.ID, track.Label)
	for i, step := range track.Steps {
		if step == 0 {
			str += "-"
		} else {
			str += "x"
		}

		if i%4 == 3 {
			str += "|"
		}
	}
	str += "\n"
	return str
}

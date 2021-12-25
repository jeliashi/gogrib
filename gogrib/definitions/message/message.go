package message

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strconv"

	"github.com/jeliashi/gogrib/gogrib/definitions"
)

const Grib2Divider = 1000000.0

var tables map[string]definitions.Table = generateTables("tables/27")

func generateTables(prefix string) map[string]definitions.Table {
	files, _ := filepath.Glob(fmt.Sprintf("%s/*.table", prefix))
	tables := make(map[string]definitions.Table)
	for _, f := range files {
		tables[f] = definitions.FilenameToTable(f)
	}
	return tables
}

func GetSectionLengthAndNumber(r io.Reader) (uint32, int, error) {
	b, err := ReadNBytes(r, 4)
	if err != nil {
		return 0, 0, errors.New("unable to read first 4 bytes of section")
	}
	if strconv.QuoteToASCII(string(b)) == "7777" {
		return 4, 8, nil
	}
	b, err = ReadNBytes(r, 5)
	if err != nil {
		return 0, 0, errors.New("unable to read first 5 bytes of section")
	}
	length := binary.BigEndian.Uint32(b[:4])
	section := int(b[4])
	return length, section, nil
}

func ReadNBytes(r io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ScaleAndValueToReal(scale int, val int) float64 {
	return float64(val) / math.Pow10(scale)
}

func DegToRad(v float64) float64 {
	return math.Pi * v / 180.0
}

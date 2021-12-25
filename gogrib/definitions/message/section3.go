package message

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/jeliashi/gogrib/gogrib/definitions"
)

// var srcGridDefTable definitions.Table = definitions.FilenameToTable("tables/27/3.0.table")
// var interpNPointsTable definitions.Table = definitions.FilenameToTable("tables/27/3.11.table")

type Section3 struct {
	SourceOfGridDefinition definitions.TableEntry
	NPts                   uint32
	NOctetsPerPt           uint8
	InterpOfNPts           definitions.TableEntry
	GridDefInt             uint16
	GridDefTemplateNumber  definitions.TableEntry
	GridDefinition         interface{}
}

func ParseSection3(r io.Reader) Section3 {
	buf, err := ReadNBytes(r, 5)
	if err != nil {
		fmt.Printf("failure to launch section 3")
		return Section3{}
	}
	buf, err = ReadNBytes(r, 1)
	if err != nil {
		fmt.Printf("failure to launch section 3")
		return Section3{}
	}
	src := definitions.GetEntryFromTable(int(buf[0]), tables["3.0.table"])
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("failure to launch section 3")
		return Section3{}
	}
	nPts := binary.BigEndian.Uint32(buf)
	buf, err = ReadNBytes(r, 1)
	if err != nil {
		fmt.Printf("failure to launch section 3")
		return Section3{}
	}
	nOctetsPerPt := uint8(buf[0])
	buf, err = ReadNBytes(r, 1)
	if err != nil {
		fmt.Printf("failure to launch section 3")
		return Section3{}
	}
	interpNPoints := definitions.GetEntryFromTable(int(buf[0]), tables["3.11.table"])
	buf, err = ReadNBytes(r, 2)
	if err != nil {
		fmt.Printf("failure to launch section 3")
		return Section3{}
	}
	gridDefInt := binary.BigEndian.Uint16(buf)
	gridDefDesc := definitions.GetEntryFromTable(int(gridDefInt), tables["3.1.table"])
	var grid interface{}
	switch gridDefInt {
	case 30:
		grid = ParseLambertConformal(r)
	default:
		grid = EmptyStruct{Empty: true}
	}
	return Section3{src, nPts, nOctetsPerPt, interpNPoints, gridDefInt, gridDefDesc, grid}

}

type EmptyStruct struct {
	Empty bool
}

package message

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/jeliashi/gogrib/gogrib/definitions"
)

var centreTable definitions.Table = definitions.FilenameToTable("definitions/common/c-11.table")
var versionTable definitions.Table = tables["1.0.table"]
var sigOfRefTimeTable definitions.Table = tables["1.2.table"]
var productionStatusTable definitions.Table = tables["1.3.table"]
var typeOfDataTable definitions.Table = tables["1.4.table"]

type Section1 struct {
	Centre                          definitions.TableEntry //'common/c-11.table'
	SubCentre                       uint16                 //dump
	TableVersion                    definitions.TableEntry //'grib2/tables/1.0.table'
	LocalTableVersion               uint8                  //'grib2/tables/local/[centreForLocal]/1.1.table' ;
	SignificanceOfReferenceTime     definitions.TableEntry //'1.2.table'
	Year                            uint16                 //last 4 digits
	Month                           uint8
	Day                             uint8
	Hour                            uint8
	Minute                          uint8
	Second                          uint8
	ProductionStatusOfProcessedData definitions.TableEntry //'1.3.table'
	TypeOfProcessedData             definitions.TableEntry //'1.4.table'
}

func ParseSection1(r io.Reader) Section1 {
	_, err := ReadNBytes(r, 5)
	if err != nil {
		fmt.Println("unable to make section 1")
		return Section1{}
	}
	b, _ := ReadNBytes(r, 2)
	centre := definitions.GetEntryFromTable(int(binary.BigEndian.Uint16(b)), centreTable)
	b, _ = ReadNBytes(r, 2)
	subCentre := binary.BigEndian.Uint16(b)
	b, _ = ReadNBytes(r, 1)
	tableVersion := definitions.GetEntryFromTable(int(uint8(b[0])), versionTable)
	b, _ = ReadNBytes(r, 1)
	localVersion := uint8(b[0])
	b, _ = ReadNBytes(r, 1)
	sigOfRefTime := definitions.GetEntryFromTable(int(uint8(b[0])), sigOfRefTimeTable)
	b, _ = ReadNBytes(r, 2)
	year := binary.BigEndian.Uint16(b)
	b, _ = ReadNBytes(r, 1)
	month := uint8(b[0])
	b, _ = ReadNBytes(r, 1)
	day := uint8(b[0])
	b, _ = ReadNBytes(r, 1)
	hour := uint8(b[0])
	b, _ = ReadNBytes(r, 1)
	minute := uint8(b[0])
	b, _ = ReadNBytes(r, 1)
	second := uint8(b[0])
	b, _ = ReadNBytes(r, 1)
	prodStatus := definitions.GetEntryFromTable(int(uint8(b[0])), productionStatusTable)
	b, _ = ReadNBytes(r, 1)
	typeOfData := definitions.GetEntryFromTable(int(uint8(b[0])), typeOfDataTable)
	return Section1{centre, subCentre, tableVersion, localVersion, sigOfRefTime, year, month, day, hour, minute, second, prodStatus, typeOfData}
}

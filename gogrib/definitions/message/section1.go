package message

import (
	"encoding/binary"
	"fmt"
	"gogrib/definitions"
	"io"
)

var centreTable definitions.Table = definitions.FilenameToTable("../common/c-11.table")
const versionTable := tables["1.0.table"]
const sigOfRefTimeTable := tables["1.2.table"
const productionStatusTable := tables["1.3.table"]
const typeOfDataTable := tables["1.4.table"]

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
	centre := definitions.GetEntryFromTable(binary.BigEndian.Uint16(ReadNBytes(r, 2)), centreTable)
	subCentre := binary.BigEndian.Uint16(ReadNBytes(r, 2))
	tableVersion := definitions.GetEntryFromTable(uint8(ReadNBytes(r,1)), versionTable)
	localVersion := uint8(ReadNBytes(r,1))
	sigOfRefTime := definitions.GetEntryFromTable(uint8(ReadNBytes(r,1)), sigOfRefTimeTable)
	year := binary.BigEndian.Uint16(ReadNBytes(r,2))
	month := uint8(ReadNBytes(r,1))
	day := uint8(ReadNBytes(r,1))
	hour := uint8(ReadNBytes(r,1))
	minute := uint8(ReadNBytes(r,1))
	second := uint8(ReadNBytes(r,1))
	prodStatus := definitions.GetEntryFromTable(uint8(ReadNBytes(r,1)), productionStatusTable)
	typeOfData := definitions.GetEntryFromTable(uint8(ReadNBytes(r,1)), typeOfDataTable)
	return Section1{centre, subCentre, tableVersion, localVersion, sigOfRefTime, year, month, day, hour, minute, second, prodStatus, typeOfData}
}

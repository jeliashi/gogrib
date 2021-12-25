package message

import (
	"encoding/binary"
	"io"

	"github.com/jeliashi/gogrib/gogrib/definitions"
)

type ShapeOfTheEarth struct {
	Shape             definitions.TableEntry
	EarthMeasurements map[string]float64
}

func ParseShapeOfTheEarth(r io.Reader) ShapeOfTheEarth {
	measurements := make(map[string]float64)
	buf, err := ReadNBytes(r, 1)
	if err != nil {
		return ShapeOfTheEarth{}
	}
	shape_int := uint8(buf[0])
	shape := definitions.GetEntryFromTable(int(shape_int), tables["3.2.table"])

	buf, err = ReadNBytes(r, 1)
	if err != nil {
		return ShapeOfTheEarth{}
	}
	scale := uint8(buf[0])
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		return ShapeOfTheEarth{}
	}
	radius_m := binary.BigEndian.Uint32(buf)
	buf, err = ReadNBytes(r, 1)
	if err != nil {
		return ShapeOfTheEarth{}
	}
	major_scale := uint8(buf[0])
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		return ShapeOfTheEarth{}
	}
	major_m := binary.BigEndian.Uint32(buf)
	buf, err = ReadNBytes(r, 1)
	if err != nil {
		return ShapeOfTheEarth{}
	}
	minor_scale := uint8(buf[0])
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		return ShapeOfTheEarth{}
	}
	minor_m := binary.BigEndian.Uint32(buf)
	switch shape_int {
	case 0:
		measurements["radius"] = 6367470.0
		measurements["radiusOfTheEarth"] = measurements["radius"]
		measurements["radiusInMetres"] = measurements["radius"]
		measurements["earthIsOblate"] = 0
	case 1:
		measurements["radius"] = ScaleAndValueToReal(int(scale), int(radius_m))
		measurements["radiusOfTheEarth"] = measurements["radius"]
		measurements["radiusInMetres"] = measurements["radius"]
		measurements["earthIsOblate"] = 0
	case 6:
		measurements["radius"] = 6371229.0
		measurements["radiusOfTheEarth"] = measurements["radius"]
		measurements["radiusInMetres"] = measurements["radius"]
		measurements["earthIsOblate"] = 0
	case 8:
		measurements["radius"] = 6371200.0
		measurements["radiusOfTheEarth"] = measurements["radius"]
		measurements["radiusInMetres"] = measurements["radius"]
		measurements["earthIsOblate"] = 0
	case 2:
		measurements["earthMajorAxis"] = 6378160.0
		measurements["earthMinorAxis"] = 6356775.0
		measurements["earthMajorAxisInMetres"] = measurements["earthMajorAxis"]
		measurements["earthMinorAxisInMetres"] = measurements["earthMinorAxis"]
		measurements["earthIsOblate"] = 1
	case 3:
		measurements["earthMajorAxis"] = ScaleAndValueToReal(int(major_scale), int(major_m))
		measurements["earthMinorAxis"] = ScaleAndValueToReal(int(minor_scale), int(minor_m))
		measurements["earthMajorAxisInMetres"] = measurements["earthMajorAxis"] * 1000
		measurements["earthMinorAxisInMetres"] = measurements["earthMinorAxis"] * 1000
		measurements["earthIsOblate"] = 1
	case 7:
		measurements["earthMajorAxis"] = ScaleAndValueToReal(int(major_scale), int(major_m))
		measurements["earthMinorAxis"] = ScaleAndValueToReal(int(minor_scale), int(minor_m))
		measurements["earthMajorAxisInMetres"] = measurements["earthMajorAxis"]
		measurements["earthMinorAxisInMetres"] = measurements["earthMinorAxis"]
		measurements["earthIsOblate"] = 1
	case 4:
		measurements["earthMajorAxis"] = 6378137.0
		measurements["earthMinorAxis"] = 6356752.314
		measurements["earthMajorAxisInMetres"] = measurements["earthMajorAxis"]
		measurements["earthMinorAxisInMetres"] = measurements["earthMinorAxis"]
		measurements["earthIsOblate"] = 1
	case 5:
		measurements["earthMajorAxis"] = 6378137.0
		measurements["earthMinorAxis"] = 6356752.314
		measurements["earthMajorAxisInMetres"] = measurements["earthMajorAxis"]
		measurements["earthMinorAxisInMetres"] = measurements["earthMinorAxis"]
		measurements["earthIsOblate"] = 1
	case 9:
		measurements["earthMajorAxis"] = 6377563.396
		measurements["earthMinorAxis"] = 6377563.396
		measurements["earthMajorAxisInMetres"] = measurements["earthMajorAxis"]
		measurements["earthMinorAxisInMetres"] = measurements["earthMinorAxis"]
		measurements["earthIsOblate"] = 1
	}
	return ShapeOfTheEarth{shape, measurements}

}

package message

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type LambertConformal struct {
	Shape                     ShapeOfTheEarth
	Nx                        uint32
	Ny                        uint32
	LatitudeOfFirstGridPoint  float64
	LongitudeOfFirstGridPoint float64
	lats                      []float64
	lons                      []float64
	scanning                  ScanningMode
}

func ParseLambertConformal(r io.Reader) LambertConformal {
	shape := ParseShapeOfTheEarth(r)
	buf, err := ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	nx := binary.BigEndian.Uint32(buf)
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	ny := binary.BigEndian.Uint32(buf)
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	firstLat := float64(int32(binary.BigEndian.Uint32(buf))) / Grib2Divider
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	firstLon := float64(int32(binary.BigEndian.Uint32(buf))) / Grib2Divider
	_, _ = ReadNBytes(r, 1)
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	LaD := float64(int32(binary.BigEndian.Uint32(buf))) / Grib2Divider
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	LoV := float64(int32(binary.BigEndian.Uint32(buf))) / Grib2Divider
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	Dx := float32(binary.BigEndian.Uint32(buf)) / 1000.0
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	Dy := float32(binary.BigEndian.Uint32(buf)) / 1000.0
	buf, err = ReadNBytes(r, 1)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	// flags := definitions.GetEntryFromTable(int(binary.BigEndian.Uint32(buf)), tables["3.5.table"])
	scanningMode := ParseScanningMode(r)
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	Latin1 := float64(int32(binary.BigEndian.Uint32(buf))) / Grib2Divider
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	Latin2 := float64(int32(binary.BigEndian.Uint32(buf))) / Grib2Divider
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	// LatitudeOfSouthernPole := float64(int(binary.BigEndian.Uint32(buf))) / Grib2Divider
	buf, err = ReadNBytes(r, 4)
	if err != nil {
		fmt.Printf("unable to parse LC")
		return LambertConformal{}
	}
	// LongitudeOfSouthernPole := float64(int(binary.BigEndian.Uint32(buf))) / Grib2Divider
	lats, lons := IterLambertConformal(shape, nx, ny, LoV, LaD, Latin1, Latin2, firstLat, firstLon, Dx, Dy, scanningMode)
	return LambertConformal{shape, nx, ny, firstLat, firstLon, lats, lons, scanningMode}
}

func IterLambertConformal(shape ShapeOfTheEarth, nx uint32, ny uint32, LoV float64, Lad float64, sLat1 float64, sLat2 float64, sLatFirst float64, sLonFirst float64, sDx float32, sDy float32, scanning ScanningMode) ([]float64, []float64) {
	var lats []float64
	var lons []float64
	if shape.EarthMeasurements["earthIsOblate"] == 1 {
		lats, lons = IterOblate(shape, nx, ny, LoV, Lad, sLat1, sLat2, sLatFirst, sLonFirst, sDx, sDy)
	} else {
		lats, lons = IterSphere(shape, nx, ny, LoV, Lad, sLat1, sLat2, sLatFirst, sLonFirst, sDx, sDy)
	}
	lats, lons = TransformLatLon(lats, lons, scanning, int(nx), int(ny))
	return lats, lons
}

func IterOblate(shape ShapeOfTheEarth, nx uint32, ny uint32, LoV float64, Lad float64, sLat1 float64, sLat2 float64, sLatFirst float64, sLonFirst float64, sDx float32, sDy float32) ([]float64, []float64) {
	return []float64{1.0}, []float64{1.0}
}

func IterSphere(shape ShapeOfTheEarth, nx uint32, ny uint32, LoV float64, Lad float64, sLat1 float64, sLat2 float64, sLatFirst float64, sLonFirst float64, sDx float32, sDy float32) ([]float64, []float64) {
	npts := nx * ny
	lat_out := make([]float64, npts)
	lon_out := make([]float64, npts)
	sLat1_r := DegToRad(sLat1)
	sLat2_r := DegToRad(sLat2)
	var n float64
	if math.Abs(sLat1_r-sLat2_r) < 1e-9 {
		n = math.Sin(sLat1_r)
	} else {
		n = math.Log(math.Cos(sLat1_r)/math.Cos(sLat2_r)) / math.Log(math.Tan(math.Pi/4.0+sLat2_r/2.0)/math.Tan((math.Pi/4.0+sLat1_r/2.0)))
	}
	f := (math.Cos(sLat1_r) * math.Pow(math.Tan(math.Pi/4.0+sLat1_r/2.0), n)) / n
	rho := shape.EarthMeasurements["radius"] * f * math.Pow(math.Tan(math.Pi/4+sLat1_r/2.0), -n)
	rho0 := shape.EarthMeasurements["radius"] * f * math.Pow(math.Tan(math.Pi/4+DegToRad(Lad)/2.0), -n)
	if n < 0 {
		rho0 = -1.0 * rho0
	}
	lonDiff := DegToRad(sLonFirst) - DegToRad(LoV)
	if lonDiff > math.Pi {
		lonDiff = lonDiff - 2*math.Pi
	}
	if lonDiff < -math.Pi {
		lonDiff = lonDiff + 2*math.Pi
	}
	angle := n * lonDiff
	x0 := rho * math.Sin(angle)
	y0 := rho - rho*math.Cos(angle)
	for j := 0; j < int(ny); j++ {
		y := y0 + float64(j)*float64(sDy)
		if n < 0 {
			y = -y
		}
		tmp := rho0 - y
		tmp2 := tmp * tmp
		for i := 0; i < int(nx); i++ {
			index := i + j*int(nx)
			x := x0 + float64(i)*float64(sDx)
			if n < 0 {
				x = -x
			}
			angle = math.Atan2(x, tmp)
			rho = math.Sqrt(x*x + tmp2)
			if n <= 0 {
				rho = -rho
			}
			lonDeg := LoV + (angle/n)*180.0/math.Pi
			latDeg := (2.0*math.Atan(math.Pow(shape.EarthMeasurements["radius"]*f/rho, 1.0/n)) - math.Pi/2) * 180.0 / math.Pi
			if lonDeg >= 360.0 {
				lonDeg -= 360.0
			} else if lonDeg < 0 {
				lonDeg += 360.0
			}
			lat_out[index] = latDeg
			lon_out[index] = lonDeg
		}
	}
	return lat_out, lon_out
}

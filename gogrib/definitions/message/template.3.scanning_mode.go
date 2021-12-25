package message

import (
	"fmt"
	"io"
)

type ScanningMode struct {
	iScansNegatively       bool
	jScansPositively       bool
	jPointsAreConsecutive  bool
	alternativeRowScanning bool

	scanningMode5 bool
	scanningMode6 bool
	scanningMode7 bool
	scanningMode8 bool
}

func ParseScanningMode(r io.Reader) ScanningMode {
	buf, err := ReadNBytes(r, 1)
	if err != nil {
		fmt.Println("Unable to parse scanning")
	}

	t := make([]bool, 8)
	for i, x := range buf {
		for j := 0; j < 8; j++ {
			if (x<<uint(j))&0x80 == 0x80 {
				t[8*i+j] = true
			}
		}
	}
	return ScanningMode{t[0], t[1], t[2], t[3], t[4], t[5], t[6], t[7]}
}

func TransformLatLon(lats []float64, lons []float64, scanning ScanningMode, nx int, ny int) ([]float64, []float64) {
	if !scanning.iScansNegatively && scanning.jScansPositively && !scanning.jPointsAreConsecutive && !scanning.alternativeRowScanning {
		return lats, lons
	}
	if scanning.alternativeRowScanning {
		lats = alternativeRowScanning(lats, nx, ny)
		lons = alternativeRowScanning(lons, nx, ny)
	}
	if scanning.iScansNegatively {
		lats = FlipX(lats, nx, ny)
		lons = FlipX(lons, nx, ny)
	}
	if scanning.jScansPositively {
		lats = FlipY(lats, nx, ny)
		lons = FlipY(lons, nx, ny)
	}
	if scanning.jPointsAreConsecutive {
		lats = Transpose(lats, nx, ny)
		lons = Transpose(lons, nx, ny)
	}
	return lats, lons
}

func Transpose(data []float64, nx int, ny int) []float64 {
	n := ny * nx
	out := make([]float64, n)
	for index := 0; index < n; index++ {
		i := index % nx
		j := index / nx
		newIndex := j + i*ny
		out[newIndex] = data[index]
	}
	return out
}
func FlipX(data []float64, nx int, ny int) []float64 {
	n := nx * ny
	out := make([]float64, n)
	for index := 0; index < n; index++ {
		i := index % nx
		j := index / nx
		newIndex := j*nx + nx - i - 1
		out[newIndex] = data[index]
	}
	return out
}
func FlipY(data []float64, nx int, ny int) []float64 {
	n := nx * ny
	out := make([]float64, n)
	for index := 0; index < n; index++ {
		i := index % nx
		j := index / nx
		newIndex := (ny-j-1)*nx + i
		out[newIndex] = data[index]
	}
	return out
}

func alternativeRowScanning(data []float64, nx int, ny int) []float64 {
	n := nx * ny
	out := make([]float64, n)
	for index := 0; index < n; index++ {
		i := index % nx
		j := index / nx
		var newIndex int
		if j%2 == 1 {
			newIndex = j*nx + i
		} else {
			newIndex = j*nx + nx - i - 1
		}
		out[newIndex] = data[index]
	}
	return out
}

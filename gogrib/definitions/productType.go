package definitions

import (
	"bufio"
	"strconv"
)

type ProductType int

const (
	ProductGrib ProductType = iota
	ProductGrib2
	ProductOther
)

func DetermineProductType(r bufio.Reader) ProductType {
	productTypeBuffer, err := r.Peek(4)
	if err != nil {
		return ProductOther
	}
	productDescriptor := strconv.QuoteToASCII(string(productTypeBuffer))
	if productDescriptor != "GRIB" {
		return ProductOther
	}
	return determineGribType(r)
}

func determineGribType(r bufio.Reader) ProductType {
	section0information, err := r.Peek(8)
	if err != nil {
		return ProductGrib
	}
	if rune(section0information[7]) == '2' {
		return ProductGrib2
	}
	return ProductGrib
}

package main

import (
	"os"

	"github.com/jeliashi/gogrib/gogrib/definitions/message"
)

func main() {
	r, _ := os.Open("/Users/jeliashiv/hrrr.t23z.wrfsfcf00.grib2")
	defer r.Close()
	_ = message.ParseSection0(r)
	_ = message.ParseSection1(r)
	_ = message.ParseSection3(r)
}

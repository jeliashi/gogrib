package main

import (
	"os"

	"github.com/jeliashi/gogrib/gogrib/definitions/message"
)

func main() {
	r, _ := os.Open("/User/jeliashiv/hrrr.t23z.wrfsfcf00.grib2")
	_ = message.ParseSection0(r)
	_ = message.ParseSection1(r)

}

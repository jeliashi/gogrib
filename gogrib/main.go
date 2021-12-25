package main

import (
	"fmt"
	"os"

	"github.com/jeliashi/gogrib/gogrib/definitions/message"
)

func main() {
	r, _ := os.Open("/Users/jeliashiv/hrrr.t23z.wrfsfcf00.grib2")
	defer r.Close()
	s0 := message.ParseSection0(r)
	fmt.Println(s0)
	s1 := message.ParseSection1(r)
	fmt.Println(s1)
	s3 := message.ParseSection3(r)
	fmt.Println(s3)
}

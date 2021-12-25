package message

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/jeliashi/gogrib/gogrib/definitions"
)

var discipleTable definitions.Table = tables["0.0.table"]

type Section0 struct {
	MessageSize uint64
	Discipline  definitions.TableEntry
}

func ParseSection0(r io.Reader) Section0 {
	buf := make([]byte, 8)
	_, err := r.Read(buf)
	if err != nil {
		fmt.Println(r)
		fmt.Println(err.Error())
		fmt.Println("unable to parse section 0")
		return Section0{}
	}
	disciple := definitions.GetEntryFromTable(int(buf[6]), discipleTable)
	_, err = r.Read(buf)
	if err != nil {
		fmt.Println("unable to parse section 0, but I tried!")
	}
	return Section0{
		MessageSize: binary.BigEndian.Uint64(buf),
		Discipline:  disciple}
}

package definitions

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TableEntry struct {
	Id          int
	ShortName   string
	Description string
}

type Table struct {
	Name    string
	Entries []TableEntry
}

func FilenameToTable(f string) Table {
	file, err := os.Open(f)
	if err != nil {
		fmt.Printf("Unable to read table %s", f)
		return Table{}
	}
	defer file.Close()
	var entries []TableEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line_values := strings.Split(scanner.Text(), " ")
		if len(line_values) < 3 || line_values[0] == string('#') || line_values[0][0] == '#' {
			continue
		}
		id, _ := strconv.Atoi(line_values[0])
		desc := strings.Join(line_values[2:], " ")
		entries = append(entries, TableEntry{Id: id, ShortName: line_values[1], Description: desc})
	}
	return Table{f, entries}
}

func GetEntryFromTable(key int, t Table) TableEntry {
	for _, e := range t.Entries {
		if e.Id == key {
			return e
		}
	}
	return TableEntry{}
}

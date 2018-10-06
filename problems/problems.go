package problems

import (
	"flag"
	"go/token"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

var count, limit uint
var table = newTable()

func init() {
	flag.UintVar(&limit, `limit`, 0, `limit the max problems to check.`)
	flag.Parse()
}

func Add(position token.Position, desc, rule string) {
	table.Append([]string{positionString(position), desc, rule})

	count++
	if limit > 0 && count >= limit {
		table.Render()
		os.Exit(1)
	}
}

func Render() {
	table.Render()
}

func Count() uint {
	return count
}

func Clear() {
	table = newTable()
}

func newTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{`position`, `problem`, `rule`})
	table.SetColWidth(100) // set max column width
	return table
}

func positionString(position token.Position) string {
	pos := position.Filename
	if position.Line > 0 {
		pos += `:` + strconv.Itoa(position.Line)
		if position.Column > 0 {
			pos += `:` + strconv.Itoa(position.Column)
		}
	}
	return pos
}

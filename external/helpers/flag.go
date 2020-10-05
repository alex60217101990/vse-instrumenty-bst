package helpers

import (
	"flag"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func PrintFlags() {
	fmt.Fprintf(os.Stderr, "Usage: service (%s) [options] param>:\n", os.Args[0])

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Flag name", "Flag Value", "Flag description"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	counter := 0
	flag.VisitAll(func(f *flag.Flag) {
		table.Append([]string{f.Name, fmt.Sprintf("%v", f.Value), fmt.Sprintf("%v", f.Usage)})
		if counter%3 == 0 {
			table.SetRowLine(true)
		}
		counter++
	})
	table.Render()
}

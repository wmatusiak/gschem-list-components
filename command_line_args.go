package main

import (
	"flag"
)

type commandLineArgs struct {
	InFiles     []string
	SortBy      sortByValue
	ReverseSort bool
	Merge       bool
}

func NewCommandLineArgs() commandLineArgs {
	var result commandLineArgs
	flag.BoolVar(&result.Merge, "merge", false, "merge same components to single output line with count added")
	// SortBy
	result.SortBy = NewSortByValue(
		map[string]bool{
			"Device":    true,
			"Value":     true,
			"Footprint": true,
			"Refdes":    true,
			"Count":     true,
		},
	)

	result.SortBy.Set("Device")
	flag.Var(&result.SortBy, "sort-by", "column name to sort (default: Device)")
	//END of SortBy

	flag.BoolVar(&result.ReverseSort, "reverse", false, "Revers sort")
	flag.Parse()
	result.InFiles = flag.Args()
	return result
}

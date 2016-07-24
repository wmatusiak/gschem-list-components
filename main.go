package main

func main() {
	args := NewCommandLineArgs()
	components := ParseFiles(args.InFiles)
	if args.Merge {
		components.Merge()
	}

	components.Sort(args.SortBy.String(), args.ReverseSort)

	if args.Merge {
		components.Print([]string{"Device", "Value", "Footprint", "Refdes", "Count"})
	} else {
		components.Print([]string{"Device", "Value", "Footprint", "Refdes"})
	}
}

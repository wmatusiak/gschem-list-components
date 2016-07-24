package main

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
)

type componentArray []component

func NewComponentArray(len int) componentArray {
	return make([]component, 0, len)
}

func (a componentArray) Len() int {
	return len(a)
}

func (a componentArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a *componentArray) Remove(i int) {
}

func (a componentArray) Print(columns []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)
	for _, comp := range a {
		table.Append(comp.GetFieldsValues(columns))
	}

	table.Render()
}

func (a *componentArray) Merge() {
	tmp := make(map[string]component)
	for _, c := range *a {
		key := c.Device + "_" + c.Footprint + "_" + c.Value
		if _, exists := tmp[key]; exists {
			cmp := tmp[key]
			cmp.AddRefdes(c.Refdes[0])
			tmp[key] = cmp
		} else {
			tmp[key] = c
		}
	}

	*a = (*a)[:0]
	for _, c := range tmp {
		*a = append(*a, c)
	}
}

func (a *componentArray) Sort(columnName string, reverse bool) {
	var s sort.Interface
	switch columnName {
	case "Device":
		s = componentSortByDevice{*a}
	case "Value":
		s = componentSortByValue{*a}
	case "Footprint":
		s = componentSortByFootprint{*a}
	case "Refdes":
		s = componentSortByRefdes{*a}
	case "Count":
		s = componentSortByCount{*a}
	}

	if reverse {
		s = sort.Reverse(s)
	}

	sort.Sort(s)
}

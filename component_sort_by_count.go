package main

type componentSortByCount struct {
	componentArray
}

func (a componentSortByCount) Less(i, j int) bool {
	return len(a.componentArray[i].Refdes) < len(a.componentArray[j].Refdes)
}

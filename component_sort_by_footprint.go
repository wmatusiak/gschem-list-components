package main

type componentSortByFootprint struct {
	componentArray
}

func (a componentSortByFootprint) Less(i, j int) bool {
	return a.componentArray[i].Footprint < a.componentArray[j].Footprint
}

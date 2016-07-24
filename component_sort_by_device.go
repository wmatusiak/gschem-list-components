package main

type componentSortByDevice struct {
	componentArray
}

func (a componentSortByDevice) Less(i, j int) bool {
	return a.componentArray[i].Device < a.componentArray[j].Device
}

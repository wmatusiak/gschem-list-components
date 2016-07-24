package main

type componentSortByRefdes struct {
	componentArray
}

func (a componentSortByRefdes) Less(i, j int) bool {
	return a.componentArray[i].GetRefdesAsString() < a.componentArray[j].GetRefdesAsString()
}

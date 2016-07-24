package main

type componentSortByValue struct {
	componentArray
}

func (a componentSortByValue) Less(i, j int) bool {
	return a.componentArray[i].Value < a.componentArray[j].Value
}

package main

import (
	"fmt"
)

type sortByValue struct {
	Name             string
	ValidColumnNames validColumnNames
}

func NewSortByValue(validColumnNames map[string]bool) sortByValue {
	return sortByValue{
		Name:             "",
		ValidColumnNames: validColumnNames,
	}
}

func (this *sortByValue) Set(name string) error {
	if !this.ValidColumnNames.IsValid(name) {
		this.Name = ""
		return ColumnNameParserError{
			message: fmt.Sprintf("%s is not valid column name. Valid names are: %s", name, this.ValidColumnNames),
		}
	}

	this.Name = name
	return nil
}

func (this sortByValue) String() string {
	return this.Name
}

package main

import (
	"sort"
	"strconv"
	"strings"
)

type component struct {
	Refdes    []string
	Device    string
	Footprint string
	Value     string
}

func NewComponent(in []string) component {
	var out component
	for _, line := range in {
		if strings.Contains(line, "=") {
			switch line[:strings.Index(line, "=")] {
			case "refdes":
				out.Refdes = make([]string, 0, 1)
				out.Refdes = append(out.Refdes, line[strings.Index(line, "=")+1:])
			case "device":
				out.Device = line[strings.Index(line, "=")+1:]
			case "footprint":
				out.Footprint = line[strings.Index(line, "=")+1:]
			case "value":
				out.Value = line[strings.Index(line, "=")+1:]
			}
		}
	}

	return out
}

func (c *component) AddRefdes(refdes string) {
	c.Refdes = append(c.Refdes, refdes)
}

func (c component) GetRefdesAsString() string {
	sort.Strings(c.Refdes)
	res := ""
	for _, r := range c.Refdes {
		if res != "" {
			res += ", "
		}

		res += r
	}

	return res
}

func (c component) GetFieldsValues(fieldNames []string) []string {
	res := make([]string, 0, len(fieldNames))
	for _, f := range fieldNames {
		switch f {
		case "Device":
			res = append(res, c.Device)
		case "Refdes":
			res = append(res, c.GetRefdesAsString())
		case "Footprint":
			res = append(res, c.Footprint)
		case "Value":
			res = append(res, c.Value)
		case "Count":
			res = append(res, strconv.Itoa(len(c.Refdes)))
		}
	}

	return res
}

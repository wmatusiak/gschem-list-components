package main

type validColumnNames map[string]bool

func (this validColumnNames) String() string {
	var result string
	for k := range this {
		if result != "" {
			result += ", "
		}

		result += k
	}

	return result
}

func (this validColumnNames) IsValid(name string) bool {
	return this[name]
}

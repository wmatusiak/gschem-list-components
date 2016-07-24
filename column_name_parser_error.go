package main

type ColumnNameParserError struct {
	message string
}

func (this ColumnNameParserError) Error() string {
	return this.message
}

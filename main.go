package main

import "bufio"
import "flag"
import "fmt"
import "github.com/olekukonko/tablewriter"
import "io"
import "log"
import "os"
import "sort"
import "strconv"
import "strings"
import "sync"

// CommandLineArgs type used to parse command line arguments
type CommandLineArgs struct {
	InFiles     []string
	SortBy      SortByValue
	ReverseSort bool
	Merge       bool
}

func (this *CommandLineArgs) Init() {
	flag.BoolVar(&this.Merge, "merge", false, "merge same components to single output line with count added")
	// SortBy
	this.SortBy = SortByValue{
		ValidColumnNames: map[string]bool{
			"Device":    true,
			"Value":     true,
			"Footprint": true,
			"Refdes":    true,
			"Count":     true,
		},
	}

	this.SortBy.Set("Device")
	flag.Var(&this.SortBy, "sort-by", "column name to sort (default: Device)")
	//END of SortBy

	flag.BoolVar(&this.ReverseSort, "reverse", false, "Revers sort")
	flag.Parse()
	this.InFiles = flag.Args()
}

// END of CommandLineArgs

// Component type containging atributes of component from schematic
type Component struct {
	Refdes    []string
	Device    string
	Footprint string
	Value     string
}

func (out *Component) Parse(in []string) {
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
}

func (c *Component) AddRefdes(refdes string) {
	c.Refdes = append(c.Refdes, refdes)
}

func (c Component) GetRefdesAsString() string {
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

func (c Component) GetFieldsValues(fieldNames []string) []string {
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

// END of Component

// ComponentArray array of componenets
type ComponentArray []Component

func (a ComponentArray) Len() int {
	return len(a)
}

func (a ComponentArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a *ComponentArray) Remove(i int) {
}

func (a ComponentArray) Print(columns []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)
	for _, comp := range a {
		table.Append(comp.GetFieldsValues(columns))
	}

	table.Render()
}

// END of ComponentArray

// ComponentSortByDevice type used to sort ComponentArray by Device column
type ComponentSortByDevice struct {
	ComponentArray
}

func (a ComponentSortByDevice) Less(i, j int) bool {
	return a.ComponentArray[i].Device < a.ComponentArray[j].Device
}

// END of ComponentSortByDevice

// ComponentSortByValue type used to sort ComponentArray by Value column
type ComponentSortByValue struct {
	ComponentArray
}

func (a ComponentSortByValue) Less(i, j int) bool {
	return a.ComponentArray[i].Value < a.ComponentArray[j].Value
}

// END of ComponentSortByValue

// ComponentSortByFootprint type used to sort ComponentArray by Footprint column
type ComponentSortByFootprint struct {
	ComponentArray
}

func (a ComponentSortByFootprint) Less(i, j int) bool {
	return a.ComponentArray[i].Footprint < a.ComponentArray[j].Footprint
}

// END of ComponentSortByFootprint

// ComponentSortByRefdes type used to sort ComponentArray by Refdes column
type ComponentSortByRefdes struct {
	ComponentArray
}

func (a ComponentSortByRefdes) Less(i, j int) bool {
	return a.ComponentArray[i].GetRefdesAsString() < a.ComponentArray[j].GetRefdesAsString()
}

// END of ComponentSortByRefdes

// ComponentSortByCount type used to sort ComponentArray by Count column
type ComponentSortByCount struct {
	ComponentArray
}

func (a ComponentSortByCount) Less(i, j int) bool {
	return len(a.ComponentArray[i].Refdes) < len(a.ComponentArray[j].Refdes)
}

// END of ComponentSortByCount

// ValidColumnNames store valid colument names
type ValidColumnNames map[string]bool

func (this ValidColumnNames) String() string {
	var result string
	for k := range this {
		if result != "" {
			result += ", "
		}

		result += k
	}

	return result
}

func (this ValidColumnNames) IsValid(name string) bool {
	return this[name]
}

// END of ValidColumnNames

// ColumnNameParserError type used to indicate column name parser error
type ColumnNameParserError struct {
	message string
}

func (this ColumnNameParserError) Error() string {
	return this.message
}

// End of ColumnNameParserError type

// SortByValue - value used by flag package to set name of field to sort
type SortByValue struct {
	Name             string
	ValidColumnNames ValidColumnNames
}

func (this *SortByValue) Set(name string) error {
	if !this.ValidColumnNames.IsValid(name) {
		this.Name = ""
		return ColumnNameParserError{
			message: fmt.Sprintf("%s is not valid column name. Valid names are: %s", name, this.ValidColumnNames),
		}
	}

	this.Name = name
	return nil
}

func (this SortByValue) String() string {
	return this.Name
}

// END of SortByValue

func ReadAllComponents(in io.ReadCloser, out chan Component, wait *sync.WaitGroup) {
	scanner := bufio.NewScanner(bufio.NewReader(in))
	var lines []string
	var inComponent bool
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == 'C' {
			lines = make([]string, 0, 10)
			inComponent = true
			lines = append(lines, line)
		} else if inComponent {
			lines = append(lines, line)
			if line[0] == '}' {
				inComponent = false
				comp := Component{}
				comp.Parse(lines)
				out <- comp
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	in.Close()
	wait.Done()
}

func StartParserOnIOReaders(in []io.ReadCloser) chan Component {
	wait := new(sync.WaitGroup)
	out := make(chan Component)
	for _, i := range in {
		wait.Add(1)
		go ReadAllComponents(i, out, wait)
	}

	go func() {
		wait.Wait()
		close(out)
	}()

	return out
}

func ParseFiles(inFileNames []string) ComponentArray {
	inReaders := make([]io.ReadCloser, 0, len(inFileNames))
	if len(inFileNames) == 0 {
		inReaders = append(inReaders, os.Stdin)
	} else {
		for _, fileName := range inFileNames {
			inFile, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			} else {
				inReaders = append(inReaders, inFile)
			}
		}
	}

	outChan := StartParserOnIOReaders(inReaders)
	components := make([]Component, 0, 100)
	for c := range outChan {
		components = append(components, c)
	}

	return components
}

func main() {
	args := new(CommandLineArgs)
	args.Init()

	// Read components from all files
	components := ParseFiles(args.InFiles)

	// Merge if needed
	if args.Merge {
		tmp := make(map[string]Component)
		for _, c := range components {
			key := c.Device + "_" + c.Footprint + "_" + c.Value
			if _, exists := tmp[key]; exists {
				cmp := tmp[key]
				cmp.AddRefdes(c.Refdes[0])
				tmp[key] = cmp
			} else {
				tmp[key] = c
			}
		}

		components = make([]Component, 0, len(tmp))
		for _, c := range tmp {
			components = append(components, c)
		}
	}

	// Sort
	var s sort.Interface
	switch args.SortBy.String() {
	case "Device":
		s = ComponentSortByDevice{components}
	case "Value":
		s = ComponentSortByValue{components}
	case "Footprint":
		s = ComponentSortByFootprint{components}
	case "Refdes":
		s = ComponentSortByRefdes{components}
	case "Count":
		s = ComponentSortByCount{components}
	}

	if args.ReverseSort {
		s = sort.Reverse(s)
	}

	sort.Sort(s)

	// Display
	if args.Merge {
		components.Print([]string{"Device", "Value", "Footprint", "Refdes", "Count"})
	} else {
		components.Print([]string{"Device", "Value", "Footprint", "Refdes"})
	}
}

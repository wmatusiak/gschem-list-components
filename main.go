package main

import "bufio"
import "flag"
import "github.com/olekukonko/tablewriter"
import "io"
import "log"
import "os"
import "strconv"
import "strings"
import "sync"

// CommandLineArgs type used to parse command line arguments
type CommandLineArgs struct {
	InFiles []string
	Merge   bool
}

func (this *CommandLineArgs) Init() {
	flag.BoolVar(&this.Merge, "merge", false, "merge same components to single output line with count added")
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

// END of Component

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

func ParseFiles(inFileNames []string) []Component {
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

	// Display
	table := tablewriter.NewWriter(os.Stdout)
	if args.Merge {
		table.SetHeader([]string{"Device", "Value", "Footprint", "Refdes", "Count"})
	} else {
		table.SetHeader([]string{"Device", "Value", "Footprint", "Refdes"})
	}

	for _, c := range components {
		refdes := ""
		for _, r := range c.Refdes {
			if refdes != "" {
				refdes += ", "
			}

			refdes += r
		}

		if args.Merge {
			table.Append([]string{c.Device, c.Value, c.Footprint, refdes, strconv.Itoa(len(c.Refdes))})
		} else {
			table.Append([]string{c.Device, c.Value, c.Footprint, refdes})
		}
	}

	table.Render()
}

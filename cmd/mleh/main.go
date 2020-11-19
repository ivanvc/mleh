package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
)

type TemplateInput struct {
	Values map[interface{}]interface{}
}

type valuesFlag []string

func (v valuesFlag) String() string {
	return strings.Join(v, ", ")
}

func (v *valuesFlag) Set(value string) error {
	*v = append(*v, value)
	return nil
}

var (
	valuesFile = flag.String("values", "", "The values.yaml file")
	outputDir  = flag.String("output-dir", "", "The output directory")
	dryMode    = flag.Bool("dry", false, "Run in dry mode")
	values     valuesFlag
)

func init() {
	flag.Var(&values, "value", "Options to be used as values, as key=value where value is in YAML format")
	flag.Parse()
}

func main() {
	if len(flag.Args()) < 1 {
		log.Fatal("Input directory not specified")
	}
	inputDir := flag.Arg(0) + "/"
	ti := new(TemplateInput)
	dat, err := ioutil.ReadFile(filepath.Join(inputDir, "values.yaml"))
	if err != nil {
		log.Fatal("Error reading chart values.yaml: ", err.Error())
	}
	if err = yaml.Unmarshal(dat, &ti.Values); err != nil {
		log.Fatal("Error parsing chart values.yaml: ", err.Error())
	}

	if len(*valuesFile) > 0 {
		dat, err = ioutil.ReadFile(*valuesFile)
		if err != nil {
			log.Fatal("Error reading input values yaml file: ", err.Error())
		}
		if err = yaml.Unmarshal(dat, &ti.Values); err != nil {
			log.Fatal("Error parsing input values yaml file: ", err.Error())
		}
	}

	for _, v := range values {
		var value interface{}
		key := strings.Split(v, "=")[0]
		v = strings.TrimPrefix(v, key+"=")
		if err = yaml.Unmarshal([]byte(v), &value); err != nil {
			log.Fatal("Error parsing ", key, " with value ", v, ": ", err.Error())
		}
		ti.Values[key] = value
	}

	tpl := template.Must(template.New("base").Funcs(sprig.TxtFuncMap()).ParseGlob(filepath.Join(inputDir, "templates/*")))

	if *dryMode {
		log.Println("Dry mode, skipping creating directory ", *outputDir)
	} else {
		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			log.Fatal("Error creating output diirectory: ", err.Error())
		}
	}

	for _, t := range tpl.Templates() {
		if _, err := os.Stat(filepath.Join(inputDir, "templates", t.Name())); os.IsNotExist(err) || strings.HasPrefix(t.Name(), "_") {
			continue
		}

		fn := filepath.Join(*outputDir+"/", t.Name())
		if *dryMode {
			fn = os.DevNull
		}
		f, err := os.Create(fn)
		if err != nil {
			log.Fatal("Error creating file ", fn, ": ", err.Error())
		}
		defer f.Close()

		if err := t.Execute(f, ti); err != nil {
			log.Fatal("Error executing template: ", err.Error())
		}
	}
}

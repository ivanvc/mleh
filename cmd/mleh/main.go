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

var (
	values    = flag.String("values", "", "The values.yaml file")
	outputDir = flag.String("output-dir", "", "The output directory")
)

func init() {
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
	yaml.Unmarshal(dat, &ti.Values)

	dat, err = ioutil.ReadFile(*values)
	if err != nil {
		log.Fatal("Error reading chart values.yaml: ", err.Error())
	}
	yaml.Unmarshal(dat, &ti.Values)

	tpl := template.Must(template.New("base").Funcs(sprig.TxtFuncMap()).ParseGlob(filepath.Join(inputDir, "templates/[^.]*")))

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatal("Error creating output diirectory: ", err.Error())
	}

	for _, t := range tpl.Templates() {
		if _, err := os.Stat(filepath.Join(inputDir, "templates", t.Name())); os.IsNotExist(err) || strings.HasPrefix(t.Name(), "_") {
			continue
		}

		fn := filepath.Join(*outputDir+"/", t.Name())
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

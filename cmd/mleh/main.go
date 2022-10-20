package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	log "github.com/sirupsen/logrus"
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
	silent     = flag.Bool("silent", false, "Prints only Error level logs or higher")
	values     valuesFlag
)

func init() {
	flag.Var(&values, "value", "Options to be used as values, as key=value where value is in YAML format")
	flag.Parse()

	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	if *silent {
		log.SetLevel(log.ErrorLevel)
	}
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

	t := template.New("base")
	funcMap := map[string]interface{}{
		"include": func(name string, data interface{}) (string, error) {
			buf := bytes.NewBuffer(nil)
			if err := t.ExecuteTemplate(buf, name, data); err != nil {
				return "", err
			}
			return buf.String(), nil
		},
	}

	tpl := template.Must(t.Funcs(sprig.TxtFuncMap()).Funcs(funcMap).ParseGlob(filepath.Join(inputDir, "templates/*")))

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

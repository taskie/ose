package main

//go:generate statik -f -src templates

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rakyll/statik/fs"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	_ "github.com/taskie/ose/cmd/colivec/statik"
)

type templateData struct {
	CmdNames []string
}

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	pout := pflag.StringP("out", "o", ".", "output directory")
	pforce := pflag.BoolP("force", "f", false, "force mode")
	pflag.Parse()
	cmdNames := pflag.Args()

	staticFS, err := fs.New()
	if err != nil {
		log.Fatal("can't load fs", zap.Error(err))
	}
	m := map[string]string{
		"_gitignore":  ".gitignore",
		"GNUmakefile": "",
	}
	tmpl := template.New("")
	for k := range m {
		f, err := staticFS.Open("/" + k)
		if err != nil {
			log.Fatal("can't open", zap.Error(err))
		}
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("can't read", zap.Error(err))
		}
		_, err = tmpl.New(k).Parse(string(bs))
		if err != nil {
			log.Fatal("can't parse", zap.Error(err))
		}
	}
	for k, v := range m {
		filename := v
		if filename == "" {
			filename = k
		}
		outpath := filepath.Join(*pout, filename)
		if !*pforce {
			_, err = os.Stat(outpath)
			if err == nil {
				log.Info("exists", zap.String("outpath", outpath))
				continue
			}
		}
		of, err := os.Create(outpath)
		if err != nil {
			log.Fatal("can't create", zap.Error(err))
		}
		defer of.Close()
		err = tmpl.ExecuteTemplate(of, k, &templateData{
			CmdNames: cmdNames,
		})
		if err != nil {
			log.Fatal("can't execute", zap.Error(err))
		}
	}
}

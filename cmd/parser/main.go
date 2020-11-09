package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/surmus/iso-country-parser/internal"
	"github.com/surmus/iso-country-parser/internal/parser"
	"github.com/surmus/iso-country-parser/internal/wiki"
	"github.com/surmus/iso-country-parser/internal/writer"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"runtime/debug"
)

const (
	version                    = "v1.0.2"
	countriesWikiPageIDFlag    = "page-id"
	defaultCountriesWikiPageID = "List_of_ISO_3166_country_codes"
	resultTemplateFlag         = "template"
	outputFilePathFlag         = "file"
	verboseFlag                = "verbose"
	wikiURL                    = "https://en.wikipedia.org"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:  countriesWikiPageIDFlag,
		Value: defaultCountriesWikiPageID,
		Usage: "Wikipedia page ID from " + wikiURL,
	},
	&cli.StringFlag{
		Name: resultTemplateFlag,
		Usage: `Template string for formatting the result, example:
				'({CODE}, {NAME}) ' will produce:
				'(USA, The United States of America) (EST, Estonia)'`,
	},
	&cli.StringFlag{
		Name:  outputFilePathFlag,
		Usage: fmt.Sprintf(`filepath the results are written into, example: --%s C:\hello.txt`, outputFilePathFlag),
	},
	&cli.BoolFlag{
		Name:  verboseFlag,
		Usage: "Debug application",
	},
}

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Name = "ISO 3166 countries Wikipedia page parser"
	app.Usage = "Parses ISO 3166 alpha-3 country codes from WIKI page into structured output"
	app.Flags = flags
	app.Action = execute

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func execute(c *cli.Context) {
	logrus.SetFormatter(&internal.CLIFormatter{})
	var output io.Writer = os.Stdout
	var resultWriter writer.ResultWriter
	debugApplication := c.Bool(verboseFlag)
	writeToFile := c.String(outputFilePathFlag)

	defer errorHandler(debugApplication, writeToFile)()

	if debugApplication {
		logrus.SetLevel(logrus.TraceLevel)
	}

	if writeToFile != "" {
		file := mustCreateFile(writeToFile)
		output = file
		defer closeOutputFile(file)()
	}

	resultWriter = setupResultWriter(c.String(resultTemplateFlag), resultWriter, output)
	page, e := wiki.NewClient(wikiURL).GetPage(c.String(countriesWikiPageIDFlag))

	if e != nil {
		panic(e)
	}

	countries := parser.NewWikiPageParser(page).Parse()
	logrus.Info("Writing result to the output")
	if e = resultWriter.Write(countries); e != nil {
		panic(e)
	}
}

func setupResultWriter(resultTemplate string, resultWriter writer.ResultWriter, output io.Writer) writer.ResultWriter {
	if resultTemplate != "" {
		resultWriter = mustCreateResultWriter(resultTemplate, output)
	} else {
		resultWriter = writer.NewJSONResultWriter(output)
	}

	return resultWriter
}

func mustCreateResultWriter(resultTemplate string, output io.Writer) writer.ResultWriter {
	resWriter, e := writer.NewTemplateResultWriter(output, resultTemplate)

	if e != nil {
		panic(e)
	}

	return resWriter
}

func closeOutputFile(file *os.File) func() {
	return func() {
		if e := file.Close(); e != nil {
			panic(e)
		}
	}
}

func mustCreateFile(filePath string) *os.File {
	file, e := os.Create(filePath)

	if e != nil {
		panic(fmt.Errorf("failed to create output file at: %s %w", filePath, e))
	}

	return file
}

func errorHandler(debugApplication bool, outputFilePath string) func() {
	return func() {
		if r := recover(); r != nil {
			if outputFilePath != "" {
				_ = os.Remove(outputFilePath)
			}

			if debugApplication {
				logrus.Fatalf("ERROR:\n%v\n%s", r, debug.Stack())
			} else {
				logrus.Fatalf("ERROR:\n%v", r)
			}
		}
	}
}

# iso-country-parser
[![CircleCI](https://circleci.com/gh/Surmus/iso-country-parser.svg?style=svg)](https://circleci.com/gh/Surmus/iso-country-parser)
[![codecov](https://codecov.io/gh/Surmus/iso-country-parser/branch/master/graph/badge.svg)](https://codecov.io/gh/Surmus/iso-country-parser)
[![Go Report Card](https://goreportcard.com/badge/github.com/surmus/iso-country-parser)](https://goreportcard.com/report/github.com/surmus/iso-country-parser)
[![Release](https://img.shields.io/github/release/surmus/iso-country-parser.svg?style=flat-square)](https://github.com/surmus/iso-country-parser/releases/latest)

Parses ISO 3166 alpha-3 country codes from WIKI page into structured output

## Usage
### Download
Choose one of the following options:

#### Download Github release from https://github.com/Surmus/iso-country-parser/releases
##### When running Windows
1. Extract win64 folder contents from downloaded release.tar.gz
2. Run application:
     ```sh
     parser.exe
     ```
##### When running Linux
1. Extract linux64 folder contents from downloaded release.tar.gz
2. Run application:
     ```sh
     ./parser
     ```
#### Compile and install from source code
1. Install Golang https://golang.org/
2. Download source and compile binaries:
    ```sh
    $ go get -u github.com/surmus/iso-country-parser/cmd/parser
    ```
3.  Run application (NB: check that GOBIN env variable is set and added to the PATH)
    Linux
    ```sh
    $ parser
    ```
   
## CLI options
```sh
$ ./parser help
NAME:
   ISO 3166 countries Wikipedia page parser - Parses ISO 3166 alpha-3 country codes from WIKI page into structured output

USAGE:
   parser [global options] command [command options] [arguments...]

VERSION:
   v1.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --page-id value   Wikipedia page ID from https://en.wikipedia.org (default: "List_of_ISO_3166_country_codes")
   --template value  Template string for formatting the result, example:
                           '({CODE}, {NAME}) ' will produce:
                           '(USA, The United States of America) (EST, Estonia)'
   --file value      filepath the results are written into, example: --file C:\hello.txt
   --verbose         Debug application
   --help, -h        show help
   --version, -v     print the version
```

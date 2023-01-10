package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type flags struct {
	stdin   bool
	verbose bool
	inFile  string
}

var fg flags
var logger *zerolog.Logger
var headers []string

func init() {
	flag.BoolVar(&fg.verbose, "v", false, "prints debug logs")
	flag.StringVar(&fg.inFile, "f", "", "input file path")
	flag.BoolVar(&fg.stdin, "i", false, "read from standard input instead of file")
	flag.Parse()
}

func main() {

	logger = GetLogger(fg.verbose)
	input := getInputFile()
	headers = getHeaders(input)
	buf := bufio.NewWriter(os.Stdout)

	for {
		record, err := input.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			logger.Fatal().Err(err).Msg("could not read CSV")
		}

		encoded := convertRowToJSON(record)
		buf.Write(encoded)
	}

	buf.Flush()

}

func convertRowToJSON(record []string) []byte {

	m := make(map[string]string, len(record))

	for i := 0; i < len(record); i++ {
		m[headers[i]] = record[i]
	}

	enc, err := json.Marshal(m)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not encode to JSON")
	}

	return enc

}

func getHeaders(input *csv.Reader) []string {
	record, err := input.Read()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not read csv headers or file is empty")
	}

	return record
}

func getInputFile() *csv.Reader {

	if fg.stdin {
		return csv.NewReader(os.Stdin)
	}

	if fg.inFile == "" {
		logger.Fatal().Msg("please pass input file or use -i for reading from stdin")
	}

	fh, err := os.Open(fg.inFile)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not open input file")
	}

	return csv.NewReader(fh)

}

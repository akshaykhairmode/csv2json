package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type flags struct {
	stdin   bool
	verbose bool
	help    bool
	inFile  string
	jsonCol string
	keepCol string
}

var fg flags
var logger *zerolog.Logger
var headers []string
var jsonCols = map[uint]struct{}{}
var keepCols = map[uint]struct{}{}

func init() {
	flag.StringVar(&fg.jsonCol, "j", "", "mention the columns that contain json to convert them to objects instead of treating them as string.\nColumn starts from 0")
	flag.StringVar(&fg.keepCol, "k", "", "mention the columns that should be considered and skip others.\nColumn starts from 0")
	flag.BoolVar(&fg.verbose, "v", false, "prints debug logs")
	flag.StringVar(&fg.inFile, "f", "", "input file path")
	flag.BoolVar(&fg.stdin, "i", false, "read from standard input instead of file")
	flag.BoolVar(&fg.help, "h", false, "print help")
	flag.Parse()

	if fg.help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	logger = GetLogger(fg.verbose)

	parseUintFlags(fg.jsonCol, jsonCols)
	parseUintFlags(fg.keepCol, keepCols)
}

func parseUintFlags(val string, m map[uint]struct{}) {
	if val != "" {
		vals := strings.Split(val, ",")
		for _, v := range vals {
			uintVal, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				logger.Fatal().Msg("jsonCol should be a number")
			}
			m[uint(uintVal)] = struct{}{}
		}
	}
}

func main() {

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

	m := mapAnyPool.Get().(map[string]any)

	for i := 0; i < len(record); i++ {
		m[headers[i]] = getValue(record[i], uint(i))
	}

	enc, err := json.Marshal(m)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not encode to JSON")
	}

	mapAnyPool.Put(m)

	return enc

}

func getValue(rec string, i uint) any {

	if _, ok := keepCols[i]; !ok && len(keepCols) > 0 {
		return ""
	}

	if rec == "" {
		return rec
	}

	if _, ok := jsonCols[i]; !ok {
		return rec
	}

	temp := anyPool.Get()
	defer anyPool.Put(temp)

	if err := json.Unmarshal([]byte(rec), &temp); err != nil {
		logger.Fatal().Err(err).Str("record", rec).Msg("invalid json object")
	}

	return temp

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

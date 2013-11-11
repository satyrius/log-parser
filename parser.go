package main

import (
	"./stat"
	"encoding/json"
	"fmt"
	"github.com/droundy/goopt"
	"github.com/satyrius/gonx"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

var debug *bool
var format *string
var nginxConfig *string
var nginxFormat *string
var aggField *string
var groupBy *string
var groupByReqexp *string
var jsonOutput *string

func init() {
	format = goopt.String([]string{"--fmt", "--format"}, "",
		"Log format (e.g. '$remote_addr [$time_local] \"$request\"')")
	nginxConfig = goopt.String([]string{"--nginx"}, "",
		"Nginx config to look for 'log_format' directive. You also should specify --nginx-format")
	nginxFormat = goopt.String([]string{"--nginx-format"}, "",
		"Name of nginx 'log_format', should be passed with --nginx option")
	aggField = goopt.String([]string{"-a", "--aggregate"}, "request_time",
		"Nginx access log variable to aggregate")
	groupBy = goopt.String([]string{"-g", "--group-by"}, "request",
		"Nginx access log variable to group by")
	groupByReqexp = goopt.String([]string{"-r", "--regexp"}, "",
		"You can specify regular expression to extract exact data from group by data. "+
			"For example, you might want to group by a path inside $request, so you should "+
			"set this option to '^\\S+\\s(.*)(?:\\?.*)?$'.")
	debug = goopt.Flag([]string{"--debug"}, []string{"--no-debug"},
		"Log debug information", "Do not log debug information")
	jsonOutput = goopt.String([]string{"-o", "--json"}, "",
		"Save result as json encoded file")
}

func badUsage() {
	fmt.Println(goopt.Usage())
	os.Exit(1)
}

func getReader(file io.Reader) (reader *gonx.Reader, err error) {
	if *format != "" {
		reader = gonx.NewReader(file, *format)
	} else {
		if *nginxConfig == "" || *nginxFormat == "" {
			badUsage()
		}
		cfg, err := os.Open(*nginxConfig)
		if err != nil {
			return nil, err
		}
		defer cfg.Close()
		reader, err = gonx.NewNginxReader(file, cfg, *nginxFormat)
	}
	return
}

func main() {
	goopt.Parse(nil)

	var logs []*os.File
	if len(goopt.Args) != 0 {
		fmt.Println("We are going to parse those files:")
		for _, log := range goopt.Args {
			log, _ := filepath.Abs(log)
			fmt.Println(log)
			file, err := os.Open(log)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			logs = append(logs, file)
		}
	} else {
		fmt.Println("Read from STDIN")
		logs = append(logs, os.Stdin)
	}

	aggregator := func(i *stat.Item, entry *gonx.Entry) (val float64, err error) {
		if strVal, ok := (*entry)[*aggField]; ok {
			v, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return 0, err
			}
			if i.Count == 1 {
				val = v
			} else {
				val = (i.AggValue*float64(i.Count-1) + v) / float64(i.Count)
			}
		} else {
			err = fmt.Errorf("Invalid entry data")
		}
		return
	}
	st := stat.NewStat(aggregator, *groupBy, *groupByReqexp)

	for _, file := range logs {
		st.AddLog(file.Name())
		reader, err := getReader(file)
		if err != nil {
			panic(err)
		}
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			if *debug {
				fmt.Println("==============================")
				for name, value := range record {
					fmt.Printf("LR $%v = '%v'\n", name, value)
				}
			}
			if err := st.Add(&record); err != nil {
				panic(err)
			}
		}
	}
	fmt.Printf("Gratz! You've parsed %v log entries, it took %v\n", st.EntriesParsed, st.Stop())
	sort.Sort(st)
	for _, item := range st.Data {
		fmt.Printf("%7.3f %6d %v\n", item.AggValue, item.Count, item.Name)
	}

	if *jsonOutput != "" {
		jsonFile, err := os.OpenFile(*jsonOutput, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()
		jsData, err := json.MarshalIndent(st.Data, "", "  ")
		if err != nil {
			panic(err)
		}
		jsonFile.Write(jsData)
		fmt.Printf("Result was saved to JSON encoded file '%v'\n", jsonFile.Name())
	}
}

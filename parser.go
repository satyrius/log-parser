package main

import (
	"./stat"
	"fmt"
	"github.com/droundy/goopt"
	"github.com/satyrius/gonx"
	"io"
	"os"
	"path/filepath"
)

var debug *bool
var format *string
var nginxConfig *string
var nginxFormat *string
var groupBy *string
var groupByReqexp *string

func init() {
	format = goopt.String([]string{"--fmt", "--format"}, "",
		"Log format (e.g. '$remote_addr [$time_local] \"$request\"')")
	nginxConfig = goopt.String([]string{"--nginx"}, "",
		"Nginx config to look for 'log_format' directive. You also should specify --nginx-format")
	nginxFormat = goopt.String([]string{"--nginx-format"}, "",
		"Name of nginx 'log_format', should be passed with --nginx option")
	groupBy = goopt.String([]string{"-g", "--group-by"}, "request",
		"Nginx access log variable to group by")
	groupByReqexp = goopt.String([]string{"-r", "--regexp"}, "",
		"You can specify regular expression to extract exact data from group by data. "+
			"For example, you might want to group by a path inside $request, so you should "+
			"set this option to '^\\S+\\s(.*)(?:\\?.*)?$'.")
	debug = goopt.Flag([]string{"--debug"}, []string{"--no-debug"},
		"Log debug information", "Do not log debug information")
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

	st := stat.NewStat(*groupBy, *groupByReqexp)
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
	for req, count := range st.Data {
		fmt.Printf("%v %v\n", count, req)
	}
}

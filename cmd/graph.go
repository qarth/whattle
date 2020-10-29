package cmd

import (
	"fmt"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/qarth/whattle/optimization"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "writes the graph of a block model using specified slope",
	Long:  `Graph outputs the DIMACs format graph for the maxflow/mincut problem. upit is a CLI program to optimise an ore reserves ultimate pit limit during feasibility study phase. This application is a tool to generate the needed precedence files to find the ultimate pit.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGraph(cmd, args)
		fmt.Println("graph called")
	},
}

func init() {
	fmt.Println("graph init")
}

func runGraph(cmd *cobra.Command, args []string) {
	viper.BindPFlags(cmd.Flags())

	logfile := viper.GetString("log")
	infile := viper.GetString("input")
	outfile := viper.GetString("output")
	jsonFile := viper.GetString("params")
	outputDest := "<console/>"

	if len(infile) == 0 || len(outfile) == 0 || len(jsonFile) == 0 {
		cmd.Usage()
		return
	}

	if len(logfile) > 0 {
		outputDest = fmt.Sprintf(log_file_tmpl, logfile)
	}
	log_cfg := strings.Replace(log_cfg_tmpl, log_out_dest, outputDest, -1)
	logger, _ := log.LoggerFromConfigAsString(log_cfg)

	if logger != nil {
		log.ReplaceLogger(logger)
	}

	param := optimization.RunCtx{
		InputFile:  infile,
		OutputFile: outfile,
		ParamFile:  jsonFile,
	}
	optimization.StartRead(param)
}

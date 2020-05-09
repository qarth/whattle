package cmd

import (
	"fmt"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/qarth/whattle/optimization"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	log_cfg_tmpl = `<seelog minlevel="info">
		<outputs formatid="detail">
			{{OutputDest}}
		</outputs>
		<formats>
			<format id="detail" format="[%File:%Line][%Date(2006-01-02 15:04:05.000)] %Msg%n" />
		</formats>
	</seelog>`
	log_out_dest  = "{{OutputDest}}"
	log_file_tmpl = `<rollingfile filename="%s" type="size" maxsize="10247680" maxrolls="10"/>`
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run opt",
	Long:  "run a mining optimization task",
	Run: func(cmd *cobra.Command, args []string) {
		runOpt(cmd, args)
	},
}

func init() {

	RootCmd.AddCommand(runCmd)

	flagset := runCmd.PersistentFlags()
	flagset.StringP("input", "i", "", "The input file")
	flagset.StringP("output", "o", "", "The output file")
	flagset.StringP("log", "l", "", "Log information to a file")
	flagset.StringP("params", "p", "", "Grid parameter json file")
}

func runOpt(cmd *cobra.Command, args []string) {

	viper.BindPFlags(cmd.Flags())

	logfile := viper.GetString("log")
	infile := viper.GetString("input")
	outfile := viper.GetString("output")
	jsonFile := viper.GetString("params")

	if len(infile) == 0 || len(outfile) == 0 || len(jsonFile) == 0 {
		cmd.Usage()
		return
	}

	//-------

	outputDest := "<console/>"

	if len(logfile) > 0 {
		outputDest = fmt.Sprintf(log_file_tmpl, logfile)
	}

	log_cfg := strings.Replace(log_cfg_tmpl, log_out_dest, outputDest, -1)

	logger, _ := log.LoggerFromConfigAsString(log_cfg)

	if logger != nil {
		log.ReplaceLogger(logger)
	}

	//-------

	param := optimization.RunCtx{
		InputFile:  infile,
		OutputFile: outfile,
		ParamFile:  jsonFile,
	}

	log.Info("miningopt begin")

	optimization.StartRead(param)

	log.Info("miningopt finished")
	log.Flush()
}

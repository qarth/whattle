package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	default_params = `
{
// input
//   1 (GEOEAS (GSLIB) format grid file. Pre-calculated EBV)
//     grid (The grid definition)
//       min_x, min_y, min_z (The lower left centroid)
//       num_x, num_y, num_z (The number of blocks)
//       siz_x, siz_y, siz_z (The size of a block)
//     ebv_column (Economic block value column, 1 indexed)
//   2 (GZIP .gz file, only ebv, one column, no header)
//     grid (as above)
\"input\" : {
  \"type\" : 1,

  \"grid\" : {
    \"num_x\": 60, \"min_x\": 810.0, \"siz_x\": 20.0,
    \"num_y\": 60, \"min_y\": 110.0, \"siz_y\": 20.0,
    \"num_z\": 13, \"min_z\": 110.0, \"siz_z\": 20.0
  },
  \"ebv_column\" : 1
},

// precedence
//   1 (Benches)
//     slope (The slope (in degrees))
//     benches (The number of benches)
\"precedence\" : {
  \"method\" : 1,

  \"slope\" : 45.0,
  \"num_benches\": 8
},

// optimization_engine
//   1 (Lerchs Grossmann)
//   2 (Dimacs program)
//     dimacs_path (Path to engine)
\"optimization\" : {
  \"engine\" : 1
}
}`
)

var paramsCmd = &cobra.Command{
	Use:   "params",
	Short: "Output the default params",
	Long:  "Output the default parameters",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(default_params)
	},
}

func init() {
	RootCmd.AddCommand(paramsCmd)
}

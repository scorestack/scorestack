package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const saveShort = "Write the current configuration to a file."
const saveLong = saveShort + `

The configuration format will be automatically determined by the extension of
the destination file specified. If the file already exists, it will be
overwritten.`

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save [output file]",
	Short: saveShort,
	Long:  saveLong,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.WriteConfigAs(args[0])
		if err != nil {
			fmt.Printf("The following extensions are supported:\n%+q\n", viper.SupportedExts)
		}
		cobra.CheckErr(err)
	},
}

func init() {
	configCmd.AddCommand(saveCmd)
}

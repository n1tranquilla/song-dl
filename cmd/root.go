package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "song-dl",
	Short: "song-dl is a song downloader written in go",
	Long: `A command line application for downloading songs from youtube. 
	For complete documentation, see https://www.github.com/n1tranquilla/song-dl`,
	Run: func(cmd *cobra.Command, args []string) {
	  // Do Stuff Here
	},
}
  
func Execute() {
	//rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	
	if err := rootCmd.Execute(); err != nil {
	  fmt.Println(err)
	  os.Exit(1)
	}
}
/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/openwurl/hw-pragma-scan/pkg/hwscan"
	"github.com/spf13/cobra"
)

var (
	scanTargetHost string
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans the target host (simple GET)",
	Long:  `Scans the target host (simple GET).`,
	Run: func(cmd *cobra.Command, args []string) {
		hw := &hwscan.Scanner{
			Target: scanTargetHost,
		}
		err := hw.Scan()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("complete")
			hw.Report()
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&scanTargetHost, "url", "u", scanTargetHost, "The full target URL (domain and uri) website.com/path/to/file.suffix")
	scanCmd.MarkFlagRequired("url")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

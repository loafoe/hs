/*
Copyright © 2021 Andy Lo-A-Foe <andy.lo-a-foe@philips.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cheynewallace/tabby"

	"github.com/spf13/cobra"
)

// servicesRegionsCmd represents the services command
var servicesRegionsCmd = &cobra.Command{
	Use:     "regions",
	Aliases: []string{"r"},
	Short:   "Retrieve region information",
	Long:    `Retrieves region information.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getConfig(cmd)
		if err != nil {
			fmt.Printf("failed to get config: %v\n", err)
			return
		}
		regions := config.Regions()
		if len(regions) == 0 {
			if jsonOut {
				fmt.Printf("[]\n")
				return
			}
			fmt.Printf("no regions found\n")
			return
		}
		if jsonOut {
			data, _ := json.Marshal(regions)
			fmt.Printf("%s\n", data)
			return
		}
		t := tabby.New()
		t.AddHeader("regions")
		for _, r := range regions {
			t.AddLine(r)
		}
		t.Print()

		fmt.Printf("\n")
	},
}

func init() {
	servicesCmd.AddCommand(servicesRegionsCmd)
}

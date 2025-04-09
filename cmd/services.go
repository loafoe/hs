/*
Copyright © 2020 Andy Lo-A-Foe <andy.lo-a-foe@philips.com>

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

	"github.com/dip-software/go-dip-api/config"

	"github.com/spf13/cobra"
)

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:     "services",
	Aliases: []string{"s"},
	Short:   "Retrieve service information",
	Long:    `Retrieves service information.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getConfig(cmd)
		if err != nil {
			fmt.Printf("failed to get config: %v\n", err)
			return
		}
		services := config.Services()
		if len(services) == 0 {
			if jsonOut {
				fmt.Printf("[]\n")
				return
			}
			fmt.Printf("no services found\n")
			return
		}
		if jsonOut {
			data, _ := json.Marshal(services)
			fmt.Printf("%s\n", data)
			return
		}
		t := tabby.New()
		t.AddHeader("service", "host", "url", "domain")
		for _, i := range services {
			host := config.Service(i).Host
			url := config.Service(i).URL
			domain := config.Service(i).Domain
			t.AddLine(i,
				host,
				url,
				domain)
		}
		t.Print()

		fmt.Printf("\n")
	},
}

func getConfig(cmd *cobra.Command) (*config.Config, error) {
	region, _ := cmd.Flags().GetString("region")
	environment, _ := cmd.Flags().GetString("env")
	c, err := config.New(config.WithRegion(region), config.WithEnv(environment))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func init() {
	rootCmd.AddCommand(servicesCmd)
	servicesCmd.PersistentFlags().StringP("region", "r", "us-east", "Select region")
	servicesCmd.PersistentFlags().StringP("env", "e", "client-test", "Select environment")
}

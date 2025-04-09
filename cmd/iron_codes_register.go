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
	"fmt"
	"github.com/dip-software/go-dip-api/iron"
	"github.com/spf13/cobra"
	"strings"
)

// ironRegisterCmd represents the register command
var ironRegisterCmd = &cobra.Command{
	Use:     "register some/image[:tag]",
	Aliases: []string{"r"},
	Short:   "Register a docker image as an Iron code",
	Long:    `Registers a docker image as an Iron code`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		config, err := readIronConfig()
		if err != nil {
			fmt.Printf("error reading iron config: %v\n", err)
			return
		}
		if len(config.ClusterInfo) == 0 {
			fmt.Printf("missing required cluster info in config")
		}
		config.Debug = debug
		client, err := iron.NewClient(config)
		if err != nil {
			fmt.Printf("error configuring iron client: %v\n", err)
			return
		}
		name := strings.Split(args[0], ":")[0]
		code, resp, err := client.Codes.CreateOrUpdateCode(iron.Code{
			Name:  name,
			Image: args[0],
		})
		if err != nil {
			fmt.Printf("error registering code: %v\n", err)
			return
		}
		if code != nil && code.Name != "" {
			fmt.Printf("registered %s, revision %d\n", code.Name, code.Rev)
		} else {
			fmt.Printf("unexpected error registering code: %d\n", resp.StatusCode)
		}
		fmt.Printf("\n")
	},
}

func init() {
	codesCmd.AddCommand(ironRegisterCmd)
}

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

	"github.com/spf13/cobra"
)

// hasConfigCmd represents the config command
var hasConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure HAS",
	Long:  `Configure HAS settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		orgID, _ := cmd.Flags().GetString("orgid")
		if url == "" && orgID == "" {
			_ = cmd.Help()
			return
		}
		if url != "" {
			currentWorkspace.HASConfig.HASURL = url
		}
		if orgID != "" {
			currentWorkspace.HASConfig.OrgID = orgID
		}
		if err := currentWorkspace.save(); err == nil {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("failed to store config: %v\n", err)
		}
	},
}

func init() {
	hasCmd.AddCommand(hasConfigCmd)
}

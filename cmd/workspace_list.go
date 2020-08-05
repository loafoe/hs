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
	"os"

	"github.com/cheynewallace/tabby"

	"github.com/spf13/cobra"
)

// ironTasksListCmd represents the list command
var workspaceListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List available workspaces",
	Long:    `Lists available workspaces.`,
	Run: func(cmd *cobra.Command, args []string) {
		list, current, err := currentWorkspace.list()
		if err != nil {
			fmt.Printf("failed to list workspaces: %v\n", err)
			os.Exit(1)
		}
		t := tabby.New()
		t.AddHeader("workspace", "region", "environment")
		for _, w := range list {
			if config, err := loadWorkspaceConfig(w); err == nil {
				name := "  " + config.Name
				if w == current {
					name = "✓ " + config.Name
				}
				t.AddLine(name, config.DefaultRegion, config.DefaultEnvironment)
			}
		}
		t.Print()
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceListCmd)
}

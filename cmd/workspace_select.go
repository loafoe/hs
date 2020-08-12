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

// workspaceSelectCmd represents the set command
var workspaceSelectCmd = &cobra.Command{
	Use:     "select <workspace>",
	Aliases: []string{"s"},
	Short:   "Select a different workspace",
	Long:    `Selects a different workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		workspace := args[0]
		if err := currentWorkspace.setDefault(workspace); err != nil {
			fmt.Printf("failed to select workspace %s: %v\n", workspace, err)
			return
		}
		fmt.Printf("switched to workspace %s\n", workspace)
		fmt.Printf("\n")
		var err error
		currentWorkspace, err = loadWorkspaceConfig(workspace)
		if err != nil {
			fmt.Printf("failed to load workspace: %v\n", err)
			return
		}
		workspaceInfoCmd.Run(cmd, args)
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceSelectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workspaceSelectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workspaceSelectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

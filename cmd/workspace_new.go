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

// workspaceNewCmd represents the new command
var workspaceNewCmd = &cobra.Command{
	Use:     "new <workspace>",
	Aliases: []string{"n"},
	Short:   "Create a new workspace",
	Long:    `Creates a new workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		region, _ := cmd.Flags().GetString("region")
		environment, _ := cmd.Flags().GetString("environment")

		workspace := args[0]

		if currentWorkspace == nil {
			fmt.Printf("WTF\n")
			return
		}
		currentWorkspace.DefaultEnvironment = environment
		currentWorkspace.DefaultRegion = region
		currentWorkspace.Name = workspace
		if err := currentWorkspace.save(); err != nil {
			fmt.Printf("failed to create new workspace: %v\n", err)
			return
		}
		fmt.Printf("workspace %s created\n", workspace)
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceNewCmd)
	workspaceNewCmd.Flags().StringP("region", "r", "us-east", "Default region to use")
	workspaceNewCmd.Flags().StringP("environment", "e", "client-test", "Default environment to use")
}

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

// workspaceInfoCmd represents the info command
var workspaceInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information on current workspace",
	Long:  `Shows detailed information on current workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Workspace name:       %s\n", currentWorkspace.Name)
		fmt.Printf("Default region:       %s\n", currentWorkspace.DefaultRegion)
		fmt.Printf("Default environment:  %s\n", currentWorkspace.DefaultEnvironment)
		loginStatus := "never logged in"
		if expired := currentWorkspace.iamExpireTime(); expired != nil {
			if currentWorkspace.iamLoginExpired() {
				loginStatus = fmt.Sprintf("refresh required (expired at %v)", expired)
				if currentWorkspace.IAMRefreshToken == "" {
					loginStatus = "login required"
				}
			} else {
				loginStatus = fmt.Sprintf("active (expires at %v)", expired)
			}
		}
		fmt.Printf("IAM Login status:     %s\n", loginStatus)
		fmt.Printf("IAM Region:           %s\n", currentWorkspace.IAMRegion)
		fmt.Printf("IAM Environment:      %s\n", currentWorkspace.IAMEnvironment)
		if currentWorkspace.HASConfig.HASURL != "" {
			fmt.Printf("HAS Region:           %s\n", currentWorkspace.HASRegion)
			fmt.Printf("HAS URL:              %s\n", currentWorkspace.HASConfig.HASURL)

		}
		fmt.Printf("\n")
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceInfoCmd)
}

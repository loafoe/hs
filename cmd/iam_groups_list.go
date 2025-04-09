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
	"github.com/dip-software/go-dip-api/iam"
	"github.com/spf13/cobra"
)

// iamListGroupsCmd represents the listGroups command
var iamListGroupsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List IAM groups",
	Long:    `List all IAM groups.`,
	Run: func(cmd *cobra.Command, args []string) {
		iamClient, err := getIAMClient(cmd)
		if err != nil {
			fmt.Printf("error initalizing IAM client: %v\n", err)
			return
		}
		opts := &iam.GetGroupOptions{}
		if name, err := cmd.Flags().GetString("name"); err == nil && name != "" {
			opts.Name = &name
		}
		if org, err := cmd.Flags().GetString("org"); err == nil && org != "" {
			opts.OrganizationID = &org
		}
		if opts.OrganizationID == nil || *opts.OrganizationID == "" {
			if currentWorkspace.IAMSelectedOrg == "" {
				fmt.Printf("Select an organization first\n")
				return
			}
			opts.OrganizationID = &currentWorkspace.IAMSelectedOrg
		}
		groups, _, err := iamClient.Groups.GetGroups(opts)
		if err != nil {
			fmt.Printf("error retrieving groups: %v\n", err)
			return
		}
		if jsonOut {
			data, _ := json.Marshal(*groups)
			fmt.Printf("%s\n", string(data))
			return
		}
		t := tabby.New()
		t.AddHeader("group", "id", "description")
		for _, group := range *groups {
			t.AddLine(group.GroupName,
				group.ID,
				group.GroupDescription)
		}
		t.Print()
		_ = currentWorkspace.saveWithIAM(iamClient)
	},
}

func init() {
	iamGroupsCmd.AddCommand(iamListGroupsCmd)
	iamListGroupsCmd.Flags().String("name", "", "Filter by name")
	iamListGroupsCmd.Flags().String("org", "", "Filter by orgID")
}

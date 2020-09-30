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
	"github.com/philips-software/go-hsdp-api/iam"

	"github.com/spf13/cobra"
)

// iamRolesListCmd represents the list command
var iamRolesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List roles",
	Long:    `Lists IAM roles.`,
	Run: func(cmd *cobra.Command, args []string) {
		iamClient, err := getIAMClient(cmd)
		if err != nil {
			fmt.Printf("error initalizing IAM client: %v\n", err)
			return
		}
		opts := &iam.GetRolesOptions{}
		if name, err := cmd.Flags().GetString("name"); err == nil && name != "" {
			opts.Name = &name
		}
		if org, err := cmd.Flags().GetString("org"); err == nil && org != "" {
			opts.OrganizationID = &org
		}
		roles, _, err := iamClient.Roles.GetRoles(opts)
		if err != nil {
			fmt.Printf("error retrieving roles: %v\n", err)
			return
		}
		if jsonOut {
			data, _ := json.Marshal(*roles)
			fmt.Printf("%s\n", string(data))
			return
		}
		t := tabby.New()
		t.AddHeader("role", "id")
		for _, role := range *roles {
			t.AddLine(role.Name,
				role.ID)
		}
		t.Print()
		_ = currentWorkspace.saveWithIAM(iamClient)
	},
}

func init() {
	iamRolesCmd.AddCommand(iamRolesListCmd)
	iamRolesListCmd.Flags().String("name", "", "Filter by name")
	iamRolesListCmd.Flags().String("org", "", "Filter by orgID")

}

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

	"github.com/spf13/cobra"
)

// iamUsersListCmd represents the list command
var iamUsersLookupCmd = &cobra.Command{
	Use:     "lookup",
	Aliases: []string{"look"},
	Short:   "Lookup user",
	Long:    `Looks up a user based on a GUID.`,
	Run: func(cmd *cobra.Command, args []string) {
		iamClient, err := getIAMClient(cmd)
		if err != nil {
			fmt.Printf("error initalizing IAM client: %v\n", err)
			return
		}
		if currentWorkspace.IAMSelectedOrg == "" {
			fmt.Printf("please select an organization first using: hs iam orgs select\n")
			return
		}
		if len(args) == 0 {
			fmt.Printf("please specify the GUID as only argument\n")
			return
		}
		guid := args[0]

		user, _, err := iamClient.Users.LegacyGetUserByUUID(guid)
		if err != nil {
			fmt.Printf("error performing IAM introspect: %v\n", err)
			return
		}
		if err != nil || user == nil {
			fmt.Printf("{}\n")
			return
		}
		data, _ := json.Marshal(user)
		fmt.Printf("%s\n", pretty(data))
		return
		_ = currentWorkspace.saveWithIAM(iamClient)
	},
}

func init() {
	iamUsersCmd.AddCommand(iamUsersLookupCmd)
}

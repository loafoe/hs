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

	"github.com/philips-software/go-hsdp-api/iam"

	"github.com/spf13/cobra"
)

// iamUsersChangeLoginCmd represents the change command
var iamUsersChangeLoginCmd = &cobra.Command{
	Use:     "change-login",
	Aliases: []string{"cl"},
	Short:   "Change login ID",
	Long:    `Changes the login ID.`,
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
		oldLogin, _ := cmd.Flags().GetString("old")
		newLogin, _ := cmd.Flags().GetString("new")
		id, _, err := iamClient.Users.GetUserIDByLoginID(oldLogin)
		if err != nil {
			fmt.Printf("error retrieving user %s: %v\n", oldLogin, err)
			return
		}
		ok, _, err := iamClient.Users.ChangeLoginID(iam.Person{ID: id}, newLogin)
		if err != nil || !ok {
			fmt.Printf("error changing loginID: %v\n", err)
			return
		}
		fmt.Printf("OK\n")
	},
}

func init() {
	iamUsersCmd.AddCommand(iamUsersChangeLoginCmd)
	iamUsersChangeLoginCmd.Flags().String("old", "", "Old login ID")
	iamUsersChangeLoginCmd.Flags().String("new", "", "New login ID")
}

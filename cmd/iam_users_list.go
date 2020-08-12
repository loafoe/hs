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

// iamUsersListCmd represents the list command
var iamUsersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Long:  `Lists users in the selected organization.`,
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
		pageSize := "50" // TODO: implement paging
		users, _, err := iamClient.Users.GetUsers(&iam.GetUserOptions{
			OrganizationID: &currentWorkspace.IAMSelectedOrg,
			PageSize:       &pageSize,
		})
		if err != nil {
			fmt.Printf("error performing IAM introspect: %v\n", err)
			return
		}
		if !jsonOut {
			fmt.Printf("Users in Organization: %s\n\n", currentWorkspace.IAMSelectedOrgName)
		}
		if len(users.UserUUIDs) == 0 {
			if jsonOut {
				fmt.Printf("[]\n")
				return
			}
			fmt.Printf("no users found or not enough permissions\n")
			return
		}
		numWorkers := 10
		numUsers := len(users.UserUUIDs)
		if numWorkers > numUsers {
			numWorkers = numUsers
		}
		done := make(chan bool)
		queue := make(chan string)
		result := make(chan *iam.User)
		// Start workers
		for i := 0; i < numWorkers; i++ {
			go fetchUser(iamClient, i, queue, result, done)
		}
		// Fill queue
		go func() {
			for _, uuid := range users.UserUUIDs {
				queue <- uuid
			}
		}()

		if jsonOut {
			fmt.Printf("[")
			for i := 0; i < numUsers; i++ {
				user := <-result
				data, _ := json.Marshal(user)
				fmt.Printf("%s", string(data))
				if i < numUsers-1 { // Last one
					fmt.Printf(",")
				}
			}
			fmt.Printf("]\n")
			return
		}
		t := tabby.New()
		t.AddHeader("loginID", "first name", "last name", "email")
		for i := 0; i < numUsers; i++ {
			user := <-result
			t.AddLine(user.LoginID, user.Name.Given, user.Name.Family, user.EmailAddress)
		}
		t.Print()
		// Clean up
		for i := 0; i < numWorkers; i++ {
			done <- true
		}
	},
}

func fetchUser(client *iam.Client, workerID int, queue chan string, result chan *iam.User, done chan bool) {
	errorUser := &iam.User{
		ID: "error",
	}
	for {
		select {
		case id := <-queue:
			user, _, _ := client.Users.GetUserByID(id)
			if user != nil {
				result <- user
			} else {
				result <- errorUser
			}
		case <-done:
			return
		}
	}
}

func init() {
	iamUsersCmd.AddCommand(iamUsersListCmd)
}

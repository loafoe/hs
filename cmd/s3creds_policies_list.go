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
	"strings"

	"github.com/cheynewallace/tabby"
	"github.com/dip-software/go-dip-api/s3creds"
	"github.com/spf13/cobra"
)

// s3credsPoliciesListCmd represents the list command
var s3credsPoliciesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li"},
	Short:   "List S3 Credentials policies",
	Long:    `Lists know S3 Credentials policies.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getCredentialsClient(cmd, args)
		if err != nil {
			fmt.Printf("error initializing S3 Credentials client: %v\n", err)
			return
		}
		policies, _, err := client.Policy.GetPolicy(&s3creds.GetPolicyOptions{
			ProductKey: &currentWorkspace.S3CredsProductKey,
		})
		if err != nil {
			fmt.Printf("error retrieving policies list: %v\n", err)
			return
		}
		if jsonOut {
			data, _ := json.Marshal(policies)
			fmt.Printf("%s\n", data)
			return
		}
		t := tabby.New()
		t.AddHeader("policy id", "managing orgs", "groups", "resources")
		for _, r := range policies {
			t.AddLine(r.ID,
				strings.Join(r.Conditions.ManagingOrganizations, ","),
				strings.Join(r.Conditions.Groups, ","),
				strings.Join(r.Allowed.Resources, ","))
		}
		t.Print()
		if len(policies) == 0 {
			fmt.Printf("no policies found\n")
		}
	},
}

func init() {
	s3credsPoliciesCmd.AddCommand(s3credsPoliciesListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3credsPoliciesListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3credsPoliciesListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

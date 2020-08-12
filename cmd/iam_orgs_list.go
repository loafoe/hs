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
	"github.com/spf13/cobra"
)

// iamOrgsListCmd represents the list command
var iamOrgsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List IAM organizations",
	Long:    `Lists IAM organizations you have access to.`,
	Run: func(cmd *cobra.Command, args []string) {
		iamClient, err := getIAMClient(cmd)
		if err != nil {
			fmt.Printf("error initalizing IAM client: %v\n", err)
			return
		}
		introspect, _, err := iamClient.Introspect()
		if err != nil {
			fmt.Printf("error performing IAM introspect: %v\n", err)
			return
		}
		if jsonOut {
			data, _ := json.Marshal(introspect.Organizations.OrganizationList)
			fmt.Printf("%s\n", string(data))
			return
		}
		t := tabby.New()
		t.AddHeader("organization", "id")
		for _, org := range introspect.Organizations.OrganizationList {
			t.AddLine(org.OrganizationName,
				org.OrganizationID)
		}
		t.Print()
	},
}

func init() {
	iamOrgsCmd.AddCommand(iamOrgsListCmd)
}

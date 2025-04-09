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

	"github.com/dip-software/go-dip-api/s3creds"

	"github.com/spf13/cobra"
)

// s3credsGetCmd represents the get command
var s3credsGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	Short:   "Get S3 Credentials",
	Long:    `Gets S3 Credentials for the given configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getCredentialsClient(cmd, args)
		if err != nil {
			fmt.Printf("error initializing S3 Credentials client: %v\n", err)
			return
		}
		access, _, err := client.Access.GetAccess(&s3creds.GetAccessOptions{
			ProductKey: &currentWorkspace.S3CredsProductKey,
		})
		if err != nil {
			fmt.Printf("Error retrieving credentials: %v\n", err)
		}
		data, _ := json.Marshal(access)
		fmt.Println(pretty(data))
	},
}

func init() {
	s3credsCmd.AddCommand(s3credsGetCmd)
}

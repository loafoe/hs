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
	"net/http"

	"github.com/philips-software/go-hsdp-api/iam"

	"github.com/spf13/cobra"
)

// iamIntrospectCmd represents the introspect command
var iamIntrospectCmd = &cobra.Command{
	Use:     "introspect",
	Aliases: []string{"in", "intro"},
	Short:   "Introspect using current token",
	Long:    `Does an introspect call with the current active token`,
	Run: func(cmd *cobra.Command, args []string) {
		iamClient, err := iam.NewClient(http.DefaultClient, &iam.Config{
			Region:         currentWorkspace.IAMRegion,
			Environment:    currentWorkspace.IAMEnvironment,
			OAuth2ClientID: clientID,
			OAuth2Secret:   clientSecret,
			Debug:          true,
			DebugLog:       "/tmp/hs_iam_introspect.log",
		})
		if err != nil {
			fmt.Printf("error initializing IAM client: %v\n", err)
			return
		}
		iamClient.SetTokens(currentWorkspace.IAMAccessToken,
			currentWorkspace.IAMRefreshToken,
			currentWorkspace.IAMIDToken,
			currentWorkspace.IAMAccessTokenExpires)
		introspect, _, err := iamClient.Introspect()
		if err != nil {
			fmt.Printf("error performing introspect: %v\n", err)
			return
		}
		data, err := json.Marshal(introspect)
		if err != nil {
			fmt.Printf("error marshalling introspect result: %v\n", err)
			return
		}
		if introspect.Expires > currentWorkspace.IAMAccessTokenExpires {
			currentWorkspace.IAMAccessToken = iamClient.Token()
			currentWorkspace.IAMRefreshToken = iamClient.RefreshToken()
			currentWorkspace.IAMAccessTokenExpires = iamClient.Expires()
			_ = currentWorkspace.save()
		}
		fmt.Println(pretty(data))
	},
}

func init() {
	iamCmd.AddCommand(iamIntrospectCmd)

}

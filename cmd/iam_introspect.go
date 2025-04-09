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
	"time"

	"github.com/dip-software/go-dip-api/iam"

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
		})
		if err != nil {
			fmt.Printf("error initializing IAM client: %v\n", err)
			return
		}
		tempToken := false
		useExpires := time.Now().Add(30 * time.Minute).Unix()
		useToken, _ := cmd.Flags().GetString("token")
		if useToken != "" {
			tempToken = true
		} else {
			useToken = currentWorkspace.IAMAccessToken
			useExpires = currentWorkspace.IAMAccessTokenExpires
		}

		iamClient.SetTokens(useToken,
			currentWorkspace.IAMRefreshToken,
			currentWorkspace.IAMIDToken,
			useExpires)
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
		if !tempToken && introspect.Expires > currentWorkspace.IAMAccessTokenExpires {
			token, _ := iamClient.Token()
			currentWorkspace.IAMAccessToken = token
			currentWorkspace.IAMRefreshToken = iamClient.RefreshToken()
			currentWorkspace.IAMAccessTokenExpires = iamClient.Expires()
			currentWorkspace.IAMIDToken = iamClient.IDToken()
			_ = currentWorkspace.save()
		}
		fmt.Println(pretty(data))
	},
}

func init() {
	iamCmd.AddCommand(iamIntrospectCmd)
	iamIntrospectCmd.Flags().String("token", "", "Introspect this token")

}

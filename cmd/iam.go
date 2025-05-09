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
	"net/http"

	"github.com/dip-software/go-dip-api/iam"
	"github.com/spf13/cobra"
)

// iamCmd represents the iam command
var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Interact with HSDP IAM resources",
	Long:  `Interact with HSDP IAM resources`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(iamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// iamCmd.PersistentFlags().String("foo", "", "A help for foo")
	iamCmd.PersistentFlags().StringP("region", "r", "", "HSDP region to use (default: us-east)")
	iamCmd.PersistentFlags().StringP("environment", "e", "", "HSDP environment to use (default: client-test)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// iamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getIAMClient(_ *cobra.Command) (*iam.Client, error) {
	iamClient, err := iam.NewClient(http.DefaultClient, &iam.Config{
		Region:         currentWorkspace.IAMRegion,
		Environment:    currentWorkspace.IAMEnvironment,
		OAuth2ClientID: clientID,
		OAuth2Secret:   clientSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("iam client: %w", err)
	}
	iamClient.SetTokens(currentWorkspace.IAMAccessToken,
		currentWorkspace.IAMRefreshToken,
		currentWorkspace.IAMIDToken,
		currentWorkspace.IAMAccessTokenExpires)
	return iamClient, nil
}

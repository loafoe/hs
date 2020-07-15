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

	"github.com/philips-software/go-hsdp-api/config"
	"github.com/philips-software/go-hsdp-api/has"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/spf13/cobra"
)

// ironTasksListCmd represents the list command
var hasImageListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li"},
	Short:   "List available has images",
	Long:    `Lists the available list of HAS machine images`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		orgID, _ := cmd.Flags().GetString("orgid")
		if url == "" {
			c, err := config.New(config.WithRegion(currentWorkspace.IAMRegion),
				config.WithEnv(currentWorkspace.IAMEnvironment))
			if err != nil {
				fmt.Printf("failed to autoconfig HAS backend URL: %v\n", err)
				return
			}
			url, err = c.Service("has").GetString("url")
			if err != nil {
				fmt.Printf("need a HAS backend URL: %v\n", err)
				return
			}
		}
		currentWorkspace.HASConfig.HASURL = url
		iamClient, err := iam.NewClient(http.DefaultClient, &iam.Config{
			Region:         currentWorkspace.IAMRegion,
			Environment:    currentWorkspace.IAMEnvironment,
			OAuth2ClientID: clientID,
			OAuth2Secret:   clientSecret,
			Debug:          true,
			DebugLog:       "/tmp/hs_has_iam.log",
		})
		if err != nil {
			fmt.Printf("error initializing IAM client: %v\n", err)
		}
		iamClient.SetToken(currentWorkspace.IAMAccessToken)
		if orgID == "" {
			introspect, _, err := iamClient.Introspect()
			if err != nil {
				fmt.Printf("failed to fetch organization: %v\n", err)
				return
			}
			orgID = introspect.Organizations.ManagingOrganization
			currentWorkspace.HASConfig.OrgID = orgID
		}
		client, err := has.NewClient(iamClient, &has.Config{
			HASURL:   url,
			OrgID:    orgID,
			Debug:    true,
			DebugLog: "/tmp/hs_has.log",
		})
		if err != nil {
			fmt.Printf("error initializing HAS client: %v\n", err)
			return
		}
		images, _, err := client.Images.GetImages()
		if err != nil {
			fmt.Printf("error retrieving image list: %v\n", err)
			return
		}
		for _, image := range *images {
			fmt.Printf("%s -- %s\n", image.ID, image.Name)
		}
	},
}

func init() {
	hasImagesCmd.AddCommand(hasImageListCmd)
}

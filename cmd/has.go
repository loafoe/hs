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

	"github.com/manifoldco/promptui"

	"github.com/philips-software/go-hsdp-api/config"
	"github.com/philips-software/go-hsdp-api/has"
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/spf13/cobra"
)

// hasCmd represents the has command
var hasCmd = &cobra.Command{
	Use:   "has",
	Short: "Manage Hosted Appstream resources",
	Long:  `Manage Hosted Appstream resources`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var resourceSelectTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   "\U0001F4E6 {{ .ID | cyan }} ({{ .ImageID | red }})",
	Inactive: "  {{ .ID | cyan }} ({{ .ImageID | red }})",
	Selected: "\U0001F4E6 {{ .ID | red | cyan }}",
}

var imageSelectTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   "\U0001F5BC {{ .Name | cyan }} ({{ .Regions | red }})",
	Inactive: "  {{ .Name | cyan }} ({{ .Regions | red }})",
	Selected: "\U0001F5BC {{ .Name | red | cyan }}",
}

var resourceTypeSelectTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   "\U0001F5A5 {{ .Name | cyan }}",
	Inactive: "  {{ .Name | cyan }}",
	Selected: "\U0001F5A5 {{ .Name | red | cyan }}",
}

var sessionSelectTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   "\U0001F5A5 {{ .SessionID | cyan }}",
	Inactive: "  {{ .SessionID | cyan }}",
	Selected: "\U0001F5A5 {{ .SessionID | red | cyan }}",
}

type hasImage struct {
	ID      string
	Name    string
	Regions string
}

func init() {
	rootCmd.AddCommand(hasCmd)

	hasCmd.PersistentFlags().StringP("url", "u", "", "The HAS backend server to use")
	hasCmd.PersistentFlags().StringP("orgid", "o", "", "The organization ID (tenant) to use")
	hasCmd.PersistentFlags().StringP("has-region", "R", "", "Use the specified HAS region for operations")
	hasCmd.Flags().StringP("region", "r", "", "Use the specified IAM region for operations")

}

func contains(a []string, b string) bool {
	for _, r := range a {
		if r == b {
			return true
		}
	}
	return false
}

func getHASClient(cmd *cobra.Command, args []string) (*has.Client, error) {
	url, _ := cmd.Flags().GetString("url")
	orgID, _ := cmd.Flags().GetString("orgid")
	region, _ := cmd.Flags().GetString("region")
	hasRegion, _ := cmd.Flags().GetString("has-region")
	if hasRegion == "" {
		if currentWorkspace.HASRegion == "" {
			return nil, fmt.Errorf("no HAS region configured")
		}
	} else {
		currentWorkspace.HASRegion = hasRegion
	}
	if url == "" {
		if currentWorkspace.HASConfig.HASURL == "" {
			c, err := config.New(config.WithRegion(currentWorkspace.IAMRegion),
				config.WithEnv(currentWorkspace.IAMEnvironment))
			if err != nil {
				return nil, fmt.Errorf("autoconfig: %w", err)
			}
			if region != "" {
				c = c.Region(region)
			}
			url, err = c.Service("has").GetString("url")
			if err != nil {
				return nil, fmt.Errorf("service: %w", err)
			}
			currentWorkspace.HASConfig.HASURL = url
		}
	} else {
		currentWorkspace.HASConfig.HASURL = url
	}
	iamClient, err := iam.NewClient(http.DefaultClient, &iam.Config{
		Region:         currentWorkspace.IAMRegion,
		Environment:    currentWorkspace.IAMEnvironment,
		OAuth2ClientID: clientID,
		OAuth2Secret:   clientSecret,
		Debug:          true,
		DebugLog:       "/tmp/hs_has_iam.log",
	})
	if err != nil {
		return nil, fmt.Errorf("iam client: %w", err)
	}
	iamClient.SetTokens(currentWorkspace.IAMAccessToken,
		currentWorkspace.IAMRefreshToken,
		currentWorkspace.IAMIDToken,
		currentWorkspace.IAMAccessTokenExpires)
	if orgID == "" {
		if currentWorkspace.HASConfig.OrgID == "" {
			introspect, _, err := iamClient.Introspect()
			if err != nil {
				return nil, fmt.Errorf("organization: %w", err)
			}
			orgID = introspect.Organizations.ManagingOrganization
			currentWorkspace.HASConfig.OrgID = orgID
		}
	} else {
		currentWorkspace.HASConfig.OrgID = orgID
	}
	return has.NewClient(iamClient, &has.Config{
		HASURL:   currentWorkspace.HASConfig.HASURL,
		OrgID:    currentWorkspace.HASConfig.OrgID,
		Debug:    true,
		DebugLog: "/tmp/hs_has.log",
	})
}

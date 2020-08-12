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

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type iamOrg struct {
	ID   string
	Name string
}

var orgSelectTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}",
	Active:   "\U0001F3E0{{ .Name | cyan }}",
	Inactive: "  {{ .Name | cyan }}",
	Selected: " {{ .Name | red | cyan }}",
}

// iamOrgsSelectCmd represents the select command
var iamOrgsSelectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select active organization",
	Long:  `Selects the active organization.`,
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
		if len(introspect.Organizations.OrganizationList) == 0 {
			fmt.Printf("No organizations found\n")
			return
		}
		orgs := make([]iamOrg, 0)
		for _, o := range introspect.Organizations.OrganizationList {
			orgs = append(orgs, iamOrg{
				Name: o.OrganizationName,
				ID:   o.OrganizationID,
			})
		}
		prompt := promptui.Select{
			Label:     "Select active organization",
			Items:     orgs,
			HideHelp:  true,
			Templates: orgSelectTemplate,
			IsVimMode: false,
			Stdout:    &bellSkipper{},
		}
		i, _, err := prompt.Run()
		if err != nil {
			return
		}
		currentWorkspace.IAMSelectedOrg = orgs[i].ID
		_ = currentWorkspace.save()
	},
}

func init() {
	iamOrgsCmd.AddCommand(iamOrgsSelectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// iamOrgsSelectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// iamOrgsSelectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

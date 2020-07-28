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
	"github.com/philips-software/go-hsdp-api/has"

	"github.com/spf13/cobra"
)

// hasResourcesDeleteCmd represents the delete command
var hasResourcesDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d", "rm"},
	Short:   "Delete a HAS resource",
	Long:    `Deletes a HAS resource.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getHASClient(cmd, args)
		if err != nil {
			fmt.Printf("error initializing HAS client: %v\n", err)
			return
		}
		resources, _, err := client.Resources.GetResources(&has.ResourceOptions{})
		if err != nil {
			fmt.Printf("error retrieving image list: %v\n", err)
			return
		}
		if len(*resources) == 0 {
			fmt.Printf("no resources found\n")
			return
		}
		prompt := promptui.Select{
			Label:     "Select resource",
			Items:     *resources,
			HideHelp:  true,
			Templates: resourceSelectTemplate,
			IsVimMode: false,
		}
		i, _, err := prompt.Run()
		if err != nil {
			return
		}
		resourceID := (*resources)[i].ResourceID
		res, resp, err := client.Resources.DeleteResource(resourceID)
		if err != nil {
			fmt.Printf("error deleting resource: %v\n", err)
			return
		}
		if res != nil {
			fmt.Printf("resource deleted: %s\n", resourceID)
			return
		}
		fmt.Printf("unexpected error deleting resource: %v\n", resp)
	},
}

func init() {
	hasResourcesCmd.AddCommand(hasResourcesDeleteCmd)

}

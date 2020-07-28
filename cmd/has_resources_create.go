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

	"github.com/philips-software/go-hsdp-api/has"

	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
)

// hasResourcesCreateCmd represents the create command
var hasResourcesCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a HAS resource",
	Long:    `Creates a HAS resource.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getHASClient(cmd, args)
		if err != nil {
			fmt.Printf("error initializing HAS client: %v\n", err)
			return
		}
		images, _, err := client.Images.GetImages()
		if err != nil {
			fmt.Printf("error retrieving image list: %v\n", err)
			return
		}
		prompt := promptui.Select{
			Label:     "Select Image",
			Items:     *images,
			HideHelp:  true,
			Templates: imageSelectTemplate,
			IsVimMode: false,
		}
		i, _, err := prompt.Run()
		if err != nil {
			return
		}
		image := (*images)[i].ID

		var resourceTypes = []struct {
			Name string
		}{
			{"g3s.xlarge"},
			{"g3.4xlarge"},
			{"g3.8xlarge"},
			{"g3.16xlarge"},
		}

		prompt = promptui.Select{
			Label:     "Select resource type",
			Items:     resourceTypes,
			Templates: resourceTypeSelectTemplate,
			HideHelp:  true,
			IsVimMode: false,
		}
		i, _, err = prompt.Run()
		if err != nil {
			return
		}
		resourceType := resourceTypes[i].Name
		res, resp, err := client.Resources.CreateResource(has.Resource{
			ImageID:      image,
			ResourceType: resourceType,
			Region:       currentWorkspace.HASRegion,
			Count:        1,
			ClusterTag:   "created-with-hs",
			EBS: has.EBS{
				DeleteOnTermination: true,
				VolumeSize:          50,
				VolumeType:          "standard",
			},
		})
		if err != nil {
			fmt.Printf("failed to create resources: %v\n", err)
			return
		}
		if res == nil {
			fmt.Printf("failed to create resource: %v\n", resp)
			return
		}
		fmt.Printf("resource %s creation started, state: %s\n", res.ID, res.State)
	},
}

func init() {
	hasResourcesCmd.AddCommand(hasResourcesCreateCmd)

	// hasResourcesCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

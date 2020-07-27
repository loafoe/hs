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

	"github.com/cheynewallace/tabby"
	"github.com/philips-software/go-hsdp-api/has"

	"github.com/spf13/cobra"
)

// hasResourcesListCmd represents the list command
var hasResourcesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li"},
	Short:   "List HAS resources",
	Long:    `Lists HAS resources`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getHASClient(cmd, args)
		if err != nil {
			fmt.Printf("error initializing HAS client: %v\n", err)
			return
		}
		if err != nil {
			fmt.Printf("error initializing HAS client: %v\n", err)
			return
		}
		resources, _, err := client.Resources.GetResources(&has.ResourceOptions{
			Region: &currentWorkspace.DefaultRegion,
		})
		if err != nil {
			fmt.Printf("error retrieving resources list: %v\n", err)
			return
		}
		t := tabby.New()
		t.AddHeader("resource id", "region", "state", "image", "dns")
		for _, r := range *resources {
			t.AddLine(r.ID,
				r.Region,
				r.State,
				r.ImageID,
				r.DNS)
		}
		t.Print()
		if len(*resources) == 0 {
			fmt.Printf("no resources found\n")
		}
		for _, resource := range *resources {
			fmt.Printf("%s -- %s\n", resource.ID, resource.State)
		}
	},
}

func init() {
	hasResourcesCmd.AddCommand(hasResourcesListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hasResourcesListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	hasResourcesListCmd.Flags().StringP("region", "r", "", "List images in this region")
}

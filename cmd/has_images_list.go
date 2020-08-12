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
	"strings"

	"github.com/cheynewallace/tabby"

	"github.com/spf13/cobra"
)

// ironTasksListCmd represents the list command
var hasImageListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List available has images",
	Long:    `Lists the available list of HAS machine images`,
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
		if jsonOut {
			data, _ := json.Marshal(images)
			fmt.Printf("%s\n", data)
			return
		}
		t := tabby.New()
		t.AddHeader("image id", "name", "regions")
		for _, i := range *images {
			t.AddLine(i.ID,
				i.Name,
				strings.Join(i.Regions, ","))
		}
		t.Print()
		if len(*images) == 0 {
			fmt.Printf("no images found\n")
		}
		fmt.Printf("\n")
	},
}

func init() {
	hasImagesCmd.AddCommand(hasImageListCmd)
}

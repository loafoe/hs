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
	"github.com/philips-software/go-hsdp-api/iam"
	"github.com/spf13/cobra"
	"net/http"
)

// ironTasksListCmd represents the list command
var imageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available has images",
	Long: `Lists the available list of HAS machine images`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.PersistentFlags().GetString("url")
		orgID, _ := cmd.PersistentFlags().GetString("orgid")
		if url == "" {
			fmt.Printf("need a HAS backend URL\n")
			return
		}
		if orgID == "" {
			fmt.Printf("need an organization ID\n")
			return
		}
		iamClient, err := iam.NewClient(http.DefaultClient, &iam.Config{
		})
		if err != nil {
			fmt.Printf("error initializing IAM client: %v\n", err)
		}
		client, err := has.NewClient(iamClient, &has.Config{
			HASURL: url,
			OrgID: orgID,
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
	imagesCmd.AddCommand(imageListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ironTasksListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ironTasksListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

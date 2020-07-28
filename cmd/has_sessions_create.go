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

// hasSessionsCreateCmd represents the create command
var hasSessionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a HAS session",
	Long:  `Creates a HAS session.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getHASClient(cmd, args)
		if err != nil {
			fmt.Printf("error initializing HAS client: %v\n", err)
			return
		}
		r := "eu-west-1"
		images, _, err := client.Images.GetImages()
		if err != nil {
			fmt.Printf("error retrieving images list: %v\n", err)
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
		imageID := (*images)[i].ID
		sessions, _, err := client.Sessions.CreateSession(currentWorkspace.IAMUserUUID, has.Session{
			ImageID:    imageID,
			Region:     r,
			ClusterTag: "created-with-hs",
		})
		if err != nil {
			fmt.Printf("failed to create session: %v\n", err)
			return
		}
		fmt.Printf("%v\n", sessions)
	},
}

func init() {
	hasSessionsCmd.AddCommand(hasSessionsCreateCmd)

}

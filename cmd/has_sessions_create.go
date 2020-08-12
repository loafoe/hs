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
	"strings"

	"github.com/pkg/browser"

	"github.com/manifoldco/promptui"
	"github.com/philips-software/go-hsdp-api/has"

	"github.com/spf13/cobra"
)

// hasSessionsCreateCmd represents the create command
var hasSessionsCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c", "n", "new"},
	Short:   "Create a HAS session",
	Long:    `Creates a HAS session.`,
	Run: func(cmd *cobra.Command, args []string) {
		devSession, _ := cmd.Flags().GetBool("dev")
		client, err := getHASClient(cmd, args)
		if err != nil {
			fmt.Printf("error initializing HAS client: %v\n", err)
			return
		}
		r := currentWorkspace.HASRegion
		images, _, err := client.Images.GetImages()
		if err != nil {
			fmt.Printf("error retrieving images list: %v\n", err)
			return
		}
		if len(*images) == 0 {
			fmt.Printf("No images found\n")
			return
		}
		hasImages := make([]hasImage, 0)
		for _, i := range *images {
			if !contains(i.Regions, currentWorkspace.HASRegion) { // Skip if no region matches
				continue
			}
			hasImages = append(hasImages, hasImage{
				Name:    i.Name,
				ID:      i.ID,
				Regions: strings.Join(i.Regions, ","),
			})
		}
		prompt := promptui.Select{
			Label:     "Select image to use",
			Items:     hasImages,
			HideHelp:  true,
			Templates: imageSelectTemplate,
			IsVimMode: false,
			Stdout:    &bellSkipper{},
		}
		i, _, err := prompt.Run()
		if err != nil {
			return
		}
		imageID := hasImages[i].ID
		session := has.Session{
			ImageID:    imageID,
			Region:     r,
			ClusterTag: "created-with-hs",
		}
		if devSession {
			session.SessionType = "DEV"
		}
		sessions, _, err := client.Sessions.CreateSession(currentWorkspace.IAMUserUUID, session)
		if err != nil {
			fmt.Printf("failed to create session: %v\n", err)
			return
		}
		if len(sessions.Sessions) > 0 {
			session := sessions.Sessions[0]
			fmt.Printf("Started new session %s\n", session.SessionID)
			if session.SessionURL != "" {
				// Open in browser
				_ = browser.OpenURL(session.SessionURL)
			}
			if session.AccessToken != "" {
				fmt.Printf("UserID: %s\nAccessToken: %s\n", session.UserID, session.AccessToken)
			}
		}
	},
}

func init() {
	hasSessionsCmd.AddCommand(hasSessionsCreateCmd)
	hasSessionsCreateCmd.Flags().BoolP("dev", "d", false, "Start a dev session")
}

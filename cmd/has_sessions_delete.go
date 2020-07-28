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

// hasSessionsDeleteCmd represents the delete command
var hasSessionsDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d", "rm"},
	Short:   "Delete a HAS session",
	Long:    `Deletes a HAS session.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getHASClient(cmd, args)
		if err != nil {
			fmt.Printf("error initializing HAS client: %v\n", err)
			return
		}
		sessions, _, err := client.Sessions.GetSessions()
		if err != nil {
			fmt.Printf("error retrieving session list: %v\n", err)
			return
		}
		if sessions == nil || len(sessions.Sessions) == 0 {
			fmt.Printf("no sessions found\n")
			return
		}
		prompt := promptui.Select{
			Label:     "Select session",
			Items:     sessions.Sessions,
			HideHelp:  true,
			Templates: sessionSelectTemplate,
			IsVimMode: false,
		}
		i, _, err := prompt.Run()
		if err != nil {
			return
		}
		sessionID := sessions.Sessions[i].SessionID
		ok, resp, err := client.Sessions.DeleteSession(currentWorkspace.IAMUserUUID)
		if err != nil {
			fmt.Printf("error deleting session %v: %v\n", sessionID, err)
			return
		}
		if ok {
			fmt.Printf("session deleted: %v\n", sessionID)
			return
		}
		fmt.Printf("unexpected error deleting session: %v\n", resp)
	},
}

func init() {
	hasSessionsCmd.AddCommand(hasSessionsDeleteCmd)

}

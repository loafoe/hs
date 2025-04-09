/*
Copyright © 2021 Andy Lo-A-Foe <andy.lo-a-foe@philips.com>

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

	"github.com/dip-software/go-dip-api/iron"
	"github.com/spf13/cobra"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:     "dump",
	Aliases: []string{"d"},
	Short:   "Dump a scheduled task",
	Long:    `Dumps all known information about a scheduled task.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		config, err := readIronConfig()
		if err != nil {
			fmt.Printf("error reading iron config: %v\n", err)
			return
		}
		client, err := iron.NewClient(config)
		if err != nil {
			fmt.Printf("error configuring iron client: %v\n", err)
			return
		}
		if !jsonOut {
			fmt.Printf("retrieving scheduled task details...\n\n")
		}
		schedules, _, err := client.Schedules.GetSchedules()
		if err != nil {
			fmt.Printf("error retrieving schedules: %v\n", err)
			return
		}
		scheduleID := ""
		for _, s := range *schedules {
			if s.CodeName == args[0] {
				scheduleID = s.ID
				break
			}
		}
		if scheduleID == "" {
			fmt.Printf("schedule not found: %s\n", args[0])
		}
		schedule, _, err := client.Schedules.GetSchedule(scheduleID)
		if err != nil {
			fmt.Printf("error retrieving schedule: %v\n", err)
			return
		}
		data, err := json.Marshal(schedule)
		if err != nil {
			fmt.Printf("error marshalling schedule result: %v\n", err)
			return
		}
		fmt.Println(pretty(data))
	},
}

func init() {
	ironTasksCmd.AddCommand(dumpCmd)
}

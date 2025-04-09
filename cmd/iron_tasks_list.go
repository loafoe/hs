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
	"time"

	"github.com/cheynewallace/tabby"
	"github.com/dip-software/go-dip-api/iron"
	"github.com/spf13/cobra"
)

// ironTasksListCmd represents the list command
var ironTasksListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List tasks on Iron",
	Long:    `Lists task on Iron`,
	Run: func(cmd *cobra.Command, args []string) {
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
			fmt.Printf("retrieving tasks and schedules...\n\n")
		}
		tasks, _, err := client.Tasks.GetTasks()
		if err != nil {
			fmt.Printf("error getting tasks: %v\n", err)
			return
		}
		if tasks == nil {
			if jsonOut {
				fmt.Printf("[]\n")
				return
			}
			fmt.Printf("no tasks found.\n")
			return
		}
		if jsonOut {
			data, _ := json.Marshal(tasks)
			fmt.Printf("%s\n", data)
			return
		}
		t := tabby.New()
		type taskStats struct {
			Queued    int
			Preparing int
			Timeout   int
			Running   int
			Cancelled int
			Error     int
			Completed int
		}
		taskEntry := map[string]taskStats{}

		t.AddHeader("code name", "queued", "preparing", "running", "error", "cancelled", "timeout", "complete")
		for _, task := range *tasks {
			entry := taskStats{}
			if existing, found := taskEntry[task.CodeName]; found {
				entry = existing
			}
			switch task.Status {
			case "running":
				entry.Running++
			case "preparing":
				entry.Preparing++
			case "timeout":
				entry.Timeout++
			case "cancelled":
				entry.Cancelled++
			case "complete":
				entry.Completed++
			case "error":
				entry.Error++
			case "queued":
				entry.Queued++
			}
			taskEntry[task.CodeName] = entry
		}
		for code, stats := range taskEntry {
			t.AddLine(code, stats.Queued, stats.Preparing, stats.Running, stats.Error, stats.Cancelled, stats.Timeout, stats.Completed)
		}
		t.Print()
		fmt.Printf("\n")
		schedules, _, err := client.Schedules.GetSchedules()
		if err != nil {
			fmt.Printf("error retrieving schedules: %v\n", err)
			return
		}
		t = tabby.New()
		t.AddHeader("schedule", "every", "status", "last", "next", "runs")
		for _, s := range *schedules {
			t.AddLine(s.CodeName,
				s.RunEvery,
				s.Status,
				s.LastRunTime.Format(time.RFC3339),
				s.NextStart.Format(time.RFC3339),
				s.RunTimes)
		}
		t.Print()
		if len(*schedules) == 0 {
			fmt.Printf("no scheduled tasks found\n")
		}
		fmt.Printf("\n")
	},
}

func init() {
	ironTasksCmd.AddCommand(ironTasksListCmd)
}

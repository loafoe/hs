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
	"github.com/philips-software/go-hsdp-api/iron"

	"github.com/spf13/cobra"
)

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule <code>",
	Aliases: []string{"s"},
	Short: "Schedule a task on a cluster",
	Long: `Schedule a task on a cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		code := args[0]
		config, err := readIronConfig()
		if err != nil {
			fmt.Printf("error reading iron config: %v\n", err)
			return
		}
		if len(config.ClusterInfo) == 0 {
			fmt.Printf("missing required cluster info in config")
		}
		config.Debug = debug
		client, err := iron.NewClient(config)
		if err != nil {
			fmt.Printf("error configuring iron client: %v\n", err)
			return
		}
		payload, _ := cmd.Flags().GetString("payload")
		if payload == "" {
			fmt.Printf("payload is required\n")
			return
		}
		encryptedPayload, err := config.ClusterInfo[0].Encrypt([]byte(payload))
		if err != nil {
			fmt.Printf("error encrypting payload: %v\n", err)
			return
		}
		timeout, _ := cmd.Flags().GetInt("timeout")
		runEvery, _ := cmd.Flags().GetInt("every")
		runTimes, _ := cmd.Flags().GetInt("times")
		cluster, _ := cmd.Flags().GetString("cluster")
		if cluster == "" {
			cluster = config.ClusterInfo[0].ClusterID
		}
		schedule, resp, err := client.Schedules.CreateSchedule(iron.Schedule{
			CodeName: code,
			Payload: encryptedPayload,
			Timeout: timeout,
			RunEvery: runEvery,
			RunTimes: runTimes,
			Cluster: cluster,
		})
		if err != nil {
			fmt.Printf("error scheduling task: %v\n", err)
			return
		}
		if schedule != nil {
			fmt.Printf("scheduled as: %s\n", schedule.ID)
			return
		}
		fmt.Printf("error status: %d\n", resp.StatusCode)
	},
}

func init() {
	ironCmd.AddCommand(scheduleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scheduleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	scheduleCmd.Flags().StringP("cluster", "c", "", "Cluster to schedule task on")
	scheduleCmd.Flags().StringP("payload", "p", "", "Payload to use")
	scheduleCmd.Flags().IntP("timeout", "t", 3600, "Timeout to use in seconds")
	scheduleCmd.Flags().IntP("every", "f", 0, "Time between runs in seconds")
	scheduleCmd.Flags().IntP("times", "r", 0, "Number of times the task will run")
}

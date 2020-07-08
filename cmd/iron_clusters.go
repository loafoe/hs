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
	"github.com/philips-software/go-hsdp-api/iron"

	"github.com/spf13/cobra"
)

// clustersCmd represents the clusters command
var clustersCmd = &cobra.Command{
	Use:   "clusters",
	Aliases: []string{"cl"},
	Short: "List available clusters",
	Long: `Lists the available Iron clusters.`,
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
		fmt.Printf("retrieving clusters...\n\n")
		clusters, _, err := client.Clusters.GetClusters()
		if err != nil {
			fmt.Printf("error retrieving clusters: %v\n", err)
			return
		}
		cl, _, _ := client.Clusters.GetCluster(config.ClusterInfo[0].ClusterID)
		if cl != nil {
			*clusters = append(*clusters, *cl)
		}
		t := tabby.New()
		t.AddHeader("cluster id", "name", "available", "total", "cpu", "memory", "disk")
		if clusters != nil {
			for _, cl := range *clusters {
				t.AddLine(cl.ID, cl.Name, cl.RunnersAvailable, cl.RunnersTotal, cl.CPUShare, cl.Memory, cl.DiskSpace)
			}
		}
		t.Print()
		fmt.Printf("\n")
	},
}

func init() {
	ironCmd.AddCommand(clustersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clustersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clustersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

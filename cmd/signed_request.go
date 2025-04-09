package cmd

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

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	signer "github.com/dip-software/go-dip-signer"
	"github.com/spf13/cobra"
)

// signedRequestCmd represents the signedRequest command
var signedRequestCmd = &cobra.Command{
	Use:     "signed-request -k key -s secret -X method url",
	Short:   "Perform a signed request",
	Long:    `Perform a request that is protected by the HSDP API signing algorithm`,
	Aliases: []string{"sr"},
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		secret, _ := cmd.Flags().GetString("secret")
		method, _ := cmd.Flags().GetString("method")
		contentType, _ := cmd.Flags().GetString("content-type")
		data, _ := cmd.Flags().GetString("data")
		debug, _ := cmd.Flags().GetBool("debug")
		apiVersion, _ := cmd.Flags().GetString("api-version")
		if key == "" || secret == "" {
			fmt.Println("required key or secret missing")
			return
		}
		s, err := signer.New(key, secret)
		if err != nil {
			fmt.Printf("error creating signer: %v\n", err)
			return
		}
		var bodyReader io.Reader
		if _, err := os.Stat(data); os.IsNotExist(err) {
			bodyReader = strings.NewReader(data)
		} else {
			body, _ := os.ReadFile(data)
			bodyReader = bytes.NewReader(body)
		}
		req, _ := http.NewRequest(method, args[0], bodyReader)
		req.Header.Set("Content-Type", contentType)
		if apiVersion != "" {
			req.Header.Set("Api-Version", apiVersion)
		}
		err = s.SignRequest(req)
		if err != nil {
			fmt.Printf("error signing: %v\n", err)
			return
		}
		client := http.DefaultClient
		if debug {
			dumped, _ := httputil.DumpRequest(req, true)
			fmt.Printf("%s\n", string(dumped))
		}
		resp, err := client.Do(req)
		if debug && resp != nil {
			dumped, _ := httputil.DumpResponse(resp, true)
			fmt.Printf("%s\n", string(dumped))
		}
		if err != nil {
			fmt.Printf("error on request: %v\n", err)
			return
		}
		if resp == nil || resp.Body == nil {
			fmt.Printf("response error\n")
		}
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error reading body: %v\n", err)
			return
		}
		if !debug {
			fmt.Printf("%v", string(respData))
		}
	},
}

func init() {
	rootCmd.AddCommand(signedRequestCmd)
	signedRequestCmd.Flags().StringP("api-version", "A", "", "Use API version")
	signedRequestCmd.Flags().Bool("debug", false, "Turn on debugging")
	signedRequestCmd.Flags().StringP("key", "k", "", "Signing key to use")
	signedRequestCmd.Flags().StringP("secret", "s", "", "Signing secret to use")
	signedRequestCmd.Flags().StringP("method", "X", "GET", "HTTP method to use")
	signedRequestCmd.Flags().StringP("content-type", "t", "application/json", "Content type to use")
	signedRequestCmd.Flags().StringP("data", "d", "", "String or file to send as body")
}

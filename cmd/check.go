// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
    "strings"

    "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks if last build is equal to the stable tag",
	Long: `Checks the image tags ':$(-tag)' against ':$(-tag)-$(-rev)' of 'qnib/GO_PIPELINE_NAME'
	When the hashes are equal the return code is 1 (FAIL), to stop the build pipeline.`,
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := os.Getenv("DOCKER_HOST")
	    if endpoint == "" {
	        endpoint = "unix:///var/run/docker.sock"
	    }

	    client, _ := docker.NewClient(endpoint)
		if os.Getenv("DOCKER_TLS_VERIFY") == "1" {
			path := os.Getenv("DOCKER_CERT_PATH")
    		ca := fmt.Sprintf("%s/ca.pem", path)
    		cert := fmt.Sprintf("%s/cert.pem", path)
    		key := fmt.Sprintf("%s/key.pem", path)
    		client, _ = docker.NewTLSClient(endpoint, cert, key, ca)
		}
	    imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	    imgName := fmt.Sprintf("qnib/%s", os.Getenv("GO_PIPELINE_NAME"))
	    nameArr := strings.Split(os.Getenv("GO_PIPELINE_NAME"), "_")
	    if len(nameArr) == 2 {
	        imgName = fmt.Sprintf("qnib/%s", nameArr[0])
	    } else if len(nameArr) > 2 {
	        fmt.Printf("Weird image name '%s'", os.Getenv("GO_PIPELINE_NAME"))
	        os.Exit(1)
	    }
	    imgTag := fmt.Sprintf("%s-%s", compTag, compRev)
	    var latest string
	    var current string
	    for _, img := range imgs {
	        for _, repotag := range img.RepoTags {
	            if strings.HasPrefix(repotag, imgName) {
	                arr := strings.Split(repotag, ":")
	                repo := arr[0]
	                tag := arr[1]
	                if repo == imgName {
	                    if tag == compTag {
	                        latest = img.ID
	                    } else if tag == imgTag {
	                        current = img.ID
	                    }
	                    fmt.Println("> RepoTag: ", repotag)
	                    fmt.Println("ID: ", img.ID)
	                    fmt.Println("Size: ", img.Size)
	                }
	            }
	        }
	    }
		if current == "" {
			fmt.Printf("current sha is empty?!\n")
			os.Exit(2)
		} else if latest == current {
	        fmt.Printf("FAIL > '%s' :  '%s'==%s !! Therefore we do not need to go on...\n", imgName, imgTag, compTag)
	        os.Exit(1)
	    } else {
	        fmt.Printf("PASS > '%s' :  (%s)'%s'!=%s(%s)\n", imgName, latest, imgTag, compTag, current)
	        os.Exit(0)
	    }
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

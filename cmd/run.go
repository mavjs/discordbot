/*
Copyright © 2021 Ang Chin Han <ang.chin.han@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

// Temporarily here
type StoicResponse struct {
	Id       int    `json:"id"` // cannot unmarshal number into Go struct field StoicResponse.id of type string
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
	Author   string `json:"author"`
}

/*
{"id":21,"body":"The soul becomes dyed with the color of its thoughts.","author_id":1,"author":"Marcus Aurelius"}
*/

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "hello" {
		_, err := s.ChannelMessageSend(m.ChannelID, "World!")
		if err != nil {
			log.Println(err)
		}
	}

	if m.Content == "o/" {
		_, err := s.ChannelMessageSend(m.ChannelID, "\\o")
		if err != nil {
			log.Println(err)
		}
	}

	if m.Content == "!stoic" {
		url := "https://stoicquotesapi.com/v1/api/quotes/random"

		resp, err := http.Get(url)

		if err != nil {
			fmt.Println("error retrieving stoicquotesapi", err)
			return
		}

		defer resp.Body.Close()

		var respBody StoicResponse

		// This is when you don't want a stream, so you have a copy you can debug
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(body, &respBody)
		if err != nil {
			log.Println("error decoding stoicquotesapi response", err, string(body))
			return
		}

		message := respBody.Body + " — " + respBody.Author
		_, err = s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			log.Println(err)
		}
	}
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the discordbot",
	Long:  `Run the discordbot`,
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("TOKEN")
		dg, err := discordgo.New("Bot " + token)
		if err != nil {
			fmt.Println("error creating Discord session,", err)
			return
		}

		dg.AddHandler(messageCreate)
		dg.Identify.Intents = discordgo.IntentsGuildMessages

		err = dg.Open()
		if err != nil {
			fmt.Println("error opening connection,", err)
			return
		}

		fmt.Println("Bot is now running.  Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc

		dg.Close()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
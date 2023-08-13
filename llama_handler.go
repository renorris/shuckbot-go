package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func LLamaHandler(s *discordgo.Session, m *discordgo.MessageCreate, llamaJobs chan<- LlamaJob) {

	// Create job
	var job LlamaJob
	job.prompt = strings.Join(strings.Split(m.Content, " ")[1:], " ")
	fmt.Println("Prompt: " + job.prompt)
	job.output = make(chan string)
	job.stopSignal = make(chan struct{})

	// Send off to llama worker
	llamaJobs <- job

	// Create message
	msg, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hmmmm... (number %d in queue)", len(llamaJobs)))
	if err != nil {
		fmt.Println("There was an error creating an AI response message: ", err)
	}

	startTime := time.Now()
	output := ""

	for token := range job.output {
		output += token

		if len(output) > 2000 {
			close(job.stopSignal)
			break
		}

		if (time.Now().UnixMilli() - startTime.UnixMilli()) >= 1000 {
			startTime = time.Now()
			_, err := s.ChannelMessageEdit(msg.ChannelID, msg.ID, output)
			if err != nil {
				fmt.Println("There was an error editing an AI response message:", err)
			}
		}
	}

	// Ensure all tokens are flushed to discord
	{
		_, err := s.ChannelMessageEdit(msg.ChannelID, msg.ID, output)
		if err != nil {
			fmt.Println("There was an error editing an AI response message:", err)
		}
	}
}

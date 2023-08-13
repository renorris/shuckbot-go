package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/go-skynet/go-llama.cpp"
)

type LlamaJob struct {
	// Initial prompt
	prompt string
	// Channel to output each response token one-by-one
	output chan string
	// Channel to signal inference to stop
	stopSignal chan struct{}
}

func LlamaWorker(jobs <-chan LlamaJob) {
	// Configure
	threads64, err := strconv.ParseInt(os.Getenv("INFERENCE_THREADS"), 10, 32)
	if err != nil {
		panic(err)
	}
	threads := int(threads64)

	parsedURL, err := url.Parse(os.Getenv("MODEL_URL"))
	check(err)
	modelFilePath := "models/" + path.Base(parsedURL.Path)

	l, err := llama.New(modelFilePath, llama.EnableF16Memory, llama.SetContext(700))
	defer l.Free()
	if err != nil {
		fmt.Println("Loading the llama model failed:", err)
	}
	log.Println("AI model loaded")

	// Handle requests one-by-one
	for {
		job := <-jobs
		prompt := os.Getenv("AI_INSTRUCTIONS") + "\n\nUSER:" + job.prompt + "\nASSISTANT:"

		_, err := l.Predict(prompt, llama.SetTokenCallback(func(token string) bool {
			fmt.Print(token)

			// Write our token to the channel output
			job.output <- token

			// Check if we need to stop early
			select {
			case <-job.stopSignal:
				fmt.Println("Received stop signal")
				close(job.output)
				return false
			default:
				return true
			}
		}), llama.SetTokens(650), llama.SetThreads(threads), llama.SetTemperature(1.0), llama.SetTopK(30), llama.SetTopP(0.18), llama.SetFrequencyPenalty(1.15), llama.EnablePromptCacheAll, llama.SetPathPromptCache("./cache"), llama.Debug)
		if err != nil {
			fmt.Println("There was an error during AI inference:", err)
		}

		// Signal job that we're done
		close(job.output)
	}
}

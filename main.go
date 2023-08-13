package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/renorris/shuckbot-go/commands"
)

var (
	Token     string = os.Getenv("DISCORD_TOKEN")
	Prefix    string = os.Getenv("COMMAND_PREFIX")
	LlamaJobs        = make(chan LlamaJob, 100)
)

func main() {
	log.Println("Starting shuckbot-go")
	config()

	log.Println("Creating discord session")
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("Error creating session: ", err)
		return
	}

	// Register MessageCreate handler func
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open the connection
	err = dg.Open()
	check(err)

	// Update bot's listening status
	dg.UpdateListeningStatus(Prefix + "help")

	// Spin up llama worker
	go LlamaWorker(LlamaJobs)

	log.Println("Running. Press CTRL-C to exit.")

	// Wait for kill signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("\nShutting down...")
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore messages from self
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore if bot isn't mentioned xor command prefix isn't present
	if (len(m.Mentions) > 0 && m.Mentions[0].ID == s.State.User.ID) != (!strings.HasPrefix(m.Content, Prefix)) {
		return
	}

	// Signal typing
	s.ChannelTyping(m.ChannelID)

	// Log message
	fmt.Println("Got command \"", m.Content, "\" from", m.Author.ID, m.Author.Username)

	// Do AI if bot is mentioned
	if len(m.Mentions) > 0 && m.Mentions[0].ID == s.State.User.ID {
		LLamaHandler(s, m, LlamaJobs)
		return
	}

	// Process command
	handlerResponse := processCommand(s, m)

	// Send replies
	for _, msg := range handlerResponse.GetReplyMessages() {
		_, err := s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			fmt.Println("There was an error sending a message: ", err)
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Configure bot
func config() {
	required_envs := []string{"DISCORD_TOKEN", "COMMAND_PREFIX", "MODEL_URL", "INFERENCE_THREADS", "AI_INSTRUCTIONS"}

	for _, env := range required_envs {
		_, present := os.LookupEnv(env)
		if !present {
			panic("Required environment variable not found: " + env)
		}
	}

	// Create models directory if not exist
	if _, err := os.Stat("models"); os.IsNotExist(err) {
		err := os.Mkdir("models", 0755)
		check(err)
	}

	// Parse GGML model url
	parsedURL, err := url.Parse(os.Getenv("MODEL_URL"))
	check(err)
	fileName := path.Base(parsedURL.Path)

	// Check if GGML model exists locally, download if necessary
	if _, err := os.Stat("models/" + fileName); os.IsNotExist(err) {
		log.Println("GGML model not found. Downloading...")
		out, err := os.Create("models/" + fileName)
		check(err)
		defer out.Close()
		resp, err := http.Get(os.Getenv("MODEL_URL"))
		check(err)
		defer resp.Body.Close()
		n, err := io.Copy(out, resp.Body)
		check(err)
		log.Println(n)
	}

	log.Println("Config OK")
}

// Delegate a handler & return its response
func processCommand(s *discordgo.Session, m *discordgo.MessageCreate) *commands.HandlerResponse {
	splitMessage := strings.Split(m.Content, " ")
	cmd := strings.TrimPrefix(splitMessage[0], Prefix)
	params := splitMessage[1:]

	// Delegate command to handler
	switch cmd {
	case "ping":
		return commands.PingHandler(s, m)
	case "avatar", "a":
		return commands.AvatarHandler(s, m, params)
	default:
		return commands.UnknownCommandHandler(s, m)
	}
}

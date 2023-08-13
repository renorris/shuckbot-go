module github.com/renorris/shuckbot-go

go 1.20

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/go-skynet/go-llama.cpp v0.0.0-20230802220037-50cee7712066
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.10.0 // indirect
)

replace github.com/go-skynet/go-llama.cpp => ./go-llama.cpp/

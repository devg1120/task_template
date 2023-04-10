package main

import (
	"context"
	"log"

	"golang.org/x/crypto/ssh"

	//issh "github.com/jlandowner/go-interactive-ssh"
	 "issh/libs"
)

func main() {
	ctx := context.Background()
/*
	config := &ssh.ClientConfig{
		User:            "gusa1120",
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password("sakiko1120"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
*/
    config := &ssh.ClientConfig{
        //User: "devg1120",
        User: "root",
        Auth: []ssh.AuthMethod{
		//ssh.Password("sakiko1120"),
            ssh.Password("Hyogo#2o2z"),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // password認証は設定
    }

	//client := issh.NewClient(config, "127.0.0.1", "22", []issh.Prompt{issh.DefaultPrompt})
	client := issh.NewClient(config, "127.0.0.1", "22", []issh.Prompt{issh.DefaultRootPrompt})

	err := client.Run(ctx, commands())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OK")
}

func commands() []*issh.Command {
	return []*issh.Command{
		//issh.CheckUser("devg1120"),
		issh.ChangeDirectory("/tmp"),
		issh.NewCommand("sleep 2"),
		issh.NewCommand("ls -l", issh.WithOutputLevelOption(issh.Output)),
	}
}

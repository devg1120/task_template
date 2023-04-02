package main

import (
    "bytes"
    "fmt"
    "log"

    "golang.org/x/crypto/ssh"
)

func main() {
    config := &ssh.ClientConfig{
        User: "devg1120",
        Auth: []ssh.AuthMethod{
            ssh.Password("sakiko1120"),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // password認証は設定
    }
    client, err := ssh.Dial("tcp", "localhost:22", config)
    if err != nil {
        log.Fatal("Failed to dial: ", err)
    }
    defer client.Close()

    session, err := client.NewSession()
    if err != nil {
        log.Fatal("Failed to create session: ", err)
    }
    defer session.Close()

    var b bytes.Buffer
    session.Stdout = &b
    if err := session.Run("/usr/bin/whoami"); err != nil {
        log.Fatal("Failed to run: " + err.Error())
    }
    fmt.Println(b.String())
}

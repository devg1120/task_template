package main

import (
    "bytes"
    "fmt"
    "log"
    "time"

    "golang.org/x/crypto/ssh"
)

func main() {
    config := &ssh.ClientConfig{
        User: "root",
        Auth: []ssh.AuthMethod{
            ssh.Password("Hyogo#2o2z"),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // password認証は設定
    }
    //client, err := ssh.Dial("tcp", "localhost:22", config)
    client, err := ssh.Dial("tcp", "10.0.2.101:22", config)
    if err != nil {
        log.Fatal("Failed to dial: ", err)
    }
    defer client.Close()

    // scp: send content remotely
    fmt.Println("scp -----")
    {
        session, err := client.NewSession()
        if err != nil {
            log.Fatal("Failed to create session: ", err)
        }
        defer session.Close()

        go func() {
            w, _ := session.StdinPipe()
            defer w.Close()
            t := time.Now()

            content := fmt.Sprintf("%d\n%s", t.Unix(), "hoge\nfuga\npiyo\n")

            fmt.Fprintln(w, "C0644", len(content), "hoge.txt")
            fmt.Fprint(w, content)
            fmt.Fprint(w, "\x00")
        }()

        //if err := session.Run("/usr/bin/scp -qrt ./"); err != nil {
        if err := session.Run("/usr/bin/scp -qrt /var/tmp"); err != nil {
            log.Fatal("Failed to run: " + err.Error())
        }
    }

    // wait
    time.Sleep(1 * time.Second)

    // view sent file
    fmt.Println("cat -----")
    {
        session, err := client.NewSession()
        if err != nil {
            log.Fatal("Failed to create session: ", err)
        }
        defer session.Close()

        var b bytes.Buffer
        session.Stdout = &b
        if err := session.Run("cat /var/tmp/hoge.txt"); err != nil {
            log.Fatal("Failed to run: " + err.Error())
        }
        fmt.Println(b.String())
    }
/*
    // add to file
    fmt.Println("add -----")
    {
        session, err := client.NewSession()
        if err != nil {
            log.Fatal("Failed to create session: ", err)
        }
        defer session.Close()

        var b bytes.Buffer
        session.Stdout = &b
        if err := session.Run(`perl -e 'printf qq/%s\n/, time' >> hoge.txt && cat hoge.txt`); err != nil {
            log.Fatal("Failed to run: " + err.Error())
        }
        fmt.Println(b.String())
    }
*/
}


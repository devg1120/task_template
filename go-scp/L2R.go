package main

import (
	"fmt"
	scp "go-scp/scp"
	"go-scp/scp/auth"
	"golang.org/x/crypto/ssh"
	"os"
        "context"
)

func main() {
	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library

	//clientConfig, _ := auth.PrivateKey("username", "/path/to/rsa/key", ssh.InsecureIgnoreHostKey())
	clientConfig, _ := auth.PasswordKey("root", "Hyogo#2o2z", ssh.InsecureIgnoreHostKey())

	/*
    clientConfig := &ssh.ClientConfig{
        User: "root",
        Auth: []ssh.AuthMethod{
            ssh.Password("Hyogo#2o2z"),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // password認証は設定
    }
*/
	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient("10.0.2.101:22", &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}


	/*
	  SCP
               local     =>  remote
	*/
	
	//src_filename := "hoge.txt"
	src_filename := "/etc/passwd"
	dst_filename := "/var/tmp/test4.txt"
	f, err := os.Open(src_filename)

        if err != nil {
            fmt.Println(err)
	    fmt.Printf("fail to open file:%s\n", src_filename)
	    os.Exit(3)

        }

	defer client.Close()

	defer f.Close()

	err = client.CopyFromFile(context.Background(), *f, dst_filename, "0655")

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}

	/*
	  SCP
               remote     =>  local
	*/
	
/*
	//src_filename := "hoge.txt"
	src_filename := "/etc/passwd"
	dst_filename := "/var/tmp/test.txt"
	//f, err := os.Open(dst_filename)
	f, err := os.Create(dst_filename)

        if err != nil {
            fmt.Println(err)
	    fmt.Printf("fail to open file:%s\n", dst_filename)
	    os.Exit(3)

        }

	defer client.Close()

	defer f.Close()

	//err = client.CopyFromFile(context.Background(), *f, dst_filename, "0655")
	err = client.CopyFromRemote(context.Background(), f, src_filename )

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}
*/

}


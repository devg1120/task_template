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

	clientConfig, _ := auth.PasswordKey("root", "Hyogo#2o2z", ssh.InsecureIgnoreHostKey())


	client := scp.NewClient("10.0.2.101:22", &clientConfig)

	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	clientConfig2, _ := auth.PasswordKey("root", "Hyogo#2o2z", ssh.InsecureIgnoreHostKey())


	client2 := scp.NewClient("10.0.2.103:22", &clientConfig2)

	err = client2.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	/*
	  SCP
               local     =>  remote
	*/
	/*
	//src_filename := "hoge.txt"
	src_filename := "/etc/passwd"
	dst_filename := "/var/tmp/test.txt"
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
*/
	/*
	  SCP
               remote     =>  local
	*/
	
	//src_filename := "hoge.txt"
	src_filename := "/etc/passwd"
	mid_filename := "/var/tmp/mid.txt"
	dst_filename := "/var/tmp/test.txt"
	//f, err := os.Open(dst_filename)
	f, err := os.Create(mid_filename)


	defer client.Close()

	defer f.Close()

	err = client.CopyFromRemote(context.Background(), f, src_filename )

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}



	f2, err := os.Open(mid_filename)

        if err != nil {
            fmt.Println(err)
            fmt.Printf("fail to open file:%s\n", src_filename)
            os.Exit(3)

        }

	err = client2.CopyFromFile(context.Background(), *f2, dst_filename, "0655")

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}
        if err != nil {
            fmt.Println(err)
	    fmt.Printf("fail to open file:%s\n", dst_filename)
	    os.Exit(3)

        }

	defer client2.Close()
}


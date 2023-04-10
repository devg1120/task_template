package main

import (
	"context"
	"fmt"
	scp "go-scp/scp"
	"go-scp/scp/auth"
	"golang.org/x/crypto/ssh"
	"os"
	//"flag"
	"github.com/jessevdk/go-flags"

)

type Target struct {
	ip       string
	port     string
	username string
	passwd   string
}

type Operation string

const (
	RemoteToLocal Operation = "RemoteToLocal"
	LocalToRemote Operation = "LocalToRemote"
)

type Session struct {
	op  Operation
	src string
	dst string
}

func session(t Target , s Session) {
/*
	t := Target{
		ip:       "127.0.0.1",
		port:     "22",
		username: "devg1120",
		passwd:   "sakiko1120",
	}
*/
/*
	s := Session{
		//op: RemoteToLocal,
		op:  LocalToRemote,
		src: "/var/tmp/test3.txt",
		dst: "/var/tmp/test4.txt",
	}
*/

	clientConfig, _ := auth.PasswordKey(t.username, t.passwd, ssh.InsecureIgnoreHostKey())

	// Create a new SCP client
	//client := scp.NewClient("10.0.2.101:22", &clientConfig)
	client := scp.NewClient(t.ip+":"+t.port, &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	switch s.op {
	case LocalToRemote:
		/*
			  SCP
		               local     =>  remote
		*/

		src_filename := s.src
		dst_filename := s.dst
		f, err := os.Open(src_filename)

		if err != nil {
			fmt.Println(err)
			fmt.Printf("fail to open file:%s\n", src_filename)
			os.Exit(3)

		}

		defer client.Close()

		defer f.Close()

		err = client.CopyFromFile(context.Background(), *f, dst_filename, "0666")
		//err = client.CopyFromFile(context.Background(), *f, dst_filename)

		if err != nil {
			fmt.Println("Error while copying file ", err)
		}

	case RemoteToLocal:
		/*
			  SCP
		               remote     =>  local
		*/

		src_filename := s.src
		dst_filename := s.dst
		f, err := os.Create(dst_filename)

		if err != nil {
			fmt.Println(err)
			fmt.Printf("fail to open file:%s\n", dst_filename)
			os.Exit(3)

		}

		defer client.Close()

		defer f.Close()

		//err = client.CopyFromFile(context.Background(), *f, dst_filename, "0655")
		err = client.CopyFromRemote(context.Background(), f, src_filename)

		if err != nil {
			fmt.Println("Error while copying file ", err)
		}
	}
}
/*
func main_() {
	var (
		op = flag.String("op", "", "OPRERATON REMOTE/LOCL")
		src = flag.String("s", "", "source file path")
		dst = flag.String("d", "", "source file path")
	)
flag.Parse()
fmt.Printf("-op: %s\n", *op)
fmt.Printf("-s: %s\n", *src)
fmt.Printf("-d: %s\n", *dst)

    fmt.Println(flag.NArg(), flag.NFlag())

}
*/

type Options struct {
	Remote  bool   `short:"r" long:"operation" description:"A name" `
	Src     string `short:"s" long:"sourcs" description:"A name" required:"true"`
	Dst     string `short:"d" long:"destination" description:"A name" required:"true"`
}

var options Options


func main() {

        var parser = flags.NewParser(&options, flags.Default)
	fmt.Printf("start....\n")
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	} else {
   	    fmt.Printf("ok....\n")
            fmt.Printf("Remote: %v\n", options.Remote)
            fmt.Printf("Src: %v\n",    options.Src)
            fmt.Printf("Dst: %v\n",    options.Dst)
        }
   
	t := Target{
		ip:       "127.0.0.1",
		port:     "22",
		username: "devg1120",
		passwd:   "sakiko1120",
	}
	s := Session {
                 op:  LocalToRemote,
                 //op:  RemoteToLocal,
                 src: "/var/tmp/test.txt",
                 dst: "/var/tmp/test1.txt",
	}

	session(t,s)
}




/*
func main() {
   
   _, err := flags.ParseArgs(&options, os.Args)
   fmt.Printf("start....\n")
   
   if err != nil {
   		switch flagsErr := err.(type) {
   		case flags.ErrorType:
   			if flagsErr == flags.ErrHelp {
   				os.Exit(0)
   			}
   			os.Exit(1)
   		default:
   			os.Exit(1)
   		}
   	os.Exit(1)
   }
   fmt.Printf("ok....\n")
   
   fmt.Printf("Remote: %v\n", options.Renote)
   fmt.Printf("Src: %v\n",    options.Src)
   fmt.Printf("Dst: %v\n",    options.Dst)

}
*/

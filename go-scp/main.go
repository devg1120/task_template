package main

import (
	"context"
	"fmt"
	scp "go-scp/scp"
	"go-scp/scp/auth"
	"golang.org/x/crypto/ssh"
	"os"
	"regexp"
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
	//RemoteToRemote Operation = "RemoteToRemote"
)

type Session struct {
	op  Operation
	src string
	dst string
}

func session(t Target, s Session) {

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
		fmt.Printf("... scp local => remote\n")

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
		fmt.Printf("... scp remote  => local\n")

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

func check_regexp(reg, str string) bool {
	//fmt.Println(regexp.MustCompile(reg).Match([]byte(str)))
	return regexp.MustCompile(reg).Match([]byte(str))
}

func is_remote(param string) bool {

	return check_regexp(`@`, param)

}

func make_data(direction string, src string, dst string) (Target, Session) {

	//fmt.Printf(".... %s   %s  %s\n", direction, src, dst)

	var user string
	var pass string
	var ip string
	var path_dst string
	var path_src string
	var mode Operation

	if direction == "LocalToRemote" {
		// dst
		bs := []byte(dst)
		// bbb#pppp@10.1.1.1:/path
		assined := regexp.MustCompile("(.*)%(.*)@(.*):(.*)")

		group := assined.FindSubmatch(bs)

		if len(group) == 5 {
			//fmt.Printf("group len:%d\n", len(group))

			//fmt.Printf("0.... %s\n", group[0])
			//fmt.Printf("1.... %s\n", group[1])
			//fmt.Printf("2.... %s\n", group[2])
			//fmt.Printf("3.... %s\n", group[3])
			//fmt.Printf("4.... %s\n", group[4])

			user = string(group[1])
			pass = string(group[2])
			ip = string(group[3])
			path_dst = string(group[4])
			path_src = src
			mode = LocalToRemote

		} else {
			fmt.Printf(" dst parse error")

		}

	} else if direction == "RemoteToLocal" {
		// src
		bs := []byte(src)
		// bbb#pppp@10.1.1.1:/path
		assined := regexp.MustCompile("(.*)%(.*)@(.*):(.*)")

		group := assined.FindSubmatch(bs)

		if len(group) == 5 {
			//fmt.Printf("group len:%d\n", len(group))

			//fmt.Printf("0.... %s\n", group[0])
			//fmt.Printf("1.... %s\n", group[1])
			//fmt.Printf("2.... %s\n", group[2])
			//fmt.Printf("3.... %s\n", group[3])
			//fmt.Printf("4.... %s\n", group[4])

			user = string(group[1])
			pass = string(group[2])
			ip = string(group[3])
			path_src = string(group[4])
			path_dst = dst
			mode = RemoteToLocal

		} else {
			fmt.Printf(" dst parse error")

		}

	} else {

	}

	/*
		t := Target{
			ip:       "127.0.0.1",
			port:     "22",
			username: "devg1120",
			passwd:   "sakiko1120",
		}
		s := Session{
			op: LocalToRemote,
			//op:  RemoteToLocal,
			src: "/var/tmp/test.txt",
			dst: "/var/tmp/test1.txt",
		}
	*/

	t := Target{
		ip:       ip,
		port:     "22",
		username: user,
		passwd:   pass,
	}
	s := Session{
		op: mode,
		//op: LocalToRemote,
		//op:  RemoteToLocal,
		src: path_src,
		dst: path_dst,
	}

	return t, s

}

func main() {

	if len(os.Args) != 3 {

		fmt.Printf("usage : %s <src> <dst> \n", os.Args[0])

		os.Exit(1)

	}

	src := os.Args[1]
	dst := os.Args[2]

	//fmt.Printf("src : %s\n", src)
	//fmt.Printf("dst : %s\n", dst)
	// Fst judgment
	src_is_remote := is_remote(src)

	// Snd judgment
	dst_is_remote := is_remote(dst)

	/*
		if (src_is_remote && dst_is_remote) ||
			(!src_is_remote && !dst_is_remote) {

			fmt.Printf("error\nusage : %s <src> <dst> \n", os.Args[0])
			os.Exit(1)

		}
	*/

	if !src_is_remote && !dst_is_remote {

		fmt.Printf("error\nusage : %s <src> <dst> \n", os.Args[0])
		os.Exit(1)

	}

	/*
		t := Target{
			ip:       "127.0.0.1",
			port:     "22",
			username: "devg1120",
			passwd:   "sakiko1120",
		}
		s := Session{
			op: LocalToRemote,
			//op:  RemoteToLocal,
			src: "/var/tmp/test.txt",
			dst: "/var/tmp/test1.txt",
		}
	*/

	if src_is_remote && dst_is_remote {
		// REMOTE TO REMOTE

		var t Target
		var s Session

		t, s = make_data("RemoteToLocal", src, "/var/tmp/temp")
		session(t, s)

		t, s = make_data("LocalToRemote", "/var/tmp/temp", dst)
		session(t, s)

	} else {

		var t Target
		var s Session

		if src_is_remote {
			// REMOTE TO LOCAL
			t, s = make_data("RemoteToLocal", src, dst)
			//_, _ = make_data(src, dst)

		} else if dst_is_remote {
			// LOCAL TO REMOTE
			//t, s := make_data(src, dst)
			t, s = make_data("LocalToRemote", src, dst)

		} else {

			os.Exit(1)
		}

		fmt.Printf("%+v\n", t)
		fmt.Printf("%+v\n", s)

		session(t, s)
	}
}

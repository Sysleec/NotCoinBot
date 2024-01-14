package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
)

type User struct {
	APIID   int    `json:"api_id"`
	APIHash string `json:"api_hash"`
	Proxy   string `json:"proxy"`
}

func commandCreate(cfg *Config, args ...string) error {
	if len(args) < 1 {
		return errors.New("no session name provided")
	}
	fmt.Println("Creating session...")

	sessionName := args[0]

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter AppID (from https://my.telegram.org/apps ): ")
	scanner.Scan()
	appID := scanner.Text()
	fmt.Print("Enter ApiHash (from https://my.telegram.org/apps ): ")
	scanner.Scan()
	apiHash := scanner.Text()

	fmt.Print("Enter your tg phone number: ")
	scanner.Scan()
	phone := scanner.Text()

	fmt.Print("Enter Proxy (example SOCKS5://login:pass@127.0.0.1:8080, can be empty): ")
	scanner.Scan()
	proxy := scanner.Text()

	appidINT, err := strconv.Atoi(appID)
	if err != nil {
		return errors.New("wrong appID, (example - 28378932)")
	}

	clientType := gotgproto.ClientType{
		Phone: phone,
	}

	_, err = gotgproto.NewClient(
		// Get AppID from https://my.telegram.org/apps
		appidINT,
		// Get ApiHash from https://my.telegram.org/apps
		apiHash,
		// ClientType, as we defined above
		clientType,
		// Optional parameters of client
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqliteSession(fmt.Sprintf("./sessions//%v", sessionName)),
		},
	)

	if err != nil {
		return errors.New("can't authorize")
	}

	newUser := User{
		APIID:   appidINT,
		APIHash: apiHash,
		Proxy:   proxy,
	}

	file, err := os.OpenFile("accounts.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	users := make(map[string]User)
	json.NewDecoder(file).Decode(&users)

	users[sessionName] = newUser

	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(users); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Session created successfully!")
	return nil
}
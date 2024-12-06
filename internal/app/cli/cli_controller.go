package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
	"yandex_GophKeeper_client/config"
	grpc_requests "yandex_GophKeeper_client/internal/app/requesters/gRPC"
	http_requesters "yandex_GophKeeper_client/internal/app/requesters/http"
	"yandex_GophKeeper_client/pkg/gophKeeperErrors"
)

// Please use these constants only in code, don`t ask user to input these numbers.
// Example:
//
//	Bad:
//	 fmt.Printf("Type '%d' to sign in and '%d' to sign up\n", commandLogin, commandRegister)
//	Good:
//	 fmt.Printf("Type '1' to sign in and '2' to sign up\n")
//	 //and then translate 1 and 2 to these constants yourself
//
// Explanation: its strange for user to see in 1 his step commands with numbers '1' and '2', and then '32' '33' on a second step.
const (
	commandRegister = iota
	commandLogin
	commandSaveData
	commandGetData
	commandBankCard
	commandTextData
	commandLoginAndPassword
	commandBinFile
)

// CommandsController is a command-line listener.
type CommandsController struct {
	Conf          config.AppConfig
	HTTPRequester *http_requesters.Requester
	GRPCRequester *grpc_requests.GRPCRequester
	Logger        *zap.SugaredLogger
}

func NewCommandsController(conf config.AppConfig, HTTPRequester *http_requesters.Requester, GRPCRequester *grpc_requests.GRPCRequester, logger *zap.SugaredLogger) *CommandsController {
	return &CommandsController{Conf: conf, HTTPRequester: HTTPRequester, GRPCRequester: GRPCRequester, Logger: logger}
}

// Run serves user`s commands until ctx is done.
func (c *CommandsController) Run(ctx context.Context) {
	fmt.Println("Auth:")
	for {
		err := c.auth()
		errWithHttpCode := &gophKeeperErrors.ErrWithHTTPCode{}
		if errors.As(err, &errWithHttpCode) {
			c.Logger.Errorf("Auth error: %v", err)
			if errWithHttpCode.Code() == http.StatusUnauthorized {
				fmt.Println("Wrong username or password")
			} else {
				fmt.Printf("Something went wrong, error: %v\n", errWithHttpCode)
			}
		} else if err != nil {
			c.Logger.Errorf("Auth error: %v", err)
			fmt.Printf("Auth error: %v\n", err)
		} else {
			break
		}
	}

	for {
		select {
		case <-ctx.Done():
			c.Logger.Info("Stopping listening commands")
			return
		default:
			c.runMainMenu()
		}
	}
}

// auth asks user to sign in or sign up and process it.
func (c *CommandsController) auth() error {
	for {
		//auth menu
		fmt.Println("1 - sign in")
		fmt.Println("2 - sign up")
		fmt.Print("Select: ")

		var command int
		_, err := fmt.Scanf("%d", &command)
		if err != nil {
			c.Logger.Errorf("Cant scan command, err: %v", err)
			continue
		}

		// prepare user data to command
		var userName, password string

		fmt.Println("Enter username:")
		_, err = fmt.Scanf("%s", &userName)
		if err != nil {
			return fmt.Errorf("cant scan username, err: %w", err)
		}

		fmt.Println("Enter password:")
		_, err = fmt.Scanf("%s", &password)
		if err != nil {
			return fmt.Errorf("cant scan password, err: %w", err)
		}

		// process command
		switch command {
		case 1:
			return c.login(userName, password)
		case 2:
			return c.register(userName, password)
		default:
			fmt.Println("Invalid command, try again.")
		}
	}
}

// login helps user to sign in.
func (c *CommandsController) login(userName, password string) error {
	jwt, err := c.HTTPRequester.Login(userName, password)
	errWithHttpCode := &gophKeeperErrors.ErrWithHTTPCode{}
	if errors.As(err, &errWithHttpCode) {
		if errWithHttpCode.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("wrong username or password (or user not exists)")
		} else {
			return fmt.Errorf("login error: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("login error: %w", err)
	}

	if jwt == "" {
		return fmt.Errorf("login error: empty jwt")
	}

	c.HTTPRequester.JWT = jwt
	c.GRPCRequester.JWT = jwt
	fmt.Println("Success!")
	return nil
}

// register helps user to signup.
func (c *CommandsController) register(userName, password string) error {
	jwt, err := c.HTTPRequester.RegisterUser(userName, password)
	errWithHttpCode := &gophKeeperErrors.ErrWithHTTPCode{}
	if errors.As(err, &errWithHttpCode) {
		if errWithHttpCode.StatusCode == http.StatusConflict {
			return fmt.Errorf("this user already exists")
		} else {
			return fmt.Errorf("auth error: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("auth error: %w", err)
	}

	if jwt == "" {
		return fmt.Errorf("auth error: empty jwt")
	}

	c.HTTPRequester.JWT = jwt
	c.GRPCRequester.JWT = jwt
	fmt.Println("Success!")
	return nil
}

// runMainMenu shows to user main menu and returns a command`s const.
func (c *CommandsController) runMainMenu() {
	for {
		fmt.Println("Main menu:")
		fmt.Println("1 - save data")
		fmt.Println("2 - get data")
		fmt.Print("Select: ")

		var userCmd int
		_, err := fmt.Scan(&userCmd)
		if err != nil {
			c.Logger.Errorf("Error reading command: %v", err)
			continue
		}

		switch userCmd {
		case 1:
			c.runSaveDataMenu()
		case 2:
			c.runGetDataMenu()
		default:
			fmt.Println("Invalid command, try again.")
		}
	}
}

// runSaveDataMenu shows data-save menu and serves a user`s next command (to save data on the server).
func (c *CommandsController) runSaveDataMenu() {
	if c.Conf.UseHTTPS == false {
		fmt.Printf("WARNING!\n" +
			"Connection to the server IS NOT PROTECTED.\n" +
			"You might have run this app in debug mode.\n" +
			"DON'T USE REAL SENSITIVE DATA!\n" +
			"(Actually, I don't care if your data is stolen even when using the release version,\n" +
			"since it's just an educational and pet project. So it's better not to store real data in GophKeeper.\n" +
			"It`s backend follows PCI DSS 4.0, but it's still just an educational project).\n")
	}
	for {
		fmt.Println("Choose data type:")
		fmt.Println("1 - Bank card")
		fmt.Println("2 - Login and password")
		fmt.Println("3 - Text data")
		fmt.Println("4 - Bin file")
		fmt.Println("5 - Exit")
		fmt.Println("Select:")

		var userCmd int
		_, err := fmt.Scan(&userCmd)
		if err != nil {
			c.Logger.Errorf("Error reading command: %v", err)
			fmt.Println("Please enter a valid number")
			continue
		}

		switch userCmd {
		case 1:
			c.saveBankCardData()
		case 2:
			c.saveLoginAndPasswordData()
		case 3:
			c.saveTextData()
		case 4:
			c.saveBinFileData()
		case 5:
			return
		default:
			fmt.Println("Invalid command, try again.")
		}
		return
	}
}

// runGetDataMenu shows the data retrieval menu and processes the user's command.
func (c *CommandsController) runGetDataMenu() {
	for {
		// Display the data retrieval menu
		fmt.Println("Choose data type to get:")
		fmt.Println("1 - Bank card")
		fmt.Println("2 - Login and password")
		fmt.Println("3 - Text data")
		fmt.Println("4 - Binary file")
		fmt.Println("5 - Exit")
		fmt.Print("Select: ")

		var userCmd int
		_, err := fmt.Scan(&userCmd)
		if err != nil {
			c.Logger.Errorf("Error reading command: %v", err)
			fmt.Println("Please enter a valid number.")
			continue
		}

		// Map user input to internal constants and execute the corresponding action
		switch userCmd {
		case 1:
			c.getBankCardData()
		case 2:
			c.getLoginAndPasswordData()
		case 3:
			c.getTextData()
		case 4:
			c.getBinFileData()
		case 5:
			return
		default:
			fmt.Println("Invalid command, please try again.")
		}
		return
	}
}

// saveBankCardData - asks user for a bank card data and sends it to a backend.
func (c *CommandsController) saveBankCardData() {
	//read data
	fmt.Println("Enter your bank card data:")
	var pan, ownerFirstName, ownerLastName, expiresDate string

	//read PAN with spaces
	fmt.Println("PAN:")
	reader := bufio.NewReader(os.Stdin)
	pan, err := reader.ReadString('\n')
	if err != nil {
		c.Logger.Errorf("Error reading input: %v\n", err)
		fmt.Println("Incorrect input.")
		return
	}
	pan = strings.TrimSpace(pan)

	fmt.Println("OwnerFirstName:")
	_, err = fmt.Scan(&ownerFirstName)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		fmt.Println("Incorrect input.")
		return
	}

	fmt.Println("OwnerLastName:")
	_, err = fmt.Scan(&ownerLastName)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		fmt.Println("Incorrect input.")
		return

	}

	fmt.Println("ExpiresDate:")
	_, err = fmt.Scan(&expiresDate)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		fmt.Println("Incorrect input.")
		return
	}

	//send data
	err = c.HTTPRequester.SendBankCard(pan, ownerFirstName, ownerLastName, expiresDate)
	errWithHttpCode := &gophKeeperErrors.ErrWithHTTPCode{}
	if errors.As(err, &errWithHttpCode) {
		if errWithHttpCode.Code() == http.StatusConflict {
			fmt.Println("This card has already been saved.")
		} else {
			fmt.Println("Something went wrong.")
		}
		c.Logger.Errorf("Error sending bank card: %v", err)
	} else if err != nil {
		c.Logger.Errorf("Error sending bank card: %v", err)
	}
}

// saveLoginAndPasswordData asks user for a login and password and sends them to a server.
func (c *CommandsController) saveLoginAndPasswordData() {
	//read data
	var login, password string

	fmt.Println("Enter login:")
	_, err := fmt.Scan(&login)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		return
	}

	fmt.Println("Enter password:")
	_, err = fmt.Scan(&password)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		return
	}

	//send data
	err = c.HTTPRequester.SendLoginAndPassword(login, password)
	errWithHttpCode := &gophKeeperErrors.ErrWithHTTPCode{}
	if errors.As(err, &errWithHttpCode) {
		if errWithHttpCode.Code() == http.StatusConflict {
			fmt.Println("This login and password have already been saved.")
		} else {
			fmt.Println("Something went wrong.")
		}
		c.Logger.Errorf("Error sending login and password : %v", err)
	} else if err != nil {
		c.Logger.Errorf("Error sending lgin and password: %v", err)
	}
}

// saveTextData - asks user for a text and sends it to a backend.
func (c *CommandsController) saveTextData() {
	//read data
	var textName, text string

	fmt.Println("Enter a text name:")
	_, err := fmt.Scan(&textName)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		fmt.Println("Incorrect input.")
	}

	//read text with spaces
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a text: ")
	text, err = reader.ReadString('\n')
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		fmt.Println("Incorrect input.")
		return
	}
	text = strings.TrimSpace(text)

	//send data
	err = c.HTTPRequester.SendText(textName, text)
	errWithHttpCode := &gophKeeperErrors.ErrWithHTTPCode{}
	if errors.As(err, &errWithHttpCode) {
		if errWithHttpCode.Code() == http.StatusConflict {
			fmt.Println("This text has already been saved.")
		} else {
			fmt.Println("Something went wrong.")
		}
		c.Logger.Errorf("Error sending text data: %v", err)
	} else if err != nil {
		c.Logger.Errorf("Error sending text data: %v", err)
	}
}

// saveBinFileData - asks user to enter a path to a file and sends this file to a server using gRPC stream.
func (c *CommandsController) saveBinFileData() {
	//read data
	var path string
	fmt.Println("Enter a path to a file:")
	_, err := fmt.Scan(&path)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		return
	}
	splits := strings.Split(path, string(os.PathSeparator))

	//ask api
	err = c.GRPCRequester.SendBinFile(path, splits[len(splits)-1])
	if err != nil {
		c.Logger.Errorf("Error sending the file: %v", err)
		fmt.Println("Something went wrong.")
		return
	}
	fmt.Printf("File has been saved. Use it`s name '%s' to get it (not a full path).\n", splits[len(splits)-1])
}

// getBankCardData retrieves bank card data from the server.
func (c *CommandsController) getBankCardData() {
	var lastFourDigits string
	fmt.Println("Enter a last four digits:")
	_, err := fmt.Scan(&lastFourDigits)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		return
	}

	//ask API
	data, err := c.HTTPRequester.GetBankCard(lastFourDigits)
	errWithHttpCode := &gophKeeperErrors.ErrWithHTTPCode{}
	if errors.As(err, &errWithHttpCode) {
		if errWithHttpCode.Code() == http.StatusNoContent {
			fmt.Println("This card hasn't been saved.")
		} else {
			fmt.Println("Something went wrong.")
		}
		c.Logger.Errorf("Error retrieving bank card: %v", err)
		return
	} else if err != nil {
		c.Logger.Errorf("Error retrieving bank card: %v", err)
		return
	}
	fmt.Printf("Bank Card Data:\n PAN: %v\n Owner last name: %v\n Owner first name: %v\n Expires at: %v\n", data.PAN, data.OwnerLastname, data.OwnerFirstname, data.ExpiresAt)
}

// getLoginAndPasswordData retrieves login and password data from the server.
func (c *CommandsController) getLoginAndPasswordData() {
	var login string
	fmt.Println("Enter a login:")
	_, err := fmt.Scan(&login)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		return
	}

	//ask api
	password, err := c.HTTPRequester.GetLoginAndPassword(login)
	errWithHttpCode := &gophKeeperErrors.ErrWithHTTPCode{}
	if errors.As(err, &errWithHttpCode) {
		if errWithHttpCode.Code() == http.StatusNoContent {
			fmt.Println("This pair of login and password hasn't been saved.")
		} else {
			fmt.Println("Something went wrong.")
		}
		c.Logger.Errorf("Error retrieving login and password: %v", err)
		return
	} else if err != nil {
		c.Logger.Errorf("Error retrieving login and password: %v", err)
		return
	}
	fmt.Printf(" Password is: %v\n", password)
}

// getTextData retrieves text data from the server.
func (c *CommandsController) getTextData() {
	var textName string
	fmt.Println("Enter a text name:")
	_, err := fmt.Scan(&textName)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		return
	}

	//ask api
	data, err := c.HTTPRequester.GetText(textName)
	errWithHttpCode := &gophKeeperErrors.ErrWithHTTPCode{}
	if errors.As(err, &errWithHttpCode) {
		if errWithHttpCode.Code() == http.StatusNoContent {
			fmt.Println("This text data hasn't been saved.")
		} else {
			fmt.Println("Something went wrong.")
		}
		c.Logger.Errorf("Error retrieving text data: %v", err)
		return
	} else if err != nil {
		c.Logger.Errorf("Error retrieving text data: %v", err)
		return
	}
	fmt.Printf("Text Data:\n%v\n", data)
}

// getBinFileData retrieves binary file data from the server.
func (c *CommandsController) getBinFileData() {
	var fileName, outputPath string
	fmt.Println("Enter a file name:")
	_, err := fmt.Scan(&fileName)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		return
	}
	fmt.Println("Enter an output path:")
	_, err = fmt.Scan(&outputPath)
	if err != nil {
		c.Logger.Errorf("Error reading input: %v", err)
		return
	}

	//ask api
	err = c.GRPCRequester.GetBinFile(fileName, outputPath)
	if err != nil {
		c.Logger.Errorf("Error retrieving binary file data: %v", err)
		fmt.Println("Failed to retrieve binary file data.")
		return
	}
	fmt.Printf("Bin data was received successfully: %v\n", outputPath)
}

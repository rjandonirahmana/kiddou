package usecase

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"kiddou/domain"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func SendEMail(ctx context.Context, templateStr string, user *domain.Users) error {
	client, err := Auth(ctx)
	if err != nil {
		return err
	}
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	t, err := template.ParseFiles(templateStr)
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, user); err != nil {
		return err
	}

	body := buffer.String()
	From := "From: rjandonirahmana@gmail.com" + "\n"
	to := "To: " + user.Email + "\n"
	cc := "Cc: " + "\n"
	bcc := "Bcc: " + "\n"
	replyTo := "Reply-to: " + "\n"
	subject := "Subject: Kiddou registrations" + "\n"
	var message gmail.Message
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"

	msg := []byte(to + cc + bcc + From + replyTo + subject + mime + "\n" + body)
	message.Raw = base64.URLEncoding.EncodeToString(msg)

	resultMSg, err := srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("error disni nih wakti=u send email %s", err.Error())
		return err
	}

	log.Println(resultMSg)
	return nil
}

func Auth(ctx context.Context) (*http.Client, error) {
	b, err := ioutil.ReadFile("." + "/credentials.json")
	if err != nil {
		log.Printf("failed to get credentials %v", err.Error())
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		log.Printf("failed to get credentials %v", err.Error())
		return nil, err
	}
	client := GetClient(config, ctx)

	return client, nil
}

func GetClient(config *oauth2.Config, ctx context.Context) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Printf("failed to get credentials %v", err.Error())
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(ctx, tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

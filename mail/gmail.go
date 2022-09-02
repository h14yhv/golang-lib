package mail

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/h14yhv/golang-lib/log"
)

type Gmail struct {
	name           string
	logger         log.Logger
	service        *gmail.Service
	credentialFile string
	tokenFile      string
}

func NewGmailService(credentialFile, tokenFile string) (Service, error) {
	handler := &Gmail{name: ModuleGmail}
	handler.logger, _ = log.New(handler.name, log.DebugLevel, true, os.Stdout)
	if credentialFile == "" {
		credentialFile = DefaultCredentialFile
	}
	handler.credentialFile = credentialFile
	if tokenFile == "" {
		tokenFile = DefaultTokenFile
	}
	handler.tokenFile = tokenFile
	service, err := handler.getService()
	if err != nil {
		return nil, err
	}
	handler.service = service
	// Success
	return handler, nil
}

func (s *Gmail) Send(email *Email) error {
	from := mail.Address{Name: email.Name, Address: email.From}
	receives := make([]string, 0)
	for _, to := range email.To {
		addr := mail.Address{Address: to}
		receives = append(receives, addr.String())
	}
	msg := fmt.Sprintf(""+
		"From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html;charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s", from.String(), strings.Join(receives, ", "), email.Subject, email.Message)
	_, err := s.service.Users.Messages.Send("me", &gmail.Message{Raw: base64.URLEncoding.EncodeToString([]byte(msg))}).Do()
	if err != nil {
		s.logger.Errorf("unable to send message (%s->%s), reason: %v", email.From, strings.Join(email.To, ","), err)
		return err
	}
	s.logger.Infof("send (%s->%s) success", email.From, strings.Join(email.To, ","))
	// Success
	return nil
}

func (s *Gmail) getClient(config *oauth2.Config) (*http.Client, error) {
	token, err := s.tokenFromFile(s.tokenFile)
	if err != nil {
		token, err = s.getTokenFromWeb(config)
		if err = s.saveToken(s.tokenFile, token); err != nil {
			return nil, err
		}
	}
	// Success
	return config.Client(context.Background(), token), nil
}

func (s *Gmail) getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	s.logger.Infof("Visit: %v\nEnter authorization code:", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		s.logger.Errorf("unable to read authorization code %v", err)
		return nil, err
	}
	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		s.logger.Errorf("unable to retrieve token from web %v", err)
		return nil, err
	}
	// Success
	return token, nil
}

func (s *Gmail) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	// Success
	return token, err
}

func (s *Gmail) saveToken(path string, token *oauth2.Token) error {
	s.logger.Infof("saving credential file to: %s", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		s.logger.Errorf("unable to cache oauth token: %v", err)
		return err
	}
	defer f.Close()
	if err = json.NewEncoder(f).Encode(token); err != nil {
		s.logger.Errorf("unable to save oauth token: %v", err)
		return err
	}
	// Success
	return nil
}

func (s *Gmail) getService() (*gmail.Service, error) {
	ctx := context.Background()
	b, err := ioutil.ReadFile(s.credentialFile)
	if err != nil {
		s.logger.Errorf("unable to read %s file, reason: %v", s.credentialFile, err)
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailComposeScope, gmail.GmailSendScope, gmail.GmailModifyScope)
	if err != nil {
		s.logger.Errorf("unable to init config from json, reason: %v", err)
		return nil, err
	}
	client, err := s.getClient(config)
	if err != nil {
		return nil, err
	}
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		s.logger.Errorf("unable to retrieve Gmail client, reason: %v", err)
		return nil, err
	}
	// Success
	return srv, nil
}

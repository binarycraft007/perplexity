package perplexity

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const AskVersion = "Ask/1.0.16/260016"
const AndroidVersion = "13"
const SdkVersion = "33"
const DefaultUserAgent = AskVersion + " " + "(Android; Version 13;" +
	" " + "Xiaomi M2011K2G/TQ1A.230205.002) SDK " + SdkVersion

type Session struct {
	Sid          string
	Wss          *websocket.Conn
	Client       *http.Client
	FrontendUUID uuid.UUID
	Token        string
	DeviceID     string
	UserAgent    string
	BaseApiURI   *url.URL
}

func NewSession() (*Session, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	rawURL := "https://www.perplexity.ai/socket.io/"
	baseApiURI, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	session := Session{
		Client: &http.Client{
			Jar: jar,
		},
		FrontendUUID: uuid.New(),
		Token:        getToken(),
		DeviceID:     getDeviceID(),
		UserAgent:    DefaultUserAgent,
		BaseApiURI:   baseApiURI,
	}

	params := url.Values{}
	params.Add("EIO", "4")
	params.Add("transport", "polling")
	baseApiURI.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseApiURI.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-encoding", "gzip")
	req.Header.Add("user-agent", session.UserAgent)
	req.Header.Add("x-app.version", "1.0.16")
	req.Header.Add("x-client-version", "1.0.16")
	req.Header.Add("x-client-name", "Perplexity-Android")
	req.Header.Add("x-app-apiclient", "android")
	req.Header.Add("x-app-apiversion", "1.0")

	resp, err := session.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {

		return nil, err
	}
	defer reader.Close()

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var result GetSidResponse
	if err := json.Unmarshal(bodyBytes[1:], &result); err != nil {
		return nil, err
	}

	session.Sid = result.Sid

	return &session, nil
}

func (s *Session) Check() error {
	params := url.Values{}
	params.Add("EIO", "4")
	params.Add("transport", "polling")
	params.Add("sid", s.Sid)
	s.BaseApiURI.RawQuery = params.Encode()

	body := bytes.NewBufferString("40")
	req, err := http.NewRequest("POST", s.BaseApiURI.String(), body)
	if err != nil {
		return err
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-encoding", "gzip")
	req.Header.Add("content-type", "text/plain;charset=UTF-8")
	req.Header.Add("user-agent", s.UserAgent)
	req.Header.Add("x-app.version", "1.0.16")
	req.Header.Add("x-client-version", "1.0.16")
	req.Header.Add("x-client-name", "Perplexity-Android")
	req.Header.Add("x-app-apiclient", "android")
	req.Header.Add("x-app-apiversion", "1.0")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {

		return err
	}
	defer reader.Close()

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if string(bodyBytes) != "OK" {
		return errors.New("Session Check Failed")
	}

	return nil
}

func (s *Session) GetSid() error {
	params := url.Values{}
	params.Add("EIO", "4")
	params.Add("transport", "polling")
	params.Add("sid", s.Sid)
	s.BaseApiURI.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", s.BaseApiURI.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-encoding", "gzip")
	req.Header.Add("user-agent", s.UserAgent)
	req.Header.Add("x-app.version", "1.0.16")
	req.Header.Add("x-client-version", "1.0.16")
	req.Header.Add("x-client-name", "Perplexity-Android")
	req.Header.Add("x-app-apiclient", "android")
	req.Header.Add("x-app-apiversion", "1.0")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {

		return err
	}
	defer reader.Close()

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if !strings.Contains(string(bodyBytes), "40{\"sid\":") {
		return errors.New("Get Sid Failed")
	}

	return nil
}

func (s *Session) InitWss() error {
	s.BaseApiURI.Scheme = "wss"
	defer func(u *url.URL) { u.Scheme = "https" }(s.BaseApiURI)

	cookie, err := concatenateCookies(s.Client.Jar)
	if err != nil {
		return err
	}
	header := http.Header{}
	header.Add("accept", "*/*")
	header.Add("accept-encoding", "gzip")
	header.Add("user-agent", s.UserAgent)
	header.Add("x-app.version", "1.0.16")
	header.Add("x-client-version", "1.0.16")
	header.Add("x-client-name", "Perplexity-Android")
	header.Add("x-app-apiclient", "android")
	header.Add("x-app-apiversion", "1.0")
	header.Add("Cookie", *cookie)

	params := url.Values{}
	params.Add("EIO", "4")
	params.Add("transport", "websocket")
	params.Add("sid", s.Sid)
	s.BaseApiURI.RawQuery = params.Encode()

	wssURI := s.BaseApiURI.String()
	conn, resp, err := websocket.DefaultDialer.Dial(wssURI, header)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		log.Println("Error switching protocols:", resp.StatusCode)
		return err
	}

	s.Wss = conn

	conn.WriteMessage(websocket.TextMessage, []byte("2probe"))
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			continue
		}

		if string(message) == "6" {
			break
		}

		switch string(message) {
		case "2":
			conn.WriteMessage(websocket.TextMessage, []byte("3"))
		case "3probe":
			conn.WriteMessage(websocket.TextMessage, []byte("5"))
		}
	}

	return nil
}

func concatenateCookies(jar http.CookieJar) (*string, error) {
	url, err := url.Parse("https://www.perplexity.ai")
	if err != nil {
		return nil, err
	}
	cookies := jar.Cookies(url)

	var cookieStrings []string
	for _, cookie := range cookies {
		cookieStrings = append(cookieStrings, cookie.String())
	}

	cookieString := strings.Join(cookieStrings, "; ")

	return &cookieString, nil
}

func (s *Session) Ask(question string) error {
	code := 421

	askReq := AskRequest{
		Source:                "android",
		Version:               "1.0",
		Token:                 s.Token,
		FrontendUUID:          s.FrontendUUID.String(),
		UseInhouseModel:       false,
		ConversationalEnabled: true,
		AndroidDeviceID:       s.DeviceID,
	}

	marshalled, _ := json.Marshal(askReq)

	params := []string{"perplexity_ask", question, string(marshalled)}

	q := fmt.Sprintf("%d[%q,%q,%v]", code, params[0], params[1], params[2])

	err := s.Wss.WriteMessage(websocket.TextMessage, []byte(q))
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) ReadAnswer() (*AnswerDetails, error) {
	_, message, err := s.Wss.ReadMessage()
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(string(message), "42[\"query_progress\"") {
		return s.ReadAnswer()
	}

	if !strings.HasPrefix(string(message), "42[\"query_answered\"") {
		return nil, errors.New("No answer found")
	}

	respBytes := message[len("42[\"query_answered\",") : len(message)-1]

	var result AskResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, err
	}

	var answer AnswerDetails
	if err := json.Unmarshal([]byte(result.Text), &answer); err != nil {
		return nil, err
	}

	return &answer, nil
}

func getToken() string {
	bytes := make([]byte, 3)
	_, _ = rand.Read(bytes)

	// Encode the bytes in hexadecimal format
	return hex.EncodeToString(bytes)
}

func getDeviceID() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)

	// Set the first byte to 0xdd
	bytes[0] = 0xdd

	// Encode the byte slice as a hexadecimal string
	return hex.EncodeToString(bytes)
}
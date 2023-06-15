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

const (
	AppVersion    = "1.0.23"
	ClientVersion = "1.0.23"
	ApiVersion    = "1.0"
	AskVersion    = "Ask/1.0.23/260023"

	// Android related info
	AndroidVersion = 13
	SdkVersion     = 33
	VendorName     = "Xiaomi"
	Model          = "M2011K2G"
	BuildID        = "TQ1A.230205.002"
)

type SearchFocus string

const (
	Internet     SearchFocus = "internet"
	Writing      SearchFocus = "writing"
	Academic     SearchFocus = "scholar"
	WolframAlpha SearchFocus = "wolfram"
	YouTube      SearchFocus = "youtube"
	Reddit       SearchFocus = "reddit"
)

type SearchSource string

const (
	Android SearchSource = "android"
	Default SearchSource = "default"
)

type SearchMode string

const (
	Concise SearchMode = "concise"
	Copilot SearchMode = "copilot"
)

type Session struct {
	Sid               string
	Wss               *websocket.Conn
	Client            *http.Client
	FrontendUUID      uuid.UUID
	FrontendSessionID uuid.UUID
	Token             string
	DeviceID          string
	UserAgent         string
	BaseApiURI        *url.URL
	AskSeqNum         int
	LastBackendUUID   string
	ReadWriteToken    string
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
		FrontendUUID:      uuid.New(),
		FrontendSessionID: uuid.New(),
		Token:             getToken(),
		DeviceID:          getDeviceID(),
		UserAgent:         getUserAgent(),
		BaseApiURI:        baseApiURI,
		AskSeqNum:         1,
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
	req.Header.Add("x-app.version", AppVersion)
	req.Header.Add("x-client-version", ClientVersion)
	req.Header.Add("x-client-name", "Perplexity-Android")
	req.Header.Add("x-app-apiclient", "android")
	req.Header.Add("x-app-apiversion", ApiVersion)

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
	req.Header.Add("x-app.version", AppVersion)
	req.Header.Add("x-client-version", ClientVersion)
	req.Header.Add("x-client-name", "Perplexity-Android")
	req.Header.Add("x-app-apiclient", "android")
	req.Header.Add("x-app-apiversion", ApiVersion)

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
	req.Header.Add("x-app.version", AppVersion)
	req.Header.Add("x-client-version", ClientVersion)
	req.Header.Add("x-client-name", "Perplexity-Android")
	req.Header.Add("x-app-apiclient", "android")
	req.Header.Add("x-app-apiversion", ApiVersion)

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
	defer resp.Body.Close()

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

		if string(message) == "3probe" {
			conn.WriteMessage(websocket.TextMessage, []byte("5"))
		}
	}

	return nil
}

func (s *Session) Close() {
	s.Wss.WriteMessage(websocket.CloseMessage, []byte{})
	s.Wss.Close()
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
	code := s.AskSeqNum
	defer func(s *Session) { s.AskSeqNum += 1 }(s)

	askReq := AskRequest{
		Source:                Default,
		Token:                 s.Token,
		FrontendUUID:          s.FrontendUUID.String(),
		ConversationalEnabled: true,
		Language:              "en-US",
		Timezone:              "Asia/Shanghai",
		LastBackendUUID:       s.LastBackendUUID,
		FrontendSessionID:     s.FrontendSessionID.String(),
		ReadWriteToken:        s.ReadWriteToken,
		SearchFocus:           Writing,
		Gpt4:                  false,
		Mode:                  Concise,
	}

	marshalled, err := json.Marshal(askReq)
	if err != nil {
		return err
	}

	q := fmt.Sprintf(
		"%d%d[%q,%q,%v]",
		42, code, "perplexity_ask", question, string(marshalled),
	)

	err = s.Wss.WriteMessage(websocket.TextMessage, []byte(q))
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

	if string(message) == "2" {
		s.Wss.WriteMessage(websocket.TextMessage, []byte("3"))
		return s.ReadAnswer()
	}

	if !strings.HasPrefix(string(message), "42[\"query_answered\"") {
		return nil, errors.New("No answer found")
	}

	var result AskResponse
	if err := parseMessage(message, &result); err != nil {
		return nil, err
	}
	s.LastBackendUUID = result.BackendUUID
	s.ReadWriteToken = result.ReadWriteToken

	var answer AnswerDetails
	if err := json.Unmarshal([]byte(result.Text), &answer); err != nil {
		return nil, err
	}

	_, _, _ = s.Wss.ReadMessage()

	return &answer, nil
}

func parseMessage(message []byte, v any) error {
	start := strings.Index(string(message), ",") + 1
	respBytes := message[start : len(message)-1]

	if err := json.Unmarshal(respBytes, v); err != nil {
		return err
	}

	return nil
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

	// Encode the byte slice as a hexadecimal string
	return hex.EncodeToString(bytes)
}

func getUserAgent() string {
	return fmt.Sprintf(
		"%s (Android; Version %d; %s %s/%s) SDK %d",
		AskVersion,
		AndroidVersion,
		VendorName,
		Model,
		BuildID,
		SdkVersion,
	)
}

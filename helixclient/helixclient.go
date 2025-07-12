package helixclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
	"ttv-statistics/constants"
)

const (
	getHelixAuthEndpoint string = "https://id.twitch.tv/oauth2/token" // currently hardcoded constant, this could be parsed as a CLI flag for flexibility

	helixLoginURLParam  string = "login"
	helixUserIDURLParam string = "user_id"
	helixFirstURLParam  string = "first"

	HelixUsersEndpoint  string = "/users"
	HelixVideosEndpoint string = "/videos"

	authorisationHeaderKey string = "Authorization"
	clientIDHeaderKey      string = "Client-ID"
)

var (
	ClientID     string
	ClientSecret string
	HelixHost    string

	helixAccessToken string

	helixClient = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     time.Minute * 2,
		},
	}
)

func InitHelixClientAuth(ctx context.Context) error {
	response, err := getHelixAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("message=%q error=%v", "failed to get authorisation for helix client", err)
	}

	helixAccessToken = response.AccessToken
	return nil
}

func GetUserData(ctx context.Context, userName string) (responseBody UsersResponseBody, err error) {

	endpoint, err := url.Parse(HelixHost)
	if err != nil {
		return responseBody, err
	}

	endpoint.Path = path.Join(endpoint.Path, HelixUsersEndpoint)

	headers := generateHeaders(false)

	queryParams := map[string]string{
		helixLoginURLParam: userName,
	}

	return executeRequest[UsersResponseBody](ctx, http.MethodGet, endpoint, queryParams, headers, nil)
}

func GetStreamerFirstNVideoStatistics(ctx context.Context, userID string, n int) (responseBody VideosResponseBody, err error) {
	endpoint, err := url.Parse(HelixHost)
	if err != nil {
		return responseBody, err
	}

	endpoint.Path = path.Join(endpoint.Path, HelixVideosEndpoint)

	headers := generateHeaders(false)

	queryParams := map[string]string{
		helixUserIDURLParam: userID,
		helixFirstURLParam:  strconv.Itoa(n),
	}

	return executeRequest[VideosResponseBody](ctx, http.MethodGet, endpoint, queryParams, headers, nil)
}

func generateHeaders(setContentType bool) map[string]string {
	return map[string]string{
		clientIDHeaderKey:      ClientID,
		authorisationHeaderKey: fmt.Sprintf("Bearer %s", helixAccessToken),
	}
}

func getHelixAccessToken(ctx context.Context) (responseBody TokenResponse, err error) {

	endpoint, err := url.Parse(getHelixAuthEndpoint)
	if err != nil {
		return responseBody, err
	}

	headers := map[string]string{
		constants.ContentTypeHeaderKey: constants.ContentTypeFormURLEndcoded,
	}

	params := url.Values{}
	params.Set("client_id", ClientID)
	params.Set("client_secret", ClientSecret)
	params.Set("grant_type", "client_credentials")

	body := strings.NewReader(params.Encode())

	return executeRequest[TokenResponse](ctx, http.MethodPost, endpoint, nil, headers, body)
}

// Template functions 😈
func executeRequest[T ClientResponseModels](
	ctx context.Context, method string, endpoint *url.URL, queryParams, headers map[string]string, body io.Reader,
) (responseBody T, err error) {

	query := endpoint.Query()

	for paramName, queryValue := range queryParams {
		query.Set(paramName, queryValue)
	}

	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), body)
	if err != nil {
		return responseBody,
			fmt.Errorf("message=%q url=%q error=%v", "failed to create request", endpoint.String(), err)
	}

	for headerName, headerValue := range headers {
		req.Header.Set(headerName, headerValue)
	}

	response, err := helixClient.Do(req)
	if err != nil {
		return responseBody, fmt.Errorf("message=%s url=%s error=%v", "failed to execute http request", endpoint.String(), err)
	}

	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Printf("message=%s url=%s error=%s\n ", "warning: failed to close response body", endpoint.String(), closeErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return responseBody,
			fmt.Errorf("message=%s url=%s status_code=%d", "received unexpected status code", endpoint.String(), response.StatusCode)
	}

	responseBuffer, err := io.ReadAll(response.Body)
	if err != nil {
		return responseBody,
			fmt.Errorf("message=%s url=%s error=%s", "failed to read response body", endpoint.String(), err)
	}

	err = json.Unmarshal(responseBuffer, &responseBody)

	return responseBody, err
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type UsersResponseBody struct {
	Data []struct {
		ID              string `json:"id"`
		Login           string `json:"login"`
		DisplayName     string `json:"display_name"`
		ProfileImageURL string `json:"profile_image_url"`
	} `json:"data"`
}

type VideoInfo struct {
	ID            string    `json:"id"`
	StreamID      *string   `json:"stream_id"`
	UserID        string    `json:"user_id"`
	UserLogin     string    `json:"user_login"`
	UserName      string    `json:"user_name"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	PublishedAt   time.Time `json:"published_at"`
	URL           string    `json:"url"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	Viewable      string    `json:"viewable"`
	ViewCount     int       `json:"view_count"`
	Language      string    `json:"language"`
	Type          string    `json:"type"`
	Duration      string    `json:"duration"`
	MutedSegments []struct {
		Duration int `json:"duration"`
		Offset   int `json:"offset"`
	} `json:"muted_segments,omitempty"`
}

type VideosResponseBody struct {
	Data       []VideoInfo       `json:"data"`
	Pagination map[string]string `json:"pagination"`
}

type ClientResponseModels interface {
	TokenResponse |
		UsersResponseBody |
		VideosResponseBody
}

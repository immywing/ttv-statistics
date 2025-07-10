package helixclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	getHelixAuthEndpoint string = "https://id.twitch.tv/oauth2/token"

	helixUsersEndpoint string = "https://api.twitch.tv/helix/users"

	authorisationHeaderKey     string = "Authorization"
	clientIDHeaderKey          string = "Client-ID"
	contentTypeHeaderKey       string = "Content-Type"
	contentTypeFormURLEndcoded string = "application/x-www-form-urlencoded"
	contentTypeApplicationJson string = "application/json"
)

var (
	ClientID     string
	ClientSecret string

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

type tokenResponse struct {
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

type clientResponseModels interface {
	tokenResponse |
		UsersResponseBody
}

func InitHelixClientAuth(ctx context.Context) error {
	response, err := getHelixAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("message=%q error=%v", "failed to get authorisation for helix client", err)
	}

	helixAccessToken = response.AccessToken
	return nil
}

func GetUserData(ctx context.Context, userName string) (responseBody UsersResponseBody, err error) {

	endpoint, err := url.Parse(helixUsersEndpoint)
	if err != nil {
		return responseBody, err
	}

	headers := generateHeaders(false)

	queryParams := map[string]string{
		"login": userName,
	}

	return executeRequest[UsersResponseBody](ctx, http.MethodGet, endpoint, queryParams, headers, nil)
}

func generateHeaders(setContentType bool) map[string]string {

	headers := map[string]string{
		clientIDHeaderKey:      ClientID,
		authorisationHeaderKey: fmt.Sprintf("Bearer %s", helixAccessToken),
	}

	if setContentType {
		headers[contentTypeHeaderKey] = contentTypeApplicationJson
	}

	return headers
}

func getHelixAccessToken(ctx context.Context) (responseBody tokenResponse, err error) {

	endpoint, err := url.Parse(getHelixAuthEndpoint)
	if err != nil {
		return responseBody, err
	}

	headers := map[string]string{
		contentTypeHeaderKey: contentTypeFormURLEndcoded,
	}

	params := url.Values{}
	params.Set("client_id", ClientID)
	params.Set("client_secret", ClientSecret)
	params.Set("grant_type", "client_credentials")

	body := strings.NewReader(params.Encode())

	return executeRequest[tokenResponse](ctx, http.MethodPost, endpoint, nil, headers, body)
}

// Template functions 😈
func executeRequest[T clientResponseModels](
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
		return responseBody, fmt.Errorf("message=%q url=%q error=%v", "failed to execute http request", endpoint.String(), err)
	}

	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Printf("message=%q url=%q error=%v\n ", "warning: failed to close response body", endpoint.String(), closeErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return responseBody,
			fmt.Errorf("message=%q url=%q status_code=%d", "received unexpected status code", endpoint.String(), response.StatusCode)
	}

	responseBuffer, err := io.ReadAll(response.Body)
	if err != nil {
		return responseBody,
			fmt.Errorf("message=%q url=%q error=%v", "failed to read response body", endpoint.String(), err)
	}

	err = json.Unmarshal(responseBuffer, &responseBody)

	return responseBody, err
}

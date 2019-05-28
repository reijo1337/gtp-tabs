package clients

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AuthClientInterface -
type AuthClientInterface interface {
	Register(user *User) (*Tokens, error)
	RegisterVk(user *VkUser) (*Tokens, error)
	GenToken(user *User) (*Tokens, error)
	GenTokenVk(user *VkUser) (*Tokens, error)
	RefreshToken(refreshToken string) (*Tokens, error)
}

// AuthClient -
type AuthClient struct {
	url string
}

// MakeAuthClient -
func MakeAuthClient(host string, port string) AuthClientInterface {
	return &AuthClient{
		url: fmt.Sprintf("http://%s:%s", host, port),
	}
}

func (ac *AuthClient) Register(user *User) (*Tokens, error) {
	jsonStr, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("parsing user: %v", err)
	}
	url := fmt.Sprintf("%s/register", ac.url)
	return ac.registerRequest(jsonStr, url)
}

func (ac *AuthClient) RegisterVk(user *VkUser) (*Tokens, error) {
	jsonStr, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("parsing user: %v", err)
	}
	url := fmt.Sprintf("%s/register/vk", ac.url)
	return ac.registerRequest(jsonStr, url)
}

func (ac *AuthClient) registerRequest(jsonStr []byte, url string) (*Tokens, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		errMsg := &ErrorResponse{}
		json.Unmarshal(body, errMsg)
		return nil, fmt.Errorf("status no ok: %s", errMsg.Error)
	}
	tokens := &Tokens{}
	if err = json.Unmarshal(body, tokens); err != nil {
		return nil, fmt.Errorf("parsing response: %v", err)
	}
	return tokens, nil
}

func (ac *AuthClient) GenToken(user *User) (*Tokens, error) {
	jsonStr, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("parsing user: %v", err)
	}
	return ac.registerRequest(jsonStr, ac.url)
}

func (ac *AuthClient) GenTokenVk(user *VkUser) (*Tokens, error) {
	jsonStr, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("parsing user: %v", err)
	}
	url := fmt.Sprintf("%s/vk", ac.url)
	return ac.registerRequest(jsonStr, url)
}

func (ac *AuthClient) RefreshToken(refreshToken string) (*Tokens, error) {
	url := fmt.Sprintf("%s?refresh_token=%s", ac.url, refreshToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("send refresh request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		errMsg := &ErrorResponse{}
		json.Unmarshal(body, errMsg)
		return nil, fmt.Errorf("status no ok: %s", errMsg.Error)
	}
	tokens := &Tokens{}
	if err = json.Unmarshal(body, tokens); err != nil {
		return nil, fmt.Errorf("parsing response: %v", err)
	}
	return tokens, nil
}

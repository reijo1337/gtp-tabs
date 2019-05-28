package clients

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ProfileClientInterface -
type ProfileClientInterface interface {
	GetProfile(userID int) (*ProfileInfo, error)
	GetProfileByAcc(userID int) (*ProfileInfo, error)
	SetProfile(user *ProfileInfo) error
}

// ProfileClient -
type ProfileClient struct {
	url string
}

// MakeProfileClient -
func MakeProfileClient(host string, port string) ProfileClientInterface {
	return &ProfileClient{
		url: fmt.Sprintf("http://%s:%s", host, port),
	}
}

func (pc *ProfileClient) GetProfileByAcc(userID int) (*ProfileInfo, error) {
	url := fmt.Sprintf("%s/profile/user/%d", pc.url, userID)
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
		return nil, fmt.Errorf("status not ok: %s", errMsg.Error)
	}
	pi := &ProfileInfo{}
	if err = json.Unmarshal(body, pi); err != nil {
		return nil, fmt.Errorf("parsing response: %v", err)
	}
	return pi, nil
}

func (pc *ProfileClient) GetProfile(userID int) (*ProfileInfo, error) {
	url := fmt.Sprintf("%s/profile/%d", pc.url, userID)
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
		return nil, fmt.Errorf("status not ok: %s", errMsg.Error)
	}
	pi := &ProfileInfo{}
	if err = json.Unmarshal(body, pi); err != nil {
		return nil, fmt.Errorf("parsing response: %v", err)
	}
	return pi, nil
}

func (pc *ProfileClient) SetProfile(user *ProfileInfo) error {
	jsonStr, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("parsing user: %v", err)
	}
	url := fmt.Sprintf("%s/register", pc.url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		errMsg := &ErrorResponse{}
		json.Unmarshal(body, errMsg)
		return fmt.Errorf("status no ok: %s", errMsg.Error)
	}
	if err = json.Unmarshal(body, user); err != nil {
		return fmt.Errorf("parsing response: %v", err)
	}
	return nil
}

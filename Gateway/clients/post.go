package clients

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type PostClientInterface interface {
	GetPost(tabID int) (*Post, error)
	SetPost(post *Post) error
	UpdateRating(postID int, rating int) error
	AddComment(postID int, comment *Comment) error
}

type PostClient struct {
	url string
}

func MakePostClient(url string) PostClientInterface {
	return &PostClient{
		url: url,
	}
}

func (pc *PostClient) GetPost(tabID int) (*Post, error) {
	url := fmt.Sprintf("%s/post/%d", pc.url, tabID)
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
	pi := &Post{}
	if err = json.Unmarshal(body, pi); err != nil {
		return nil, fmt.Errorf("parsing response: %v", err)
	}
	return pi, nil
}

func (pc *PostClient) SetPost(post *Post) error {
	jsonStr, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("parsing post: %v", err)
	}
	url := fmt.Sprintf("%s/post", pc.url)
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
	if err = json.Unmarshal(body, post); err != nil {
		return fmt.Errorf("parsing response: %v", err)
	}
	return nil
}

func (pc *PostClient) UpdateRating(postID int, rating int) error {
	reqStruct := &UpdateRatingRequest{
		PostID: postID,
		Rating: rating,
	}
	jsonStr, err := json.Marshal(reqStruct)
	if err != nil {
		return fmt.Errorf("parsing post: %v", err)
	}
	url := fmt.Sprintf("%s/post", pc.url)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
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
	return nil
}

func (pc *PostClient) AddComment(postID int, comment *Comment) error {
	jsonStr, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("parsing post: %v", err)
	}
	url := fmt.Sprintf("%s/comment/%s", pc.url, strconv.Itoa(postID))
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
	if err = json.Unmarshal(body, comment); err != nil {
		return fmt.Errorf("parsing response: %v", err)
	}
	return nil
}

package clients

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// StorageClientInterface интерфейс для работы с сервисом хранилища
type StorageClientInterface interface {
	GetAuthorsByLetter(letter string) ([]MusiciansWithCount, error)
	GetAuthorsByName(search string) ([]MusiciansWithCount, error)
	FindTabsByName(search string) ([]TabWithSize, error)
	GetAuthorsByCategory(name string) ([]MusiciansWithCount, error)
	UploadFile(upload *FileUploadRequest) (*TabInfo, error)
	DownloadFile(name string) (FileDownloadResponse, error)
	GetTab(tabID int) (*TabInfo, error)
}

type StorageClient struct {
	url string
}

func MakeStorageClient(url string) StorageClientInterface {
	return &StorageClient{
		url: url,
	}
}

func (sc *StorageClient) GetTab(tabID int) (*TabInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/tab/%d", sc.url, tabID))
	if err != nil {
		log.Println("Can't get tab from service", err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read body", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errResp.Error)
	}

	ret := &TabInfo{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (sc *StorageClient) GetAuthorsByLetter(letter string) ([]MusiciansWithCount, error) {
	log.Println("GetAuthorsByLetter", letter)
	resp, err := http.Get(fmt.Sprintf("%s/letter/musicians/%s", sc.url, letter))
	if err != nil {
		log.Println("Can't get musicians from service", err)
		return nil, err
	}
	return sc.returnMusicians(resp)
}

func (sc *StorageClient) GetAuthorsByName(search string) ([]MusiciansWithCount, error) {
	log.Println("GetAuthorsByName", search)
	resp, err := http.Get(fmt.Sprintf("%s/musicians/%s", sc.url, search))
	if err != nil {
		log.Println("Can't get musicians from service", err)
		return nil, err
	}
	return sc.returnMusicians(resp)
}

func (sc *StorageClient) returnMusicians(resp *http.Response) ([]MusiciansWithCount, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read body", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errResp.Error)
	}

	ret := make([]MusiciansWithCount, 0)
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (sc *StorageClient) FindTabsByName(search string) ([]TabWithSize, error) {
	log.Println("GetAuthorsByName", search)
	resp, err := http.Get(fmt.Sprintf("%s/tabs/%s", sc.url, search))
	if err != nil {
		log.Println("Can't get musicians from service", err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read body", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errResp.Error)
	}

	ret := make([]TabWithSize, 0)
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (sc *StorageClient) GetAuthorsByCategory(name string) ([]MusiciansWithCount, error) {
	log.Println("GetAuthorsByName", name)
	resp, err := http.Get(fmt.Sprintf("%s/category/%s", sc.url, name))
	if err != nil {
		log.Println("Can't get musicians from service", err)
		return nil, err
	}
	return sc.returnMusicians(resp)
}

func (sc *StorageClient) UploadFile(upload *FileUploadRequest) (*TabInfo, error) {
	stringsUpload, err := json.Marshal(upload)
	if err != nil {
		return nil, fmt.Errorf("marshal upload request: %v", err)
	}
	reader := bytes.NewReader(stringsUpload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/file", sc.url), reader)
	if err != nil {
		return nil, fmt.Errorf("making upload request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending upload request: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return nil, fmt.Errorf("getting error upload response: %v", err)
		}
		return nil, fmt.Errorf("error upload response: %v", err)
	}
	ret := &TabInfo{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, fmt.Errorf("getting response: %v", err)
	}
	return ret, nil
}

func (sc *StorageClient) DownloadFile(name string) (FileDownloadResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/file?name=%s", sc.url, name))
	if err != nil {
		return FileDownloadResponse{}, fmt.Errorf("getting file: %v", err)
	}
	var ret FileDownloadResponse
	ret.FileContent = resp.Body
	ret.ContentLength = resp.ContentLength
	ret.ContentType = resp.Header.Get("Content-Type")
	ret.ExtraHeaders = map[string]string{
		"Content-Disposition": resp.Header.Get("Content-Disposition"),
	}
	return ret, nil
}

package clients

import (
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
	GetAuthorsByLetter(letter string) ([]*MusiciansWithCount, error)
	// GetAuthorsByName(search string) ([]*MusiciansWithCount, error)
	// FindTabsByName(search string) ([]*TabWithSize, error)
	// GetAuthorsByCategory(name string) ([]*MusiciansWithCount, error)
}

type StorageClient struct {
	url string
}

func MakeStorageClient(host string, port string) *StorageClient {
	return &StorageClient{
		url: fmt.Sprintf("http://%s:%s", host, port),
	}
}

func (sc *StorageClient) GetAuthorsByLetter(letter string) ([]*MusiciansWithCount, error) {
	log.Println("GetAuthorsByLetter", letter)
	resp, err := http.Get(fmt.Sprintf("%s/authors_by_letter?letter=%s", sc.url, letter))
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

	ret := make([]*MusiciansWithCount, 0)
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

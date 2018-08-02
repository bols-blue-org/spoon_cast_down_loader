package cast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
)

func Download(ctx *cli.Context) error {
	id := ctx.String("id")

	data, err := LoadMetaData(id)
	if err != nil {
		log.Println("Error, " + err.Error() + "!")
		return err
	}
	fmt.Printf("%v", data)
	err = data.Download()
	if err != nil {
		log.Println("Error, " + err.Error() + "!")
		return err
	}

	log.Println("Download, " + id + "!")
	return err
}

var client = &http.Client{Timeout: time.Duration(10) * time.Second}

func LoadMetaData(id string) (MetaData, error) {
	url := "https://jp-api.spooncast.net/casts/" + id + "/"
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)

	jsonBytes := ([]byte)(byteArray)
	data := new(MetaData)

	if err := json.Unmarshal(jsonBytes, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return *data, err
	}
	fmt.Printf("%v\n", data)
	return *data, nil
}

type MetaData struct {
	StatusCode int    `json:"status_code"`
	Detail     string `json:"detail"`
	Results    []struct {
		ID         int           `json:"id"`
		Title      string        `json:"title"`
		ImgURL     string        `json:"img_url"`
		VoiceURL   string        `json:"voice_url"`
		SpoonCount int           `json:"spoon_count"`
		Duration   float64       `json:"duration"`
		Reporters  []interface{} `json:"reporters"`
		IsDonated  bool          `json:"is_donated"`
		IsLike     bool          `json:"is_like"`
		Created    time.Time     `json:"created"`
		IsStorage  bool          `json:"is_storage"`
	} `json:"results"`
}

func (a MetaData) downloadFromURL(url string, retryMaxCount int) (*bytes.Buffer, error) {
	// Download from download URL.
	retryCounter := 0
	var resp *http.Response
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for retryMaxCount >= retryCounter {
		resp, err = client.Do(req)
		if a.isRequestSuccess(resp, err) {
			break
		}
		retryCounter++
	}

	if err != nil {
		return nil, err
	}

	defer func() {
		deferErr := resp.Body.Close()
		if deferErr != nil {
			log.Print(deferErr)
		}
	}()

	bodyBuf := &bytes.Buffer{}
	if _, err = io.Copy(bodyBuf, resp.Body); err != nil {
		log.Printf("failed to copy download file. err = %s", err)
	}
	return bodyBuf, err
}

func (a MetaData) Download() error {
	retry := 2
	buff, err := a.downloadFromURL(a.Results[0].VoiceURL, retry)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Base(a.Results[0].VoiceURL), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}

	defer func() {
		deferErr := file.Close()
		if deferErr != nil {
			log.Print(deferErr)
		}
	}()
	_, err = io.Copy(file, buff)
	if err != nil {
		return err
	}
	return nil
}

func (a MetaData) isRequestSuccess(resp *http.Response, err error) bool {
	if err != nil {
		if resp == nil {
			return false
		} else if resp.StatusCode >= 500 {
			return false
		} else {
			return true
		}
	}
	return true
}

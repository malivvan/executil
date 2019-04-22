package chromaprint

import (
	"encoding/json"
	"fmt"
	"github.com/satnamram/executil"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"runtime"
)

type AcoustIDRequest struct {
	Fingerprint string `json:"fingerprint"`
	Duration    int    `json:"duration"`
	ApiKey      string `json:"client"`
	Metadata    string `json:"meta"`
}

type Result struct {
	ID string `json:"id"`

	Recordings []struct {
		Artists []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"artists"`

		ReleaseGroups []struct {
			Type           string   `json:"type"`
			ID             string   `json:"id"`
			Title          string   `json:"title"`
			SecondaryTypes []string `json:"secondarytypes"`
		} `json:"releasegroups"`

		Duration float64 `json:"duration"`
		ID       string  `json:"id"`
		Title    string  `json:"title"`
	} `json:"recordings"`

	Score float64 `json:"score"`
}

type AcoustIDResponse struct {
	Results []Result `json:"results"`
	Status  string   `json:"status"`
}

func (a *AcoustIDRequest) Do() (*AcoustIDResponse, error) {
	client := http.Client{}
	postValues, err := a.PostValues()
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	response, err := client.PostForm("http://api.acoustid.org/v2/lookup", postValues)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	aidResp := &AcoustIDResponse{}
	err = json.Unmarshal(body, &aidResp)
	if err != nil {
		return nil, err
	}

	return aidResp, nil
}

func (a *AcoustIDRequest) PostValues() (url.Values, error) {
	query := fmt.Sprintf(
		"client=%s&duration=%d&meta=%s&fingerprint=%s",
		a.ApiKey,
		a.Duration,
		a.Metadata,
		a.Fingerprint)

	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}
	return values, nil
}

type AudioFingerprint struct {
	fingerprint string
	duration    int
}

func (afp *AudioFingerprint) Fingerprint() string {
	return afp.fingerprint
}

func (afp *AudioFingerprint) Duration() int {
	return afp.duration
}

func (afp *AudioFingerprint) Lookup(apikey string) (*AcoustIDResponse, error) {
	request := AcoustIDRequest{
		ApiKey:      apikey,
		Duration:    afp.duration,
		Fingerprint: afp.fingerprint,
		Metadata:    "recordings",
	}
	return request.Do()
}

func NewAudioFingerprint(filePath string) (*AudioFingerprint, error) {
	out, err := executil.Command("fpcalc", filePath).Output()
	if err != nil {
		return nil, err
	}

	fp := &AudioFingerprint{}
	lineSeperator := "\n"
	if runtime.GOOS == "windows" {
		lineSeperator = "\r\n"
	}
	outstrs := strings.Split(string(out), lineSeperator)

	for _, s := range outstrs {
		if strings.Index(s, "DURATION=") == 0 {
			ds := strings.Split(s, "=")[1]
			fp.duration, _ = strconv.Atoi(ds)
		} else if strings.Index(s, "FINGERPRINT=") == 0 {
			fp.fingerprint = strings.Split(s, "=")[1]
		}
	}

	return fp, nil
}

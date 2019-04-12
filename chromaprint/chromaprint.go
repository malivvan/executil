package chromaprint

import (
	"encoding/json"
	"fmt"
	"github.com/mikkyang/id3-go"
	"github.com/satnamram/shtool"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
)

func init() {
	shtool.Register(shtool.Package{
		Name:      "chromaprint",
		OS:        "linux",
		Arch:      "amd64",
		Commands:  []string{"fpcalc"},
		Download:  "https://github.com/acoustid/chromaprint/releases/download/v1.4.3/chromaprint-fpcalc-1.4.3-linux-x86_64.tar.gz",
		Checksum:  "a84425fccb43faa11b5bdc9d5b6101d6810b3b74876916191d42d31f7d73f4ce",
		Installer: shtool.ArchiveInstaller("chromaprint-fpcalc-1.4.3-linux-x86_64"),
		Activator: shtool.PathActivator,
	})
}

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

	err := shtool.Ensure("fpcalc")
	if err != nil {
		return nil, err
	}

	out, err := exec.Command("fpcalc", filePath).Output()
	if err != nil {
		return nil, err
	}

	fp := &AudioFingerprint{}
	outstrs := strings.Split(string(out), "\n")

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

func TagID3(title string, artist string, path string) error {
	mp3File, err := id3.Open(path)
	defer mp3File.Close()
	if err != nil {
		return err
	}
	mp3File.SetArtist(artist)
	mp3File.SetTitle(title)
	return nil
}

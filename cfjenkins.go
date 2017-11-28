package main

import (
	"fmt"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	log "github.com/sirupsen/logrus"
)

type (
	// Jenkins Job parameters
	JenkinksJobParams struct {
		Username string
		Token    string
		Host	string
		Job	string
		JobParams	string
	}

)

func main() {

	host:= os.Getenv("JENKINS_HOST")
	token:= os.Getenv("TOKEN")
	user:= os.Getenv("USER")
	job:= os.Getenv("JOB")

	jenkins := NewJenkinksJobParams(user, token, host,job,"")
	jenkins.trigger()

}


func NewJenkinksJobParams(user string, token string, host string, job string, jobparams string) *JenkinksJobParams {
	host = strings.TrimRight(host, "/")
	return &JenkinksJobParams{
		Username:    user,
		Token:	token,
		Host: host,
		Job:job,
		JobParams:jobparams,
	}
}
func (jenkins *JenkinksJobParams) trigger() {

	path := fmt.Sprintf("%s/job/%s/%s", jenkins.Host, jenkins.Job,"build")
	log.Info(fmt.Sprintf("Going to trigger %s job on %s", jenkins.Job, jenkins.Host))
	requestURL := jenkins.buildURL(path, url.Values{})
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	resp, err := jenkins.sendRequest(req)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info(resp.Status)
	log.Warn("Done")
}


func (jenkins *JenkinksJobParams) buildURL(path string, params url.Values) (requestURL string) {
	requestURL = path
	if params != nil {
		queryString := params.Encode()
		if queryString != "" {
			requestURL = requestURL + "?" + queryString
		}
	}
	return
}

func (jenkins *JenkinksJobParams) sendRequest(req *http.Request) (*http.Response, error) {

	req.SetBasicAuth(jenkins.Username, jenkins.Token)
	return http.DefaultClient.Do(req)
}

func (jenkins *JenkinksJobParams) parseResponse(resp *http.Response, body interface{}) (err error) {
	defer resp.Body.Close()

	if body == nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(data, body)
}
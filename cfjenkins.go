package main

import (
	"fmt"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
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
func (jenkins *JenkinksJobParams) trigger() error{

	path := fmt.Sprintf("%s/job/%s/%s", jenkins.Host, jenkins.Job,"build")
	fmt.Println("SHalom: "+jenkins.Token,jenkins.Host, jenkins.Username, jenkins.Job,path)
	return jenkins.post(path, url.Values{}, jenkins.JobParams)


	fmt.Println(jenkins.Token,jenkins.Host, jenkins.Username, jenkins.Job)
	return nil
}

func (jenkins *JenkinksJobParams) runJob(job string) error {
	path := fmt.Sprintf("%s/job/%s/%s", jenkins.Host, job,"build")
	fmt.Println("SHalom: "+jenkins.Token,jenkins.Host, jenkins.Username, job,path)
	return jenkins.post(path, url.Values{}, jenkins.JobParams)

}

func (jenkins *JenkinksJobParams) post(path string, params url.Values, body string) (err error) {
	requestURL := jenkins.buildURL(path, params)
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := jenkins.sendRequest(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	return jenkins.parseResponse(resp, body)
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
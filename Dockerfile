# ----- Go Builder ------
#
FROM golang:1.9 AS go

# set working directory
RUN mkdir -p /go/src/github.com/codefresh-io/cf-run-jenkins-job-plugin
WORKDIR /go/src/github.com/codefresh-io/cf-run-jenkins-job-plugin

# copy sources
COPY . .

RUN CGO_ENABLED=0 go build -o /cfjenkins cfjenkins.go

#
# ------ Release ------
#
FROM alpine:3.6

COPY --from=go /cfjenkins /usr/local/bin/

ENV JENKINS_HOST ${JENKINS_HOST}
ENV USER ${USER}
ENV TOKEN ${TOKEN}
ENV JOB ${JOB}

# Run locally docker run -it --rm -e JENKINS_URL="http://192.168.1.73:8080" -e JENKINS_USER="vadim" -e JENKINS_TOKEN="testpassword" -e JENKINS_JOB="testjob" local/cfjenkins
CMD ["/usr/local/bin/cfjenkins"]
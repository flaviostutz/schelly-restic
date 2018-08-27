FROM golang:1.10 AS BUILD

#doing dependency build separated from source build optimizes time for developer, but is not required
#install external dependencies first
# ADD go-plugins-helpers/Gopkg.toml $GOPATH/src/go-plugins-helpers/
ADD /main.go $GOPATH/src/restic-api/main.go
RUN go get -v restic-api

#now build source code
ADD restic-api $GOPATH/src/restic-api
RUN go get -v restic-api


FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y restic

VOLUME [ "/backup-source" ]
VOLUME [ "/backup-repo" ]

ENV RESTIC_PASSWORD ''
ENV LISTEN_PORT 8080
ENV LISTEN_IP '0.0.0.0'
ENV LOG_LEVEL 'debug'
ENV PRE_BACKUP_COMMAND ''
ENV POST_BACKUP_COMMAND ''

COPY --from=BUILD /go/bin/* /bin/
ADD startup.sh /
ADD pre-test.sh /
ADD post-test.sh /

CMD [ "/startup.sh" ]

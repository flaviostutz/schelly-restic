FROM golang:1.10 AS BUILD

#doing dependency build separated from source build optimizes time for developer, but is not required
#install external dependencies first
# ADD go-plugins-helpers/Gopkg.toml $GOPATH/src/go-plugins-helpers/
ADD /main.go $GOPATH/src/restic-api/main.go
RUN go get -v restic-api

#now build source code
ADD restic-api $GOPATH/src/
RUN go get -v restic-api


FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -f restic

VOLUME [ "/backup-source" ]
VOLUME [ "/backup-repo" ]

ENV RESTIC_PASSWORD ''

COPY --from=BUILD /go/bin/* /bin/
ADD startup.sh /

CMD [ "/startup.sh" ]

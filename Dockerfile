FROM golang:1.10 AS BUILD

#doing dependency build separated from source build optimizes time for developer, but is not required
#install external dependencies first
# ADD go-plugins-helpers/Gopkg.toml $GOPATH/src/go-plugins-helpers/
ADD /main.go $GOPATH/src/schelly-restic/main.go
RUN go get -v schelly-restic

#now build source code
ADD schelly-restic $GOPATH/src/schelly-restic
RUN go get -v schelly-restic


FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y restic

VOLUME [ "/backup-source" ]
VOLUME [ "/backup-repo" ]

EXPOSE 7070

ENV RESTIC_PASSWORD ''
ENV LISTEN_PORT 7070
ENV LISTEN_IP '0.0.0.0'
ENV LOG_LEVEL 'debug'

ENV PRE_POST_TIMEOUT '7200'
ENV PRE_BACKUP_COMMAND ''
ENV POST_BACKUP_COMMAND ''
ENV SOURCE_DATA_PATH '/backup-source'
ENV TARGET_DATA_PATH '/backup-repo'

COPY --from=BUILD /go/bin/* /bin/
ADD startup.sh /

CMD [ "/startup.sh" ]

FROM golang:1.10 AS BUILD

#doing dependency build separated from source build optimizes time for developer, but is not required
#install external dependencies first
# ADD go-plugins-helpers/Gopkg.toml $GOPATH/src/go-plugins-helpers/
ADD /main.go $GOPATH/src/backy2-api/main.go
RUN go get -v backy2-api

#now build source code
ADD backy2-api $GOPATH/src/
RUN go get -v backy2-api


FROM ubuntu:18.04

VOLUME [ "/data" ]

RUN apt-get update
RUN apt-get install wget
RUN wget https://github.com/wamdam/backy2/releases/download/v2.9.17/backy2_2.9.17_all.deb -O /tmp/backy2.deb
RUN dpkg -i /tmp/backy2.deb
RUN apt-get -f install
RUN rm /tmp/backy2.deb

COPY --from=BUILD /go/bin/* /bin/
ADD startup.sh /

CMD [ "/startup.sh" ]

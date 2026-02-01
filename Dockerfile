FROM golang:1.25rc1-alpine3.22

ARG VERSION
LABEL com.bluffpointtech.bpti18next.version="${VERSION}"

RUN mkdir /app
RUN mkdir /app/proto
RUN mkdir /app/resourcelocalization
RUN mkdir /app/util
RUN mkdir /opt/certs
WORKDIR /app

COPY proto/ /app/proto/
COPY resourcelocalization/ /app/resourcelocalization/
COPY util/ /app/util/
COPY go.sum /app/go.sum
COPY go.mod /app/go.mod 
COPY main.go /app/main.go

RUN go mod download
RUN go build -o bpti18next main.go

RUN rm -rf /app/proto
RUN rm -rf /app/resourcelocalization
RUN rm -rf /app/util
RUN rm go.sum
RUN rm go.mod
RUN rm main.go

RUN echo "bptnext:x:1014:1014::/bin/sh" >> /etc/passwd \
    && echo "bptnext:x:1014:bptnext" >> /etc/group
    
USER 1014

EXPOSE 8080 8081 9080 9081

CMD ["/app/bpti18next"]
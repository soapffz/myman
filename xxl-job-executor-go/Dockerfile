FROM golang:1.18

WORKDIR $GOPATH/src/xxl-job-go-executor
COPY . $GOPATH/src/xxl-job-go-executor
RUN rm -rf $GOPATH/src/xxl-job-go-executor/go.mod && rm -rf $GOPATH/src/xxl-job-go-executor/go.sum

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod init xxl-job-go-executor && go mod tidy

RUN go build .
ENTRYPOINT ["./xxl-job-go-executor"]
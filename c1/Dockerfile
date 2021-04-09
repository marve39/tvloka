FROM golang:1.16.3-alpine3.13 as test-base
RUN apk update && apk add build-base
RUN go get -u golang.org/x/lint/golint && \
    go get github.com/axw/gocov/... && \
    go get github.com/AlekSi/gocov-xml && \
    go get -u github.com/jstemmer/go-junit-report

FROM test-base as test 
ARG GITHUB_TOKEN
WORKDIR /app 
ADD . . 
RUN go get -a && make run-ci
RUN if [[ $(grep -c "FAIL:" junit.out) != 0 ]];then exit 2;fi
RUN make build 


FROM alpine:3.13.4 as base 
WORKDIR /app 
COPY --from=test /app/bin/c1-bin c1-bin

CMD ["/app/c1-bin"]

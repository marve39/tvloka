FROM golang:1.16.3-alpine3.13 as test-base
RUN apk update && apk add build-base
RUN go get -u golang.org/x/lint/golint

FROM test-base as test 
ARG GITHUB_TOKEN
WORKDIR /app 
ADD . . 
RUN go get -a && make run-ci
RUN if [[ $(grep -c "\"fail"\" test-result.json) != 0 ]];then exit 2;fi
RUN make build 


FROM alpine:3.13.4 as base 
WORKDIR /app 
COPY --from=test /app/bin/c1-bin c1-bin

CMD ["/app/c1-bin"]

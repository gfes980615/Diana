FROM golang:1.18.1-alpine
WORKDIR /diana
ADD . /diana
RUN cd /diana
RUN set GO111MODULE=on
RUN go build -o diana main.go
EXPOSE 8000
ENTRYPOINT ./diana
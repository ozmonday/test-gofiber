FROM golang:1.19.3-alpine3.15

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN  go mod download && go mod verify
COPY . .

RUN go build -o /usr/local/bin

ENV PORT=:3030
CMD [ "testfiber" ]

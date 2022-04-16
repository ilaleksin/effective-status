FROM golang:1.18

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
CMD ["./effective-status"]

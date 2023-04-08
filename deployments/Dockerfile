FROM golang

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY configs ./configs
COPY reports ./reports
COPY Makefile ./

RUN CGO_ENABLED=0 GOOS=linux make build

CMD ["./server", "configs/config.json"]
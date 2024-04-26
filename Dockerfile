FROM golang:1.22

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY data/ ./data/
COPY exceptions/ ./exceptions/
COPY helpers/ ./helpers/
COPY res/ ./res/
COPY routes/ ./routes/
COPY security/ ./security/
COPY services/ ./services/
COPY *.go ./
RUN go build -o ./export.out

EXPOSE 8088/tcp
CMD ["/app/export.out"]

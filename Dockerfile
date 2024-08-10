FROM golang:1.22.6-bullseye

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY appliance/ ./appliance/
COPY device/ ./device/
COPY runner/ ./runner/
COPY signal/ ./signal/
COPY slack/ ./slack/
COPY temp_controller/ ./temp_controller/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /home_controller

CMD ["/home_controller"]

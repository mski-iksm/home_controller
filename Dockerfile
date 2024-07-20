FROM golang:1.22

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

RUN GOOS=linux GOARCH=amd64 go build -o /home_controller

CMD ["/home_controller"]

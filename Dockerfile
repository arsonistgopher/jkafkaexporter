FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

#RUN go build -o main .

CMD ["app", "-identity=vmx01", "-kafkahost=http://192.168.10.200:9092", "-kafkaperiod=5", "-username=jet", "-password=Passw0rd", "-target=10.47.0.132"]

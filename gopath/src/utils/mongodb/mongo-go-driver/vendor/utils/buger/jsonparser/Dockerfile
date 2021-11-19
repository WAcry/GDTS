FROM golang:1.6

RUN go get utils/Jeffail/gabs
RUN go get utils/bitly/go-simplejson
RUN go get utils/pquerna/ffjson
RUN go get utils/antonholmquist/jason
RUN go get utils/mreiferson/go-ujson
RUN go get -tags=unsafe -u utils/ugorji/go/codec
RUN go get utils/mailru/easyjson

WORKDIR /go/src/utils/buger/jsonparser
ADD . /go/src/utils/buger/jsonparser
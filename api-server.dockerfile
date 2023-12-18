FROM golang:1.21.5

WORKDIR /app

COPY . .

RUN go mod download &&\
    go build -o /api-server ./api/cmd/server

EXPOSE 8030

ENTRYPOINT [ "/api-server" ]
CMD [ "-addr", "0.0.0.0:8030" ]

FROM golang:1.18 as builder

WORKDIR /src

COPY . .

RUN go build -o /bin/ilogin cmd/ilogin/*.go

FROM golang:1.18 as app

WORKDIR /app

COPY --from=builder /bin/ilogin /bin/ilogin
COPY web /app/web

ENTRYPOINT [ "/bin/ilogin" ]
CMD [ "-server" ]
FROM golang as builder

WORKDIR /backend

COPY . /backend/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Second stage

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /backend

RUN ls

COPY --from=builder . /backend/

EXPOSE 8000

CMD ["./main"]
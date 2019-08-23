FROM golang:1.12.0-stretch as builder
ADD . /go-cron-schedules/
ADD ./src/app/main.go /go-cron-schedules/
ADD ./src/types/types.go /go-cron-schedules/
WORKDIR /go-cron-schedules/
ENV GO111MODULE=on
ADD go.mod go.mod
RUN rm go.sum

RUN CGO_ENABLED=0 go build /go-cron-schedules/main.go

FROM scratch
COPY --from=builder /go-cron-schedules/main .
EXPOSE 80
ENTRYPOINT ["./main"]
FROM cimg/go:1.23.1 AS builder
USER 0

ENV PROJECT=jachobs-mind

RUN mkdir /${PROJECT}
RUN mkdir /botdata
RUN touch /botdata/data.db
RUN chown -R circleci /botdata
RUN chmod -R 777 /botdata

#COPY ./.env /${PROJECT}/.env
COPY ./libs /${PROJECT}/libs
COPY ./go.* /${PROJECT}/
COPY ./static /${PROJECT}/static
COPY ./bot /${PROJECT}/bot

RUN mkdir /${PROJECT}/bin

WORKDIR /${PROJECT}/
#ENV GOPATH=/${PROJECT}/libs
RUN go build -o ./bin/${PROJECT} ./bot 
#USER circleci

FROM alpine:latest

COPY --from=builder ./bin/${PROJECT} .


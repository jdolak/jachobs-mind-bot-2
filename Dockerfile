FROM cimg/go:1.23.1 AS builder
USER 0

ENV PROJECT=jachobs-mind

RUN mkdir /${PROJECT}
RUN mkdir /botdata
RUN touch /botdata/data.db
RUN chown -R circleci /botdata
RUN chmod -R 777 /botdata

WORKDIR /${PROJECT}/

COPY ./go.* /${PROJECT}/
RUN go mod download

COPY ./bot /${PROJECT}/bot
RUN mkdir /${PROJECT}/bin

RUN go build -o /${PROJECT}/bin/${PROJECT} ./bot 
COPY ./static /${PROJECT}/static

CMD ["/jachobs-mind/bin/jachobs-mind"]
#FROM ubuntu:latest

#ENV PROJECT=jachobs-mind
#COPY --from=builder /${PROJECT}/bin/${PROJECT} .
#CMD ["/jachobs-mind"]


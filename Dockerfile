FROM cimg/go:1.23.1
USER 0

ENV PROJECT=jachobs-mind

RUN mkdir /${PROJECT}
RUN mkdir /botdata
RUN touch /botdata/data.db
RUN chown -R circleci /botdata
RUN chmod -R 777 /botdata

COPY ./.env /${PROJECT}/.env
COPY ./libs /${PROJECT}/libs
COPY ./go.* /${PROJECT}/
COPY ./static /${PROJECT}/static
COPY ./src /${PROJECT}/src

WORKDIR /${PROJECT}
ENV GOPATH=/${PROJECT}/libs
RUN go build ./src/main.go
#USER circleci
CMD ["./main"]

FROM golang:1.6.2 

RUN go get github.com/tools/godep
RUN go get github.com/mjibson/esc
RUN CGO_ENABLED=0 go install -a std
ENV APP_DIR $GOPATH/src/github.com/arkenio/s3provisioner

# Set the entrypoint as the binary, so `docker run <image>` will behave as the binary
EXPOSE 8788
##ENTRYPOINT ["/s3provisioner"]
ADD . $APP_DIR
# Compile the binary and statically link
RUN cd $APP_DIR && \
    godep update && \
    esc -o server/static.go -prefix static -pkg server static && \
    CGO_ENABLED=0 godep restore && \
    godep go build -o s3provisioner -ldflags '-w -s'
CMD ["src/github.com/arkenio/s3provisioner/s3provisioner"]  
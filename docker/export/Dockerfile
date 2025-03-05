FROM golang:1.24.0 as builder
WORKDIR /app
RUN go env -w GOMODCACHE=/root/.cache/go-build
RUN \
--mount=type=bind,src=./,target=/app/,rw \
--mount=type=cache,target=/root/.cache/go-build \
go get -u ./... && go mod download
COPY --chmod=755 ./build.sh /app/build.sh
COPY . .
RUN \ 
--mount=type=cache,target=/root/.cache/go-build \
/app/build.sh && tar -czf rpkengine.tar.gz ./bin/

FROM scratch AS export
COPY --from=builder /app/rpkengine.tar.gz ./rpkgengine.tar.gz
ENTRYPOINT ["./rpkgengine.tar.gz"]
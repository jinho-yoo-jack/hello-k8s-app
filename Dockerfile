# 1단계: 테스트 — `docker build --target test .` 로 이 단계까지만 실행할 수 있다
FROM golang:1.27-alpine AS test
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go test -v ./...

# 2단계: 빌드 — 테스트 단계를 이어받으므로 테스트가 실패하면 이미지도 만들어지지 않는다
FROM test AS build
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=${VERSION}" -o /out/hello .

# 3단계: 실행 이미지 — 셸도 패키지 관리자도 없는 아주 작은 이미지에 바이너리만 담는다
FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/hello /hello
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/hello"]

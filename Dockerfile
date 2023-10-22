FROM --platform=$BUILDPLATFORM golang:1.21.1-alpine3.18 as build
WORKDIR /src
ARG TARGETOS TARGETARCH
RUN --mount=target=. --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o /out/frontend -ldflags "-s -w"
COPY public /out/public
COPY views /out/views

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=build /out .
ENTRYPOINT ["./frontend", "-backend=http://backend:5001", "-hostname=https://demo.metronero.com", "-token-secret=xufeg6uoth4oM7CohK2D"]

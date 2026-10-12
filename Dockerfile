FROM golang:1.27.2-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
COPY web ./web
COPY openapi.yaml ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/calendarops-server ./cmd/calendarops-server && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/calendarops-export ./cmd/calendarops-export
FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build --chown=nonroot:nonroot /out/calendarops-server /calendarops-server
COPY --from=build --chown=nonroot:nonroot /out/calendarops-export /calendarops-export
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/calendarops-server"]

FROM golang:1.24-alpine AS build

# variables build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy app source (inclut ./public)
COPY . .

# build static, optimisé
RUN go build -ldflags="-s -w" -o /app/server ./...

# Final image minimal (scratch)
FROM scratch

# copy binary
COPY --from=build /app/server /server

# copy static files
COPY --from=build /src/public /public

# si ton binaire attend public dans ./public, ajuste WORKDIR
WORKDIR /

EXPOSE 8081

# si ton binaire lit la variable STATIC_DIR, tu peux la fournir
ENV STATIC_DIR=/public

ENTRYPOINT ["/server"]

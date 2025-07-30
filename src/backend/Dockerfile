FROM golang:1.24.4

# Install wireguard-tools (adds 'wg' command)
RUN apt-get update && apt-get install -y wireguard-tools

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o backend .

CMD ["./backend"]

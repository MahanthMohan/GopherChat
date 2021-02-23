# Base Image for golang (latest)
FROM golang:latest

# Make a chat directory
RUN mkdir /chat

# Add all files in current directory into chat
ADD . /chat

# Set it as the Working Directory ($PWD = /chat)
WORKDIR /chat

# Build the main.go file into an executable called chat
RUN go build -o chat

# Run the executable using CMD
CMD ["./chat"]

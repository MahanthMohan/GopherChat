echo "** Building GopherChat **"

echo "-- (Step 1) Building MacOS/Darwin binaries --"
gobin=$PWD/bin/darwin/amd64 goos=darwin goarch=amd64 go install -ldflags="-s -w" *.go
echo "-- Finished MacOS/Darwin build --"

echo "-- (Step 2) Building Linux binaries --"
gobin=$PWD/bin/linux/amd64 goos=linux goarch=amd64 go install -ldflags="-s -w" *.go
gobin=$PWD/bin/linux/arm64 goos=linux goarch=arm64 go install -ldflags="-s -w" *.go
echo "-- Finished Linux build --"

echo "-- (Step 3) Building Windows binaries --"
gobin=$PWD/bin/windows/amd64 goos=windows goarch=amd64 go install -ldflags="-s -w" *.go
gobin=$PWD/bin/windows/arm goos=windows goarch=arm go install -ldflags="-s -w" *.go
echo "-- Finished Windows build --"

echo "** Finished GopherChat build **"

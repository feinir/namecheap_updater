rm -rf bin/update_cheap-mac-inter
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags '-w -s' -o bin/update_cheap-mac-inter main.go

rm -rf bin/update_cheap-lin
rm -rf bin/update_cheap-lin-x86
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags '-w -s' -o bin/update_cheap-lin main.go
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -trimpath -ldflags '-w -s' -o bin/update_cheap-lin-x86 main.go

rm -rf bin/update_cheap.exe
rm -rf bin/update_cheap-x86.exe
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-H windowsgui -w -s" -o bin/update_cheap.exe main.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -trimpath -ldflags="-H windowsgui -w -s" -o bin/update_cheap-x86.exe main.go
all:
	mkdir -p dist
	gox -ldflags="-s -w" -tags="full" -osarch="darwin/amd64 linux/386 linux/amd64 linux/arm freebsd/amd64" -output="dist/jiq_{{.OS}}_{{.Arch}}" github.com/fiatjaf/jiq/cmd/jiq

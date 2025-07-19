build-parse:
	GOOS=linux GOARCH=amd64 go build -o bootstrap ./parseUpload
	zip parseUpload.zip bootstrap
	rm bootstrap

build-presign:
	GOOS=linux GOARCH=amd64 go build -o bootstrap ./getPresignedUrl
	zip getPresignedUrl.zip bootstrap
	rm bootstrap

build-all: build-parse build-presign

# Wrote by yijian on 2024/09/15

ifeq ($(OS),Windows_NT)
	target=cretime_renamer.exe
else
	target=cretime_renamer
endif

all: ${target}

${target}: main.go
ifeq ($(OS),Windows_NT)
	set GOOS=windows
	set GOARCH=amd64
endif
	go mod tidy && go build -o $@ $<

.PHONY: clean

clean:
	rm -f ${target

install: ${target}
ifeq ($(OS),Windows_NT)
	copy ${target} %GOPATH%\bin\
else
	cp ${target} $$GOPATH/bin/
endif

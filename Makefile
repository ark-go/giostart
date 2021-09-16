SHELL := /bin/bash
arklibgo := ~/ProjectsGo/arkAlias.sh
version = ~/ProjectsGo/arkAlias.sh getlastversion
.PHONY: check

.SILENT: build getlasttag buildzip buildwin


buildlinux:
	$(info +Компиляция Linux)
	@echo $$($(version))
	go build -ldflags "-s -w -X 'main.versionProg=$$($(version))'" -o ./bin/main/giostart cmd/main/main.go
buildzip:
	$(info +Компиляция с жатием)
	go build -ldflags "-s -w" -o ./bin/main/giostart cmd/main/main.go
buildwin:
	$(info +Компиляция windows)
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o ./bin/main/giostart.exe -tags static -ldflags "-s -w -X 'main.versionProg=$$($(version))'" cmd/main/main.go

run: buildlinux buildwin
	$(info +Запуск)
	./bin/main/giostart

build: buildlinux buildwin
	$(info +Сборка)

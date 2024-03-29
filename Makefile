ranssibot: $(shell find -name "*.go")
	go install

compilepie: $(shell find -name "*.go")
	env GOOS=linux GOARCH=arm GOARM=7 go build -v

packagepie: compilepie ranssibot $(shell ls lang/*.lang)
	tar cvfJ ranssibot.tar.xz ranssibot lang/*.lang

pie: packagepie
	rm -f ranssibot

clean:
	rm -f ranssibot
	rm -f ranssibot.tar.xz

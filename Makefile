
ranssibot: $(shell ls */*.go *.go)
	go install

compilepie: $(shell ls */*.go *.go)
	env GOOS=linux GOARCH=arm GOARM=7 go build -v

packagepie: compilepie ranssibot lang/en_US.lang lang/fi_FI.lang
	tar -zcvf ranssibot.tar.gz ranssibot lang/en_US.lang lang/fi_FI.lang

pie: packagepie
	rm -f ranssibot

clean:
	rm -f ranssibot
	rm -f ranssibot.tar.gz

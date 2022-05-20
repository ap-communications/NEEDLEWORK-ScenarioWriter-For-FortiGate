test:
	go test -v tests/*.go
	# go test -v tests/service_test.go
	# go test -v tests/other_test.go
build:
	GOOS=windows GOARCH=386 go build -o bin/NEEDLEWORK-ScenarioWriter-For-FortiGate_windows_386.exe github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate
	GOOS=windows GOARCH=amd64 go build -o bin/NEEDLEWORK-ScenarioWriter-For-FortiGate.exe github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate
	GOOS=darwin GOARCH=386 go build -o bin/NEEDLEWORK-ScenarioWriter-For-FortiGate_darwin_386 github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate
	GOOS=darwin GOARCH=amd64 go build -o bin/NEEDLEWORK-ScenarioWriter-For-FortiGate_darwin_amd64 github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate
	GOOS=linux GOARCH=386 go build -o bin/NEEDLEWORK-ScenarioWriter-For-FortiGate_linux_386 github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate
	GOOS=linux GOARCH=amd64 go build -o bin/NEEDLEWORK-ScenarioWriter-For-FortiGate_linux_amd64 github.com/ap-communications/NEEDLEWORK-ScenarioWriter-For-FortiGate

SRV_00 = server00
SRV_01 = server01
SRV_02 = server02

CLI_01 = candy-client01
CLI_02 = candy-client02

ex00:
	cd ex00/server/go && go build -o $(SRV_00)
	mv ex00/server/go/$(SRV_00) .

ex01:
	cd ex01/server/go && go build -o $(SRV_01)
	mv ex01/server/go/$(SRV_01) .
	cd ex01/client && go build -o $(CLI_01)
	mv ex01/client/$(CLI_01) .

ex02:
	cd ex02/server/go && go build -o $(SRV_02)
	mv ex02/server/go/$(SRV_02) .
	cd ex02/client && go build -o $(CLI_02)
	mv ex02/client/$(CLI_02) .

build: ex00 ex01 ex02

clean:
	go clean
	rm -f $(SRV_00) $(SRV_01) $(SRV_02) $(CLI_01) $(CLI_02)

.PHONY: ex00 ex01 ex02
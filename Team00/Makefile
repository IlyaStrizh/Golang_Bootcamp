SRV_00 = transmitter
SRV_01 = transmitterQuick

CLI_00 = receiver
CLI_01 = receiverDetection
CLI_02 = receiverReport
CLI_03 = receiverReportDetection

ex00:
	cd ex00/transmitter && go build -o $(SRV_00)
	mv ex00/transmitter/$(SRV_00) .
	cd ex00/receiver && go build -o $(CLI_00)
	mv ex00/receiver/$(CLI_00) .

ex01:
	cd ex01/transmitter && go build -o $(SRV_01)
	mv ex01/transmitter/$(SRV_01) .
	cd ex01/receiver && go build -o $(CLI_01)
	mv ex01/receiver/$(CLI_01) .

ex02:
	cd ex02 && go build -o $(CLI_02)
	mv ex02/$(CLI_02) .

ex03:
	cd ex03/receiver && go build -o $(CLI_03)
	mv ex03/receiver/$(CLI_03) .

build: ex00 ex01 ex02 ex03

clean:
	go clean
	rm -f $(SRV_00) $(SRV_01) $(CLI_00) $(CLI_01) $(CLI_02) $(CLI_03)

.PHONY: ex00 ex01 ex02 ex03
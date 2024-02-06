EXECUTABLE_0 = readDB
EXECUTABLE_1 = compareDB
EXECUTABLE_2 = compareFS

ex00:
	cd ex00 && go build -o $(EXECUTABLE_0)
	mv ex00/$(EXECUTABLE_0) ./

ex01:
	cd ex01 && go build -o $(EXECUTABLE_1)
	mv ex01/$(EXECUTABLE_1) ./

ex02:
	cd ex02 && go build -o $(EXECUTABLE_2)
	mv ex02/$(EXECUTABLE_2) ./

test: clean ex00 ex01 ex02
	./readDB -f ex00/original_database.xml > ex00_xml_to_json.json
	./readDB -f ex00/stolen_database.json > ex00_json_to_xml.xml
	./compareDB --old ex00/original_database.xml --new ex00/stolen_database.json > ex01_compareDB.txt
	./compareFS --old ex02/snapshot1.txt --new ex02/snapshot2.txt > ex02_compareFS.txt
	./compareFS --new ex02/snapshot1.txt --old ex02/snapshot2.txt > ex02_compareFS_2.txt

clean:
	go clean
	rm -rf $(EXECUTABLE_0) $(EXECUTABLE_1) $(EXECUTABLE_2)
	rm -rf *json *xml *txt

.PHONY: ex00 ex01 ex02
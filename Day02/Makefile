EXECUTABLE_0 = myFind
EXECUTABLE_1 = myWc
EXECUTABLE_2 = myXargs
EXECUTABLE_3 = myRotate

ex00:
	cd ex00 && go build -o $(EXECUTABLE_0)
	mv ex00/$(EXECUTABLE_0) ./

ex01:
	cd ex01 && go build -o $(EXECUTABLE_1)
	mv ex01/$(EXECUTABLE_1) ./

ex02:
	cd ex02 && go build -o $(EXECUTABLE_2)
	mv ex02/$(EXECUTABLE_2) ./

ex03:
	cd ex03 && go build -o $(EXECUTABLE_3)
	mv ex03/$(EXECUTABLE_3) ./

test: clean ex00 ex01 ex02 ex03
	./myFind -f -ext 'log' . | ./myXargs ./myWc -l > test.txt
	./myRotate testFile/test.log
	./myRotate -a testMyRotate testFile/test.log

clean:
	go clean
	rm -rf $(EXECUTABLE_0) $(EXECUTABLE_1) $(EXECUTABLE_2) $(EXECUTABLE_3)
	rm -rf *txt testMyRotate testFile/*.gz

.PHONY: ex00 ex01 ex02 ex03
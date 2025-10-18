SRC_DIR := ./src
DST_DIR := ./dst

GOLANG := go

build: $(DST_DIR)/uspy_arm64.exe $(DST_DIR)/uspy_arm.exe $(DST_DIR)/uspy_amd64.exe $(DST_DIR)/uspy_386.exe

$(DST_DIR)/uspy_arm64.exe: $(SRC_DIR)/main.go
	cd $(SRC_DIR) && GOOS=windows GOARCH=arm64 $(GOLANG) build -o ../$(DST_DIR)/uspy_arm64.exe .

$(DST_DIR)/uspy_arm.exe: $(SRC_DIR)/main.go
	cd $(SRC_DIR) && GOOS=windows GOARCH=arm $(GOLANG) build -o ../$(DST_DIR)/uspy_arm.exe .

$(DST_DIR)/uspy_amd64.exe: $(SRC_DIR)/main.go
	cd $(SRC_DIR) && GOOS=windows GOARCH=amd64 $(GOLANG) build -o ../$(DST_DIR)/uspy_amd64.exe .

$(DST_DIR)/uspy_386.exe: $(SRC_DIR)/main.go
	cd $(SRC_DIR) && GOOS=windows GOARCH=386 $(GOLANG) build -o ../$(DST_DIR)/uspy_386.exe .

test:
	cd $(SRC_DIR) && go test .

clean:
	rm -r ./dst/*.exe
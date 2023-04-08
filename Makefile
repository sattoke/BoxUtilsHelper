# [Prerequisite]
# - You must be using WSL
# - NSIS must be installed.
# - Go must be installed

MAKENSIS = "/mnt/c/Program Files (x86)/NSIS/makensis.exe"
GO = go
NSI = Install.nsi
INSTALLER = Install.exe
EXE = boxutils-helper.exe
SRC = boxutils-helper.go
BUILD_OPTIONS = GOOS=windows GOARCH=amd64

.PHONE: all
all: $(INSTALLER)

$(INSTALLER): $(NSI) $(EXE)
	$(MAKENSIS) $(NSI)

$(EXE): $(SRC)
	$(BUILD_OPTIONS) $(GO) build -o $@ $^

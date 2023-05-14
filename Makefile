# [Prerequisite]
# - You must be using WSL
# - NSIS must be installed.
# - Go must be installed

ARCH = x86_64
VERSION_FILE = VERSION
VERSION = $(shell cat $(VERSION_FILE))

SED = sed
MAKENSIS = "/mnt/c/Program Files (x86)/NSIS/makensis.exe"
GO = go
NSI_SRC = Install.tmpl.nsi
NSI_DST = Install.nsi
INSTALLER = Install-$(VERSION).$(ARCH).exe
EXE = boxutils-helper.exe
SRC = boxutils-helper.go
BUILD_OPTIONS = GOOS=windows GOARCH=amd64

.PHONE: all
all: $(INSTALLER)

$(INSTALLER): $(NSI_DST) $(EXE)
	$(MAKENSIS) $(NSI_DST)

$(EXE): $(SRC)
	$(BUILD_OPTIONS) $(GO) build -o $@ $^

$(NSI_DST): $(NSI_SRC) $(VERSION_FILE)
	$(SED)  -e s/__VERSION__/$(VERSION)/g -e s/__ARCH__/$(ARCH)/g < $(NSI_SRC) > $(NSI_DST)

.PHONE: clean
clean:
	-rm $(EXE) $(NSI_DST) $(INSTALLER)

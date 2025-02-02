GO_FILES = main.go
C_SOURCES = $(wildcard ccode/OBJ/*.c)
C_OBJECTS = $(C_SOURCES:.c=.o)
LIB_NAME = liblkh.a

#all: go

#go: $(GO_FILES) $(LIB_NAME)
#	go build -o lkh.exe

$(LIB_NAME): $(C_OBJECTS)
	ar rcs $@ $(C_OBJECTS)

#clean:
#	del $(C_OBJECTS) $(LIB_NAME) lkh.exe

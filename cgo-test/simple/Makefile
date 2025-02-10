GO_FILES = main.go
C_SOURCES = $(wildcard ccode/*.c)
C_OBJECTS = $(C_SOURCES:.c=.o)
LIB_NAME = libmycodelib.a

all: go

go: $(GO_FILES) $(LIB_NAME)
	go build -o myprogram.exe

$(LIB_NAME): $(C_OBJECTS)
	ar rcs $@ $(C_OBJECTS)

%.o: %.c
	gcc -c $< -o $@ -Iccode

clean:
	del $(C_OBJECTS) $(LIB_NAME) myprogram.exe

#GO_FILES = main.go
#C_FILES = ccode/a.c ccode/b.c
#C_OBJECTS = $(C_FILES:.c=.o)

#all: go

#go: $(GO_FILES) $(C_OBJECTS)
#	go build -o simple.exe

#%.o: %.c
#	gcc -c $< -o $@ -Iccode

#clean:
#	del $(C_OBJECTS) simple.exe
##	rm -f $(C_OBJECTS) simple.exe

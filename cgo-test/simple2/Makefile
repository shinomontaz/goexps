GO_FILES = main.go
C_SOURCES = $(wildcard ccode/*.c)
C_OBJECTS = $(C_SOURCES:.c=.o)
LIB_NAME = libmy.a

all: go

go: $(GO_FILES) $(LIB_NAME)
	go build -o myprogram.exe

$(LIB_NAME): $(C_OBJECTS)
	ar rcs $@ $(C_OBJECTS)

%.o: %.c
	gcc -c $< -o $@ -Iccode

clean:
	del $(C_OBJECTS) $(LIB_NAME) myprogram.exe

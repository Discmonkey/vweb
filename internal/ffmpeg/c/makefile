CC = gcc  # C compiler
CFLAGS = -fPIC -Wall -Wextra -O2 -g  # C flags
LDFLAGS = -shared   # linking flags
RM = rm -f   # rm command
TARGET_LIB = libtarget.so  # target lib
TARGET_TEST = build/test_ffmpeg
SOURCE = src/test/test.c  # source files

# where our local install of ffmpeg is located
FFMPEG_ROOT = ../../../third_party/ffmpeg
FR=-L$(FFMPEG_ROOT)/
# the include path relative to the install
INCLUDE = -I$(FFMPEG_ROOT)

# the libraries we are interested in linking to
LIB_PATH =  $(FR)libavformat $(FR)libavcodec $(FR)libswscale $(FR)libavutil $(FR)libavdevice $(FR)libswresample

# the actual libraries
LIBRARIES = -lavformat -lavcodec -lswscale -lavutil -lavdevice -lswresample -lm -lz -ldl -lpthread

.PHONY: all
all: $(TARGET_TEST)

$(TARGET_TEST): $(SOURCE)
	$(CC) $(INCLUDE) $(LIB_PATH) -o $@ $^ $(LIBRARIES)
# $(TARGET_LIB): $(OBJS)
#	$(CC) ${LDFLAGS} -o $@ $^

# $(SRCS:.c=.d):%.d:%.c
#	$(CC) $(CFLAGS) -MM $< >$@

# include $(SRCS:.c=.d)

.PHONY: clean
clean:
	-${RM} ${TARGET_LIB} ${OBJS} $(SOURCE:.c=.d)
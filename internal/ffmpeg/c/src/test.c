#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>

struct Stream {
    
};
int main() {
    AVFormatContext *pFormatCtx = NULL;

    char *f = "/home/max/go/src/vweb/test/data/big_buck_bunny_1080_10s_1mb_h264.mp4\0";
    if(avformat_open_input(&pFormatCtx, f,
                           NULL, NULL)!=0) {
        return -1; // Couldn't open file
    }

    if(avformat_find_stream_info(pFormatCtx, NULL)<0) {
        return -1; // Couldn't find stream information
    }

    av_dump_format(pFormatCtx, 0, f, 0);


}
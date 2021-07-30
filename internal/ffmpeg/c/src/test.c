#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>

int main() {
    AVFormatContext *pFormatCtx = NULL;

    if(avformat_open_input(&pFormatCtx, "/home/max/go/src/vweb/test/data/big_buck_bunny_1080_10s_1mb_h264.mp4",
                           NULL, NULL)!=0) {
        return -1; // Couldn't open file
    }

    if(avformat_find_stream_info(pFormatCtx, NULL)<0) {
        return -1; // Couldn't find stream information
    }


}
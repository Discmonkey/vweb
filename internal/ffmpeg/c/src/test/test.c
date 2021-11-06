#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>

#define BUFFER_SIZE 4096;

const int32_t kNotFound = -1;

struct Stream {
    
};
int main() {
    AVFormatContext *pFormatCtx = NULL;

    AVFrame *frame;
    AVPacket  *pkt;

    char *f = "../../../test/data/big_buck_bunny_1080_10s_1mb_h264.mp4";
    if(avformat_open_input(&pFormatCtx, f,
                           NULL, NULL)!=0) {
	    printf("could not open file\n");
        return -1; // Couldn't open file
    }
    if(avformat_find_stream_info(pFormatCtx, NULL)<0) {
        printf("could not open file\n");
	    return -1; // Couldn't find stream information
    }
    if (!(pkt = av_packet_alloc())) {
        printf("could not allocate pkt\n");
        return -1;
    }
    av_dump_format(pFormatCtx, 0, f, 0);

    AVCodecContext *pCodecCtxOrig = NULL;
    AVCodecContext *pCodecCtx = NULL;

    int videoStreamIdx = kNotFound;
    for (int i = 0; i < pFormatCtx->nb_streams; i++) {
        if (pFormatCtx->streams[i]->codecpar->codec_type ==
            AVMEDIA_TYPE_VIDEO) {
            printf("found AVMEDIA_TYPE_VIDEO position: %d\n", i);
            videoStreamIdx = i;
            break;
        }
    }

    if (videoStreamIdx == kNotFound) {
        printf("stream not found");
        return -1;
    }

    AVStream * video = pFormatCtx->streams[videoStreamIdx];
    AVCodec *p_codec = avcodec_find_decoder(video->codecpar->codec_id);
    if (!p_codec) {
        printf("codec not found");
        return -1;
    }


}

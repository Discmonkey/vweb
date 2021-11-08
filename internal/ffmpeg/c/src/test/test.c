#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>

const int32_t k_not_found = -1;

int main() {
    AVFormatContext *p_format_context = NULL;
    AVPacket  *pkt;

    char *f = "../../../test/data/big_buck_bunny_1080_10s_1mb_h264.mp4";
    if(avformat_open_input(&p_format_context, f,NULL, NULL)!=0) {
	    printf("could not open file\n");
        return -1; // Couldn't open file
    }

    if(avformat_find_stream_info(p_format_context, NULL) < 0) {
        printf("could not open file\n");
	    return -1; // Couldn't find stream information
    }

    if (!(pkt = av_packet_alloc())) {
        printf("could not allocate pkt\n");
        return -1;
    }

    av_dump_format(p_format_context, 0, f, 0);

    int video_stream_index = k_not_found;
    for (int i = 0; i < p_format_context->nb_streams; i++) {
        if (p_format_context->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_VIDEO) {
            video_stream_index = i;
            break;
        }
    }

    if (video_stream_index == k_not_found) {
        printf("stream not found");
        return -1;
    }

    printf("video index is %d\n", video_stream_index);
    while (av_read_frame(p_format_context, pkt) == 0) {
        printf("%d\n", pkt->stream_index);
        printf("%d\n", pkt->size);
        printf("key frame? %d\n", (pkt->flags & AV_PKT_FLAG_KEY) > 0);
        printf("\n");
    }
}

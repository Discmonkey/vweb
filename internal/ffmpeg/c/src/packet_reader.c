#include <string.h>

#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>

#include "packet_reader.h"

// safely free an error struct
void free_error(Error *error) {
    if (error != NULL) {
        if (error->reason != NULL) {
            free(error->reason);
        }
        free(error);
    }
}

// safely free a stream struct
void free_stream(Stream *stream) {
    if (stream != NULL) {
        if (stream->context != NULL) {
            avformat_free_context((AVFormatContext *) stream->context);
        }

        if (stream->packet != NULL) {
            av_packet_unref(((AVPacket *) stream->packet));
        }

        free(stream);
    }
}

// utility function for constructing an error
Error* new_error(char *message, bool is_recoverable) {
    Error *e = malloc(sizeof(Error));
    e->reason = av_strdup(message);
    e->can_recover = is_recoverable;

    return e;
}

// attempt to open video resource (file or video stream)
// either the stream or the error will be set in the returned structure
StreamOrError open_stream(char *url) {
    AVFormatContext  *p_format_context = NULL;
    AVPacket *pkt;
    StreamOrError stream_or_error = {NULL, NULL};
    if (avformat_open_input(&p_format_context, url, NULL, NULL) < 0) {
        stream_or_error.error = new_error("could not open stream", false);
        return stream_or_error;
    }

    if (avformat_find_stream_info(p_format_context, NULL) < 0) {
        stream_or_error.error = new_error("could not find stream info", false);
        return stream_or_error;
    }

    if (!(pkt = av_packet_alloc())) {
        stream_or_error.error = new_error("could not allocate packet", true);
        return stream_or_error;
    }

    int video_stream_index = -1;
    for (int i = 0; i < p_format_context->nb_streams; i++) {
        if (p_format_context->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_VIDEO) {
            video_stream_index = i;
            break;
        }
    }

    if (video_stream_index == -1) {
        stream_or_error.error = new_error("video channel not found within mpeg stream", false);
        return stream_or_error;
    }

    stream_or_error.stream = malloc(sizeof(Stream));
    stream_or_error.stream->context = (void *)p_format_context;
    stream_or_error.stream->packet = (void *)pkt;
    stream_or_error.stream->video_index = video_stream_index;

    return stream_or_error;
}

// read from a stream and obtain the next packet or error.
// either the packet or the error will be set in the returned struct
PacketOrError next_packet(Stream *stream) {
    PacketOrError packet_or_error = {NULL, NULL};
    AVPacket *av_packet = (AVPacket *)stream->packet;
    while (av_read_frame((AVFormatContext *)stream->context, av_packet) == 0) {
       if (av_packet->stream_index == stream->video_index) {
           packet_or_error.packet = malloc(sizeof(Packet));
           packet_or_error.packet->data = av_packet->data;
           packet_or_error.packet->is_key_frame = (av_packet->flags & AV_PKT_FLAG_KEY) > 0;
           packet_or_error.packet->size = av_packet->size;
           return packet_or_error;
       }
    }

    packet_or_error.error = new_error("could not read packet from stream", true);
    return packet_or_error;
}

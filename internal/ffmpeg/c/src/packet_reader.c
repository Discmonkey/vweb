#include <string.h>

#include <libavcodec/avcodec.h>
#include <libavformat/avformat.h>

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

typedef struct Buffer {
    void *data;
    int capacity;
    int length;
} Buffer;

Buffer *new_buffer() {
    Buffer *buffer = malloc(sizeof(Buffer));
    buffer->length = 0;
    buffer->capacity = 0;
    buffer->data = NULL;

    return buffer;
}

void write_buffer(Buffer *dst, AVPacket *packet, uint8_t *side_data, int side_data_size) {
    int new_length = packet->size + side_data_size;
    if (dst->capacity < new_length) {
       if (dst->data != NULL) {
           free(dst->data);
       }

       dst->data = malloc(new_length);
       dst->capacity = new_length;
    }

    memcpy(dst->data, side_data, side_data_size);
    memcpy(dst->data + side_data_size, packet->data, packet->size);
    dst->length = new_length;
}

// safely free a stream struct
void free_stream(Stream *stream) {
    if (stream != NULL) {
        if (stream->context != NULL) {
            avformat_free_context((AVFormatContext *) stream->context);
        }

        if (stream->av_packet != NULL) {
            av_packet_unref(((AVPacket *) stream->av_packet));
        }

        if (stream->re_packet != NULL) {
            free(stream->re_packet);
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
    printf("%s\n", url);
    AVFormatContext  *p_format_context = NULL;
    AVPacket *pkt;
    StreamOrError stream_or_error = {NULL, NULL};
    if (avformat_open_input(&p_format_context, url, NULL, NULL) < 0) {
        stream_or_error.error = new_error("could not open stream", false);
        return stream_or_error;
    }
    printf("avformat_open succeeded\n");
    if (avformat_find_stream_info(p_format_context, NULL) < 0) {
        stream_or_error.error = new_error("could not find stream info", false);
        return stream_or_error;
    }

    if (!(pkt = av_packet_alloc())) {
        stream_or_error.error = new_error("could not allocate packet", true);
        return stream_or_error;
    }
    printf("allocated packet successfully\n");

    int video_stream_index = -1;
    printf("num streams %d\n", p_format_context->nb_streams);
    for (int i = 0; i < p_format_context->nb_streams; i++) {
        printf("checking %d\n", i);
        if (p_format_context->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_VIDEO) {
            video_stream_index = i;
            break;
        }
    }

    printf("found video stream %d\n", video_stream_index);

    if (video_stream_index == -1) {
        stream_or_error.error = new_error("video channel not found within mpeg stream", false);
        return stream_or_error;
    }

    stream_or_error.stream = malloc(sizeof(Stream));
    stream_or_error.stream->context = (void *)p_format_context;
    stream_or_error.stream->re_packet = malloc(sizeof(Packet));
    stream_or_error.stream->av_packet = pkt;
    stream_or_error.stream->video_index = video_stream_index;

    // buffer for packet data
    stream_or_error.stream->buffer = new_buffer();

    printf("returning\n");
    return stream_or_error;
}

// read from a stream and obtain the next packet or error.
// either the packet or the error will be set in the returned struct
PacketOrError next_packet(Stream *stream) {
    PacketOrError packet_or_error = {NULL, NULL};
    AVPacket *av_packet = (AVPacket *)stream->av_packet;
    AVFormatContext *context = (AVFormatContext *)stream->context;

    while (av_read_frame(context, av_packet) == 0) {
        if (av_packet->stream_index == stream->video_index) {
            packet_or_error.packet = stream->re_packet;
            packet_or_error.packet->is_key_frame = (av_packet->flags & AV_PKT_FLAG_KEY) > 0;

            // For key frames we need to copy the sps and pps buffers at the start of the frame
            if (packet_or_error.packet->is_key_frame) {
                Buffer *buffer = (Buffer *)stream->buffer;
                write_buffer(buffer,
                             av_packet,
                             context->streams[stream->video_index]->codecpar->extradata,
                             context->streams[stream->video_index]->codecpar->extradata_size);
                packet_or_error.packet->data = buffer->data;
                packet_or_error.packet->size = buffer->length;
            } else {
                packet_or_error.packet->data = av_packet->data;
                packet_or_error.packet->size = av_packet->size;
            }

            return packet_or_error;
        }
    }

    packet_or_error.error = new_error("could not read packet from stream", true);
    return packet_or_error;
}

//
// Created by max on 11/7/21.
//

#ifndef C_READ_PACKETS_H
#define C_READ_PACKETS_H

typedef struct Stream {
    void *context;
    int video_index;
} Stream;

typedef struct Error {
    char *reason;
    bool can_recover;
} Error;

typedef struct StreamOrError {
    Stream *stream;
    Error *error;
} StreamOrError;

typedef struct Packet {
    bool is_key_frame;
    void *data;
} Packet;

typedef struct PacketOrError {
    Packet *packet;
    Error *error;
} Error;

// safely free an error struct
void free_error(Error *error);

// safely free a stream struct
void free_stream(Stream *stream);

// attempt to open video resource (file or video stream)
// either the stream or the error will be set in the returned structure
StreamOrError open_stream(char *url);

// read from a stream and obtain the next packet or error.
// either the packet or the error will be set in the returned struct
PacketOrError next_packet(Stream *stream);

#endif //C_READ_PACKETS_H
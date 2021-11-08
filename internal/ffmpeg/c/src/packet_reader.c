//
// Created by max on 11/7/21.
//

#include "packet_reader.h"

// safely free an error struct
void free_error(Error *error) {
}

// safely free a stream struct
void free_stream(Stream *stream) {

}

// attempt to open video resource (file or video stream)
// either the stream or the error will be set in the returned structure
StreamOrError open_stream(char *url) {
    return StreamOrError{};
}

// read from a stream and obtain the next packet or error.
// either the packet or the error will be set in the returned struct
PacketOrError next_packet(Stream *stream) {
    return PacketOrError{};
}


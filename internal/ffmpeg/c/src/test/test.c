#include "../packet_reader.h"
#include <stdio.h>

int main() {

    char *f = "/home/max/go/src/vweb/test/data/output.ts";

    StreamOrError stream_or_error = open_stream(f);

    if (stream_or_error.error != NULL) {
        printf("error opening %s", stream_or_error.error->reason);
    }
    while (true) {
        PacketOrError packet_or_error = next_packet(stream_or_error.stream);

        if (packet_or_error.error != NULL) {
            printf("%s\n", packet_or_error.error->reason);
            break;
        }

        printf("is key? %d, size: %d\n", packet_or_error.packet->is_key_frame, packet_or_error.packet->size);
    }
}

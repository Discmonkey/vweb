//
// Created by max on 12/20/21.
//

#include <stdio.h>
#include <libavutil/avutil.h>
int main() {
    char *buf = malloc(300);

    av_strerror(-1094995529, buf, 300);

    printf("%s\n", buf);
}

cmake_minimum_required(VERSION 3.15)

project(Rewinder)

find_package(PkgConfig REQUIRED)
pkg_check_modules(LIBAV REQUIRED IMPORTED_TARGET
        libavdevice
        libavfilter
        libavformat
        libavcodec
        libswresample
        libswscale
        libavutil
        )

add_library(packet_reader src/packet_reader.c)
target_link_libraries(packet_reader PkgConfig::LIBAV)

add_executable(packet_read_test src/test/test.c)
target_link_libraries(packet_read_test packet_reader)

add_executable(print_error src/test/print_error.c)
target_link_libraries(print_error PkgConfig::LIBAV)
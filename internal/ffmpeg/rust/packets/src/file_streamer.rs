extern crate ffmpeg_next as ffmpeg;

use ffmpeg::Error;

struct FileStreamer {
    input: ffmpeg::format::Input,
    output: ffmpeg::format::Output,
}

impl FileStreamer {
    pub fn new(video_file: &str) -> Result<Self, Error> {
        let mut input = ffmpeg::format::input(video_file)?;
        let mut output = ffmpeg::format::output("test.rtp");
        output.
    }

    pub fn next(&self) -> Vec<u8> {

    }
}
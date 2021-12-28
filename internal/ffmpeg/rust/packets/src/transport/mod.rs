mod transport;

extern crate ffmpeg_next as ffmpeg;
use ffmpeg::{Error};

// The mode for the transport,
pub enum Mode {
    Input,
    Output,
}

/// Transport represents the outer most wrapper of a video asset
/// whether it be a an MPEGts container or an mp4 file
pub struct Transport {
    input: Option<ffmpeg::format::Input>,
    output: Option<ffmpeg::format::Output>,
    mode: Mode,
}
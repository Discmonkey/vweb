package com.example.rewinder.encode

import android.annotation.SuppressLint
import android.media.Image
import android.media.MediaCodec
import android.media.MediaCodecInfo
import android.media.MediaFormat
import androidx.camera.core.ImageAnalysis
import androidx.camera.core.ImageProxy

class EncoderAnalyzer(private val encoder: MediaCodec, private var ready: Boolean = false):
    ImageAnalysis.Analyzer {

    private fun setEncoderParams(width: Int, height: Int, channels: Int) {
        val bitrate = width * height * channels * 8;
        val mediaFormat = MediaFormat.createVideoFormat("video/avc", width,  height)
        mediaFormat.setInteger(
            MediaFormat.KEY_BITRATE_MODE,
            MediaCodecInfo.EncoderCapabilities.BITRATE_MODE_VBR)
        mediaFormat.setInteger(MediaFormat.KEY_BIT_RATE, bitrate)
        mediaFormat.setInteger(MediaFormat.KEY_FRAME_RATE, 10)
        mediaFormat.setInteger(
            MediaFormat.KEY_COLOR_FORMAT,
            MediaCodecInfo.CodecCapabilities.COLOR_FormatYUV420Flexible)
        mediaFormat.setInteger(MediaFormat.KEY_I_FRAME_INTERVAL, 1)
        mediaFormat.setString(
            MediaFormat.KEY_LATENCY,
            MediaCodecInfo.CodecCapabilities.FEATURE_LowLatency)
//        mediaCodec.setCallback(callback)
        encoder.configure(mediaFormat, null, null, MediaCodec.CONFIGURE_FLAG_ENCODE)
    }

    @SuppressLint("UnsafeOptInUsageError")
    private fun copyTo(from: ImageProxy, to: Image?) {
        if (from.image != null && to != null) {
            for (i in 0..from.planes.size) {
                to.planes[i].buffer.put(from.planes[i].buffer)
            }
        }
    }

    override fun analyze(imageProxy: ImageProxy) {
        val width = imageProxy.width
        val height = imageProxy.height
        val channels = imageProxy.planes.size
        if (!ready) {
            setEncoderParams(width, height, channels)
        }

        val inputBufIdx = encoder.dequeueInputBuffer(1000);
        val inputBuf = encoder.getInputImage(inputBufIdx)
        copyTo(imageProxy, inputBuf)

        encoder.queueInputBuffer(inputBufIdx, 0, width * height * channels,
            0, 0)

        imageProxy.close()
    }
}
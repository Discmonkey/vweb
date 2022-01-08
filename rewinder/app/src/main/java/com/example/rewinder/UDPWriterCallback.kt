package com.example.rewinder

import android.media.MediaCodec
import android.media.MediaCodec.Callback
import android.media.MediaFormat
import android.util.Log
import java.io.BufferedOutputStream

class UDPWriterCallback(private val outputStream: BufferedOutputStream): Callback() {

    override fun onInputBufferAvailable(p0: MediaCodec, p1: Int) {
        Log.d("onInputBufferAvailable", p1.toString())
    }

    override fun onOutputBufferAvailable(p0: MediaCodec, p1: Int, p2: MediaCodec.BufferInfo) {
        Log.d("onOutputBufferAvailable", p2.toString())
    }

    override fun onError(p0: MediaCodec, p1: MediaCodec.CodecException) {
        Log.d("onError", p1.toString())

    }

    override fun onOutputFormatChanged(p0: MediaCodec, p1: MediaFormat) {
        Log.d("onOutputFormatChanged", "output format changed")
    }
}
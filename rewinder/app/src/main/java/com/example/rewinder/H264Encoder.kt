package com.example.rewinder

import android.media.MediaCodec
import android.media.MediaCodecInfo
import android.media.MediaFormat
import android.util.Log
import android.view.Surface
import android.view.SurfaceView
import androidx.camera.core.Preview
import androidx.camera.core.SurfaceRequest
import androidx.core.util.Consumer
import java.io.BufferedOutputStream
import java.lang.Exception
import java.util.concurrent.Executor

class H264Encoder(private val outputStream: BufferedOutputStream,
                  callback: MediaCodec.Callback
                  ) : Preview.SurfaceProvider {

    private val mediaCodec: MediaCodec = MediaCodec.createEncoderByType("video/avc")
    private var inputSurface: Surface? = null;
    init {
        val mediaFormat = MediaFormat.createVideoFormat("video/avc", 320, 240)
        mediaFormat.setInteger(MediaFormat.KEY_BIT_RATE, 125000)
        mediaFormat.setInteger(MediaFormat.KEY_FRAME_RATE, 15)
        mediaFormat.setInteger(MediaFormat.KEY_COLOR_FORMAT,
            MediaCodecInfo.CodecCapabilities.COLOR_FormatYUV420Flexible)
        mediaFormat.setInteger(MediaFormat.KEY_I_FRAME_INTERVAL, 5);
        mediaCodec.setCallback(callback)
        mediaCodec.configure(mediaFormat, null, null, MediaCodec.CONFIGURE_FLAG_ENCODE);
        inputSurface = mediaCodec.createInputSurface()
    }

    private fun start() {
        mediaCodec.start()
    }

    fun stop() {
        mediaCodec.stop()
    }

    private fun getSurface(): Surface? {
        return inputSurface
    }

    fun close() {
        try {
            mediaCodec.stop();
            mediaCodec.release();
            outputStream.flush();
            outputStream.close();
        } catch (e: Exception){
            e.printStackTrace();
        }
    }

    class SimpleExecutor: Executor {
        override fun execute(p0: Runnable?) {
            Log.d("simple consumer", "starting to execute")
            Thread(p0).start()
        }
    }

    class SimpleConsumer: Consumer<SurfaceRequest.Result> {
        override fun accept(t: SurfaceRequest.Result?) {
           Log.d("simple consumer", t.toString())
        }

    }

    override fun onSurfaceRequested(request: SurfaceRequest) {
        Log.d("onSurfaceRequested", "called")
        getSurface()?.let { request.provideSurface(it, SimpleExecutor(), SimpleConsumer()) }
        this.start()
    }
}
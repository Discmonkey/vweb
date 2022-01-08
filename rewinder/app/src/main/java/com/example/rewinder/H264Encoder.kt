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

    fun start() {
        mediaCodec.start()
    }

    fun stop() {
        mediaCodec.stop()
    }

    fun getSurface(): Surface? {
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

    fun onBuffer() {

    }
    fun encode(input: ByteArray) {
        try {
            val inputBufferIndex = mediaCodec.dequeueInputBuffer(-1);
            if (inputBufferIndex >= 0) {
                val inputBuffer = mediaCodec.getInputBuffer(inputBufferIndex)
                inputBuffer?.clear()
                inputBuffer?.put(input)
                mediaCodec.queueInputBuffer(inputBufferIndex, 0,
                    input.size, 0, 0)
            }
            val bufferInfo = MediaCodec.BufferInfo()
            var outputBufferIndex = mediaCodec.dequeueOutputBuffer(bufferInfo, 0)
            var data = ByteArray(0)

            while (outputBufferIndex >= 0) {
                val outputBuffer = mediaCodec.getOutputBuffer(outputBufferIndex)
                if (data.size < bufferInfo.size) {
                    data = ByteArray((bufferInfo.size))
                }
                outputBuffer?.get(data);
                outputStream.write(data, 0, data.size);

                mediaCodec.releaseOutputBuffer(outputBufferIndex, false);
                outputBufferIndex = mediaCodec.dequeueOutputBuffer(bufferInfo, 0);
            }
        } catch (t: Throwable) {
            t.printStackTrace();
        }

    }

    class SimpleExecutor: Executor {
        override fun execute(p0: Runnable?) {
            Log.d("executor", p0.toString())
        }
    }

    class SimpleConsumer: Consumer<SurfaceRequest.Result> {
        override fun accept(t: SurfaceRequest.Result?) {
           Log.d("simple consumer", t.toString())
        }

    }

    override fun onSurfaceRequested(request: SurfaceRequest) {
        Log.d("onSurfaceRequestd", "called")
        getSurface()?.let { request.provideSurface(it, SimpleExecutor(), SimpleConsumer()) }

    }
}
package com.example.rewinder

import android.media.MediaCodec
import android.media.MediaCodecInfo
import android.media.MediaFormat
import android.os.Build
import android.provider.MediaStore.Audio.Media
import android.util.Log
import android.view.Surface
import androidx.annotation.RequiresApi
import androidx.camera.core.Preview
import androidx.camera.core.SurfaceRequest
import androidx.core.util.Consumer
import java.io.BufferedOutputStream
import java.lang.Exception
import java.util.concurrent.Executor
import java.util.concurrent.Executors

@RequiresApi(Build.VERSION_CODES.S_V2)
class H264Encoder(private val outputStream: BufferedOutputStream,
                  callback: MediaCodec.Callback
                  ) : Preview.SurfaceProvider {

    private val mediaCodec: MediaCodec = MediaCodec.createEncoderByType("video/avc")
    private var inputSurface: Surface? = null;

    companion object {
        private const val TAG = "H264Encoder"
        private const val MIME_TYPE = "video/avc" // H.264
        private const val DEFAULT_WIDTH = 640
        private const val DEFAULT_HEIGHT = 360
        private const val DEFAULT_FPS = 30
        // Adjust I-Frame interval. 1 means every frame is an I-frame (good for seeking, bad for compression).
        // 5-10 seconds is more common for streaming. For 30 FPS, 150-300.
        // For quick recovery/seeking in a "rewinder" app, a shorter interval like 1-2 seconds (30-60 frames) might be better.
        private const val DEFAULT_I_FRAME_INTERVAL_SECONDS = 2
    }
    init {
        val mediaFormat = MediaFormat.createVideoFormat(MIME_TYPE, DEFAULT_WIDTH, DEFAULT_HEIGHT)
        mediaFormat.setInteger(MediaFormat.KEY_BITRATE_MODE,
            MediaCodecInfo.EncoderCapabilities.BITRATE_MODE_VBR)
//        mediaFormat.setInteger(MediaFormat.KEY_BIT_RATE, width * height * 3 * 8 * 30)
        mediaFormat.setInteger(MediaFormat.KEY_FRAME_RATE, DEFAULT_FPS)
        mediaFormat.setInteger(MediaFormat.KEY_COLOR_FORMAT,
            MediaCodecInfo.CodecCapabilities.COLOR_FormatSurface)
        mediaFormat.setInteger(MediaFormat.KEY_I_FRAME_INTERVAL, 150)
        mediaFormat.setString(MediaFormat.KEY_LATENCY,
            MediaCodecInfo.CodecCapabilities.FEATURE_LowLatency)

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) { // FEATURE_LowLatency constant available from API 26
            val caps = MediaCodecInfo.CodecCapabilities.FEATURE_LowLatency
            // Check if supported before setting, though setString might not throw immediately
            // mediaFormat.setInteger(MediaFormat.KEY_LOW_LATENCY, 1); // Alternative way to request low latency
            mediaFormat.setFeatureEnabled(caps, true) // Preferred way if available
        } else {
            mediaFormat.setString(MediaFormat.KEY_LATENCY, "1") // Older way, might not be as effective as FEATURE_LowLatency
        }

        mediaCodec.setCallback(callback)
        mediaCodec.configure(mediaFormat, null, null, MediaCodec.CONFIGURE_FLAG_ENCODE)
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
        private val executor =
            Executors.newSingleThreadExecutor { r -> Thread(r, "SurfaceProviderExecutor") }
        override fun execute(command: Runnable?) {
            if (command != null) {
                Log.d(TAG, "SimpleExecutor executing command.")
                executor.execute(command)
            }
        }
        fun shutdown() { // Good practice to allow shutting down the executor
            executor.shutdown()
        }
    }

    class SimpleConsumer: Consumer<SurfaceRequest.Result> {
        override fun accept(result: SurfaceRequest.Result?) {
            when (result?.resultCode) {
                SurfaceRequest.Result.RESULT_SURFACE_USED_SUCCESSFULLY ->
                    Log.d(TAG, "Surface successfully provided and used.")
                SurfaceRequest.Result.RESULT_INVALID_SURFACE ->
                    Log.e(TAG, "Surface provision failed: Invalid surface.")
                SurfaceRequest.Result.RESULT_REQUEST_CANCELLED ->
                    Log.w(TAG, "Surface provision failed: Request cancelled.")
                SurfaceRequest.Result.RESULT_SURFACE_ALREADY_PROVIDED ->
                    Log.w(TAG, "Surface provision failed: Surface already provided.")
                SurfaceRequest.Result.RESULT_WILL_NOT_PROVIDE_SURFACE ->
                    Log.e(TAG, "Surface provision failed: Will not provide surface.")
                else -> Log.d(TAG, "Surface provision result: $result")
            }
        }
    }

    override fun onSurfaceRequested(request: SurfaceRequest) {
        getSurface()?.let { request.provideSurface(it, SimpleExecutor(), SimpleConsumer()) }
        this.start()
    }
}
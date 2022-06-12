package com.example.rewinder

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import androidx.annotation.RequiresApi
import androidx.camera.core.CameraSelector
import androidx.camera.core.Preview
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.core.content.ContextCompat
import com.google.android.material.slider.Slider
import io.swagger.server.models.Address
import net.pwall.json.parseJSON
import java.io.BufferedOutputStream
import java.util.concurrent.ExecutorService


class StreamActivity : AppCompatActivity() {
    private lateinit var cameraExecutor: ExecutorService
    @RequiresApi(32)
    private var permissionManager = PermissionManager()
    @RequiresApi(32)
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_stream)

        val address: Address? = intent.getStringExtra("streamAddress")?.parseJSON()
        // Request camera permissions
        if (address != null && permissionManager.allPermissionsGranted(baseContext)) {
            val udpOutputStream = UDPOutputStream(address.ip, address.port)
            val encoder  = H264Encoder(BufferedOutputStream(udpOutputStream),
                UDPWriterCallback(BufferedOutputStream(udpOutputStream)))
            startCamera(encoder)
        } else {
            startActivity(Intent(baseContext, ConnectActivity::class.java))
        }
    }

    @RequiresApi(32)
    private fun startCamera(encoder: H264Encoder) {
        val cameraProviderFuture = ProcessCameraProvider.getInstance(this)

        cameraProviderFuture.addListener(Runnable {
            // Used to bind the lifecycle of cameras to the lifecycle owner
            val cameraProvider: ProcessCameraProvider = cameraProviderFuture.get()
            val viewFinder: androidx.camera.view.PreviewView = findViewById(R.id.viewFinder);

            // Preview
            val preview = Preview.Builder()
                .build()
                .also {
                    it.setSurfaceProvider(viewFinder.surfaceProvider)
                }

            val encoderPreview  = Preview.Builder()
                .build()
                .also {
                    it.setSurfaceProvider(
                        encoder
                    )
                }

            // Select back camera as a default
            val cameraSelector = CameraSelector.DEFAULT_BACK_CAMERA

            try {
                // Unbind use cases before rebinding
                cameraProvider.unbindAll()

                // Bind use cases to camera
                val camera = cameraProvider.bindToLifecycle(
                    this, cameraSelector, preview, encoderPreview)
                findViewById<Slider>(R.id.zoom).addOnChangeListener { _, value, _ ->
                    camera.cameraControl.setLinearZoom(value / 100)
                    println(camera.cameraInfo.sensorRotationDegrees)
                }
            } catch (exc: Exception) {
                Log.e("StreamActivity", "Use case binding failed", exc)
            }

        }, ContextCompat.getMainExecutor(this))
    }

    override fun onDestroy() {
        super.onDestroy()
        cameraExecutor.shutdown()
    }
}
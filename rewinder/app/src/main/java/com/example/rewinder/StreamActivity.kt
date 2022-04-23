package com.example.rewinder

import android.Manifest
import android.content.pm.PackageManager
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.widget.Toast
import androidx.annotation.RequiresApi
import androidx.camera.core.CameraSelector
import androidx.camera.core.Preview
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import java.io.BufferedOutputStream
import java.util.concurrent.ExecutorService
import java.util.concurrent.Executors


class StreamActivity : AppCompatActivity() {
    private lateinit var cameraExecutor: ExecutorService
    private var udpOutputStream = UDPOutputStream("192.168.1.11", 9000)
    @RequiresApi(32)
    private var encoder = H264Encoder(BufferedOutputStream(udpOutputStream),
        UDPWriterCallback(BufferedOutputStream(udpOutputStream)))
    private var permissionManager = PermissionManager()

    @RequiresApi(32)
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_stream)

        // Request camera permissions
        if (permissionManager.allPermissionsGranted(baseContext)) {
            startCamera()
        } else {
            permissionManager.requestPermissions(this);
        }
    }

    @RequiresApi(32)
    private fun startCamera() {
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
                cameraProvider.bindToLifecycle(
                    this, cameraSelector, preview, encoderPreview
                )


            } catch (exc: Exception) {
                Log.e("StreamActivity", "Use case binding failed", exc)
            }

        }, ContextCompat.getMainExecutor(this))
    }

    override fun onDestroy() {
        super.onDestroy()
        cameraExecutor.shutdown()
    }

    @RequiresApi(32)
    override fun onRequestPermissionsResult(
        requestCode: Int, permissions: Array<String>, grantResults:
        IntArray) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == PermissionManager.REQUEST_CODE_PERMISSIONS) {
            if (permissionManager.allPermissionsGranted(baseContext)) {
                startCamera()
            } else {
                Toast.makeText(this,
                    "Permissions not granted by the user.",
                    Toast.LENGTH_SHORT).show()
                finish()
            }
        }
    }


}
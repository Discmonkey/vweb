package com.example.rewinder

import android.Manifest
import android.content.Context
import android.content.pm.PackageManager
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat

class PermissionManager {

    fun allPermissionsGranted(context: Context) = REQUIRED_PERMISSIONS.all {
        ContextCompat.checkSelfPermission(
            context, it) == PackageManager.PERMISSION_GRANTED
    }

    fun requestPermissions(activity: AppCompatActivity) {
        ActivityCompat.requestPermissions(
            activity, REQUIRED_PERMISSIONS, REQUEST_CODE_PERMISSIONS)
    }

    companion object {
        private const val TAG = "CameraXBasic"
        const val REQUEST_CODE_PERMISSIONS = 7
        private val REQUIRED_PERMISSIONS = arrayOf(Manifest.permission.CAMERA)
    }
}
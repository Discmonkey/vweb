package com.example.rewinder

import android.os.Bundle
import androidx.annotation.RequiresApi
import androidx.appcompat.app.AppCompatActivity
import io.swagger.server.models.Address

class ConnectActivity : AppCompatActivity() {
    @RequiresApi(32)
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_connect)
    }

}
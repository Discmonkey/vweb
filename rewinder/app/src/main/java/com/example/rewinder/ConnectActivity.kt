package com.example.rewinder

import android.os.Bundle
import android.util.Log
import android.widget.Button
import android.widget.EditText
import android.widget.TextView
import androidx.annotation.RequiresApi
import androidx.appcompat.app.AppCompatActivity
import com.example.rewinder.network.createSource
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking

class ConnectActivity : AppCompatActivity() {
    @RequiresApi(32)
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_connect)

        val submit = findViewById<Button>(R.id.connect_button)
        val stream = findViewById<Button>(R.id.stream_button)
        val text = findViewById<EditText>(R.id.ip_field)
        val status = findViewById<TextView>(R.id.status)

        submit.setOnClickListener {
            runBlocking {
                launch {
                    val ip = text.text.toString()
                    createSource(ip)
                }
            }
        }
    }
}
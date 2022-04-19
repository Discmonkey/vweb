package com.example.rewinder

import android.os.Bundle
import android.util.Log
import android.widget.Button
import android.widget.EditText
import android.widget.TextView
import androidx.annotation.RequiresApi
import androidx.appcompat.app.AppCompatActivity
import io.swagger.server.models.Address
import io.swagger.server.models.Source

class ConnectActivity : AppCompatActivity() {
    @RequiresApi(32)
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_connect)

        val submit = findViewById<Button>(R.id.connect_button)
        val text = findViewById<EditText>(R.id.ip_field)
        val status = findViewById<TextView>(R.id.status)

        submit.setOnClickListener {
            val ip = text.text.toString()
            val source = Source()

            Log.e("test", source.toString())
        }
    }
}
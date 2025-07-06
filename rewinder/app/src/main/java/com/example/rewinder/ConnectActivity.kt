package com.example.rewinder

import android.content.Intent
import android.os.Bundle
import android.widget.Button
import android.widget.EditText
import android.widget.TextView
import android.widget.Toast
import androidx.annotation.RequiresApi
import androidx.appcompat.app.AppCompatActivity
import com.example.rewinder.network.createSource
import com.github.kittinunf.result.Result
import io.swagger.server.models.Address
import net.pwall.json.parseJSON
import net.pwall.json.stringifyJSON


class ConnectActivity : AppCompatActivity() {

    private var permissionManager = PermissionManager()
    private var address: Address? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_connect)
        if (!permissionManager.allPermissionsGranted(baseContext)) {
            permissionManager.requestPermissions(this);
        }
        val submit = findViewById<Button>(R.id.connect_button)
        val stream = findViewById<Button>(R.id.stream_button)
        val text = findViewById<EditText>(R.id.ip_field)
        val status = findViewById<TextView>(R.id.status)

        submit.setOnClickListener {
            val ip = text.text.toString()
            createSource(ip).responseString {
                _, _, result ->
                    when (result) {
                        is Result.Failure -> {
                            status.text = result.toString()
                        }

                        is Result.Success -> {
                            this.address = result.get().parseJSON()
                            this.address = Address(ip, this.address!!.port)
                        }
                    }
            }
        }

        stream.setOnClickListener {
            if (this.address == null) {
                Toast.makeText(this,
                    "server not configured",
                    Toast.LENGTH_SHORT).show()
            } else {
                val intent = Intent(baseContext, StreamActivity::class.java)
                intent.putExtra("streamAddress", this.address.stringifyJSON())
                startActivity(intent)
            }
        }
    }

    override fun onRequestPermissionsResult(
        requestCode: Int, permissions: Array<String>, grantResults:
        IntArray) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == PermissionManager.REQUEST_CODE_PERMISSIONS) {
           if (!permissionManager.allPermissionsGranted(baseContext)) {
                Toast.makeText(this,
                    "Permissions not granted by the user.",
                    Toast.LENGTH_SHORT).show()
                finish()
           }
        }
    }
}
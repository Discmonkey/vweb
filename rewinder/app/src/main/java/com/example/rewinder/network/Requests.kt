package com.example.rewinder.network

import android.util.Log
import io.swagger.server.models.Address
import io.swagger.server.models.IP
import io.swagger.server.models.Source
import org.json.JSONObject
import org.json.JSONTokener
import java.io.*
import java.net.HttpURLConnection
import java.net.URL


// add source to server, and return the address to which we want to stream
suspend fun createSource(address: String): Address {

    runCatching {
        val url  = URL("$address/source")
        val urlConnection: HttpURLConnection = url.openConnection() as HttpURLConnection
        val out = BufferedOutputStream(urlConnection.outputStream)
        val writer = BufferedWriter(OutputStreamWriter(out, "UTF-8"))
        writer.write(Source("video/H264").toString())
        writer.flush()
        writer.close()
        out.close()

        urlConnection.connect()

        val input  = BufferedInputStream(urlConnection.inputStream)
        val reader = BufferedReader(InputStreamReader(input, "UTF-8"))
        val returns = JSONObject(JSONTokener(reader.toString()))

        Log.e("returned:", returns.toString())
    }

    return Address("hello", 3000)
}
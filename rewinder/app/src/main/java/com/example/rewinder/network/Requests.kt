package com.example.rewinder.network

import io.swagger.server.models.Address
import io.swagger.server.models.IP
import io.swagger.server.models.Source
import java.io.BufferedInputStream
import java.io.InputStream
import java.net.HttpURLConnection
import java.net.URL


// add source to server, and return the address to which we want to stream
fun create_source(ip: IP, source: Source): Address {
    val url  = URL("$ip:3000")
    val urlConnection: HttpURLConnection = url.openConnection() as HttpURLConnection
    urlConnection.requestMethod = "POST"

    try {
        val response: InputStream = BufferedInputStream(urlConnection.inputStream)
    } finally {
        urlConnection.disconnect()
    }

    return Address(ip, 3000)
}
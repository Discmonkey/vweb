package com.example.rewinder.network

import android.util.Log
import io.swagger.server.models.Address
import com.github.kittinunf.fuel.httpPost
import com.github.kittinunf.result.Result
import java.io.*
import java.net.HttpURLConnection
import java.net.URL

// add source to server, and return the address to which we want to stream
fun createSource(address: String): Address {
    val bodyJson = """
        { 
            "codec" : "video/H264",
            "name" : "max's phone"
        }
    """
    val httpAsync = "http://192.168.1.10:3000/source"
        .httpPost()
        .body(bodyJson)
        .responseString { request, response, result ->
            when (result) {
                is Result.Failure -> {
                    val ex = result.getException()
                    println(ex)
                }
                is Result.Success -> {
                    val data = result.get()
                    println(data)
                }
            }
        }
//    runCatching {
//        val url = URL("http://192.168.1.10:3000/source")
//        val postData = "{\"codec\":\"video/H264\"}"
//
//        val conn = url.openConnection()
//        conn.doOutput = true
//        conn.setRequestProperty("Content-Type", "application/json")
//        conn.setRequestProperty("Content-Length", postData.length.toString())
//
//        val stream = conn.getOutputStream();
//        println("here")
//        DataOutputStream(stream).use {
//            it.writeBytes(postData)
//        }
//        println("here")
//        BufferedReader(InputStreamReader(conn.getInputStream())).use { bf ->
//            var line: String?
//            while (bf.readLine().also { line = it } != null) {
//                println(line)
//            }
//        }
//    }

    return Address("hello", 3000)
}
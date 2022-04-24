package com.example.rewinder.network

import com.github.kittinunf.fuel.core.Request
import com.github.kittinunf.fuel.httpPost

// add source to server, and return the address to which we want to stream
fun createSource(ip: String): Request {
    val bodyJson = """
        { 
            "codec" : "video/H264",
            "name" : "max's phone"
        }
    """
    return "http://$ip:3000/source"
        .httpPost()
        .body(bodyJson)
}
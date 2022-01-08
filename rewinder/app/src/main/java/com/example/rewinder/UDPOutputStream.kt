package com.example.rewinder

import android.util.Log
import java.io.OutputStream
import java.net.DatagramPacket
import java.net.DatagramSocket
import java.net.Inet4Address
import kotlin.math.min

class UDPOutputStream(address: Inet4Address, port : Int) : OutputStream() {
    private val socket = DatagramSocket(port, address)
    private val maxSendLength = socket.sendBufferSize

    override fun write(p0: Int) {
        Log.d("udp writer", "int write called, not currently supported")
    }

    override fun write(b: ByteArray?) {
        if (b != null) {
            safeSend(b, 0, b.size)
        }
    }

    override fun write(b: ByteArray?, off: Int, len: Int) {
        safeSend(b, off, len)
    }

    // no need to flush UDP socket, its basically autoflushed
    override fun flush() {}

    override fun close() {
        socket.close()
    }

    private fun safeSend(b: ByteArray?, offset: Int, length: Int) {
        var stillNeedToSend = length
        var currentOffset = offset
        while (stillNeedToSend > 0) {
            val sendLength = min(stillNeedToSend, maxSendLength)
            val packet = DatagramPacket(b, currentOffset, sendLength)
            socket.send(packet)

            stillNeedToSend -= sendLength
            currentOffset += sendLength
        }
    }
}
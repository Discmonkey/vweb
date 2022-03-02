package com.example.rewinder

import android.util.Log
import java.io.OutputStream
import java.net.DatagramPacket
import java.net.DatagramSocket
import java.net.Inet4Address
import java.util.concurrent.BlockingQueue
import java.util.concurrent.LinkedBlockingDeque
import kotlin.math.min

class UDPOutputStream(address: String, port : Int) : OutputStream() {
    private val queue: BlockingQueue<DatagramPacket> = LinkedBlockingDeque();
    private val task = ThreadedWriter(address, port, queue)
    private val maxSendLength = 256
    private val thread = Thread(task)
    init {
        thread.start()
    }

    class ThreadedWriter(private val address: String,
                         private val port : Int,
                         private val queue: BlockingQueue<DatagramPacket>): Runnable {
        @Volatile
        var running = true
        override fun run() {
            Log.d("thread", Thread.currentThread().name)
            val ipv4 = Inet4Address.getByName(address) as Inet4Address
            val socket = DatagramSocket()
            socket.connect(ipv4, port)

            while (running) {
                val next = queue.take()
                next.address = ipv4
                next.port = port
                socket.send(next)
            }

            socket.close()
        }
    }


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
        task.running = false
    }

    private fun safeSend(b: ByteArray?, offset: Int, length: Int) {
        var stillNeedToSend = length
        var currentOffset = offset
        while (stillNeedToSend > 0) {
            val sendLength = min(stillNeedToSend, maxSendLength)
            val packet = DatagramPacket(b, currentOffset, sendLength)
            queue.put(packet)
            stillNeedToSend -= sendLength
            currentOffset += sendLength
        }
    }
}
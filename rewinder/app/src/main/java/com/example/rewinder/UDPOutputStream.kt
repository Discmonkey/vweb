package com.example.rewinder

import android.util.Log
import java.io.OutputStream
import java.net.Inet4Address
import java.net.InetSocketAddress
import java.net.Socket
import java.util.concurrent.BlockingQueue
import java.util.concurrent.LinkedBlockingDeque
import kotlin.math.min

class UDPOutputStream(address: String, port : Int) : OutputStream() {
    private val queue: BlockingQueue<ByteArray> = LinkedBlockingDeque();
    private val task = ThreadedWriter(address, port, queue)
    private val maxSendLength = 1024
    private val thread = Thread(task)
    init {
        thread.start()
    }

    class ThreadedWriter(private val address: String,
                         private val port : Int,
                         private val queue: BlockingQueue<ByteArray>): Runnable {
        @Volatile
        var running = true
        override fun run() {
            val socket = Socket()
            socket.connect(InetSocketAddress(address, port))
            val out = socket.getOutputStream()
            while (running) {
                val next = queue.take()
                out.write(next)
            }

            socket.close()
        }
    }


    override fun write(p0: Int) {
        Log.d("udp writer", "int write called, not currently supported")
    }

    override fun write(b: ByteArray?) {
        if (b != null) {
            safeSend(b)
        }
    }

    override fun write(b: ByteArray?, off: Int, len: Int) {
        safeSend(b)
    }

    // no need to flush UDP socket, its basically autoflushed
    override fun flush() {}

    override fun close() {
        task.running = false
    }

    private fun safeSend(b: ByteArray?) {
        queue.put(b)
    }
}
# ByteBuffer & Channel

### 1. ByteBuffer
> ByteBuffer 用来缓冲读写数据，常见的有

    * ByteBuffer
      * MappedByteBuffer
      * DirectByteBuffer
      * HeapByteBuffer
    * ShortBuffer
    * IntBuffer
    * LongBuffer
    * FloatBuffer
    * DoubleBuffer
    * CharBuffer

### 2. Channel 
> Channel 有一点类似于 stream，它就是读写数据的**双向通道**，可以从 channel 将数据读入 buffer，也可以将 buffer 的数据写入 channel，而之前的 stream 要么是输入，要么是输出，channel 比 stream 更为底层，常见的有

    * FileChannel
    * DatagramChannel
    * SocketChannel
    * ServerSocketChannel

### 3. 使用 FileChannel 来读取文件内容
```java
// FileChannel
// 1.In/OutputStream 2. RandomAccessFile
try(var channel = new FileInputStream("file/c1/data.txt").getChannel()) {
    // 准备缓冲区
    var buffer = ByteBuffer.allocate(10);
    // 从 channel 中读取数据，向缓冲区 buffer 写入, 当len为-1时没有读取结束
    var len = channel.read(buffer);
    while (len != -1) {
        log.warn("read len: {}", len);
        // 输出 buffer 的内容
        buffer.flip(); // 切换到buffer的读模式 limit -> pos; pos -> 0;
        // 检查是否剩余
        while (buffer.hasRemaining()) {
            var b = buffer.get();
            log.debug("result: {}", (char) b);
        }
        buffer.clear(); // 切换到写模式 pos -> 0; limit -> capacity
        len = channel.read(buffer);
    }

} catch (IOException e) {
    e.printStackTrace();
}
```
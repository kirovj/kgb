> 翻译自[kingori.co](https://kingori.co/minutae/2013/04/httpurlconnection-disconnect/)

额，虽然要看情况，但是一般来说...

no.

是的，不需要。

see the [javadoc](https://docs.oracle.com/javase/7/docs/api/java/net/HttpURLConnection.html)

> Each HttpURLConnection instance is used to make a single request but the underlying network connection to the HTTP server may be transparently shared by other instances. Calling the close() methods on the InputStream or OutputStream of an HttpURLConnection after a request may free network resources associated with this instance but has no effect on any shared persistent connection. Calling the disconnect() method may close the underlying socket if a persistent connection is otherwise idle at that time.

> 每个HttpURLConnection实例都用于发出单个请求，但其他实例可以透明地共享与HTTP服务器的基础网络连接。 请求后在HttpURLConnection的InputStream或OutputStream上调用close()方法可能会释放与此实例关联的网络资源，但对任何共享的持久连接都没有影响。***如果持久连接当时处于空闲状态，则调用disconnect()方法可能会关闭基础Socket***。

简单来说：

> I don’t say it is a mistake. But, disconnect is extreme case (socket close open operations are costly), unless you really want I wouldn’t go for it. stream.close() releases most of the network resources and should be enough. Again, if your requirement is ok to create socket everytime, there is nothing wrong in calling disconnect.  ~ [Nambari](https://stackoverflow.com/questions/11056088/do-i-need-to-call-httpurlconnection-disconnect-after-finish-using-it#comment14465352_11056207)

> 我不是说这是错误的。但是断开连接是极端情况（***Socket关闭打开操作的成本很高***），除非您真的希望，而我不会这么做。stream.close()释放了大部分网络资源，应该足够了。同样，如果您的需求是每次都创建Socket，那就调用disconnect()吧。

例子：

```java
//Initialize variables
InputStream responseInputStream = null;
HttpURLConnection conn = null;
int responseCode = -1;

try {
	// Create a connection
	URL url = new URL(targetURL);
	conn = (HttpURLConnection) url.openConnection();
	
	// Starts the connection
	conn.connect();

	// Get the response code
	responseCode = conn.getResponseCode();
	if (responseCode < 400) {
		// Get the response InputStream if all is well
		responseInputStream = conn.getInputStream();
	} else if (responseCode >= 400) {
		// Get the response ErrorStream if we got an error instead
		responseInputStream = conn.getErrorStream();
	}
} catch (IOException e) {
	// YOU CAN DO WHATEVER HERE
} finally {
	// Makes sure that the InputStream is closed after
	// we are done with it
	if (responseInputStream != null) {
		try {
			responseInputStream.close();
		} catch (IOException e) {
			// YOU CAN DO WHATEVER HERE
		}
	}
}
```

***And that is enough.***

thanks:

- [Do I need to call HttpURLConnection.disconnect after finish using it](https://stackoverflow.com/questions/11056088/do-i-need-to-call-httpurlconnection-disconnect-after-finish-using-it)
- [Do we need to call HttpURLConnection.disconnect()?](https://kingori.co/minutae/2013/04/httpurlconnection-disconnect/)
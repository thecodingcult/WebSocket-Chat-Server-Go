<!DOCTYPE html>
<html>
  <head>
   
  </head>
  <body>
    <h1>WebSocket Chat Server</h1>
    <p>This is a simple WebSocket-based chat server written in Go. Users can connect to the server and send text messages to all other connected users.</p>
    <h2>Usage</h2>
<ol>
  <li>Clone the repository to your local machine.</li>
  <li>Open a terminal and navigate to the project's root directory.</li>
  <li>Run the command <code>go run main.go</code> to start the server.</li>
  <li>Open a web browser and navigate to <a href="http://localhost:8080">http://localhost:8080</a>.</li>
  <li>Enter your name and click "Connect".</li>
  <li>Start chatting with other users!</li>
</ol>

<h2>Dependencies</h2>
<p>This project relies on the following dependencies:</p>
<ul>
  <li><a href="https://github.com/gorilla/websocket">github.com/gorilla/websocket</a></li>
</ul>
<p>Make sure you have these dependencies installed before running the server.</p>

<h2>Code Overview</h2>
<p>The main logic of the chat server is implemented in the <code>Server</code> type, which holds a map of connections, a channel for broadcasting messages, and channels for registering and unregistering connections.</p>
<p>The <code>ServeHTTP</code> method of the <code>Server</code> type handles incoming HTTP requests and upgrades them to WebSocket connections. When a new connection is established, it is registered with the <code>Server</code>.</p>
<p>The <code>run</code> method of the <code>Server</code> type is the main loop for the chat server. It listens for new connections, closed connections, and incoming messages, and broadcasts messages to all connected clients.</p>
<p>The <code>Connection</code> type holds a channel for sending messages and a WebSocket connection. It also has writer and reader methods that handle sending and receiving messages, respectively.</p>

<h2>Conclusion</h2>
<p>This is a simple implementation of a WebSocket-based chat server in Go. Feel free to use it as a starting point for your own projects!</p>
  </body>
</html>

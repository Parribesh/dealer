import { useState, useEffect, useCallback, useRef } from "react";

const useWebSocket = (url) => {
  const [ws, setWs] = useState(null);
  const [lastMessage, setLastMessage] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const [userId, setUserId] = useState(null);
  const reconnectAttempts = useRef(0);
  const maxReconnectAttempts = 5;
  const reconnectTimeoutId = useRef(null);
  const tokenRef = useRef(null);

  const connect = (token) => {
    if (ws) {
      console.log("WebSocket already connected. Disconnecting first.");
      ws.close();
    }

    console.log("Connecting to WebSocket...");
    const websocket = new WebSocket(url);
    // setWs(websocket);
    tokenRef.current = token;

    websocket.onopen = () => {
      console.log("WebSocket connected successfully");
      // Send the token message immediately after connection
      const tokenMessage = JSON.stringify({ token: tokenRef.current });
      websocket.send(tokenMessage);
      setIsConnected(true);
      reconnectAttempts.current = 0;
    };

    websocket.onmessage = (event) => {
      console.log("WebSocket message received:", event.data);
      try {
        const message = JSON.parse(event.data);
        setLastMessage(message);
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    };

    websocket.onclose = (event) => {
      console.log(
        `WebSocket disconnected. Code: ${event.code}, Reason: ${event.reason}`
      );
      setIsConnected(false);
      if (event.code === 1005) {
        console.log(
          "Received close 1005 (no status). Attempting immediate reconnection."
        );
        reconnectAttempts.current += 1;
        if (reconnectAttempts.current < maxReconnectAttempts) {
          connect(tokenRef.current);
        } else {
          console.log("Max reconnect attempts reached. Stopping reconnection.");
        }
      } else {
        reconnectAttempts.current += 1;
        if (reconnectAttempts.current < maxReconnectAttempts) {
          const reconnectDelay = Math.min(
            1000 * 2 ** reconnectAttempts.current,
            30000
          );
          reconnectTimeoutId.current = setTimeout(
            () => connect(tokenRef.current),
            reconnectDelay
          );
        } else {
          console.log("Max reconnect attempts reached. Stopping reconnection.");
        }
      }
    };

    websocket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  };

  const disconnect = useCallback(() => {
    if (ws) {
      console.log("Closing WebSocket connection");
      ws.close();
    }
    if (reconnectTimeoutId.current) {
      clearTimeout(reconnectTimeoutId.current);
    }
    setWs(null);
    setIsConnected(false);
    tokenRef.current = null;
  }, [ws]);

  useEffect(() => {
    return () => {
      disconnect();
    };
  }, [disconnect]);

  const sendMessage = useCallback(
    (message) => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(message));
      } else {
        console.error("WebSocket is not open. Unable to send message.");
      }
    },
    [ws]
  );

  return {
    isConnected,
    lastMessage,
    connect,
    disconnect,
    sendMessage,
    setUserId,
    userId,
  };
};

export default useWebSocket;

import { useState, useEffect, useCallback, useRef } from "react";
import { createComponentLogger } from "../logger";

const log = createComponentLogger("useWebSocket", "info")

const useWebSocket = (url) => {
  const [ws, setWs] = useState(null);
  const [lastMessage, setLastMessage] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const [userId, setUserId] = useState(null);
  const reconnectAttempts = useRef(0);
  const maxReconnectAttempts = 5;
  const reconnectTimeoutId = useRef(null);
  const tokenRef = useRef(null);
  const isConnectingRef = useRef(false); // Track if currently connecting

  const connect = (token) => {
    // Prevent multiple connection attempts
    if (isConnectingRef.current) return;
    isConnectingRef.current = true;

    if (ws) {
      log.warn("WebSocket already connected. Disconnecting first.");
      ws.close();
    }

    log.info("Connecting to WebSocket...");
    const websocket = new WebSocket(url);
    setWs(websocket);
    tokenRef.current = token;

    websocket.onopen = () => {
      log.info("WebSocket connected successfully");
      const tokenMessage = JSON.stringify({ token: tokenRef.current });
      websocket.send(tokenMessage);
      setIsConnected(true);
      reconnectAttempts.current = 0;
      isConnectingRef.current = false;
    };

    websocket.onmessage = (event) => {
      log.debug("WebSocket message received:", event.data);
      try {
        const message = JSON.parse(event.data);
        setLastMessage(message);
      } catch (error) {
        log.error("Error parsing WebSocket message:", error);
      }
    };

    websocket.onclose = (event) => {
      log.info(
        `WebSocket disconnected. Code: ${event.code}, Reason: ${event.reason}`
      );
      setIsConnected(false);
      isConnectingRef.current = false;
      if (event.code === 1005) {
        reconnectAttempts.current += 1;
        if (reconnectAttempts.current < maxReconnectAttempts) {
          connect(tokenRef.current);
        } else {
          log.info("Max reconnect attempts reached. Stopping reconnection.");
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
          log.info("Max reconnect attempts reached. Stopping reconnection.");
        }
      }
    };

    websocket.onerror = (error) => {
      log.error("WebSocket error:", error);
      isConnectingRef.current = false;
    };
  };

  const disconnect = useCallback(() => {
    if (ws) {
      log.info("Closing WebSocket connection");
      ws.close();
    }
    if (reconnectTimeoutId.current) {
      clearTimeout(reconnectTimeoutId.current);
    }
    setWs(null);
    setIsConnected(false);
    tokenRef.current = null;
    isConnectingRef.current = false;
  }, [ws]);

  const sendMessage = useCallback(
    (message) => {
      if (ws) {
        ws.send(JSON.stringify(message));
      } else {
        console.error("WebSocket is not open. Attempting to reconnect.");
        connect(tokenRef.current);
      }
    },
    [ws]
  );

  return {
    isConnected,
    lastMessage,
    connect,
    sendMessage,
    setUserId,
    disconnect,
    userId,
  };
};

export default useWebSocket;

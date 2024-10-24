import React, { createContext, useContext, useState, useCallback } from "react";
import useWebSocket from "../hooks/useWebSocket";

// Create a context for the WebSocket connection
const WebSocketContext = createContext();

export const WebSocketProvider = ({ children }) => {
  const {
    isConnected,
    lastMessage,
    connect,
    disconnect,
    sendMessage,
    setUserId,
    userId,
  } = useWebSocket("ws://localhost:8080/ws");

  return (
    <WebSocketContext.Provider
      value={{
        isConnected,
        lastMessage,
        connect,
        disconnect,
        sendMessage,
        setUserId,
        userId,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
};

// Custom hook to use WebSocketContext
export const useWebSocketContext = () => {
  return useContext(WebSocketContext);
};

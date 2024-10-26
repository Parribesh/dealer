import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useWebSocketContext } from "../context/WebSocketContext";
import useGameState from "../hooks/useGameState";
import { useGameStateContext } from "../context/GameStateContext";
import { usePlayerContext } from "../context/PlayerContext";

const JoinLobby = () => {
  const [token, setToken] = useState(null);
  const [loading, setLoading] = useState(false);
  const [connectedPlayers, setConnectedPlayers] = useState([]);
  const [playerName, setPlayerName] = useState("");
  const { gameState, updateGameState } = useGameStateContext();
  const { player, setPlayer } = usePlayerContext();
  const {
    isConnected,
    lastMessage,
    connect,
    sendMessage,
    disconnect,
    setUserId,
    userId,
  } = useWebSocketContext();
  const navigate = useNavigate();
  const [lastReceivedMessage, setLastReceivedMessage] = useState("");

  // Function to join the lobby and retrieve the JWT token
  const joinLobby = async () => {
    setLoading(true);
    try {
      const response = await axios.post("http://localhost:8080/lobby/join");
      const jwtToken = response.data.token;
      localStorage.setItem("jwtToken", jwtToken);
      setToken(jwtToken);
      setPlayerName(response.data.playerName);
    } catch (error) {
      console.error("Error joining the lobby:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    console.log("userId set to", userId);
  }, [userId]);

  useEffect(() => {
    if (playerName && token) {
      connect(token);
      setUserId(playerName);
      setPlayer(playerName);
    }
  }, [playerName, token]);

  useEffect(() => {
    console.log("JoinLobby: WebSocket connected:", isConnected);
    console.log("JoinLobby: WebSocket last message:", lastMessage);
    if (lastMessage) {
      const { type, data } = lastMessage;
      console.log("JoinLobby: Message type:", type, "Message content:", data);

      switch (type.toLowerCase()) {
        case "gamestate":
          console.log("JoinLobby: Game started, navigating with state:", data);
          updateGameState(data);
          navigate("/game");
          break;
        case "playerlist":
          setConnectedPlayers(data);
          break;
        case "playerleft":
        case "playerjoined":
          setConnectedPlayers(data);
          break;
        case "roomcreated":
          // Handle room creation if needed
          console.log("Room created:", data);
          navigate("/game", { state: { initialGameState: data } });
          break;
        case "message":
          setLastReceivedMessage(data);
          break;
        default:
          console.log("Unhandled message type:", type);
      }
    }
  }, [lastMessage, navigate, isConnected]);

  useEffect(() => {
    return () => {
      if (!navigate) {
        disconnect();
      }
    };
  }, [disconnect]);

  const startGame = async () => {
    try {
      await axios.post("http://localhost:8080/game/start", null, {
        headers: { Authorization: `Bearer ${token}` },
      });
    } catch (error) {
      console.error("Error starting the game:", error);
    }
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        height: "100vh",
      }}
    >
      {!token ? (
        <button onClick={joinLobby} disabled={loading}>
          {loading ? (
            <>
              Loading... <span className="spinner">🔄</span>
            </>
          ) : (
            "Join Lobby"
          )}
        </button>
      ) : (
        <>
          <div>You are {playerName}</div>
          {lastReceivedMessage && (
            <div style={{ marginBottom: "10px", fontWeight: "bold" }}>
              Last message: {lastReceivedMessage}
            </div>
          )}
          <h2>Connected Players:</h2>
          {connectedPlayers && connectedPlayers.length > 0 ? (
            <ul>
              {connectedPlayers.map((player, index) => (
                <li key={index}>{player}</li>
              ))}
            </ul>
          ) : (
            <p>No players connected yet.</p>
          )}
          {connectedPlayers && connectedPlayers.length >= 4 && (
            <button onClick={startGame}>Start Game</button>
          )}
        </>
      )}
    </div>
  );
};

export default JoinLobby;

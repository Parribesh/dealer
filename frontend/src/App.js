import { Routes, Route } from "react-router-dom";
import "./App.css";
import JoinLobby from "./components/JoinLobby";
import Game from "./components/Game";
import { WebSocketProvider } from "./context/WebSocketContext";
import { GameStateProvider } from "./context/GameStateContext";

function App() {
  return (
    <div className="App">
      <main>
        <WebSocketProvider>
          <GameStateProvider>
            <Routes>
              <Route path="/" element={<JoinLobby />} />
              <Route path="/game" element={<Game />} />
            </Routes>
          </GameStateProvider>
        </WebSocketProvider>
      </main>
    </div>
  );
}

export default App;

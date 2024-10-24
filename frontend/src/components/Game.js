import React, { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";
import { useWebSocketContext } from "../context/WebSocketContext";
import styled from "styled-components";
import { useGameStateContext } from "../context/GameStateContext";

const GameContainer = styled.div`
  background: #1a1a2e;
  color: #ffffff;
  font-family: "Arial", sans-serif;
  padding: 20px;
  min-height: 100vh;
  position: relative;
`;

const GameHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
`;

const GameID = styled.h1`
  font-size: 24px;
  margin: 0;
`;

const TurnIndicator = styled.div`
  font-size: 18px;
  background: #16213e;
  padding: 10px 20px;
  border-radius: 20px;
`;

const PlayersContainer = styled.div`
  position: relative;
  width: 100%;
  height: calc(100vh);
`;

const PlayerSection = styled.div`
  background: #0f3460;
  border-radius: 10px;
  padding: 20px;
  position: absolute;
  width: 200px;
`;

const PlayerHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
`;

const PlayerName = styled.h2`
  margin: 0;
  font-size: 20px;
`;

const HealthBar = styled.div`
  background: #e94560;
  width: ${(props) => props.health}%;
  height: 10px;
  border-radius: 5px;
  transition: width 0.5s ease-in-out;
`;

const HandContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
`;

const Card = styled.div`
  background: ${(props) => (props.hidden ? "#000" : "#ffffff")};
  color: ${(props) =>
    props.hidden
      ? "#000"
      : props.suit === "H" || props.suit === "D"
      ? "#e94560"
      : "#16213e"};
  width: 40px;
  height: 60px;
  border-radius: 5px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold;
  font-size: 16px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
`;

const CentralArea = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 550px;
  height: 500px;
  background: #16213e;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 15px;
`;

const DeckCard = styled.div`
  background: #ffffff;
  width: 80px;
  height: 120px;
  border-radius: 5px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 24px;
  font-weight: bold;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
`;

const CentralBoard = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #16213e;
  width: 500px;
  height: 500px;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 0 15px rgba(0, 0, 0, 0.5);
`;

const Game = () => {
  const location = useLocation();
  // const [gameState, setGameState] = useState(
  //   location.state?.initialGameState || null
  // );
  const { gameState, updateGameState, updateHealthState } =
    useGameStateContext();
  const { lastMessage, isConnected, userId } = useWebSocketContext();

  useEffect(() => {
    if (!gameState) {
      // Request initial game state here
    }
  }, []);

  useEffect(() => {
    if (lastMessage) {
      const { type, data } = lastMessage;
      if (type.toLowerCase() === "gamestate") {
        updateGameState(data);
      }
      if (type.toLowerCase() === "healthstate") {
        updateHealthState(data);
      }
    }
  }, [lastMessage]);

  useEffect(() => {
    console.log("game state after health upate: ", gameState);
  }, gameState);

  const isCurrentPlayer = (playerId) => userId === playerId;

  if (!gameState) {
    return <div>Loading game state...</div>;
  }

  const { GameID, Players, State } = gameState;

  const getSuitSymbol = (suit) => {
    switch (suit) {
      case "H":
        return "♥";
      case "D":
        return "♦";
      case "C":
        return "♣";
      case "S":
        return "♠";
      default:
        return suit;
    }
  };

  return (
    <GameContainer>
      <GameHeader>
        <GameID>Game: {GameID}</GameID>
        {console.log("userId: ", userId)}
        <div>Your are: {userId}</div>
        <TurnIndicator>Turn: {State.turn}</TurnIndicator>
      </GameHeader>

      <PlayersContainer>
        {Players.map((playerId, index) => {
          const player = State[`player${index + 1}`];
          const playerPosition =
            index === 0
              ? { top: "5%", left: "5%" } // Top-left
              : index === 1
              ? { top: "5%", right: "5%" } // Top-right
              : index === 2
              ? { bottom: "5%", left: "5%" } // Bottom-left
              : { bottom: "5%", right: "5%" }; // Bottom-right

          return (
            <PlayerSection key={playerId} style={playerPosition}>
              <PlayerHeader>
                <PlayerName>{player.id}</PlayerName>
                <HealthBar health={player.health} />
              </PlayerHeader>
              <HandContainer>
                {player.hand.map((card, cardIndex) => (
                  <Card
                    key={cardIndex}
                    suit={card.suit}
                    hidden={!isCurrentPlayer(player.id)}
                  >
                    {!isCurrentPlayer(player.id)
                      ? "🂠"
                      : `${card.rank}${getSuitSymbol(card.suit)}`}
                  </Card>
                ))}
              </HandContainer>
            </PlayerSection>
          );
        })}

        <CentralArea>
          <DeckCard>🂠</DeckCard>
        </CentralArea>
      </PlayersContainer>
    </GameContainer>
  );
};

export default Game;

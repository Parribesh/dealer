import React, { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";
import { useWebSocketContext } from "../context/WebSocketContext";
import styled from "styled-components";
import { useGameStateContext } from "../context/GameStateContext";
import { usePlayerContext } from "../context/PlayerContext";
import { createComponentLogger } from "../logger";
import Scoreboard from "./ScoreBoard";

const log = createComponentLogger("GameContainer", "info");

const GameContainer = styled.div`
  background: #1a1a2e;
  color: #ffffff;
  font-family: "Arial", sans-serif;
  padding: 20px;
  min-height: calc(100vh - 50px);
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const GameHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  width: 100%;
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
  display: flex;
  justify-content: space-around;
  margin-bottom: 20px;
  gap: 40px;
`;

const PlayerSection = styled.div`
  background: #0f3460;
  border-radius: 10px;
  padding: 10px;
  width: 150px;
  text-align: center;
`;

const PlayerName = styled.h2`
  margin: 0;
  font-size: 18px;
`;

const HealthBar = styled.div`
  background: #e94560;
  width: ${(props) => props.health}%;
  height: 10px;
  border-radius: 5px;
  transition: width 0.5s ease-in-out;
  margin-top: 5px;
`;

const CentralArea = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: #16213e;
  border-radius: 15px;
  height: 500px;
  width: 500px;
  margin: 20px 0;
`;

const DeckGrid = styled.div`
  display: flex;

  align-items: center;
  gap: 10px;
  flex: 1;
`;

const CardContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  font-size: 20px;
`;

const Card = styled.div`
  background: ${(props) => (props.hidden ? "#000" : "#ffffff")};
  color: ${(props) =>
    props.hidden
      ? "#000"
      : props.suit === "H" || props.suit === "D"
      ? "#e94560"
      : "#16213e"};
  width: 80px;
  height: 120px;
  border-radius: 5px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold;
  font-size: 16px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
  cursor: pointer;
`;

const BottomRow = styled.div`
  display: flex;
  justify-content: center;
  padding: 10px;
  background: #16213e;
  border-radius: 15px;
  gap: 10px;
  position: absolute;
  bottom: 20px;
`;

const Game = () => {
  const { gameState, updateGameState, updateHealthState } =
    useGameStateContext();
  const { lastMessage, userId, sendMessage } = useWebSocketContext();
  const { player, setPlayer } = usePlayerContext();
  const [selectedCards, setSelectedCards] = useState({});
  const [trickWinner, setTrickWinner] = useState(null);

  useEffect(() => {
    if (lastMessage) {
      const { type, data } = lastMessage;
      if (type.toLowerCase() === "gamestate") {
        log.debug("updating game state", data);
        updateGameState(data);
        acknowledgeGameState(player);
      }
      if (type.toLowerCase() === "healthstate") {
        updateHealthState(data);
      }
      if (type.toLowerCase() === "trickwon") {
        // updateHealthState(data);
        log.info("trick won. Setting cards null...");
        setSelectedCards({});
      }
      if (type.toLowerCase() === "cardplayed") {
        // updateHealthState(data);
        log.debug("cardplayed: ", data.card);
        log.debug("player: ", data.playerId);

        setSelectedCards((prevSelectedCards) => ({
          ...prevSelectedCards,
          [data.playerId]: data.card,
        }));
      }
      if (type.toLowerCase() === "resetcardplayed") {
        setSelectedCards({});
      }
      if (type.toLowerCase() === "trickwon") {
        // updateHealthState(data);
        setSelectedCards(null);
        setTrickWinner(data.player);
      }
    }
  }, [lastMessage]);

  const acknowledgeGameState = (playerId) => {
    sendMessage({ type: "acknowledgment", playerId });
  };

  const handleCardClick = async (card, playerId) => {
    const moveData = { card };
    try {
      const response = await fetch(
        `http://localhost:8080/game/move?gameID=${gameState.GameID}&playerID=${playerId}`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(moveData),
        }
      );
      if (!response.ok) throw new Error("Network response was not ok");
    } catch (error) {
      log.error("Error sending move:", error);
    }
  };

  const { GameID, Players, State } = gameState || {};

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

  if (!gameState) {
    return <div>Loading game state...</div>;
  }

  return (
    <GameContainer>
      <GameHeader>
        <GameID>Game: {GameID}</GameID>
        <div>Player: {player.id}</div>
        <TurnIndicator>Turn: {State.turn}</TurnIndicator>
      </GameHeader>

      <PlayersContainer>
        {Players.map((playerId, index) => {
          const player = State[`player${index + 1}`];
          return (
            <PlayerSection key={playerId}>
              <PlayerName>{player.id}</PlayerName>
              <HealthBar health={player.health} />
              {player.health} <br></br>
              Score: {player.score}
            </PlayerSection>
          );
        })}
      </PlayersContainer>

      <CentralArea>
        {trickWinner && (
          <>
            <div>Winner {trickWinner.id}</div>
            <div>Score {trickWinner.score}</div>
          </>
        )}
        <DeckGrid>
          {Players.map((playerId) => {
            let selectedCard = "";
            if (selectedCards && selectedCards[playerId]) {
              selectedCard = selectedCards[playerId];
              log.debug("selectedCard: ", selectedCard);
              // Proceed with your logic
            } else {
              // console.warn(
              // selectedCards is either null or ${playerId} is not present.
              //);
            }
            return (
              <CardContainer key={playerId}>
                <div>{playerId}</div>
                {selectedCard ? (
                  <Card suit={selectedCard.suit} hidden={false}>
                    {`${selectedCard.rank}${getSuitSymbol(selectedCard.suit)}`}
                  </Card>
                ) : (
                  <Card hidden={true}>🂠</Card>
                )}
              </CardContainer>
            );
          })}
        </DeckGrid>
        <Scoreboard gameState={gameState} />
      </CentralArea>

      <BottomRow>
        {State[`player${Players.indexOf(userId) + 1}`]?.hand.map(
          (card, index) => (
            <Card
              key={index}
              suit={card.suit}
              onClick={() => handleCardClick(card, userId)}
            >
              {`${card.rank}${getSuitSymbol(card.suit)}`}
            </Card>
          )
        )}
      </BottomRow>
    </GameContainer>
  );
};

export default Game;

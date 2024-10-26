import React, { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";
import { useWebSocketContext } from "../context/WebSocketContext";
import styled from "styled-components";
import { useGameStateContext } from "../context/GameStateContext";
import { usePlayerContext } from "../context/PlayerContext";

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
  cursor: pointer; /* Change the cursor to indicate clickable cards */
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

// Styled components for the grid layout
const DeckGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr); // Adjust columns as needed
  gap: 10px; // Space between cards
`;

const CardContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100px; // Adjust height as needed
  width: 100px; // Adjust width as needed
`;
const PlayerId = styled.div`
  font-weight: bold; // Make the text bold
  text-align: center; // Center the text
  margin-bottom: 5px; // Space between player ID and card
`;

const Game = () => {
  const location = useLocation();
  const {
    gameState,
    updateGameState,
    updateHealthState,
    updatePlayedCardState,
  } = useGameStateContext();
  const { lastMessage, isConnected, userId, sendMessage } =
    useWebSocketContext();
  const { player, setPlayer } = usePlayerContext();
  const [selectedCards, setSelectedCards] = useState({}); // Step 1: State for selected card

  useEffect(() => {
    if (!gameState) {
      // Request initial game state here
    }
  }, []);

  const acknowledgeGameState = (playerId) => {
    const acknowledgment = {
      type: "acknowledgment",
      playerId: playerId,
    };
    console.log("isConnected: ", isConnected);
    console.log("Sending Ack: ");
    sendMessage(acknowledgment);
  };

  useEffect(() => {
    if (lastMessage) {
      const { type, data } = lastMessage;
      if (type.toLowerCase() === "gamestate") {
        updateGameState(data);
        acknowledgeGameState(player);
      }
      if (type.toLowerCase() === "healthstate") {
        updateHealthState(data);
      }
      if (type.toLowerCase() === "trickwon") {
        // updateHealthState(data);
        console.log("trick won. Setting cards null...");
        setSelectedCards({});
      }
      if (type.toLowerCase() === "cardplayed") {
        // updateHealthState(data);
        setSelectedCards((prevSelectedCards) => ({
          ...prevSelectedCards,
          [data.playerId]: data.card,
        }));
      }
      if (type.toLowerCase() === "trickwon") {
        // updateHealthState(data);
        setSelectedCards(null);
      }
    }
  }, [lastMessage]);

  useEffect(() => {
    console.log("game state after health update: ", gameState);
  }, [gameState]);

  const isCurrentPlayer = (playerId) => userId === playerId;

  const handleCardClick = async (card, playerId) => {
    console.log("card is being processed hold on... ");

    // Construct the move data (message body)
    const moveData = {
      card, // Include the card information (you can customize this based on your requirements)
    };

    try {
      const response = await fetch(
        `http://localhost:8080/game/move?gameID=${gameState.GameID}&playerID=${playerId}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json", // Specify the content type as JSON
          },
          body: JSON.stringify(moveData), // Convert the moveData object to a JSON string
        }
      );

      if (!response.ok) {
        throw new Error("Network response was not ok");
      }

      const responseData = await response.text();
      console.log("Response from server:", responseData);
    } catch (error) {
      console.error("Error sending move:", error);
    }
  };

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
                    onClick={() => handleCardClick(card, player.id)} // Add onClick event
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
          <DeckGrid>
            {Players.map((playerId) => {
              let selectedCard = "";
              if (selectedCards && selectedCards[playerId]) {
                selectedCard = selectedCards[playerId];
                // Proceed with your logic
              } else {
                console.warn(
                  `selectedCards is either null or ${playerId} is not present.`
                );
              }

              return (
                <CardContainer key={playerId}>
                  <PlayerId>{playerId}</PlayerId> {/* Display the player ID */}
                  {selectedCard ? (
                    <Card
                      suit={selectedCard.suit}
                      hidden={false} // Show the selected card
                    >
                      {`${selectedCard.rank}${getSuitSymbol(
                        selectedCard.suit
                      )}`}
                    </Card>
                  ) : (
                    <Card suit={null} hidden={true}>
                      {/* Placeholder for no card selected */}
                      🂠
                    </Card>
                  )}
                </CardContainer>
              );
            })}
          </DeckGrid>
        </CentralArea>
      </PlayersContainer>
    </GameContainer>
  );
};

export default Game;

import React, { useState } from "react";

import { createComponentLogger } from "../logger";

const log = createComponentLogger("useGameState", "info");

const useGameState = () => {
  const [gameState, setGameState] = useState(null);
  const rankOrder = {
    2: 2,
    3: 3,
    4: 4,
    5: 5,
    6: 6,
    7: 7,
    8: 8,
    9: 9,
    10: 10,
    J: 11,
    Q: 12,
    K: 13,
    A: 14,
  };

  const updateGameState = (data) => {
    Object.keys(data.State).map((key) => {
      if (key.includes("player")) {
        data.State[key].hand.sort((a, b) => {
          const suitComparison = a.suit.localeCompare(b.suit);
          if (suitComparison !== 0) {
            return suitComparison; // If suits are different, sort by suit
          }
          // Sort by rank (converting rank to numeric values)
          return rankOrder[a.rank] - rankOrder[b.rank];
        });
      }
    });

    setGameState(data);
  };

  const updateHealthState = (data) => {
    setGameState((prevState) => {
      // Find the key (player1, player2, etc.) where the player's id matches data.player
      const playerKey = Object.keys(prevState.State).find(
        (key) => prevState.State[key].id === data.player
      );

      log.debug("Updating health for player:", data.player);
      log.debug("Player Health:", data.health);
      log.debug("Matched player key:", playerKey);

      // If the player is found, update their health
      if (playerKey) {
        return {
          ...prevState,
          State: {
            ...prevState.State,
            [playerKey]: {
              ...prevState.State[playerKey],
              health: data.health,
            },
          },
        };
      }

      log.error("Player not found for update");
      return prevState;
    });
  };

  const updatePlayedCardState = (data) => {
    setGameState((prevState) => {
      // Find the key (player1, player2, etc.) where the player's id matches data.player
      const playerKey = Object.keys(prevState.State).find(
        (key) => prevState.State[key].id === data.player
      );

      log.debug("Updating played card for player:", data.player);
      log.debug("New played card:", data.played_card);
      log.debug("Matched player key:", playerKey);

      // If the player is found, update their played_card
      if (playerKey) {
        return {
          ...prevState,
          State: {
            ...prevState.State,
            [playerKey]: {
              ...prevState.State[playerKey],
              played_card: data.played_card, // Update the played_card value
            },
          },
        };
      }

      log.error("Player not found for update");
      return prevState;
    });
  };

  return {
    gameState,
    updateGameState,
    updateHealthState,
    updatePlayedCardState,
  };
};

export default useGameState;

//sample game state

/**
 * 
{
  "data": {
    "GameID": "game-59008",
    "Players": ["SwiftTiger", "SwiftEagle", "LuckyFalcon", "BraveTiger"],
    "State": {
      "player1": {
        "id": "SwiftTiger",
        "hand": [
          { "rank": "3", "suit": "S" },
          { "rank": "K", "suit": "H" },
          { "rank": "4", "suit": "D" },
          { "rank": "9", "suit": "H" },
          { "rank": "10", "suit": "D" },
          { "rank": "7", "suit": "S" },
          { "rank": "K", "suit": "S" },
          { "rank": "6", "suit": "C" },
          { "rank": "4", "suit": "C" },
          { "rank": "8", "suit": "C" },
          { "rank": "5", "suit": "D" },
          { "rank": "6", "suit": "D" },
          { "rank": "2", "suit": "C" }
        ],
        "health": 0
      },
      "player2": {
        "id": "SwiftEagle",
        "hand": [
          { "rank": "7", "suit": "C" },
          { "rank": "4", "suit": "H" },
          { "rank": "2", "suit": "D" },
          { "rank": "5", "suit": "C" },
          { "rank": "6", "suit": "S" },
          { "rank": "5", "suit": "S" },
          { "rank": "Q", "suit": "H" },
          { "rank": "2", "suit": "S" },
          { "rank": "8", "suit": "S" },
          { "rank": "3", "suit": "C" },
          { "rank": "J", "suit": "C" },
          { "rank": "9", "suit": "C" },
          { "rank": "2", "suit": "H" }
        ],
        "health": 0
      },
      "player3": {
        "id": "LuckyFalcon",
        "hand": [
          { "rank": "10", "suit": "C" },
          { "rank": "10", "suit": "S" },
          { "rank": "J", "suit": "H" },
          { "rank": "9", "suit": "S" },
          { "rank": "K", "suit": "D" },
          { "rank": "A", "suit": "D" },
          { "rank": "5", "suit": "H" },
          { "rank": "A", "suit": "C" },
          { "rank": "9", "suit": "D" },
          { "rank": "Q", "suit": "C" },
          { "rank": "A", "suit": "H" },
          { "rank": "4", "suit": "S" },
          { "rank": "10", "suit": "H" }
        ],
        "health": 0
      },
      "player4": {
        "id": "BraveTiger",
        "hand": [
          { "rank": "Q", "suit": "S" },
          { "rank": "7", "suit": "D" },
          { "rank": "J", "suit": "D" },
          { "rank": "6", "suit": "H" },
          { "rank": "3", "suit": "D" },
          { "rank": "K", "suit": "C" },
          { "rank": "J", "suit": "S" },
          { "rank": "8", "suit": "D" },
          { "rank": "A", "suit": "S" },
          { "rank": "Q", "suit": "D" },
          { "rank": "3", "suit": "H" },
          { "rank": "8", "suit": "H" },
          { "rank": "7", "suit": "H" }
        ],
        "health": 10
      },
      "turn": 4
    }
  },
  "type": "gamestate"
}

 */

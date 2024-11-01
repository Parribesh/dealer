import { createContext, useContext } from "react";
import useGameState from "../hooks/useGameState";

const GameStateContext = createContext();

export const GameStateProvider = ({ children }) => {
  const { gameState, updateGameState, updateHealthState, updatePlayedCardState, updateBidState } = useGameState();
  return (
    <GameStateContext.Provider
      value={{ gameState, updateGameState, updateHealthState, updatePlayedCardState, updateBidState }}
    >
      {children}
    </GameStateContext.Provider>
  );
};

export const useGameStateContext = () => {
  return useContext(GameStateContext);
};

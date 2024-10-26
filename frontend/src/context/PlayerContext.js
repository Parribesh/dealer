import { createContext, useContext } from "react";
import usePlayer from "../hooks/usePlayer";

const playerContext = createContext();
export const PlayerContextProvider = ({ children }) => {
  const { player, setPlayer } = usePlayer();
  return (
    <playerContext.Provider value={{ player, setPlayer }}>
      {children}
    </playerContext.Provider>
  );
};

export const usePlayerContext = () => {
  return useContext(playerContext);
};

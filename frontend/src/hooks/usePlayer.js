import { useState } from "react";

const usePlayer = () => {
  const [player, setPlayer] = useState(null);

  return { player, setPlayer };
};

export default usePlayer;

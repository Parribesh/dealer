import React, { useEffect } from "react";
import useGameState from "../hooks/useGameState";
import { createComponentLogger } from "../logger";

const log = createComponentLogger("Scoreboard", "info");
const Scoreboard = ({ gameState: state }) => {
  log.debug("state: ", state);
  // Create an array of players from the gamestate object
  //   const players = [""];

  const players = [
    state.State.player1,
    state.State.player2,
    state.State.player3,
    state.State.player4,
  ];

  // Inline CSS styles
  const styles = {
    container: {
      padding: "20px",
      width: "80%",
      maxWidth: "100%",
      backgroundColor: "#f4f4f9",
      borderRadius: "8px",
      boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
    },
    table: {
      width: "100%",
      borderCollapse: "collapse",
      marginTop: "10px",
    },
    header: {
      backgroundColor: "black",
      color: "white",
      textAlign: "left",
      padding: "4px",
    },
    cell: {
      padding: "2px",
      borderBottom: "1px solid #ddd",
      textAlign: "left",
      color: "black",
    },
    title: {
      fontSize: "1.2rem",
      fontWeight: "600",
      textAlign: "center",
      color: "#333",
      marginTop: "0",
    },
  };

  return (
    <div style={styles.container}>
      <p style={styles.title}>Scoreboard</p>
      <table style={styles.table}>
        <thead>
          <tr>
            <th style={styles.header}>Player Name</th>
            <th style={styles.header}>Bid</th>
            <th style={styles.header}>Score</th>
          </tr>
        </thead>
        <tbody>
          {players.map((player, index) => (
            <tr key={index}>
              <td style={styles.cell}>{player.id}</td>
              <td style={styles.cell}>{player.bid}</td>
              <td style={styles.cell}>
                {player.score != null
                  ? player.score 
                  : 0}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

// Example usage:
// Assuming `gameState` is passed in as props and contains the required fields:
// <Scoreboard state={gameState} />

export default Scoreboard;

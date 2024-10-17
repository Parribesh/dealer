# Multiplayer Turn-Based Card Game

This project is a multiplayer turn-based card game using a **Go backend** and a **JavaScript frontend**. The game will have a real-time interaction feature, with a WebSocket-based backend to handle player actions and state synchronization.

---

## Project Overview

### **Tech Stack**:
- **Backend**: Go
  - WebSocket for real-time game actions
  - REST API for initial game state
- **Frontend**: JavaScript (or React, optional)
  - WebSocket integration for real-time updates
  - HTML/CSS for UI and animations

---

## **1-Week Agile Sprint Plan**

### **Day 1: Project Setup & Architecture Planning (Sprint 1)**

#### **Tasks**:
1. **Set up version control**:
   - Initialize a Git repository.
   - Create a project board with tasks.
   
2. **Define game architecture**:
   - Plan how the backend will handle game sessions and player states.
   - Decide how WebSockets will be used for real-time communication.

3. **Set up basic Go server**:
   - Install dependencies (e.g., `github.com/gorilla/websocket`).
   - Create an HTTP server in Go with WebSocket support.

4. **Frontend setup**:
   - Create a new project (vanilla JS or React).
   - Set up a bundler like Webpack or Parcel.

#### **Deliverables**:
- GitHub repo with Go backend initialized.
- Initial Go server running with WebSocket capability.
- Basic frontend connected to the backend.

---

### **Day 2: Basic Game Logic & Card System (Sprint 2)**

#### **Tasks**:
1. **Implement card system in Go**:
   - Define card attributes (attack, defense, etc.).
   - Create functions for drawing and playing cards.
   
2. **Create deck and player system**:
   - Implement player actions and a basic deck.

3. **Mock API**:
   - Test interactions with basic REST API routes.

#### **Frontend**:
- Display static cards and connect to backend using **fetch** requests for testing.

#### **Deliverables**:
- Go card system and player logic.
- Simple REST API for game state interaction.
- Frontend displays cards with backend interaction.

---

### **Day 3: WebSocket Integration & Turn System (Sprint 3)**

#### **Tasks**:
1. **WebSocket integration in Go**:
   - Handle player actions in real-time.
   - Broadcast game state updates to players.

2. **Synchronize game state**:
   - Ensure player actions are reflected on both ends.

#### **Frontend**:
- Use **WebSocket API** in JS for real-time updates.
- Update card display based on WebSocket messages.

#### **Deliverables**:
- WebSocket communication between players.
- Real-time turn-based system.

---

### **Day 4: Frontend Game UI (Sprint 4)**

#### **Tasks**:
1. **Design card UI**:
   - Implement a dynamic card grid.
   - Add clickable cards to play them.

2. **Game board setup**:
   - Show player health, hands, and actions (e.g., play card, end turn).

#### **Frontend**:
- Cards are clickable, and player actions update in real-time.

#### **Deliverables**:
- Functional game board with dynamic card UI.
- Real-time updates in the UI.

---

### **Day 5: Game State Management & Logic Enhancements (Sprint 5)**

#### **Tasks**:
1. **Enhance game logic**:
   - Add more complex card interactions (attack, defense, spells).
   
2. **Implement win/loss conditions**:
   - Broadcast game-over events via WebSocket.

3. **Validation and error handling**:
   - Ensure players can’t play out-of-turn cards.

#### **Frontend**:
- Display win/loss notifications and turn feedback.

#### **Deliverables**:
- Full turn-based game logic with win conditions.
- Real-time UI reflecting gameplay.

---

### **Day 6: User Authentication & Matchmaking (Sprint 6)**

#### **Tasks**:
1. **Basic user login**:
   - Implement JWT authentication for players.

2. **Matchmaking**:
   - Create a lobby and pair players for matches.

#### **Frontend**:
- Login page and game lobby for player matching.

#### **Deliverables**:
- Basic authentication and matchmaking system.
- Frontend login and matchmaking flow.

---

### **Day 7: Testing, Deployment, and Polish (Sprint 7)**

#### **Tasks**:
1. **Test and fix bugs**:
   - Test edge cases and smooth gameplay.

2. **Code clean-up**:
   - Refactor and document code for maintainability.

3. **Deploy the project**:
   - Deploy the Go backend (Heroku/AWS).
   - Deploy the frontend (GitHub Pages/Netlify).

#### **Frontend**:
- Ensure responsiveness and add animations for a polished user experience.

#### **Deliverables**:
- Fully tested and deployed project.
- Clean, documented code.

---

## How to Run Locally

### **Backend**:
1. Clone the repository.
2. Install Go dependencies:
   ```bash
   go mod tidy

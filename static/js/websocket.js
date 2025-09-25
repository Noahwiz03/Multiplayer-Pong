import * as handlers from "./handlers.js"; 

let ws;
export function initWebSocket(socket){
ws = socket;
ws.onopen =() =>{
  console.log("connected succesfully to server");
  ws.send("Hi from client");
}

ws.onmessage = e => {
  console.log("Recieved:", e.data);
  let msg = JSON.parse(e.data); 
  console.log(msg.roomCreated);
  switch (msg.type){
    case "roomCreated":
      handlers.handleRoomCreated(msg);
    break;
    case  "joinedRoom":
      handlers.handleRoomJoined(msg);
      break; 
    case "joinedTeam":
      handlers.handleJoinedTeam(msg);
      break;
    case "roomLeft":
      handlers.handleRoomLeft(msg);
      break;
    case "hostReassigned":
      handlers.handleHostReassigned(msg);
      break;
    case "gameState":
      handlers.handleGameState(msg);
      break;
    case "returnToLobby":
      handlers.handleReturnToLobby();
      break;
  }
}

ws.onclose =(e) =>{
  console.log("socket closed connection:", e);
}

ws.onerror =(error) =>{
  console.log("Socket error:", error);
}

}

export function sendMessage(type, payload ={}){
  ws.send(JSON.stringify({type, ...payload}));
}
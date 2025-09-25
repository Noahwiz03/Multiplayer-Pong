import { renderGameState } from "./gameState.js";
import * as uiState from "./uiState.js";

export function handleRoomCreated(msg){
  console.log("we made it here...", msg)
  uiState.setRoomCode(msg.roomCode);
  uiState.setHost(msg.host);
  if(msg.roomCreated){
    document.getElementById("roomcodeDisplay").textContent = "RoomCode: " + uiState.getRoomCode(); 
    if(uiState.isHost()){
      document.getElementById("hostDetector").textContent = "You are the Host";
      document.getElementById("StartGameBtn").style.display = 'block';
      document.getElementById("ReturnGameToLobbyBtn").style.display = 'block';
    }

    showPage("gameLobby"); 
  }
  else{
    return;
  } 
}

export function handleRoomJoined(msg){
 if(msg.joined){
    showPage("gameLobby"); 
  }
  if(!uiState.isHost()){
      document.getElementById("roomcodeDisplay").textContent = "RoomCode: ";
      document.getElementById("hostDetector").textContent = "";
      document.getElementById("StartGameBtn").style.display = 'none';
      document.getElementById("ReturnGameToLobbyBtn").style.display = 'none';
    }
  else{
    return;
  }
}

export function handleJoinedTeam(msg){
  showPage("Game"); 
}

export function handleRoomLeft(msg){
  uiState.setHost(false);
  uiState.setRoomCode('');
  showPage("landing"); 
}

export function handleHostReassigned(msg){
  uiState.setRoomCode(msg.roomCode);
  uiState.setHost(msg.host);
  document.getElementById("hostDetector").textContent = "You are the Host";
  document.getElementById("roomcodeDisplay").textContent = "RoomCode: " + uiState.getRoomCode();
  document.getElementById("StartGameBtn").style.display = 'block';
  document.getElementById("ReturnGameToLobbyBtn").style.display = 'block';
}

export function handleGameState(msg){
  renderGameState(msg);
}

export function handleReturnToLobby(){
  if(uiState.isHost()){
    document.getElementById("roomcodeDisplay").textContent = "RoomCode: " + uiState.getRoomCode(); 
    document.getElementById("hostDetector").textContent = "You are the Host";
  }

  showPage("gameLobby");
}

function showPage(pageId){
  console.log("showing the page now")
  document.querySelectorAll(".page").forEach(p => p.classList.remove("active"));
  document.getElementById(pageId).classList.add("active"); 
}
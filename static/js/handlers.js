export function handleRoomCreated(msg){
  console.log("we made it here...", msg)
  if(msg.roomCreated){
    document.getElementById("roomcodeDisplay").textContent = "RoomCode: " + msg.roomCode;
    if(msg.host){
      document.getElementById("hostDetector").textContent = "You are the Host";
      document.getElementById("StartGameBtn").style.display = 'block';
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
  if(!msg.host){
      document.getElementById("roomcodeDisplay").textContent = "RoomCode: ";
      document.getElementById("hostDetector").textContent = "";
      document.getElementById("StartGameBtn").style.display = 'none';
    }
  else{
    return;
  }
}

export function handleJoinedTeam(msg){
  showPage("Game"); 
}

export function handleRoomLeft(msg){
  showPage("landing"); 
}

export function handleHostReassigned(msg){
  if(msg.host){
    document.getElementById("hostDetector").textContent = "You are the Host";
    document.getElementById("roomcodeDisplay").textContent = "RoomCode: " + msg.roomCode;
    document.getElementById("StartGameBtn").style.display = 'block';
  }
}

export function handleGameState(msg){
 //send info needed to render game to canvas 
}

function showPage(pageId){
  console.log("showing the page now")
  document.querySelectorAll(".page").forEach(p => p.classList.remove("active"));
  document.getElementById(pageId).classList.add("active"); }
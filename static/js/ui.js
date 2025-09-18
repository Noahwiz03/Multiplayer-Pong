
export function setupUI(sendMessage){
  console.log("setting up dom");
  document.getElementById("CreateRoomBtn").onclick = () =>{
    sendMessage("createRoom");
  }

  document.getElementById("JoinRoomBtn").onclick = () =>{
  const  roomcode = document.getElementById("roomCode").value.trim();
    sendMessage("joinRoom", {roomCode: roomcode});
  }

  document.getElementById("joinLeftBtn").onclick = () => {
    sendMessage("joinTeam", {team: "left"});
  }

  document.getElementById("joinRightBtn").onclick = () => {
    sendMessage("joinTeam", {team: "right"});
  }

  document.getElementById("goBackBtn").onclick = () => {
    sendMessage("leaveRoom");
  }

  document.getElementById("leaveGameBtn").onclick = () => {
    sendMessage("leaveRoom");
  }
  
  document.getElementById("StartGameBtn").style.display = 'none';
  document.getElementById("StartGameBtn").onclick = () => {
    sendMessage("gameStart");
  }
}
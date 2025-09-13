window.onload = () =>{
let ws = new WebSocket("http://localhost:8080/ws");

ws.onopen =() =>{
  console.log("connected succesfully to server");
  ws.send("Hi from client");
}

ws.onmessage = e => console.log("message:", e.data);
ws.onclose =(e) =>{
  console.log("socket closed connection:", e);
}
ws.onerror =(error) =>{
  console.log("Socket error:", error);
}

document.getElementById("CreateRoomBtn").onclick = () =>{
  ws.send(JSON.stringify({type: "createRoom"}));
  showPage("gameLobby"); 
}

document.getElementById("JoinRoomBtn").onclick = () =>{
const  roomcode = document.getElementById("roomCode").value.trim();
  ws.send(JSON.stringify({type: "joinRoom", roomCode: roomcode}));
  console.log("sending:",JSON.stringify({type: "joinRoom", roomCode: roomcode}));
  showPage("gameLobby"); 
}

document.getElementById("joinLeftBtn").onclick = () => {
  ws.send(JSON.stringify({type: "joinTeam", team: "left"}));
  showPage("Game"); 
}

document.getElementById("joinRightBtn").onclick = () => {
  ws.send(JSON.stringify({type: "joinTeam", team: "right"}));
  showPage("Game"); 
}

document.getElementById("goBackBtn").onclick = () => {
  ws.send(JSON.stringify({type: "leaveRoom"}));
  showPage("landing");
}

document.getElementById("leaveGameBtn").onclick = () => {
  ws.send(JSON.stringify({type: "leaveRoom"}));
  showPage("landing");
}

function showPage(pageId){
  document.querySelectorAll(".page").forEach(p => p.classList.remove("active"));
  document.getElementById(pageId).classList.add("active");
}
}



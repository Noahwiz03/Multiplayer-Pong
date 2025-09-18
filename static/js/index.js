import { initWebSocket, sendMessage } from "./websocket.js";
import { setupUI } from "./ui.js";

window.onload = () =>{
 let socket = new WebSocket("http://localhost:8080/ws");
  initWebSocket(socket);
  setupUI(sendMessage);
}



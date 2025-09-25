let Host = false;
let RoomCode = '';

export function setRoomCode(code){
  RoomCode = code;
}
export function getRoomCode(){
  return RoomCode;
}

export function setHost(host){
  Host = host;
}
export function isHost(){
  return Host;
}
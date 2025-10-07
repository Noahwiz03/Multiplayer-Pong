const canvas = document.getElementById("PongCanvas");
const ctx = canvas.getContext("2d");

const Paddle = {
	X: 0,
  Y: 0,
	Width: 0,
	Height: 0,
	Speed: 0,
}

const Ball = {
	X:0,
  Y:0,
	Radius:0,
	VelocityX:0,
	VelocityY:0,
}

const GameState = {
	LeftPaddle:   {...Paddle},
	RightPaddle: {...Paddle},
	Ball:    Ball,
	ScoreLeft: 0,
	ScoreRight: 0,
}


ctx.fillStyle = "white";
ctx.font = '30px Arial';

//render gamestate into canvas!
export function renderGameState(msg){

	let state = msg.gameState;
	Object.assign(GameState.LeftPaddle, state.LeftPaddle);
	Object.assign(GameState.RightPaddle, state.RightPaddle);
	Object.assign(GameState.Ball, state.Ball);
	GameState.ScoreLeft = state.ScoreLeft;
	GameState.ScoreRight = state.ScoreRight;

	ctx.clearRect(0,0, canvas.width, canvas.height);

	ctx.beginPath();
	ctx.setLineDash([10, 15]); 
	ctx.moveTo(canvas.width / 2, 0);
	ctx.lineTo(canvas.width / 2, canvas.height);
	ctx.strokeStyle = "#ffffff57"; 
	ctx.lineWidth = 4;
	ctx.stroke();
	ctx.setLineDash([]); 

	ctx.fillText(GameState.ScoreLeft.toString(), 250, 50);
	ctx.fillText(GameState.ScoreRight.toString(), 550, 50);

	ctx.fillRect(GameState.LeftPaddle.X, GameState.LeftPaddle.Y, GameState.LeftPaddle.Width, GameState.LeftPaddle.Height);
	ctx.fillRect(GameState.RightPaddle.X, GameState.RightPaddle.Y, GameState.RightPaddle.Width, GameState.RightPaddle.Height);
	ctx.beginPath();
	ctx.arc(GameState.Ball.X, GameState.Ball.Y, GameState.Ball.Radius,0, Math.PI *2, true);
	ctx.fill();
}
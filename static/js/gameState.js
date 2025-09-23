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


//render gamestate into canvas!
export function renderGameState(msg){

	let state = msg.gameState;
	Object.assign(GameState.LeftPaddle, state.LeftPaddle);
	Object.assign(GameState.RightPaddle, state.RightPaddle);
	Object.assign(GameState.Ball, state.Ball);
	Object.assign(GameState.ScoreLeft , state.ScoreLeft);
	Object.assign(GameState.ScoreRight , state.ScoreRight);

	ctx.fillStyle = "white";
	ctx.clearRect(0,0, canvas.width, canvas.height);

	ctx.fillRect(GameState.LeftPaddle.X, GameState.LeftPaddle.Y, GameState.LeftPaddle.Width, GameState.LeftPaddle.Height);
	ctx.fillRect(GameState.RightPaddle.X, GameState.RightPaddle.Y, GameState.RightPaddle.Width, GameState.RightPaddle.Height);
	ctx.beginPath();
	ctx.arc(GameState.Ball.X, GameState.Ball.Y, GameState.Ball.Radius,0, Math.PI *2, true);
	ctx.fill();
}
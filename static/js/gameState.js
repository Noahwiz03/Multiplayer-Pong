const canvas = document.getElementById("PongCanvas");
const ctx = canvas.getContext("2d");

const GameState = {
	LeftPaddle:  new Paddle,
	RightPaddle: new Paddle,
	Ball:  new Ball,
	ScoreLeft,
	ScoreRight,
}

const Paddle = {
	X,
  Y,
	Width,
	Height,
	Speed,
}

const Ball = {
	X,
  Y,
	Radius,
	VelocityX,
	VelocityY,
}

//render gamestate into canvas!
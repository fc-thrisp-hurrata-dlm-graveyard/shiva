package engine

import (
	"time"
)

var (
	frameDelta int64
	frameTime  int64
	frameCount int64
	FPS        int64
)

func fps(e *Engine) {
	now := time.Now().Unix()
	frameDelta = now - frameTime
	frameTime = now
	frameCount = frameCount + 1
	if frameDelta >= 1 {
		FPS = frameCount
		e.FPS = FPS
		frameCount = 0
	}
}

package domain

import "3205.team/go-mp3-converter/application/progress"

type Progress interface {
	Send(stage progress.Stage, prog int)
	SendError(stage progress.Stage)
	SendDone()
}

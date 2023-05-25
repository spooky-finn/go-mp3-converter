package progress

type (
	Stage  string
	Status string
)

const (
	StatusOk    Status = "ok"
	StatusError Status = "error"

	DownloadStage Stage = "downloading"
	ConvertStage  Stage = "converting"
)

type ProgressEventPayload struct {
	Stage  Stage  `json:"stage"`
	Prog   int    `json:"progress"`
	Status Status `json:"status"`
}

type Progress struct {
	Ch   chan ProgressEventPayload `json:",omitempty"`
	Done chan struct{}
}

func New() *Progress {
	return &Progress{
		Ch:   make(chan ProgressEventPayload),
		Done: make(chan struct{}),
	}
}

func (p *Progress) Send(stage Stage, prog int) {
	p.Ch <- ProgressEventPayload{
		Stage:  stage,
		Prog:   prog,
		Status: StatusOk,
	}
}

func (p *Progress) SendError(stage Stage) {
	p.Ch <- ProgressEventPayload{
		Stage:  stage,
		Prog:   0,
		Status: StatusError,
	}
}

func (p *Progress) SendDone() {
	p.Done <- struct{}{}
}

func (p *Progress) Close() {
	close(p.Ch)
	close(p.Done)
}

func (p *Progress) Listen() <-chan ProgressEventPayload {
	return p.Ch
}

func (p *Progress) Wait() <-chan struct{} {
	return p.Done
}

func (p *Progress) ToPercentage(progEvent ProgressEventPayload) int {
	// i have 2 stages is now stge percentage is less than 50
	if progEvent.Stage == DownloadStage {
		return progEvent.Prog / 2
	} else {
		return progEvent.Prog/2 + 50
	}
}

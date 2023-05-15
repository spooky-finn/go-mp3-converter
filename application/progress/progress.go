package progress

type Stage string

const DownloadStage Stage = "downloading"
const ConvertStage Stage = "converting"

type Status string

const (
	StatusOk    Status = "ok"
	StatusError Status = "error"
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

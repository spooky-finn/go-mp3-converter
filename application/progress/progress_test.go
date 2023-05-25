package progress

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPercentage(t *testing.T) {
	p := New()
	defer p.Close()

	go func() {
		p.Send(DownloadStage, 10)
	}()
	progEvent := <-p.Ch
	assert.Equal(t, 5, progEvent.AggragetedProgress())

	go func() {
		p.Send(DownloadStage, 100)
	}()
	progEvent = <-p.Ch
	assert.Equal(t, 50, progEvent.AggragetedProgress())

	go func() {
		p.Send(ConvertStage, 50)
	}()
	progEvent = <-p.Ch
	assert.Equal(t, 75, progEvent.AggragetedProgress())

	go func() {
		p.Send(ConvertStage, 100)
	}()
	progEvent = <-p.Ch
	assert.Equal(t, 100, progEvent.AggragetedProgress())

}

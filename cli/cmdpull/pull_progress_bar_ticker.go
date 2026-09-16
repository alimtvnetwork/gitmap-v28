package cmdpull

import (
	"fmt"
	"time"
)

func (p *PullProgressBar) initTickerLifecycle() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.startTime = time.Now()
	p.ticker = time.NewTicker(80 * time.Millisecond)
	p.renderLocked()
}

func (p *PullProgressBar) runTickerLoop() {
	defer p.tickerWg.Done()
	for {
		select {
		case <-p.stopChan:
			return
		case <-p.ticker.C:
			p.tick()
		}
	}
}

func (p *PullProgressBar) tick() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isStopped {
		return
	}
	p.spinnerIdx++
	if p.isTTY {
		p.renderLocked()
	}
}

func (p *PullProgressBar) shutdownTicker() {
	p.mu.Lock()
	p.isStopped = true
	if p.ticker != nil {
		p.ticker.Stop()
	}
	if p.stopChan != nil {
		close(p.stopChan)
	}
	p.mu.Unlock()
	p.tickerWg.Wait()
}

func (p *PullProgressBar) flushFinalOutput() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isTTY && !p.isQuiet {
		fmt.Fprintln(p.out)
	}
}

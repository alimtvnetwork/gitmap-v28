package cmdcdp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const domProbeScript = `(() => {
  const ladder = [
    'textarea[placeholder*="Ask"]',
    'textarea[placeholder*="Message"]',
    'textarea[placeholder*="Type"]',
    'textarea',
    'div[contenteditable="true"]',
    '[role="textbox"]',
    'input[type="text"]',
    'input:not([type])',
    '#prompt-textarea',
    '[data-placeholder]'
  ];
  const found = [];
  for (let i = 0; i < ladder.length; i++) {
    try {
      const el = document.querySelector(ladder[i]);
      if (el) {
        found.push({
          selector: ladder[i],
          tag: el.tagName.toLowerCase(),
          ariaLabel: el.getAttribute('aria-label') || '',
          placeholder: el.getAttribute('placeholder') || el.getAttribute('data-placeholder') || '',
          text: (el.value || el.innerText || '').slice(0, 100)
        });
      }
    } catch (e) {}
  }
  let active = '';
  if (document.activeElement) {
    const act = document.activeElement;
    active = act.tagName.toLowerCase() + (act.id ? '#' + act.id : '') + (act.className ? '.' + String(act.className).split(' ')[0] : '');
  }
  const editables = [];
  try {
    const allCE = document.querySelectorAll('[contenteditable="true"]');
    for (let j = 0; j < allCE.length; j++) {
      editables.push({
        tag: allCE[j].tagName.toLowerCase(),
        aria: allCE[j].getAttribute('aria-label') || '',
        className: allCE[j].className || ''
      });
    }
  } catch (e) {}
  return {
    found: found,
    active: active,
    allContentEditables: editables
  };
})()`

func buildEvalRequest(id int, expression string) ([]byte, error) {
	req := CDPRequest{
		ID:     id,
		Method: "Runtime.evaluate",
		Params: map[string]any{
			"expression":    expression,
			"returnByValue": true,
			"awaitPromise":  true,
		},
	}

	return json.Marshal(req)
}

func parseFrameResponse(frame []byte, targetID int) (*CDPResponse, bool) {
	var resp CDPResponse
	if err := json.Unmarshal(frame, &resp); err != nil {
		return nil, false
	}

	return &resp, resp.ID == targetID
}

func readEvalResponse(client *WSClient, targetID int) (*CDPResponse, error) {
	for {
		frame, err := client.RecvFrame()
		if err != nil {
			return nil, apperror.WrapSimple(err, "cdp.recv_eval")
		}

		resp, isMatch := parseFrameResponse(frame, targetID)
		if isMatch {
			return resp, nil
		}
	}
}

func unpackEvalResult(result json.RawMessage) json.RawMessage {
	var evalRes runtimeEvalResult
	if err := json.Unmarshal(result, &evalRes); err == nil && len(evalRes.Result.Value) > 0 {
		return evalRes.Result.Value
	}

	return result
}

func extractEvalValue(resp *CDPResponse) (json.RawMessage, error) {
	if resp.Error != nil {
		return nil, fmt.Errorf("cdp error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	return unpackEvalResult(resp.Result), nil
}

func executeEval(client *WSClient, reqBytes []byte) (json.RawMessage, error) {
	if err := client.SendFrame(reqBytes); err != nil {
		return nil, apperror.WrapSimple(err, "cdp.send_eval")
	}

	resp, err := readEvalResponse(client, 1)
	if err != nil {
		return nil, err
	}

	return extractEvalValue(resp)
}

// EvaluateJS sends a Runtime.evaluate command to the CDP target and returns the evaluated JSON result.
func EvaluateJS(wsURL string, expression string, timeout time.Duration) (json.RawMessage, error) {
	client, err := DialWebSocket(wsURL, timeout)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	reqBytes, err := buildEvalRequest(1, expression)
	if err != nil {
		return nil, err
	}

	return executeEval(client, reqBytes)
}

// RunDOMInputProbe executes the DOM selector ladder probe on the target and returns parsed findings.
func RunDOMInputProbe(wsURL string) (*ProbeResult, error) {
	raw, err := EvaluateJS(wsURL, domProbeScript, 5*time.Second)
	if err != nil {
		return nil, err
	}

	var res ProbeResult
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, apperror.WrapSimple(err, "cdp.parse_probe_result")
	}

	return &res, nil
}

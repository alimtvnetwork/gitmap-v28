package cmdagy

import (
	"encoding/json"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// AgentAPISendMessage sends a message into an existing conversation.
func AgentAPISendMessage(convID, title, content string) result.Result[string] {
	args := buildSendMessageArgs(convID, title, content)
	rawRes := executeAgentAPICmd(args)
	if rawRes.IsFailure() {
		return result.Fail[string](rawRes.Error)
	}

	return parseSendMessageResponse(rawRes.Value)
}

func buildSendMessageArgs(convID, title, content string) []string {
	hasTitle := len(strings.TrimSpace(title)) > 0
	if hasTitle {
		return []string{"send-message", "--title=" + title, convID, content}
	}

	return []string{"send-message", convID, content}
}

func parseSendMessageResponse(data []byte) result.Result[string] {
	var resp rawAgentAPIResponse
	err := json.Unmarshal(data, &resp)
	hasErr := err != nil
	if hasErr {
		return result.Fail[string](apperror.WrapSimple(err, "parse send-message response"))
	}
	hasAPIError := len(resp.Error) > 0
	if hasAPIError {
		return result.Fail[string](apperror.NewSimple("agentapi error: "+resp.Error, "E9032"))
	}
	hasMsg := resp.Response.SendMessage != nil
	if hasMsg {
		return result.Ok(resp.Response.SendMessage.RecipientID)
	}

	return result.Ok("")
}

// AgentAPINewConversation creates a new conversation with an initial prompt.
func AgentAPINewConversation(title, prompt string) result.Result[string] {
	args := buildNewConvArgs(title, prompt)
	rawRes := executeAgentAPICmd(args)
	if rawRes.IsFailure() {
		return result.Fail[string](rawRes.Error)
	}

	return parseNewConversationResponse(rawRes.Value)
}

func buildNewConvArgs(title, prompt string) []string {
	hasTitle := len(strings.TrimSpace(title)) > 0
	if hasTitle {
		return []string{"new-conversation", "--title=" + title, prompt}
	}

	return []string{"new-conversation", prompt}
}

func parseNewConversationResponse(data []byte) result.Result[string] {
	var resp rawAgentAPIResponse
	err := json.Unmarshal(data, &resp)
	hasErr := err != nil
	if hasErr {
		return result.Fail[string](apperror.WrapSimple(err, "parse new-conversation response"))
	}
	hasAPIError := len(resp.Error) > 0
	if hasAPIError {
		return result.Fail[string](apperror.NewSimple("agentapi error: "+resp.Error, "E9033"))
	}
	hasConv := resp.Response.NewConversation != nil
	if hasConv {
		return result.Ok(resp.Response.NewConversation.ConversationID)
	}

	return result.Fail[string](apperror.NewSimple("empty new conversation response", "E9034"))
}

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
		return result.Fail[string](rawRes.Err)
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
	return AgentAPINewConversationWithOptions(title, "", "", prompt)
}

// AgentAPINewConversationWithOptions creates a new conversation with model, profile, title, and initial prompt.
func AgentAPINewConversationWithOptions(title, model, profile, prompt string) result.Result[string] {
	args := buildNewConvArgsWithOptions(title, model, profile, prompt)
	rawRes := executeAgentAPICmd(args)
	if rawRes.IsFailure() {
		return result.Fail[string](rawRes.Err)
	}

	return parseNewConversationResponse(rawRes.Value)
}

func buildNewConvArgsWithOptions(title, model, profile, prompt string) []string {
	args := []string{"new-conversation"}
	if len(strings.TrimSpace(model)) > 0 {
		args = append(args, "--model="+strings.TrimSpace(model))
	}
	if len(strings.TrimSpace(title)) > 0 {
		args = append(args, "--title="+strings.TrimSpace(title))
	}
	if len(strings.TrimSpace(profile)) > 0 {
		args = append(args, "--profile="+strings.TrimSpace(profile))
	}
	args = append(args, prompt)

	return args
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

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

// Message representa un mensaje con autor y partes
type Message struct {
	Author string `json:"author"`
	Parts  []Part `json:"parts"`
}

//Parsea un json como el siguiente:
/*
{
	"author": "user",
	"parts": [
		"Hola, ¿cómo estás?",
		"¿Puedes ayudarme con algo?"
	]
}
*/

// Part representa cada parte del mensaje
type Part struct {
	Text       string      `json:"text,omitempty"`
	Transcript string      `json:"transcript,omitempty"`
	Asset      interface{} `json:"asset,omitempty"`
}

//Nota: struct field tag `json:"text, omitempty"` not compatible with reflect.StructTag.Get: suspicious space in struct tag value
// Key takeaway: Struct tag syntax is very literal—spaces matter! The format should always be key:"value,option1,option2" with no spaces around commas or colons.

// Conversations representa una conversación con ID y mensajes
type ConversationNode struct {
	Message *MessageData `json:"message,omitempty"`
	Parent  string       `json:"parent,omitempty"`
}

// MessageData contiene los datos del mensaje
type MessageData struct {
	Author   AuthorData      `json:"author"`
	Content  ContentData     `json:"content"`
	Metadata MessageMetadata `json:"metadata"`
}

// AuthorData representa el autor del mensaje
type AuthorData struct {
	Role string `json:"role"`
}

// ContentData representa el contenido del mensaje
type ContentData struct {
	ContentType string     `json:"content_type"`
	Parts       []PartData `json:"parts"`
}

// PartData representa una parte del contenido
type PartData struct {
	ContentType                string          `json:"content_type"`
	Text                       string          `json:"text"`
	AudioAssetPointer          *AssetPointer   `json:"audio_asset_pointer,omitempty"`
	VideoContainerAssetPointer *AssetPointer   `json:"video_container_asset_pointer,omitempty"`
	FrameAssetPointers         []*AssetPointer `json:"frame_asset_pointers,omitempty"`
}

// UnmarshalJSON allows PartData to accept either a plain string or a structured object in the JSON array.
func (p *PartData) UnmarshalJSON(data []byte) error {
	// If the part is a JSON string (legacy export), treat it as text-only.
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		p.Text = s
		return nil
	}

	// Otherwise, decode as the structured payload.
	type partAlias PartData
	var aux partAlias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*p = PartData(aux)
	return nil
}

// AssetPointer representa un puntero a asset
type AssetPointer struct {
	// Campos dinámicos según el asset
}

// MessageMetadata contiene metadata del mensaje
type MessageMetadata struct {
	IsUserSystemMessage bool `json:"is_user_system_message"`
}

// Conversation representa la estructura completa
type Conversation struct {
	Title       string                      `json:"title"`
	CurrentNode string                      `json:"current_node"`
	Mapping     map[string]ConversationNode `json:"mapping"`
}

// ProcessParts procesa las partes del contenido
func processParts(parts []PartData) []Part {
	var result []Part

	for _, part := range parts {
		if part.Text != "" {
			result = append(result, Part{Text: part.Text})
			continue
		}

		switch part.ContentType {
		case "audio_transcription":
			result = append(result, Part{Transcript: part.Text})
		case "audio_asset_pointer", "image_asset_pointer", "video_container_asset_pointer":
			result = append(result, Part{Asset: part})
		case "real_time_user_audio_video_asset_pointer":
			if part.AudioAssetPointer != nil {
				result = append(result, Part{Asset: part.AudioAssetPointer})
			}
			if part.VideoContainerAssetPointer != nil {
				result = append(result, Part{Asset: part.VideoContainerAssetPointer})
			}
			for _, frame := range part.FrameAssetPointers {
				if frame != nil {
					result = append(result, Part{Asset: frame})
				}
			}

		}
	}

	return result
}

/* -------------------------------------------------------------------------- */
/*                             CARGADO DE ARCHIVO                             */
/* -------------------------------------------------------------------------- */

func loadConversations(filename string) ([]Conversation, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var conversations []Conversation
	if err := json.Unmarshal(data, &conversations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return conversations, nil
}

/* -------------------------------------------------------------------------- */
/*                                  IMPORTER                                  */
/* -------------------------------------------------------------------------- */

// GetConversationMessages convierte la conversación a mensajes ordenados
func GetConversationMessages(conversation Conversation) []Message {
	var messages []Message
	currentNode := conversation.CurrentNode

	for currentNode != "" {
		if node, exists := conversation.Mapping[currentNode]; exists && node.Message != nil {
			msg := node.Message
			author := msg.Author.Role

			//Verificar condiciones para incluir el mensaje
			//msg.Content.Parts != nil es redundante porque es un slice, en go un slice nil tiene len 0
			if len(msg.Content.Parts) > 0 && (msg.Author.Role != "system" || msg.Metadata.IsUserSystemMessage) {
				if author == "assistant" || author == "tool" {
					author = "CHATGPT"
				} else if author == "system" && msg.Metadata.IsUserSystemMessage {
					author = "Custom User Info"
				}
			}

			//Procesar solo los tipos de contenido text/multimodal
			if msg.Content.ContentType == "text" || msg.Content.ContentType == "multimodal" {
				parts := processParts(msg.Content.Parts)
				if len(parts) > 0 {
					messages = append(messages, Message{
						Author: author,
						Parts:  parts,
					})
				}
			}
		}
		currentNode = conversation.Mapping[currentNode].Parent
	}

	//Reversar el orden de los mensajes
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages
}

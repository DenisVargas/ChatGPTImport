/*
Copyright © 2025 Denis Vargas <denis.vargasrivero@outlook.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/* -------------------------------------------------------------------------- */
/*                            FORMATEADO A MARKDOWN                           */
/* -------------------------------------------------------------------------- */

func composeMarkdown(conversation Conversation) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", conversation.Title))

	messages := GetConversationMessages(conversation)
	for _, msg := range messages {
		sb.WriteString(fmt.Sprintf("**%s:**\n\n", msg.Author))

		for _, part := range msg.Parts {
			if part.Text != "" {
				// sb.WriteString(fmt.Sprintf("%s\n\n", escapeMarkdown(part.Text)))  // No se deberian escapar los mensajes de texto normales
				sb.WriteString(fmt.Sprintf("%s\n\n", part.Text))
			} else if part.Transcript != "" {
				// sb.WriteString(fmt.Sprintf("[Transcript]\n%s\n\n", escapeMarkdown(part.Transcript)))
				sb.WriteString(fmt.Sprintf("[Transcript]\n%s\n\n", part.Transcript))
			} else if part.Asset != nil {
				// Simula assetsJson lookup - en producción usarías un map real
				sb.WriteString("[File]: [asset_pointer]\n\n")
			}
		}
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

/* -------------------------------------------------------------------------- */
/*                              SALIDA A ARCHIVO                              */
/* -------------------------------------------------------------------------- */

func renderMarkdown(markdownContent, path string) error {
	if path == "" {
		path = fmt.Sprintf("conversation_%d.md", time.Now().UnixNano()/1e6)
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return os.WriteFile(path, []byte(markdownContent), 0644)
}

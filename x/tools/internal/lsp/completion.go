// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
<<<<<<< HEAD
	"bytes"
=======
	"context"
>>>>>>> bd25a1f6d07d2d464980e6a8576c1ed59bb3950a
	"fmt"
	"sort"
	"strings"

	"golang.org/x/tools/internal/lsp/protocol"
	"golang.org/x/tools/internal/lsp/source"
<<<<<<< HEAD
)

func toProtocolCompletionItems(items []source.CompletionItem, prefix string, pos protocol.Position, snippetsSupported, signatureHelpEnabled bool) []protocol.CompletionItem {
	var results []protocol.CompletionItem
	sort.Slice(items, func(i, j int) bool {
		return items[i].Score > items[j].Score
	})
	insertTextFormat := protocol.PlainTextFormat
	if snippetsSupported {
		insertTextFormat = protocol.SnippetTextFormat
	}
	for i, item := range items {
		// Matching against the label.
		if !strings.HasPrefix(item.Label, prefix) {
			continue
		}
		insertText, triggerSignatureHelp := labelToProtocolSnippets(item.Label, item.Kind, insertTextFormat, signatureHelpEnabled)
		if strings.HasPrefix(insertText, prefix) {
			insertText = insertText[len(prefix):]
		}
		i := protocol.CompletionItem{
			Label:            item.Label,
			Detail:           item.Detail,
			Kind:             float64(toProtocolCompletionItemKind(item.Kind)),
			InsertTextFormat: insertTextFormat,
=======
	"golang.org/x/tools/internal/span"
)

func (s *Server) completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	uri := span.NewURI(params.TextDocument.URI)
	view := s.findView(ctx, uri)
	f, m, err := newColumnMap(ctx, view, uri)
	if err != nil {
		return nil, err
	}
	spn, err := m.PointSpan(params.Position)
	if err != nil {
		return nil, err
	}
	rng, err := spn.Range(m.Converter)
	if err != nil {
		return nil, err
	}
	items, prefix, err := source.Completion(ctx, f, rng.Start)
	if err != nil {
		s.log.Infof(ctx, "no completions found for %s:%v:%v: %v", uri, int(params.Position.Line), int(params.Position.Character), err)
		items = []source.CompletionItem{}
	}
	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        toProtocolCompletionItems(items, prefix, params.Position, s.insertTextFormat, s.usePlaceholders),
	}, nil
}

func toProtocolCompletionItems(candidates []source.CompletionItem, prefix string, pos protocol.Position, insertTextFormat protocol.InsertTextFormat, usePlaceholders bool) []protocol.CompletionItem {
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})
	items := []protocol.CompletionItem{}
	for i, candidate := range candidates {
		// Match against the label.
		if !strings.HasPrefix(candidate.Label, prefix) {
			continue
		}
		insertText := candidate.InsertText
		if insertTextFormat == protocol.SnippetTextFormat {
			if usePlaceholders && candidate.PlaceholderSnippet != nil {
				insertText = candidate.PlaceholderSnippet.String()
			} else if candidate.Snippet != nil {
				insertText = candidate.Snippet.String()
			}
		}
		// If the user has already typed some part of the completion candidate,
		// don't insert that portion of the text.
		if strings.HasPrefix(insertText, prefix) {
			insertText = insertText[len(prefix):]
		}
		// Don't filter on text that might have snippets in it.
		filterText := candidate.InsertText
		if strings.HasPrefix(filterText, prefix) {
			filterText = filterText[len(prefix):]
		}
		item := protocol.CompletionItem{
			Label:  candidate.Label,
			Detail: candidate.Detail,
			Kind:   toProtocolCompletionItemKind(candidate.Kind),
>>>>>>> bd25a1f6d07d2d464980e6a8576c1ed59bb3950a
			TextEdit: &protocol.TextEdit{
				NewText: insertText,
				Range: protocol.Range{
					Start: pos,
					End:   pos,
				},
			},
<<<<<<< HEAD
			// InsertText is deprecated in favor of TextEdit.
			InsertText: insertText,
			// This is a hack so that the client sorts completion results in the order
			// according to their score. This can be removed upon the resolution of
			// https://github.com/Microsoft/language-server-protocol/issues/348.
			SortText: fmt.Sprintf("%05d", i),
		}
		// If we are completing a function, we should trigger signature help if possible.
		if triggerSignatureHelp && signatureHelpEnabled {
			i.Command = &protocol.Command{
				Command: "editor.action.triggerParameterHints",
			}
		}
		results = append(results, i)
	}
	return results
=======
			InsertTextFormat: insertTextFormat,
			// This is a hack so that the client sorts completion results in the order
			// according to their score. This can be removed upon the resolution of
			// https://github.com/Microsoft/language-server-protocol/issues/348.
			SortText:   fmt.Sprintf("%05d", i),
			FilterText: filterText,
			Preselect:  i == 0,
		}
		// Trigger signature help for any function or method completion.
		// This is helpful even if a function does not have parameters,
		// since we show return types as well.
		switch item.Kind {
		case protocol.FunctionCompletion, protocol.MethodCompletion:
			item.Command = &protocol.Command{
				Command: "editor.action.triggerParameterHints",
			}
		}
		items = append(items, item)
	}
	return items
>>>>>>> bd25a1f6d07d2d464980e6a8576c1ed59bb3950a
}

func toProtocolCompletionItemKind(kind source.CompletionItemKind) protocol.CompletionItemKind {
	switch kind {
	case source.InterfaceCompletionItem:
		return protocol.InterfaceCompletion
	case source.StructCompletionItem:
		return protocol.StructCompletion
	case source.TypeCompletionItem:
		return protocol.TypeParameterCompletion // ??
	case source.ConstantCompletionItem:
		return protocol.ConstantCompletion
	case source.FieldCompletionItem:
		return protocol.FieldCompletion
	case source.ParameterCompletionItem, source.VariableCompletionItem:
		return protocol.VariableCompletion
	case source.FunctionCompletionItem:
		return protocol.FunctionCompletion
	case source.MethodCompletionItem:
		return protocol.MethodCompletion
	case source.PackageCompletionItem:
		return protocol.ModuleCompletion // ??
	default:
		return protocol.TextCompletion
	}
}
<<<<<<< HEAD

func labelToProtocolSnippets(label string, kind source.CompletionItemKind, insertTextFormat protocol.InsertTextFormat, signatureHelpEnabled bool) (string, bool) {
	switch kind {
	case source.ConstantCompletionItem:
		// The label for constants is of the format "<identifier> = <value>".
		// We should not insert the " = <value>" part of the label.
		if i := strings.Index(label, " ="); i >= 0 {
			return label[:i], false
		}
	case source.FunctionCompletionItem, source.MethodCompletionItem:
		var trimmed, params string
		if i := strings.Index(label, "("); i >= 0 {
			trimmed = label[:i]
			params = strings.Trim(label[i:], "()")
		}
		if params == "" || trimmed == "" {
			return label, true
		}
		// Don't add parameters or parens for the plaintext insert format.
		if insertTextFormat == protocol.PlainTextFormat {
			return trimmed, true
		}
		// If we do have signature help enabled, the user can see parameters as
		// they type in the function, so we just return empty parentheses.
		if signatureHelpEnabled {
			return trimmed + "($1)", true
		}
		// If signature help is not enabled, we should give the user parameters
		// that they can tab through. The insert text format follows the
		// specification defined by Microsoft for LSP. The "$", "}, and "\"
		// characters should be escaped.
		r := strings.NewReplacer(
			`\`, `\\`,
			`}`, `\}`,
			`$`, `\$`,
		)
		b := bytes.NewBufferString(trimmed)
		b.WriteByte('(')
		for i, p := range strings.Split(params, ",") {
			if i != 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(b, "${%v:%v}", i+1, r.Replace(strings.Trim(p, " ")))
		}
		b.WriteByte(')')
		return b.String(), false

	}
	return label, false
}
=======
>>>>>>> bd25a1f6d07d2d464980e6a8576c1ed59bb3950a

package formcollection

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

type rewriteTarget struct {
	FormKey string
	Action  string
}

func discoverForms(file UploadedHTMLFile) ([]DiscoveredForm, error) {
	root, err := html.Parse(bytes.NewReader(file.Body))
	if err != nil {
		return nil, err
	}

	forms := make([]DiscoveredForm, 0)
	position := 0
	traverse(root, func(node *html.Node) {
		if node.Type != html.ElementNode || !strings.EqualFold(node.Data, "form") {
			return
		}

		position++
		forms = append(forms, buildDiscoveredForm(file, node, position))
	})

	return forms, nil
}

func rewriteForms(pagePath string, body []byte, targets []rewriteTarget) ([]byte, error) {
	if len(targets) == 0 {
		return body, nil
	}

	root, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	targetByKey := make(map[string]string, len(targets))
	for _, target := range targets {
		targetByKey[target.FormKey] = target.Action
	}

	position := 0
	traverse(root, func(node *html.Node) {
		if node.Type != html.ElementNode || !strings.EqualFold(node.Data, "form") {
			return
		}

		position++
		key := formKey(pagePath, node, position)
		action, ok := targetByKey[key]
		if !ok {
			return
		}

		setAttr(node, "action", action)
		setAttr(node, "method", "post")
	})

	var out bytes.Buffer
	if err := html.Render(&out, root); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func buildDiscoveredForm(file UploadedHTMLFile, node *html.Node, position int) DiscoveredForm {
	method := strings.ToLower(strings.TrimSpace(attr(node, "method")))
	if method == "" {
		method = "get"
	}

	fields := formFields(node)
	unsupportedReason := unsupportedReason(node, method)

	return DiscoveredForm{
		FileName:          file.FileName,
		PagePath:          file.PagePath,
		FormKey:           formKey(file.PagePath, node, position),
		Label:             deriveLabel(file.PagePath, node, fields, position),
		OriginalAction:    strings.TrimSpace(attr(node, "action")),
		Method:            method,
		Fields:            fields,
		Supported:         unsupportedReason == "",
		UnsupportedReason: unsupportedReason,
		Active:            false,
	}
}

func formKey(pagePath string, node *html.Node, position int) string {
	method := strings.ToLower(strings.TrimSpace(attr(node, "method")))
	if method == "" {
		method = "get"
	}

	parts := []string{
		pagePath,
		fmt.Sprintf("%d", position),
		method,
		strings.TrimSpace(attr(node, "action")),
	}
	for _, field := range formFields(node) {
		parts = append(parts, field.Name, field.Type)
	}

	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(sum[:])[:24]
}

func unsupportedReason(node *html.Node, method string) string {
	if method != "get" && method != "post" {
		return "unsupported_method"
	}

	enctype := strings.ToLower(strings.TrimSpace(attr(node, "enctype")))
	if enctype == "multipart/form-data" {
		return "multipart_form"
	}

	hasFileInput := false
	traverse(node, func(child *html.Node) {
		if child.Type != html.ElementNode || !strings.EqualFold(child.Data, "input") {
			return
		}
		if strings.EqualFold(attr(child, "type"), "file") {
			hasFileInput = true
		}
	})
	if hasFileInput {
		return "file_input"
	}

	return ""
}

func formFields(node *html.Node) FieldList {
	fields := make(FieldList, 0)
	seen := make(map[string]struct{})

	traverse(node, func(child *html.Node) {
		if child.Type != html.ElementNode {
			return
		}

		switch strings.ToLower(child.Data) {
		case "input", "textarea", "select", "button":
		default:
			return
		}

		name := strings.TrimSpace(attr(child, "name"))
		if name == "" {
			return
		}

		fieldType := strings.ToLower(child.Data)
		if strings.EqualFold(child.Data, "input") || strings.EqualFold(child.Data, "button") {
			if t := strings.TrimSpace(attr(child, "type")); t != "" {
				fieldType = strings.ToLower(t)
			}
		}

		key := name + "\x00" + fieldType
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}

		fields = append(fields, Field{Name: name, Type: fieldType})
	})

	return fields
}

func deriveLabel(pagePath string, node *html.Node, fields FieldList, position int) string {
	candidates := []string{
		attr(node, "aria-label"),
		attr(node, "name"),
		attr(node, "id"),
		firstSubmitLabel(node),
		pagePath,
		fmt.Sprintf("Form %d", position),
	}

	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate != "" && candidate != "/" {
			return candidate
		}
	}

	if len(fields) > 0 {
		return fields[0].Name
	}

	return fmt.Sprintf("Form %d", position)
}

func firstSubmitLabel(node *html.Node) string {
	var label string
	traverse(node, func(child *html.Node) {
		if label != "" || child.Type != html.ElementNode {
			return
		}

		if strings.EqualFold(child.Data, "input") && strings.EqualFold(attr(child, "type"), "submit") {
			label = attr(child, "value")
			return
		}

		if strings.EqualFold(child.Data, "button") {
			t := strings.ToLower(strings.TrimSpace(attr(child, "type")))
			if t == "" || t == "submit" {
				label = nodeText(child)
			}
		}
	})

	return label
}

func nodeText(node *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(node)
	return strings.Join(strings.Fields(b.String()), " ")
}

func traverse(node *html.Node, visit func(*html.Node)) {
	visit(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		traverse(child, visit)
	}
}

func attr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if strings.EqualFold(attr.Key, key) {
			return attr.Val
		}
	}
	return ""
}

func setAttr(node *html.Node, key, value string) {
	for i := range node.Attr {
		if strings.EqualFold(node.Attr[i].Key, key) {
			node.Attr[i].Key = key
			node.Attr[i].Val = value
			return
		}
	}

	node.Attr = append(node.Attr, html.Attribute{Key: key, Val: value})
	slices.SortFunc(node.Attr, func(a, b html.Attribute) int {
		return strings.Compare(a.Key, b.Key)
	})
}

func renderNode(node *html.Node) string {
	var b bytes.Buffer
	_ = html.Render(&b, node)
	return b.String()
}

func parseHTMLFragment(source string) (*html.Node, error) {
	nodes, err := html.ParseFragment(strings.NewReader(source), nil)
	if err != nil {
		return nil, err
	}
	root := &html.Node{Type: html.ElementNode, Data: "root"}
	for _, node := range nodes {
		root.AppendChild(node)
	}
	return root, nil
}

func readAll(reader io.Reader) ([]byte, error) {
	var b bytes.Buffer
	_, err := b.ReadFrom(reader)
	return b.Bytes(), err
}

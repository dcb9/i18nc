package i18nc

import (
	"encoding/json"
	"fmt"
	"github.com/gobeam/stringy"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io"
	"regexp"
	"sort"
	"strings"
)

// FromFile generate type safe Go from locale file path
func FromFile(path string, pkg string, w io.Writer) error {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	messageFile, err := bundle.LoadMessageFile(path)
	if err != nil {
		return err
	}
	return generate(messageFile, pkg, w)
}

// FromBytes generate type safe Go from locale file content
func FromBytes(content []byte, pkg string, w io.Writer) error {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load source language
	messageFile, err := bundle.ParseMessageFileBytes(content, "en.json")
	if err != nil {
		return err
	}

	return generate(messageFile, pkg, w)
}

func generate(messageFile *i18n.MessageFile, pkg string, w io.Writer) error {
	_, err := fmt.Fprintf(w, `package %s

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var Localizer *i18n.Localizer

`, pkg)
	if err != nil {
		return err
	}

	for _, m := range messageFile.Messages {
		str := stringy.New(m.ID)
		fnName := str.CamelCase()
		reg := regexp.MustCompile(`\{\{\.(\S+)\}\}`)
		matches := reg.FindAllStringSubmatch(m.Zero+m.One+m.Two+m.Few+m.Many+m.Other, -1)

		argsMap := map[string]struct{}{}
		for _, match := range matches {
			argsMap[match[1]] = struct{}{}
		}

		args := make([]string, 0, len(argsMap))
		for k := range argsMap {
			args = append(args, k)
		}
		sort.Strings(args)

		activeArgs := make([]string, 0, len(args))
		for _, k := range args {
			if k == "PluralCount" {
				continue
			}
			activeArgs = append(activeArgs, k)
		}

		templateDataList := make([]string, 0, len(args))
		for _, k := range activeArgs {
			templateDataList = append(templateDataList, fmt.Sprintf(`"%s": %s`, k, k))
		}

		var templateData string
		if len(templateDataList) > 0 {
			templateData = fmt.Sprintf(`map[string]interface{}{ %s }`, strings.Join(templateDataList, ", \n"))
		} else {
			templateData = "nil"
		}
		activeArgsStr := strings.Join(activeArgs, ", ")
		if activeArgsStr != "" {
			activeArgsStr += " interface{}, "
		}

		_, err = fmt.Fprintf(w, `func %s(%splural ...interface{}) string {
	messageID := "%s"
	var pluralCount interface{}
	if len(plural) != 0 {
		pluralCount = plural[0]
	}
	message, _ := Localizer.Localize(&i18n.LocalizeConfig{
		MessageID:      messageID,
		TemplateData:   %s,
		PluralCount:    pluralCount,
	})
    return message
}

`, fnName, activeArgsStr, m.ID, templateData)
		if err != nil {
			return err
		}

		if len(activeArgs) != 0 {
			_, err = fmt.Fprintf(w, `func %s_WithData(templateData, pluralCount interface{}) string {
	messageID := "%s"
	message, _ := Localizer.Localize(&i18n.LocalizeConfig{
		MessageID:      messageID,
		TemplateData:   templateData,
		PluralCount:    pluralCount,
	})
    return message
}

`, fnName, m.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

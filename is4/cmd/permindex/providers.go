package main

type PrefixByTemplateProvider struct{}

func (p PrefixByTemplateProvider) Provide(template string) (string, error) {
	justWord := true
	for _, char := range template {
		if char == '*' {
			justWord = false
			break
		}
	}
	if justWord {
		return "$" + template, nil
	}
	for i := 0; i < len(template); i++ {
		if template[i] == '*' {
			return template[i+1:] + "$" + template[:i], nil
		}
	}
	return "", nil
}

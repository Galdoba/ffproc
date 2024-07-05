package translit

import "strings"

type Option string

var Title Option = "Title"

func Process(origin string, opts ...Option) string {
	result := ""
	glyphs := strings.Split(origin, "")
	for _, gl := range glyphs {
		result += change(gl)
	}
	words := strings.Split(result, "_")
	result = ""
	for _, w := range words {
		if w == "" {
			continue
		}
		result += w + "_"
	}
	result = strings.TrimSuffix(result, "_")
	if result == "" {
		return result
	}
	if haveOption(Title, opts) {
		ltrs := strings.Split(result, "")
		ltrs[0] = strings.ToUpper(ltrs[0])
		result = strings.Join(ltrs, "")
	}

	return result
}

func haveOption(opt Option, pool []Option) bool {
	for _, check := range pool {
		if opt == check {
			return true
		}
	}
	return false
}

func change(a string) string {
	a = strings.ToLower(a)
	switch a {
	default:
		return "_"
	case "а", "б", "в", "г", "д", "е", "ё", "ж", "з", "и", "й", "к", "л", "м", "н", "о", "п", "р", "с", "т", "у", "ф", "х", "ц", "ч", "ш", "щ", "ъ", "ы", "ь", "э", "ю", "я":
		lMap := make(map[string]string)
		lMap = map[string]string{
			"а": "a",
			"б": "b",
			"в": "v",
			"г": "g",
			"д": "d",
			"е": "e",
			"ё": "e",
			"ж": "zh",
			"з": "z",
			"и": "i",
			"й": "y",
			"к": "k",
			"л": "l",
			"м": "m",
			"н": "n",
			"о": "o",
			"п": "p",
			"р": "r",
			"с": "s",
			"т": "t",
			"у": "u",
			"ф": "f",
			"х": "h",
			"ц": "c",
			"ч": "ch",
			"ш": "sh",
			"щ": "sh",
			"ъ": "",
			"ы": "y",
			"ь": "",
			"э": "e",
			"ю": "yu",
			"я": "ya"}
		return lMap[a]
	case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", `/`, `\`:
		return a
	}
}

package translit

import "strings"

type Option string

var Title Option = "Title"
var KeepRegister Option = "Keep Register"

func Process(origin string, opts ...Option) string {
	result := ""
	if !haveOption(KeepRegister, opts) {
		origin = strings.ToLower(origin)
	}
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
	case "А", "Б", "В", "Г", "Д", "Е", "Ё", "Ж", "З", "И", "Й", "К", "Л", "М", "Н", "О", "П", "Р", "С", "Т", "У", "Ф", "Х", "Ц", "Ч", "Ш", "Щ", "Ъ", "Ы", "Ь", "Э", "Ю", "Я":
		lMap := make(map[string]string)
		lMap = map[string]string{
			"А": "A",
			"Б": "B",
			"В": "V",
			"Г": "G",
			"Д": "D",
			"Е": "E",
			"Ё": "E",
			"Ж": "ZH",
			"З": "Z",
			"И": "I",
			"Й": "Y",
			"К": "K",
			"Л": "L",
			"М": "M",
			"Н": "N",
			"О": "O",
			"П": "P",
			"Р": "R",
			"С": "S",
			"Т": "T",
			"У": "U",
			"Ф": "F",
			"Х": "H",
			"Ц": "C",
			"Ч": "CH",
			"Ш": "SH",
			"Щ": "SH",
			"Ъ": "",
			"Ы": "Y",
			"Ь": "",
			"Э": "E",
			"Ю": "YU",
			"Я": "YA"}
		return lMap[a]
	case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", `/`, `\`:
		return a
	}
}

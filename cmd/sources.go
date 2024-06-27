package cmd

// langCode is the code of the language
type langCode string

// enum of languages.
const (
	Arabic          langCode = "ar"
	Bengali         langCode = "bn"
	Chinese         langCode = "zh"
	English         langCode = "en"
	Spanish         langCode = "es"
	French          langCode = "fr"
	Indonesian      langCode = "id"
	Russian         langCode = "ru"
	Swedish         langCode = "sv"
	Turkish         langCode = "tr"
	Urdu            langCode = "ur"
	Transliteration langCode = "transliteration"
)

// sources is a map of languages and their corresponding sources.
var sources = map[langCode]string{
	Arabic:          "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran.json",
	Bengali:         "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_bn.json",
	Chinese:         "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_zh.json",
	English:         "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_en.json",
	Spanish:         "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_es.json",
	French:          "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_fr.json",
	Indonesian:      "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_id.json",
	Russian:         "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_ru.json",
	Swedish:         "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_sv.json",
	Turkish:         "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_tr.json",
	Urdu:            "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_ur.json",
	Transliteration: "https://cdn.jsdelivr.net/npm/quran-json@3.1.2/dist/quran_transliteration.json",
}

package kufar

const (
	gtsyBelarus  = "country-belarus"
	gtsyMinsk    = "country-belarus~province-minsk~locality-minsk"
	gtsyBrestO   = "country-belarus~province-brestskaja_oblast"
	gtsyGomelO   = "country-belarus~province-gomelskaja_oblast"
	gtsyGrodnoO  = "country-belarus~province-grodnenskaja_oblast"
	gtsyMogilevO = "country-belarus~province-mogilyovskaja_oblast"
	gtsyMinskO   = "country-belarus~province-minskaja_oblast"
	gtsyVitebskO = "country-belarus~province-vitebskaja_oblast"
)

// pathTokens are SEO path segments that are not a city/district.
var pathTokens = map[string]map[string]string{
	"snyat":           {"typ": "let"},
	"kvartiru":        {"cat": "1010"},
	"bez-posrednikov": {"cmp": "0"},
	"1k":              {"rms": "v.or:1"},
	"2k":              {"rms": "v.or:2"},
	"3k":              {"rms": "v.or:3"},
	"4k":              {"rms": "v.or:4"},
	"5k":              {"rms": "v.or:5"},
}

// placeParams maps /l/{slug}/... to Kufar API geo filters.
// Microdistricts use red; administrative districts of Minsk use ar + gtsy coder_district.
var placeParams = map[string]map[string]string{
	"belarus":             {"gtsy": gtsyBelarus},
	"minsk":               {"gtsy": gtsyMinsk, "rgn": "7"},
	"brest":               {"gtsy": gtsyBrestO + "~locality-brest", "rgn": "1", "ar": "1"},
	"gomel":               {"gtsy": gtsyGomelO + "~locality-gomel", "rgn": "2", "ar": "5"},
	"grodno":              {"gtsy": gtsyGrodnoO + "~locality-grodno", "rgn": "3", "ar": "9"},
	"mogilev":             {"gtsy": gtsyMogilevO + "~locality-mogilyov", "rgn": "4", "ar": "13"},
	"vitebsk":             {"gtsy": gtsyVitebskO + "~locality-vitebsk", "rgn": "6", "ar": "18"},
	"baranovichi":         {"gtsy": gtsyBrestO + "~locality-baranovichi", "rgn": "1", "ar": "37"},
	"brestskaya-oblast":   {"gtsy": gtsyBrestO, "rgn": "1"},
	"gomelskaya-oblast":   {"gtsy": gtsyGomelO, "rgn": "2"},
	"grodnenskaya-oblast": {"gtsy": gtsyGrodnoO, "rgn": "3"},
	"mogilevskaya-oblast": {"gtsy": gtsyMogilevO, "rgn": "4"},
	"minskaya-oblast":     {"gtsy": gtsyMinskO, "rgn": "5"},
	"vitebskaya-oblast":   {"gtsy": gtsyVitebskO, "rgn": "6"},

	"minsk-centralnyj-rajon":   minskDistrict("22"),
	"minsk-sovetskij-rajon":    minskDistrict("23"),
	"minsk-pervomajskij-rajon": minskDistrict("24"),
	"minsk-partizanskij-rajon": minskDistrict("25"),
	"minsk-zavodskoj-rajon":    minskDistrict("26"),
	"minsk-leninskij-rajon":    minskDistrict("27"),
	"minsk-oktyabrskij-rajon":  minskDistrict("28"),
	"minsk-moskovskij-rajon":   minskDistrict("29"),
	"minsk-frunzenskij-rajon":  minskDistrict("30"),

	"minsk-akademiya-nauk-mkrn":     minskMicro("5"),
	"minsk-angarskaya-mkrn":         minskMicro("10"),
	"minsk-aerodromnaya-mkrn":       minskMicro("15"),
	"minsk-borovaya-mkrn":           minskMicro("20"),
	"minsk-borovlyany-mkrn":         minskMicro("25"),
	"minsk-brilevichi-mkrn":         minskMicro("30"),
	"minsk-verhnij-gorod-mkrn":      minskMicro("35"),
	"minsk-vesnyanka-mkrn":          minskMicro("40"),
	"minsk-vostok-mkrn":             minskMicro("45"),
	"minsk-vostochnyj-mkrn":         minskMicro("50"),
	"minsk-grushevka-mkrn":          minskMicro("55"),
	"minsk-dombrovka-mkrn":          minskMicro("65"),
	"minsk-drazhnya-mkrn":           minskMicro("70"),
	"minsk-zapad-mkrn":              minskMicro("80"),
	"minsk-zaslavskaya-mkrn":        minskMicro("90"),
	"minsk-zacen-mkrn":              minskMicro("95"),
	"minsk-zelenyj-lug-mkrn":        minskMicro("100"),
	"minsk-kamennaya-gorka-mkrn":    minskMicro("105"),
	"minsk-kaskad-mkrn":             minskMicro("110"),
	"minsk-komarovka-mkrn":          minskMicro("120"),
	"minsk-krasnyj-bor-mkrn":        minskMicro("125"),
	"minsk-kuncevshchina-mkrn":      minskMicro("130"),
	"minsk-kurasovshchina-mkrn":     minskMicro("135"),
	"minsk-lebyazhij-mkrn":          minskMicro("140"),
	"minsk-loshica-mkrn":            minskMicro("145"),
	"minsk-malinovka-mkrn":          minskMicro("150"),
	"minsk-masyukovshchina-mkrn":    minskMicro("160"),
	"minsk-medvezhino-mkrn":         minskMicro("165"),
	"minsk-mihalovo-mkrn":           minskMicro("175"),
	"minsk-novinki-mkrn":            minskMicro("180"),
	"minsk-ozerishche-mkrn":         minskMicro("185"),
	"minsk-poselok-raduzhnyj-mkrn":  minskMicro("195"),
	"minsk-p-ulihova-mkrn":          minskMicro("200"),
	"minsk-rakovskoe-shosse-mkrn":   minskMicro("205"),
	"minsk-rzhavec-mkrn":            minskMicro("210"),
	"minsk-severnyj-poselok-mkrn":   minskMicro("215"),
	"minsk-selhozposelok-mkrn":      minskMicro("220"),
	"minsk-semashko-mkrn":           minskMicro("225"),
	"minsk-serebryanka-mkrn":        minskMicro("230"),
	"minsk-serova-mkrn":             minskMicro("235"),
	"minsk-slepyanka-mkrn":          minskMicro("240"),
	"minsk-sokol-mkrn":              minskMicro("245"),
	"minsk-sosny-mkrn":              minskMicro("250"),
	"minsk-stepyanka-mkrn":          minskMicro("255"),
	"minsk-storozhovka-mkrn":        minskMicro("260"),
	"minsk-suharevo-mkrn":           minskMicro("265"),
	"minsk-troickoe-predmeste-mkrn": minskMicro("275"),
	"minsk-uruche-mkrn":             minskMicro("280"),
	"minsk-harkovskaya-mkrn":        minskMicro("290"),
	"minsk-chizhovka-mkrn":          minskMicro("305"),
	"minsk-shabany-mkrn":            minskMicro("310"),
	"minsk-yugo-zapad-mkrn":         minskMicro("320"),
}

func minskDistrict(ar string) map[string]string {
	return map[string]string{
		"gtsy": gtsyMinsk + "~coder_district-" + ar,
		"rgn":  "7",
		"ar":   ar,
	}
}

func minskMicro(id string) map[string]string {
	return map[string]string{
		"gtsy": gtsyMinsk,
		"rgn":  "7",
		"red":  "v.or:" + id,
	}
}

// metroIDs maps /metro-{slug} to the mee API filter. Minsk only.
var metroIDs = map[string]string{
	"avtozavodskaya":                    "v.or:2",
	"akademiya-nauk":                    "v.or:3",
	"aerodromnaya":                      "v.or:36",
	"borisovskij-trakt":                 "v.or:4",
	"vostok":                            "v.or:5",
	"grushevka":                         "v.or:6",
	"institut-kultury":                  "v.or:7",
	"kamennaya-gorka":                   "v.or:8",
	"kovalskaya-sloboda":                "v.or:32",
	"kuncevshchina":                     "v.or:9",
	"malinovka":                         "v.or:10",
	"mihalovo":                          "v.or:11",
	"mogilevskaya":                      "v.or:12",
	"molodezhnaya":                      "v.or:13",
	"moskovskaya":                       "v.or:14",
	"nemiga":                            "v.or:15",
	"nemorshanskij-sad":                 "v.or:35",
	"oktyabrskaya-kupalovskaya":         "v.or:16",
	"park-chelyuskincev":                "v.or:17",
	"partizanskaya":                     "v.or:18",
	"pervomajskaya":                     "v.or:19",
	"petrovshchina":                     "v.or:20",
	"ploshchad-lenina":                  "v.or:21",
	"ploshchad-pobedy":                  "v.or:22",
	"ploshchad-frantishka-bogushevicha": "v.or:33",
	"ploshchad-yakuba-kolasa":           "v.or:23",
	"proletarskaya":                     "v.or:24",
	"pushkinskaya":                      "v.or:25",
	"sportivnaya":                       "v.or:26",
	"sluckij-gostinec":                  "v.or:34",
	"traktornyj-zavod":                  "v.or:27",
	"uruche":                            "v.or:28",
	"frunzenskaya":                      "v.or:29",
}

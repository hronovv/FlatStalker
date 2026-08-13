package kufar

import (
	"os"
	"strings"
	"testing"
	"time"
)

type checkURL struct {
	name string
	raw  string
	ads  bool
}

func TestLiveParseAndFetchDiverse(t *testing.T) {
	if os.Getenv("LIVE_KUFAR") == "" {
		t.Skip("set LIVE_KUFAR=1 to run against kufar.by")
	}

	client := NewClient()
	cases := liveCases()

	var parseFail, adsFail, parseOK, adsOK int
	for i, c := range cases {
		if i > 0 && mayHitKufar(c.raw) {
			time.Sleep(400 * time.Millisecond)
		}
		params, err := ParseSearchURL(c.raw)
		if err != nil {
			parseFail++
			t.Errorf("PARSE FAIL %s\n  %s\n  %v", c.name, c.raw, err)
			continue
		}
		parseOK++
		t.Logf("OK %s gtsy=%s rgn=%s ar=%s red=%s mee=%s rms=%s prc=%s cur=%s",
			c.name, params["gtsy"], params["rgn"], params["ar"], params["red"], params["mee"], params["rms"], params["prc"], params["cur"])
		if params["gtsy"] == "" || params["cat"] != "1010" || params["typ"] != "let" {
			t.Errorf("PARSE WEIRD %s gtsy=%q cat=%q typ=%q", c.name, params["gtsy"], params["cat"], params["typ"])
		}
		if !c.ads {
			continue
		}
		ads, err := client.FetchAds(c.raw)
		if err != nil {
			adsFail++
			t.Errorf("ADS FAIL %s\n  %s\n  %v", c.name, c.raw, err)
			continue
		}
		adsOK++
		t.Logf("ADS %s count=%d", c.name, len(ads))
	}
	t.Logf("summary parse ok=%d fail=%d; ads ok=%d fail=%d; total=%d", parseOK, parseFail, adsOK, adsFail, len(cases))
}

func rent(slug, name string, ads bool) checkURL {
	return checkURL{name, "https://re.kufar.by/l/" + slug + "/snyat/kvartiru", ads}
}

func liveCases() []checkURL {
	out := []checkURL{
		{"belarus", "https://re.kufar.by/l/belarus/snyat/kvartiru", true},
		{"minsk + usd + price", "https://re.kufar.by/l/minsk/snyat/kvartiru?cur=USD&prc=r%3A0%2C500", true},
		{"minsk 1k path", "https://re.kufar.by/l/minsk/snyat/kvartiru/1k", true},
		{"minsk no agents", "https://re.kufar.by/l/minsk/snyat/kvartiru/bez-posrednikov", true},
		{"kuncevshchina mkrn", "https://re.kufar.by/l/minsk-kuncevshchina-mkrn/snyat/kvartiru?cur=USD", true},
		{"uruche mkrn", "https://re.kufar.by/l/minsk-uruche-mkrn/snyat/kvartiru", true},
		{"metro kuncevshchina", "https://re.kufar.by/l/minsk/snyat/kvartiru/metro-kuncevshchina", true},
		{"metro malinovka + 1k", "https://re.kufar.by/l/minsk/snyat/kvartiru/metro-malinovka/1k", true},
		{"frunzenskij rajon", "https://re.kufar.by/l/minsk-frunzenskij-rajon/snyat/kvartiru", true},
		{"brest", "https://re.kufar.by/l/brest/snyat/kvartiru?cur=USD", true},
		{"gomel 2k", "https://re.kufar.by/l/gomel/snyat/kvartiru/2k", true},
		{"grodno", "https://re.kufar.by/l/grodno/snyat/kvartiru", true},
		{"vitebsk", "https://re.kufar.by/l/vitebsk/snyat/kvartiru", true},
		{"mogilev", "https://re.kufar.by/l/mogilev/snyat/kvartiru", true},
		{"brestskaya oblast", "https://re.kufar.by/l/brestskaya-oblast/snyat/kvartiru", true},
		{"minskaya oblast", "https://re.kufar.by/l/minskaya-oblast/snyat/kvartiru", false},
		{"vitebskaya oblast", "https://re.kufar.by/l/vitebskaya-oblast/snyat/kvartiru", false},
		{"gomelskaya oblast", "https://re.kufar.by/l/gomelskaya-oblast/snyat/kvartiru", false},
		{"grodnenskaya oblast", "https://re.kufar.by/l/grodnenskaya-oblast/snyat/kvartiru", false},
		{"mogilevskaya oblast", "https://re.kufar.by/l/mogilevskaya-oblast/snyat/kvartiru", false},
	}

	cities := [][2]string{
		{"Минский район", "minskij-rajon"},
		{"Березино", "berezino"},
		{"Борисов", "borisov"},
		{"Боровляны", "borovlyany"},
		{"Вилейка", "vilejka"},
		{"Воложин", "volozhin"},
		{"Городея", "gorodeya"},
		{"Дзержинск", "dzerzhinsk"},
		{"Ждановичи", "zhdanovichi"},
		{"Жодино", "zhodino"},
		{"Заславль", "zaslavl"},
		{"Зелёный Бор", "zelenyj-bor"},
		{"Ивенец", "ivenec"},
		{"Клецк", "kleck"},
		{"Копыль", "kopyl"},
		{"Крупки", "krupki"},
		{"Логойск", "logojsk"},
		{"Любань", "lyuban"},
		{"Марьина Горка", "marina-gorka"},
		{"Мачулищи", "machulishchi"},
		{"Молодечно", "molodechno"},
		{"Мядель", "myadel"},
		{"Несвиж", "nesvizh"},
		{"Плещеницы", "pleshchenicy"},
		{"Радошковичи", "radoshkovichi"},
		{"Руденск", "rudensk"},
		{"Слуцк", "sluck"},
		{"Смиловичи", "smilovichi"},
		{"Смолевичи", "smolevichi"},
		{"Солигорск", "soligorsk"},
		{"Старобин", "starobin"},
		{"Старые Дороги", "starye-dorogi"},
		{"Столбцы", "stolbcy"},
		{"Узда", "uzda"},
		{"Фаниполь", "fanipol"},
		{"Червень", "cherven"},
		{"Брестский район", "brestskij-rajon"},
		{"Барановичи", "baranovichi"},
		{"Белоозерск", "beloozersk"},
		{"Береза", "bereza"},
		{"Высокое", "vysokoe"},
		{"Ганцевичи", "gancevichi"},
		{"Городище", "gorodishche"},
		{"Давид-Городок", "david-gorodok"},
		{"Дрогичин", "drogichin"},
		{"Жабинка", "zhabinka"},
		{"Иваново", "ivanovo"},
		{"Ивацевичи", "ivacevichi"},
		{"Каменец", "kamenec"},
		{"Кобрин", "kobrin"},
		{"Лунинец", "luninec"},
		{"Ляховичи", "lyahovichi"},
		{"Малорита", "malorita"},
		{"Микашевичи", "mikashevichi"},
		{"Пинск", "pinsk"},
		{"Пружаны", "pruzhany"},
		{"Ружаны", "ruzhany"},
		{"Столин", "stolin"},
		{"Телеханы", "telehany"},
		{"Витебский район", "vitebskij-rajon"},
		{"Барань", "baran"},
		{"Бешенковичи", "beshenkovichi"},
		{"Браслав", "braslav"},
		{"Браславский район", "braslavskij-rajon"},
		{"Верхнедвинск", "verhnedvinsk"},
		{"Ветрино", "vetrino"},
		{"Коханово", "kohanovo"},
		{"Глубокое", "glubokoe"},
		{"Городок", "gorodok"},
		{"Докшицы", "dokshicy"},
		{"Дубровно", "dubrovno"},
		{"Лепель", "lepel"},
		{"Лиозно", "liozno"},
		{"Миоры", "miory"},
		{"Новолукомль", "novolukoml"},
		{"Новополоцк", "novopolock"},
		{"Орша", "orsha"},
		{"Полоцк", "polock"},
		{"Поставы", "postavy"},
		{"Россоны", "rossony"},
		{"Сенно", "senno"},
		{"Толочин", "tolochin"},
		{"Ушачи", "ushachi"},
		{"Чашники", "chashniki"},
		{"Шарковщина", "sharkovshchina"},
		{"Шумилино", "shumilino"},
		{"Гомельский район", "gomelskij-rajon"},
		{"Большевик", "bolshevik"},
		{"Брагин", "bragin"},
		{"Буда-Кошелево", "buda-koshelevo"},
		{"Василевичи", "vasilevichi"},
		{"Ветка", "vetka"},
		{"Добруш", "dobrush"},
		{"Ельск", "elsk"},
		{"Житковичи", "zhitkovichi"},
		{"Жлобин", "zhlobin"},
		{"Калинковичи", "kalinkovichi"},
		{"Корма", "korma"},
		{"Костюковка", "kostyukovka"},
		{"Лельчицы", "lelchicy"},
		{"Лоев", "loev"},
		{"Мозырь", "mozyr"},
		{"Наровля", "narovlya"},
		{"Октябрьский", "oktyabrskij"},
		{"Паричи", "parichi"},
		{"Петриков", "petrikov"},
		{"Речица", "rechica"},
		{"Рогачев", "rogachev"},
		{"Светлогорск", "svetlogorsk"},
		{"Тереховка", "terehovka"},
		{"Туров", "turov"},
		{"Уваровичи", "uvarovichi"},
		{"Хойники", "hojniki"},
		{"Чечерск", "chechersk"},
		{"Гродненский район", "grodnenskij-rajon"},
		{"Берёзовка", "berezovka"},
		{"Берестовица", "berestovica"},
		{"Большая Берестовица", "bolshaya-berestovica"},
		{"Волковыск", "volkovysk"},
		{"Вороново", "voronovo"},
		{"Дятлово", "dyatlovo"},
		{"Зельва", "zelva"},
		{"Ивье", "ivie"},
		{"Кореличи", "korelichi"},
		{"Красносельский", "krasnoselskij"},
		{"Лида", "lida"},
		{"Мир", "mir"},
		{"Мосты", "mosty"},
		{"Новогрудок", "novogrudok"},
		{"Островец", "ostrovec"},
		{"Ошмяны", "oshmyany"},
		{"Радунь", "radun"},
		{"Россь", "ross"},
		{"Свислочь", "svisloch"},
		{"Скидель", "skidel"},
		{"Слоним", "slonim"},
		{"Сморгонь", "smorgon"},
		{"Щучин", "shchuchin"},
		{"Могилевский район", "mogilevskij-rajon"},
		{"Белыничи", "belynichi"},
		{"Бобруйск", "bobrujsk"},
		{"Быхов", "byhov"},
		{"Глуск", "glusk"},
		{"Горки", "gorki-city"},
		{"Дрибин", "dribin"},
		{"Елизово", "elizovo"},
		{"Кировск", "kirovsk"},
		{"Климовичи", "klimovichi"},
		{"Кличев", "klichev"},
		{"Костюковичи", "kostyukovichi"},
		{"Краснополье", "krasnopole"},
		{"Кричев", "krichev"},
		{"Круглое", "krugloe"},
		{"Мстиславль", "mstislavl"},
		{"Осиповичи", "osipovichi"},
		{"Славгород", "slavgorod"},
		{"Хотимск", "hotimsk"},
		{"Чаусы", "chausy"},
		{"Чериков", "cherikov"},
		{"Шклов", "shklov"},
	}
	for _, c := range cities {
		fetchAds := c[1] == "pinsk" || c[1] == "lida" || c[1] == "bobrujsk" || c[1] == "mozyr" || c[1] == "orsha"
		out = append(out, rent(c[1], c[0], fetchAds))
	}
	return out
}

func mayHitKufar(raw string) bool {
	u := strings.ToLower(raw)
	for slug := range placeParams {
		if strings.Contains(u, "/l/"+slug+"/") && !strings.Contains(u, "/l/"+slug+"-") {
			return false
		}
	}
	return true
}

const tg = window.Telegram?.WebApp;

const params = new URLSearchParams(location.search);
const API_BASE =
  params.get("api") ||
  window.FLATSTALKER_API ||
  "http://127.0.0.1:8080";

const LANG_KEY = "flatstalker_lang";

const I18N = {
  ru: {
    hello_guest: "Привет — это твой кабинет",
    hello_user: "{name}, это твой кабинет",
    eyebrow: "Умный поиск аренды",
    hero_title: "Новая квартира — раньше других",
    hero_text:
      "FlatStalker следит за твоим поиском и присылает новые объявления в Telegram.",
    status_title: "Мониторинг активен",
    status_sub: "Проверка по тарифу",
    add_title: "Добавить поиск",
    add_hint:
      "Настрой фильтры на Kufar, скопируй ссылку поиска и вставь сюда — дальше бот сам пришлёт новые объявления.",
    link_label: "Ссылка",
    link_submit: "Добавить ссылку",
    link_note: "Только ссылка поиска аренды на Kufar.",
    my_links_title: "Мои ссылки",
    links_empty: "Пока пусто — добавь ссылку поиска выше.",
    links_open_tg: "Открой кабинет из Telegram, чтобы видеть ссылки.",
    links_need_start: "Сначала нажми /start в боте.",
    links_load_error: "Не удалось загрузить ссылки.",
    how_title: "Как это работает",
    how_1: "Открываешь бота и кабинет",
    how_2: "Настраиваешь поиск на Kufar и вставляешь ссылку",
    how_3: "Получаешь уведомление со ссылкой на новое объявление",
    caps_title: "Возможности",
    cap1_title: "Следим за объявлениями",
    cap1_text: "Новые варианты не нужно искать вручную",
    cap2_title: "Быстрые уведомления",
    cap2_text: "Ссылка приходит сразу, без долгих сводок",
    cap3_title: "Твой поиск с Kufar",
    cap3_text: "Город, цена и комнаты — в ссылке поиска",
    cap4_title: "По всей РБ",
    cap4_text: "Под поиск аренды в Беларуси",
    pricing_title: "Тарифы",
    plan_badge_current: "Твой",
    plan_current_loading: "Загружаем твой тариф…",
    plan_current: "Твой тариф: {plan}. Проверка {interval}.",
    plan_current_guest: "Открой кабинет из Telegram, чтобы увидеть тариф.",
    plan_current_need_start: "Сначала нажми /start в боте.",
    plan_current_error: "Не удалось загрузить тариф.",
    plan_interval_every: "Проверка объявлений {interval}",
    plan_links_one: "{n} ссылка",
    plan_links_few: "{n} ссылки",
    plan_links_many: "{n} ссылок",
    interval_seconds: "каждые {n} сек",
    interval_minutes: "каждые {n} мин",
    interval_short_seconds: "{n} сек",
    interval_short_minutes: "{n} мин",
    plan_price: "{amount} BYN",
    plan_price_day: "{amount} BYN/день",
    plan_price_free: "Бесплатно",
    plan_period_chip: "{n}д",
    plan_period_hint: "Цена за",
    plan_days_one: "{n} день",
    plan_days_few: "{n} дня",
    plan_days_many: "{n} дней",
    plan_buy: "Купить",
    plan_renew: "Продлить",
    pay_kicker: "Оплата тарифа",
    pay_period: "Срок: {period}",
    pay_confirm: "Оплатить {amount} BYN",
    pay_cancel: "Отмена",
    pay_currency: "К оплате в белорусских рублях, BYN.",
    pay_test_note: "В тестовом Telegram счёт может показаться в долларах — это витрина теста, цена всё равно в BYN.",
    pay_sent_chat: "Счёт отправили в чат с ботом. Закрой кабинет и оплати там.",
    pay_open_chat: "Открыть чат",
    toast_pay_guest: "Открой кабинет из Telegram",
    toast_pay_start: "Сначала нажми /start в боте",
    toast_pay_open: "Открой Mini App из Telegram",
    toast_pay_fail: "Не удалось открыть оплату",
    toast_pay_cancel: "Оплата отменена",
    toast_pay_paid: "Оплата прошла",
    pricing_note: "Цены в BYN. Оплата тестовая через Telegram.",
    faq_q: "Есть вопрос по FlatStalker?",
    faq_a: "Напиши",
    footer_meta: "РБ · АРЕНДА",
    badge_active: "Активна",
    badge_paused: "Пауза",
    action_pause: "Пауза",
    action_resume: "Возобновить",
    action_copy: "Копировать",
    action_paste: "Вставить из буфера",
    action_delete: "Удалить",
    toast_paused: "На паузе",
    toast_resumed: "Снова активна",
    toast_copied: "Ссылка скопирована",
    toast_pasted: "Вставлено из буфера",
    toast_paste_empty: "В буфере нет ссылки",
    toast_paste_fail: "Не удалось прочитать буфер",
    toast_copy_fail: "Не удалось скопировать",
    toast_deleted: "Ссылка удалена",
    toast_update_fail: "Не удалось обновить",
    toast_slow: "Слишком часто — подожди секунду",
    toast_open_tg: "Открой Mini App из Telegram",
    toast_paste_link: "Вставь ссылку",
    toast_need_kufar: "Нужна ссылка kufar.by",
    toast_need_start: "Сначала нажми /start в боте",
    toast_exists: "Такая ссылка уже есть",
    toast_added: "Ссылка добавлена",
    toast_limit: "Лимит тарифа: {used} из {limit}",
    toast_add_fail: "Не удалось добавить",
    note_kufar_only: "Пока парсим только поиск аренды на Kufar.",
    note_need_start: "Нужен /start в боте, потом снова добавь ссылку.",
    note_duplicate: "Дубликат не сохраняем.",
    note_added: "Добавлено. Управление — в списке ниже.",
    note_limit: "Удали ссылку или смени тариф, чтобы добавить новую.",
    note_api_error: "Ошибка. Проверь, что backend запущен и api= доступен.",
    confirm_delete: "Удалить эту ссылку?",
    ban_stamp: "ДОСТУП ЗАКРЫТ",
    ban_title: "Кабинет на паузе",
    ban_text:
      "Сейчас этот аккаунт не может пользоваться FlatStalker. Если так вышло по ошибке — напиши, спокойно разберёмся.",
    ban_reason_label: "Причина",
    ban_contact: "Связаться",
  },
  by: {
    hello_guest: "Прывітанне — гэта твой кабінет",
    hello_user: "{name}, гэта твой кабінет",
    eyebrow: "Разумны пошук арэнды",
    hero_title: "Новая кватэра — раней за іншых",
    hero_text:
      "FlatStalker сочыць за тваім пошукам і прысылае новыя аб'явы ў Telegram.",
    status_title: "Маніторынг актыўны",
    status_sub: "Праверка па тарыфе",
    add_title: "Дадаць пошук",
    add_hint:
      "Наладзь фільтры на Kufar, скапіюй спасылку пошуку і ўстаў сюды — далей бот сам прышле новыя аб'явы.",
    link_label: "Спасылка",
    link_submit: "Дадаць спасылку",
    link_note: "Толькі спасылка пошуку арэнды на Kufar.",
    my_links_title: "Мае спасылкі",
    links_empty: "Пакуль пуста — дадай спасылку пошуку вышэй.",
    links_open_tg: "Адкрый кабінет з Telegram, каб бачыць спасылкі.",
    links_need_start: "Спачатку націсні /start у боце.",
    links_load_error: "Не ўдалося загрузіць спасылкі.",
    how_title: "Як гэта працуе",
    how_1: "Адкрываеш бота і кабінет",
    how_2: "Наладжваеш пошук на Kufar і ўстаўляеш спасылку",
    how_3: "Атрымліваеш апавяшчэнне са спасылкай на новую аб'яву",
    caps_title: "Магчымасці",
    cap1_title: "Сочым за аб'явамі",
    cap1_text: "Новыя варыянты не трэба шукаць уручную",
    cap2_title: "Хуткія апавяшчэнні",
    cap2_text: "Спасылка прыходзіць адразу, без доўгіх зводак",
    cap3_title: "Твой пошук з Kufar",
    cap3_text: "Горад, цана і пакоі — у спасылцы пошуку",
    cap4_title: "Па ўсёй РБ",
    cap4_text: "Пад пошук арэнды ў Беларусі",
    pricing_title: "Тарыфы",
    plan_badge_current: "Твой",
    plan_current_loading: "Загружаем твой тарыф…",
    plan_current: "Твой тарыф: {plan}. Праверка {interval}.",
    plan_current_guest: "Адкрый кабінет з Telegram, каб убачыць тарыф.",
    plan_current_need_start: "Спачатку націсні /start у боце.",
    plan_current_error: "Не ўдалося загрузіць тарыф.",
    plan_interval_every: "Праверка аб'яў {interval}",
    plan_links_one: "{n} спасылка",
    plan_links_few: "{n} спасылкі",
    plan_links_many: "{n} спасылак",
    interval_seconds: "кожныя {n} сек",
    interval_minutes: "кожныя {n} хв",
    interval_short_seconds: "{n} сек",
    interval_short_minutes: "{n} хв",
    plan_price: "{amount} BYN",
    plan_price_day: "{amount} BYN/дзень",
    plan_price_free: "Бясплатна",
    plan_period_chip: "{n}д",
    plan_period_hint: "Цана за",
    plan_days_one: "{n} дзень",
    plan_days_few: "{n} дні",
    plan_days_many: "{n} дзён",
    plan_buy: "Купіць",
    plan_renew: "Падоўжыць",
    pay_kicker: "Аплата тарыфу",
    pay_period: "Тэрмін: {period}",
    pay_confirm: "Аплаціць {amount} BYN",
    pay_cancel: "Адмена",
    pay_currency: "Да аплаты ў беларускіх рублях, BYN.",
    pay_test_note: "У тэставым Telegram рахунак можа паказацца ў доларах — гэта вітрына тэсту, цана ўсё роўна ў BYN.",
    pay_sent_chat: "Рахунак адправілі ў чат з ботам. Закрый кабінет і аплаці там.",
    pay_open_chat: "Адкрыць чат",
    toast_pay_guest: "Адкрый кабінет з Telegram",
    toast_pay_start: "Спачатку націсні /start у боце",
    toast_pay_open: "Адкрый Mini App з Telegram",
    toast_pay_fail: "Не ўдалося адкрыць аплату",
    toast_pay_cancel: "Аплата адменена",
    toast_pay_paid: "Аплата прайшла",
    pricing_note: "Цэны ў BYN. Аплата тэставая праз Telegram.",
    faq_q: "Ёсць пытанне па FlatStalker?",
    faq_a: "Напішы",
    footer_meta: "РБ · АРЭНДА",
    badge_active: "Актыўная",
    badge_paused: "Паўза",
    action_pause: "Паўза",
    action_resume: "Аднавіць",
    action_copy: "Капіяваць",
    action_paste: "Уставіць з буфера",
    action_delete: "Выдаліць",
    toast_paused: "На паўзе",
    toast_resumed: "Зноў актыўная",
    toast_copied: "Спасылка скапіявана",
    toast_pasted: "Устаўлена з буфера",
    toast_paste_empty: "У буферы няма спасылкі",
    toast_paste_fail: "Не ўдалося прачытаць буфер",
    toast_copy_fail: "Не ўдалося скапіяваць",
    toast_deleted: "Спасылка выдалена",
    toast_update_fail: "Не ўдалося абнавіць",
    toast_slow: "Занадта часта — пачакай секунду",
    toast_open_tg: "Адкрый Mini App з Telegram",
    toast_paste_link: "Устаў спасылку",
    toast_need_kufar: "Патрэбна спасылка kufar.by",
    toast_need_start: "Спачатку націсні /start у боце",
    toast_exists: "Такая спасылка ўжо ёсць",
    toast_added: "Спасылка дададзена",
    toast_limit: "Ліміт тарыфу: {used} з {limit}",
    toast_add_fail: "Не ўдалося дадаць",
    note_kufar_only: "Пакуль парсім толькі пошук арэнды на Kufar.",
    note_need_start: "Патрэбны /start у боце, потым зноў дадай спасылку.",
    note_duplicate: "Дублікат не захоўваем.",
    note_added: "Дададзена. Кіраванне — у спісе ніжэй.",
    note_limit: "Выдалі спасылку або змяні тарыф, каб дадаць новую.",
    note_api_error: "Памылка. Правер, што backend запушчаны і api= даступны.",
    confirm_delete: "Выдаліць гэтую спасылку?",
    ban_stamp: "ДОСТУП ЗАЧЫНЕНЫ",
    ban_title: "Кабінет на паўзе",
    ban_text:
      "Зараз гэты акаўнт не можа карыстацца FlatStalker. Калі так выйшла памылкова — напішы, спакойна разбярэмся.",
    ban_reason_label: "Прычына",
    ban_contact: "Звязацца",
  },
};

function applyTelegramChrome() {
  if (!tg) return;
  if (tg.setHeaderColor) tg.setHeaderColor("#08111f");
  if (tg.setBackgroundColor) tg.setBackgroundColor("#08111f");
}

if (tg) {
  tg.ready();
  tg.expand();
  applyTelegramChrome();
}

const user = tg?.initDataUnsafe?.user;
const linkForm = document.getElementById("link-form");
const linkNote = document.getElementById("link-note");
const linkSubmit = document.getElementById("link-submit");
const linkUrl = document.getElementById("link-url");
const pasteLink = document.getElementById("paste-link");
const toast = document.getElementById("toast");
const linkList = document.getElementById("link-list");
const linksEmpty = document.getElementById("links-empty");
const linksCount = document.getElementById("links-count");
const langButtons = document.querySelectorAll(".lang-btn");

let lang = localStorage.getItem(LANG_KEY) === "by" ? "by" : "ru";
let toastTimer;
let links = [];
let linksEmptyKey = "links_empty";
let me = null;
let meStatusKey = "plan_current_loading";
const cabinetCacheKey = user?.id ? `flatstalker_cabinet_${user.id}` : "";
const banCacheKey = user?.id ? `flatstalker_banned_${user.id}` : "";
const busyLinks = new Set();

function t(key, vars) {
  const dict = I18N[lang] || I18N.ru;
  let value = dict[key] ?? I18N.ru[key] ?? key;
  if (vars) {
    for (const [name, raw] of Object.entries(vars)) {
      value = value.replaceAll(`{${name}}`, String(raw));
    }
  }
  return value;
}

function formatInterval(raw) {
  if (!raw && raw !== 0) return "";
  const totalMs = parseGoDuration(raw);
  if (totalMs <= 0) return String(raw);
  if (totalMs < 60000) {
    return t("interval_seconds", { n: Math.max(1, Math.round(totalMs / 1000)) });
  }
  return t("interval_minutes", { n: Math.max(1, Math.round(totalMs / 60000)) });
}

function formatIntervalShort(raw) {
  if (!raw && raw !== 0) return "—";
  const totalMs = parseGoDuration(raw);
  if (totalMs <= 0) return "—";
  if (totalMs < 60000) {
    return t("interval_short_seconds", { n: Math.max(1, Math.round(totalMs / 1000)) });
  }
  return t("interval_short_minutes", { n: Math.max(1, Math.round(totalMs / 60000)) });
}

function parseGoDuration(raw) {
  if (typeof raw === "number") return raw;
  let totalMs = 0;
  const re = /(\d+(?:\.\d+)?)(ns|us|µs|ms|s|m|h)/g;
  let match;
  while ((match = re.exec(String(raw))) !== null) {
    const amount = Number(match[1]);
    switch (match[2]) {
      case "h":
        totalMs += amount * 3_600_000;
        break;
      case "m":
        totalMs += amount * 60_000;
        break;
      case "s":
        totalMs += amount * 1_000;
        break;
      case "ms":
        totalMs += amount;
        break;
      default:
        break;
    }
  }
  return totalMs;
}

const DEFAULT_PLAN_INTERVALS = {
  free: "5m",
  plus: "2m",
  pro: "30s",
};

const DEFAULT_LINK_LIMITS = {
  free: 1,
  plus: 3,
  pro: 5,
};

const DEFAULT_PERIOD_DAYS = [1, 3, 7, 15, 30, 90, 180];
const DEFAULT_PRICES = {
  currency: "BYN",
  period_days: DEFAULT_PERIOD_DAYS,
  plus: {
    1: "0.70",
    3: "1.60",
    7: "2.90",
    15: "5.50",
    30: "9.90",
    90: "25.20",
    180: "44.50",
  },
  pro: {
    1: "1.40",
    3: "3.20",
    7: "5.80",
    15: "11.00",
    30: "19.90",
    90: "50.40",
    180: "89.00",
  },
};

let selectedPeriodDays = 15;
let paying = false;

function pluralKey(base, n) {
  const abs = Math.abs(Number(n)) % 100;
  const digit = abs % 10;
  if (abs > 10 && abs < 20) return `${base}_many`;
  if (digit === 1) return `${base}_one`;
  if (digit >= 2 && digit <= 4) return `${base}_few`;
  return `${base}_many`;
}

function linkLimitFor(planName) {
  const fromMe = me?.link_limits?.[planName];
  if (Number.isFinite(Number(fromMe))) return Number(fromMe);
  return DEFAULT_LINK_LIMITS[planName] || DEFAULT_LINK_LIMITS.free;
}

function currentLinkLimit() {
  if (Number.isFinite(Number(me?.link_limit))) return Number(me.link_limit);
  return linkLimitFor(me?.plan || "free");
}

function atLinkLimit() {
  return Boolean(me) && links.length >= currentLinkLimit();
}

function priceCatalog() {
  const fromMe = me?.prices;
  if (fromMe?.plus && fromMe?.pro) return fromMe;
  return DEFAULT_PRICES;
}

function periodDays() {
  const days = priceCatalog().period_days;
  if (Array.isArray(days) && days.length) return days.map(Number);
  return DEFAULT_PERIOD_DAYS;
}

function planAmount(planName, days) {
  const catalog = priceCatalog();
  const table = catalog[planName];
  if (!table) return "";
  return table[String(days)] || table[days] || "";
}

function perDayAmount(amount, days) {
  const total = Number(amount);
  if (!days || !Number.isFinite(total)) return "";
  return (total / days).toFixed(2);
}

function renderPeriodChips() {
  const root = document.getElementById("plan-periods");
  if (!root) return;
  const days = periodDays();
  if (!days.includes(selectedPeriodDays)) {
    selectedPeriodDays = days.includes(15) ? 15 : days[0];
  }
  root.replaceChildren(
    ...days.map((n) => {
      const btn = document.createElement("button");
      btn.type = "button";
      btn.className = "plan-period";
      btn.classList.toggle("is-active", n === selectedPeriodDays);
      btn.textContent = t("plan_period_chip", { n });
      btn.addEventListener("click", () => {
        selectedPeriodDays = n;
        renderPlans();
      });
      return btn;
    })
  );
  const value = document.getElementById("plan-period-value");
  if (value) {
    value.textContent = t(pluralKey("plan_days", selectedPeriodDays), {
      n: selectedPeriodDays,
    });
  }
}

function renderPlans() {
  const currentEl = document.getElementById("plan-current");
  const cards = document.querySelectorAll("[data-plan]");

  renderPeriodChips();

  cards.forEach((card) => {
    const name = card.dataset.plan;
    const isCurrent = Boolean(me && me.plan === name);
    card.classList.toggle("is-current", isCurrent);
    const badge = card.querySelector(".plan-badge");
    if (badge) {
      badge.hidden = !isCurrent;
      badge.textContent = t("plan_badge_current");
    }
    const priceEl = card.querySelector("[data-plan-price]");
    const dayEl = card.querySelector("[data-plan-price-day]");
    if (priceEl) {
      if (name === "free") {
        priceEl.textContent = t("plan_price_free");
        if (dayEl) dayEl.textContent = "";
      } else {
        const amount = planAmount(name, selectedPeriodDays);
        priceEl.textContent = amount ? t("plan_price", { amount }) : "";
        const perDay = perDayAmount(amount, selectedPeriodDays);
        if (dayEl) {
          dayEl.textContent = perDay ? t("plan_price_day", { amount: perDay }) : "";
        }
      }
    }
    const intervalEl = card.querySelector("[data-plan-interval]");
    if (intervalEl) {
      const raw = me?.intervals?.[name] || DEFAULT_PLAN_INTERVALS[name];
      intervalEl.textContent = raw
        ? t("plan_interval_every", { interval: formatInterval(raw) })
        : "";
    }
    const linksEl = card.querySelector("[data-plan-links]");
    if (linksEl) {
      const n = linkLimitFor(name);
      linksEl.textContent = t(pluralKey("plan_links", n), { n });
    }
    const buyEl = card.querySelector("[data-plan-buy]");
    if (buyEl) {
      buyEl.hidden = name === "free";
      buyEl.disabled = paying;
      buyEl.textContent = isCurrent ? t("plan_renew") : t("plan_buy");
    }
  });

  const statusEl = document.getElementById("status-interval");

  if (!currentEl) return;
  if (!hasTelegramAuth()) {
    currentEl.textContent = t("plan_current_guest");
    if (statusEl) statusEl.textContent = "—";
    return;
  }
  if (!me) {
    currentEl.textContent = t(meStatusKey);
    if (statusEl) statusEl.textContent = "—";
    return;
  }
  currentEl.textContent = t("plan_current", {
    plan: me.plan_label || String(me.plan || "").toUpperCase(),
    interval: formatInterval(me.interval),
  });
  if (statusEl) {
    statusEl.textContent = formatIntervalShort(me.interval);
  }
}

function applyLanguage() {
  document.documentElement.lang = lang === "by" ? "be" : "ru";

  langButtons.forEach((button) => {
    button.classList.toggle("is-active", button.dataset.lang === lang);
  });

  document.querySelectorAll("[data-i18n]").forEach((el) => {
    const key = el.dataset.i18n;
    if (!key) return;
    if (key === "hello_guest" && user?.first_name) {
      el.textContent = t("hello_user", { name: user.first_name });
      return;
    }
    el.textContent = t(key);
  });

  if (linksEmpty && !links.length) {
    linksEmpty.textContent = t(linksEmptyKey);
  }

  renderLinks();
  renderPlans();
}

function setLanguage(next) {
  if (next !== "ru" && next !== "by") return;
  lang = next;
  localStorage.setItem(LANG_KEY, lang);
  applyLanguage();
}

langButtons.forEach((button) => {
  button.addEventListener("click", () => setLanguage(button.dataset.lang));
});

document.querySelectorAll("[data-plan-buy]").forEach((button) => {
  button.addEventListener("click", () => openPaySheet(button.dataset.planBuy));
});

const payOverlay = document.getElementById("pay-overlay");
const payTitle = document.getElementById("pay-title");
const payPeriod = document.getElementById("pay-period");
const paySum = document.getElementById("pay-sum");
const payDay = document.getElementById("pay-day");
const payStatus = document.getElementById("pay-status");
const payConfirm = document.getElementById("pay-confirm");
const payCancel = document.getElementById("pay-cancel");
const payClose = document.getElementById("pay-close");
let payPlan = "";
let payBotURL = "";

function closePaySheet() {
  if (!payOverlay) return;
  payOverlay.hidden = true;
  document.documentElement.classList.remove("pay-open");
  paying = false;
  payBotURL = "";
  renderPlans();
}

function openPaySheet(planName) {
  if (paying) return;
  if (!hasTelegramAuth()) {
    showToast(t("toast_pay_guest"));
    return;
  }
  if (!me) {
    showToast(t("toast_pay_start"));
    return;
  }
  const amount = planAmount(planName, selectedPeriodDays);
  if (!amount) {
    showToast(t("toast_pay_fail"));
    return;
  }
  payPlan = planName;
  if (payTitle) payTitle.textContent = String(planName || "").toUpperCase();
  if (payPeriod) {
    payPeriod.textContent = t("pay_period", {
      period: t(pluralKey("plan_days", selectedPeriodDays), { n: selectedPeriodDays }),
    });
  }
  if (paySum) paySum.textContent = t("plan_price", { amount });
  if (payDay) {
    const perDay = perDayAmount(amount, selectedPeriodDays);
    payDay.textContent = perDay ? t("plan_price_day", { amount: perDay }) : "";
  }
  if (payStatus) {
    payStatus.hidden = true;
    payStatus.textContent = "";
  }
  if (payConfirm) {
    payConfirm.hidden = false;
    payConfirm.disabled = false;
    payConfirm.textContent = t("pay_confirm", { amount });
  }
  if (payCancel) payCancel.textContent = t("pay_cancel");
  if (payOverlay) payOverlay.hidden = false;
  document.documentElement.classList.add("pay-open");
}

async function buyPlan() {
  if (paying || !payPlan) return;
  paying = true;
  if (payConfirm) payConfirm.disabled = true;
  renderPlans();
  try {
    const res = await apiFetch("/api/pay", {
      method: "POST",
      body: JSON.stringify({ plan: payPlan, days: selectedPeriodDays }),
    });
    if (!res.ok) {
      showToast(t("toast_pay_fail"));
      paying = false;
      if (payConfirm) payConfirm.disabled = false;
      renderPlans();
      return;
    }
    const body = await res.json().catch(() => ({}));
    payBotURL = body.bot_url || "";
    showChatInvoice();
  } catch {
    showToast(t("toast_pay_fail"));
    paying = false;
    if (payConfirm) payConfirm.disabled = false;
    renderPlans();
  }
}

function showChatInvoice() {
  paying = false;
  renderPlans();
  if (payConfirm) {
    payConfirm.hidden = !payBotURL;
    payConfirm.disabled = false;
    payConfirm.textContent = t("pay_open_chat");
  }
  if (payStatus) {
    payStatus.hidden = false;
    payStatus.textContent = t("pay_sent_chat");
  }
}

function openPayChat() {
  if (payBotURL) openTelegramUrl(payBotURL);
  try {
    tg?.close?.();
  } catch {
    /* desktop webview */
  }
}

payConfirm?.addEventListener("click", () => {
  if (payStatus && !payStatus.hidden && payBotURL) {
    openPayChat();
    return;
  }
  buyPlan();
});
payCancel?.addEventListener("click", closePaySheet);
payClose?.addEventListener("click", closePaySheet);
payOverlay?.addEventListener("click", (event) => {
  if (event.target === payOverlay) closePaySheet();
});
document.addEventListener("keydown", (event) => {
  if (event.key === "Escape" && payOverlay && !payOverlay.hidden) {
    closePaySheet();
  }
});

function showToast(message) {
  if (!toast) return;
  toast.textContent = message;
  toast.hidden = false;
  requestAnimationFrame(() => toast.classList.add("show"));
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => {
    toast.classList.remove("show");
    setTimeout(() => {
      toast.hidden = true;
    }, 220);
  }, 2200);
}

function initData() {
  return tg?.initData || "";
}

function hasTelegramAuth() {
  return Boolean(initData());
}

function authHeaders(extra) {
  const headers = { ...(extra || {}) };
  const data = initData();
  if (data) {
    headers.Authorization = `tma ${data}`;
  }
  return headers;
}

async function apiFetch(path, options = {}) {
  const headers = authHeaders(options.headers);
  if (options.body && !headers["Content-Type"]) {
    headers["Content-Type"] = "application/json";
  }
  const res = await fetch(`${API_BASE}${path}`, { ...options, headers });
  if (res.status === 403) {
    const body = await res.clone().json().catch(() => ({}));
    if (body.code === "banned") {
      showBanScreen(body.support, body.reason);
    }
  }
  return res;
}

function supportHandle(raw) {
  const handle = String(raw || "@bazan_ivan").trim();
  return handle.startsWith("@") ? handle : `@${handle}`;
}

function supportUrl(handle) {
  return `https://t.me/${supportHandle(handle).replace(/^@/, "")}`;
}

function openTelegramUrl(url) {
  if (!url) return;
  try {
    if (typeof tg?.openTelegramLink === "function") {
      tg.openTelegramLink(url);
      return;
    }
  } catch {
    /* desktop webview */
  }
  try {
    if (typeof tg?.openLink === "function") {
      tg.openLink(url, { try_instant_view: false });
      return;
    }
  } catch {
    /* ignore */
  }
  window.open(url, "_blank", "noopener,noreferrer");
}

function openSupportChat(handle) {
  openTelegramUrl(supportUrl(handle));
}

function showBanScreen(support, reason) {
  const screen = document.getElementById("ban-screen");
  const contact = document.getElementById("ban-contact");
  const handleEl = document.getElementById("ban-handle");
  const reasonBlock = document.getElementById("ban-reason-block");
  const reasonEl = document.getElementById("ban-reason");
  const handle = supportHandle(support);
  const reasonText = normalizeBanReason(reason);
  writeBanCache(handle, reasonText);
  clearCabinetCache();
  document.documentElement.classList.remove("cabinet-ready");
  document.documentElement.classList.add("is-banned");
  document.body.classList.add("is-banned");
  if (handleEl) handleEl.textContent = handle;
  if (contact) contact.dataset.tg = handle.replace(/^@/, "");
  if (reasonEl) reasonEl.textContent = reasonText;
  if (reasonBlock) reasonBlock.hidden = !reasonText;
  screen?.querySelectorAll("[data-i18n]").forEach((el) => {
    const key = el.dataset.i18n;
    if (key) el.textContent = t(key);
  });
}

function hideBanScreen() {
  document.documentElement.classList.remove("is-banned");
  document.body.classList.remove("is-banned");
}

function bootFromCache() {
  const banned = readBanCache();
  if (banned) showBanScreen(banned.support, banned.reason);
}

function revealCabinet() {
  if (document.documentElement.classList.contains("is-banned")) return;
  document.documentElement.classList.add("cabinet-ready");
}

function shortURL(url) {
  try {
    const parsed = new URL(url);
    const path = `${parsed.pathname}${parsed.search}`;
    return path.length > 54 ? `${path.slice(0, 51)}…` : path;
  } catch {
    return url.length > 54 ? `${url.slice(0, 51)}…` : url;
  }
}

function renderLinks() {
  if (!linkList || !linksEmpty) return;

  if (linksCount) {
    if (me) {
      linksCount.textContent = `${links.length} / ${currentLinkLimit()}`;
    } else {
      linksCount.textContent = links.length ? String(links.length) : "";
    }
  }
  if (linkSubmit) {
    linkSubmit.disabled = atLinkLimit();
  }

  if (links.length === 0) {
    linkList.hidden = true;
    linkList.innerHTML = "";
    linksEmpty.hidden = false;
    linksEmpty.textContent = t(linksEmptyKey);
    return;
  }

  linksEmpty.hidden = true;
  linkList.hidden = false;
  linkList.innerHTML = links
    .map((link) => {
      const paused = Boolean(link.paused);
      return `
        <li class="link-item${paused ? " is-paused" : ""}" data-id="${link.id}">
          <a class="link-item-url" href="${escapeAttr(link.url)}" target="_blank" rel="noopener">
            ${escapeHTML(shortURL(link.url))}
          </a>
          <div class="link-item-meta">
            <span class="link-badge">${paused ? t("badge_paused") : t("badge_active")}</span>
            <div class="link-actions">
              <button type="button" class="link-action" data-action="copy" data-id="${link.id}">
                ${t("action_copy")}
              </button>
              <button type="button" class="link-action" data-action="toggle" data-id="${link.id}">
                ${paused ? t("action_resume") : t("action_pause")}
              </button>
              <button type="button" class="link-action is-danger" data-action="delete" data-id="${link.id}">
                ${t("action_delete")}
              </button>
            </div>
          </div>
        </li>
      `;
    })
    .join("");

  busyLinks.forEach((busyId) => {
    linkList.querySelectorAll(`button[data-id="${busyId}"]`).forEach((el) => {
      el.disabled = true;
    });
  });
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}

function escapeAttr(value) {
  return escapeHTML(value).replaceAll("'", "&#39;");
}

function readCabinetCache() {
  if (!cabinetCacheKey) return null;
  try {
    const raw = localStorage.getItem(cabinetCacheKey);
    if (!raw) return null;
    const data = JSON.parse(raw);
    if (!data || typeof data !== "object" || !Array.isArray(data.links)) return null;
    return data;
  } catch {
    return null;
  }
}

function writeCabinetCache() {
  if (!cabinetCacheKey || !me) return;
  try {
    localStorage.setItem(
      cabinetCacheKey,
      JSON.stringify({
        ...me,
        links,
      })
    );
  } catch {
    /* quota / private mode */
  }
}

function clearCabinetCache() {
  if (!cabinetCacheKey) return;
  try {
    localStorage.removeItem(cabinetCacheKey);
  } catch {
    /* ignore */
  }
}

function readBanCache() {
  if (!banCacheKey) return null;
  try {
    const raw = localStorage.getItem(banCacheKey);
    if (!raw) return null;
    const data = JSON.parse(raw);
    if (!data || typeof data !== "object") return { support: "@bazan_ivan" };
    return data;
  } catch {
    return { support: "@bazan_ivan" };
  }
}

function normalizeBanReason(raw) {
  return String(raw || "")
    .replace(/\s+/g, " ")
    .trim()
    .slice(0, 500);
}

function writeBanCache(support, reason) {
  if (!banCacheKey) return;
  try {
    localStorage.setItem(
      banCacheKey,
      JSON.stringify({
        support: supportHandle(support),
        reason: normalizeBanReason(reason),
      })
    );
  } catch {
    /* ignore */
  }
}

function clearBanCache() {
  if (!banCacheKey) return;
  try {
    localStorage.removeItem(banCacheKey);
  } catch {
    /* ignore */
  }
}

function applyCabinet(body) {
  me = body;
  meStatusKey = "plan_current_loading";
  links = Array.isArray(body.links) ? body.links : [];
  linksEmptyKey = "links_empty";
  renderPlans();
  renderLinks();
}

async function loadCabinet() {
  if (!hasTelegramAuth()) {
    me = null;
    meStatusKey = "plan_current_guest";
    links = [];
    linksEmptyKey = "links_open_tg";
    renderPlans();
    renderLinks();
    revealCabinet();
    return;
  }

  const banned = readBanCache();
  if (banned) {
    showBanScreen(banned.support, banned.reason);
  } else {
    const cached = readCabinetCache();
    if (cached) {
      applyCabinet(cached);
      revealCabinet();
    }
  }

  try {
    const res = await apiFetch("/api/me");
    const body = await res.json().catch(() => ({}));
    if (res.status === 403 && body.code === "banned") {
      return;
    }
    if (res.status === 401) {
      clearBanCache();
      hideBanScreen();
      clearCabinetCache();
      me = null;
      meStatusKey = "plan_current_guest";
      links = [];
      linksEmptyKey = "links_open_tg";
      renderPlans();
      renderLinks();
      revealCabinet();
      return;
    }
    if (res.status === 404) {
      clearBanCache();
      hideBanScreen();
      clearCabinetCache();
      me = null;
      meStatusKey = "plan_current_need_start";
      links = [];
      linksEmptyKey = "links_need_start";
      renderPlans();
      renderLinks();
      revealCabinet();
      return;
    }
    if (!res.ok) {
      throw new Error(body.error || `HTTP ${res.status}`);
    }
    clearBanCache();
    hideBanScreen();
    applyCabinet(body);
    writeCabinetCache();
    revealCabinet();
  } catch (err) {
    console.error(err);
    if (readBanCache()) return;
    if (readCabinetCache()) {
      revealCabinet();
      return;
    }
    me = null;
    meStatusKey = "plan_current_error";
    linksEmptyKey = "links_load_error";
    renderPlans();
    if (linksEmpty) {
      linksEmpty.textContent = t(linksEmptyKey);
      linksEmpty.hidden = false;
    }
    revealCabinet();
  }
}

async function setPaused(id, paused) {
  if (!hasTelegramAuth()) return;

  const res = await apiFetch(`/api/links/${id}`, {
    method: "PATCH",
    body: JSON.stringify({ paused }),
  });
  const body = await res.json().catch(() => ({}));
  if (res.status === 429) {
    const err = new Error("slow");
    err.code = "slow";
    throw err;
  }
  if (!res.ok) {
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return body;
}

async function deleteLink(id) {
  if (!hasTelegramAuth()) return;

  const res = await apiFetch(`/api/links/${id}`, {
    method: "DELETE",
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return body;
}

async function copyText(text) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text);
    return;
  }
  const area = document.createElement("textarea");
  area.value = text;
  area.setAttribute("readonly", "");
  area.style.position = "fixed";
  area.style.opacity = "0";
  document.body.appendChild(area);
  area.select();
  const ok = document.execCommand("copy");
  area.remove();
  if (!ok) throw new Error("copy failed");
}

async function readClipboardText() {
  if (navigator.clipboard?.readText) {
    return String(await navigator.clipboard.readText()).trim();
  }
  throw new Error("clipboard read unsupported");
}

function syncPasteGlow() {
  if (!pasteLink || !linkUrl) return;
  const empty = !String(linkUrl.value || "").trim();
  pasteLink.classList.toggle("is-glow", empty);
  pasteLink.classList.toggle("is-done", !empty);
}

pasteLink?.addEventListener("click", async () => {
  pasteLink.disabled = true;
  try {
    const text = await readClipboardText();
    if (!text) {
      showToast(t("toast_paste_empty"));
      return;
    }
    if (linkUrl) {
      linkUrl.value = text;
      linkUrl.focus();
      linkUrl.setSelectionRange(text.length, text.length);
    }
    syncPasteGlow();
    if (tg?.HapticFeedback?.notificationOccurred) {
      tg.HapticFeedback.notificationOccurred("success");
    }
    showToast(t("toast_pasted"));
  } catch (err) {
    console.error(err);
    showToast(t("toast_paste_fail"));
  } finally {
    pasteLink.disabled = false;
  }
});

linkUrl?.addEventListener("input", syncPasteGlow);
linkUrl?.addEventListener("change", syncPasteGlow);
syncPasteGlow();

linkList?.addEventListener("click", async (event) => {
  const button = event.target.closest("[data-action]");
  if (!button) return;

  const id = Number(button.dataset.id);
  const action = button.dataset.action;
  if (!id) return;

  const link = links.find((item) => Number(item.id) === id);
  if (!link) return;

  if ((action === "toggle" || action === "delete") && busyLinks.has(id)) {
    return;
  }

  button.disabled = true;
  if (action === "toggle" || action === "delete") {
    busyLinks.add(id);
  }
  try {
    if (action === "copy") {
      await copyText(link.url);
      if (tg?.HapticFeedback?.notificationOccurred) {
        tg.HapticFeedback.notificationOccurred("success");
      }
      showToast(t("toast_copied"));
      return;
    }

    if (action === "toggle") {
      const nextPaused = !link.paused;
      await setPaused(id, nextPaused);
      link.paused = nextPaused;
      renderLinks();
      writeCabinetCache();
      showToast(nextPaused ? t("toast_paused") : t("toast_resumed"));
      return;
    }

    if (action === "delete") {
      const confirmed =
        typeof tg?.showConfirm === "function"
          ? await new Promise((resolve) => {
              tg.showConfirm(t("confirm_delete"), resolve);
            })
          : window.confirm(t("confirm_delete"));
      if (!confirmed) return;

      await deleteLink(id);
      links = links.filter((item) => Number(item.id) !== id);
      renderLinks();
      writeCabinetCache();
      showToast(t("toast_deleted"));
    }
  } catch (err) {
    console.error(err);
    if (err?.code === "slow") {
      showToast(t("toast_slow"));
    } else {
      showToast(action === "copy" ? t("toast_copy_fail") : t("toast_update_fail"));
    }
  } finally {
    if (action === "toggle") {
      setTimeout(() => {
        busyLinks.delete(id);
        renderLinks();
      }, 1000);
    } else if (action === "delete") {
      busyLinks.delete(id);
    } else {
      button.disabled = false;
    }
  }
});

linkForm?.addEventListener("submit", async (event) => {
  event.preventDefault();

  if (!hasTelegramAuth()) {
    showToast(t("toast_open_tg"));
    return;
  }

  const url = String(new FormData(linkForm).get("url") || "").trim();
  if (!url) {
    showToast(t("toast_paste_link"));
    return;
  }
  if (!/kufar\.by/i.test(url)) {
    showToast(t("toast_need_kufar"));
    if (linkNote) linkNote.textContent = t("note_kufar_only");
    return;
  }
  if (atLinkLimit()) {
    showToast(t("toast_limit", { used: links.length, limit: currentLinkLimit() }));
    if (linkNote) linkNote.textContent = t("note_limit");
    return;
  }

  if (linkSubmit) linkSubmit.disabled = true;

  try {
    const res = await apiFetch("/api/links", {
      method: "POST",
      body: JSON.stringify({ url }),
    });
    const body = await res.json().catch(() => ({}));
    if (!res.ok) {
      if (res.status === 401) {
        showToast(t("toast_open_tg"));
        return;
      }
      if (res.status === 404) {
        showToast(t("toast_need_start"));
        if (linkNote) linkNote.textContent = t("note_need_start");
        return;
      }
      if (res.status === 409) {
        const limit = body.limit ?? currentLinkLimit();
        const used = body.used ?? links.length;
        showToast(t("toast_limit", { used, limit }));
        if (linkNote) linkNote.textContent = t("note_limit");
        return;
      }
      if (res.status === 400) {
        const msg = body.error || t("toast_need_kufar");
        showToast(msg);
        if (linkNote) linkNote.textContent = msg;
        return;
      }
      throw new Error(body.error || `HTTP ${res.status}`);
    }

    if (body.created === false) {
      showToast(t("toast_exists"));
      if (linkNote) linkNote.textContent = t("note_duplicate");
    } else {
      showToast(t("toast_added"));
      if (linkNote) linkNote.textContent = t("note_added");
      linkForm.reset();
      syncPasteGlow();
      links = [
        ...links,
        { id: body.id, url: body.url, paused: Boolean(body.paused) },
      ];
      linksEmptyKey = "links_empty";
      renderLinks();
      writeCabinetCache();
    }
    linkNote?.classList.add("is-flash");
    setTimeout(() => linkNote?.classList.remove("is-flash"), 700);
  } catch (err) {
    console.error(err);
    showToast(t("toast_add_fail"));
    if (linkNote) linkNote.textContent = t("note_api_error");
  } finally {
    if (linkSubmit) linkSubmit.disabled = atLinkLimit();
  }
});

document.getElementById("ban-contact")?.addEventListener("click", () => {
  const handle =
    document.getElementById("ban-handle")?.textContent ||
    document.getElementById("ban-contact")?.dataset.tg;
  openSupportChat(handle);
});

document.addEventListener("click", (event) => {
  const link = event.target.closest?.('a[href^="https://t.me/"]');
  if (!link) return;
  if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
  event.preventDefault();
  openTelegramUrl(link.href);
});

bootFromCache();
applyLanguage();
loadCabinet();

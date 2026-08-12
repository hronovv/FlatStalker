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
    pricing_note: "Цены и оплату добавим позже.",
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
    pricing_note: "Цэны і аплату дададзім пазней.",
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

function renderPlans() {
  const currentEl = document.getElementById("plan-current");
  const cards = document.querySelectorAll("[data-plan]");

  cards.forEach((card) => {
    const name = card.dataset.plan;
    const isCurrent = Boolean(me && me.plan === name);
    card.classList.toggle("is-current", isCurrent);
    const badge = card.querySelector(".plan-badge");
    if (badge) {
      badge.hidden = !isCurrent;
      badge.textContent = t("plan_badge_current");
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
  return fetch(`${API_BASE}${path}`, { ...options, headers });
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
    return;
  }

  const cached = readCabinetCache();
  if (cached) {
    applyCabinet(cached);
  } else {
    meStatusKey = "plan_current_loading";
    const currentEl = document.getElementById("plan-current");
    if (currentEl) currentEl.textContent = t(meStatusKey);
  }

  try {
    const res = await apiFetch("/api/me");
    const body = await res.json().catch(() => ({}));
    if (res.status === 401) {
      clearCabinetCache();
      me = null;
      meStatusKey = "plan_current_guest";
      links = [];
      linksEmptyKey = "links_open_tg";
      renderPlans();
      renderLinks();
      return;
    }
    if (res.status === 404) {
      clearCabinetCache();
      me = null;
      meStatusKey = "plan_current_need_start";
      links = [];
      linksEmptyKey = "links_need_start";
      renderPlans();
      renderLinks();
      return;
    }
    if (!res.ok) {
      throw new Error(body.error || `HTTP ${res.status}`);
    }
    applyCabinet(body);
    writeCabinetCache();
  } catch (err) {
    console.error(err);
    if (cached) return;
    me = null;
    meStatusKey = "plan_current_error";
    linksEmptyKey = "links_load_error";
    renderPlans();
    if (linksEmpty) {
      linksEmpty.textContent = t(linksEmptyKey);
      linksEmpty.hidden = false;
    }
  }
}

async function setPaused(id, paused) {
  if (!hasTelegramAuth()) return;

  const res = await apiFetch(`/api/links/${id}`, {
    method: "PATCH",
    body: JSON.stringify({ paused }),
  });
  const body = await res.json().catch(() => ({}));
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

  button.disabled = true;
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
    showToast(action === "copy" ? t("toast_copy_fail") : t("toast_update_fail"));
  } finally {
    button.disabled = false;
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

applyLanguage();
loadCabinet();

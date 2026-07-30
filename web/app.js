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
      "FlatStalker следит за объявлениями круглосуточно и сразу присылает подходящие варианты в Telegram.",
    status_title: "Мониторинг активен",
    status_sub: "Обновляем объявления",
    filters_title: "Фильтры",
    filters_hint:
      "Настроишь один раз — дальше бот будет присылать подходящие объявления по РБ.",
    city_label: "Город",
    city_placeholder: "Например, Минск",
    price_from: "Цена от",
    price_to: "Цена до",
    rooms_label: "Комнаты",
    rooms_aria: "Количество комнат",
    rooms_any: "Любое",
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
    how_2: "Задаёшь город, бюджет и комнаты",
    how_3: "Получаешь уведомление со ссылкой на новое объявление",
    caps_title: "Возможности",
    cap1_title: "Следим за объявлениями",
    cap1_text: "Новые варианты не нужно искать вручную",
    cap2_title: "Быстрые уведомления",
    cap2_text: "Ссылка приходит сразу, без долгих сводок",
    cap3_title: "Твои условия",
    cap3_text: "Город, цена, комнаты и ссылка на поиск",
    cap4_title: "По всей РБ",
    cap4_text: "Под поиск аренды в Беларуси",
    pricing_title: "Тарифы",
    plan_free: "Базовая скорость проверки объявлений",
    plan_plus: "Быстрее находит новые объявления",
    plan_pro: "Максимальная скорость проверки",
    plan_badge_current: "Твой",
    plan_current_loading: "Загружаем твой тариф…",
    plan_current: "Твой тариф: {plan}. Проверка {interval}.",
    plan_current_guest: "Открой кабинет из Telegram, чтобы увидеть тариф.",
    plan_current_need_start: "Сначала нажми /start в боте.",
    plan_current_error: "Не удалось загрузить тариф.",
    plan_interval_every: "Проверка {interval}",
    interval_seconds: "каждые {n} сек",
    interval_minutes: "каждые {n} мин",
    pricing_note: "Цены и оплату добавим позже.",
    faq_q: "Есть вопрос по FlatStalker?",
    faq_a: "Напиши",
    footer_meta: "РБ · АРЕНДА",
    badge_active: "Активна",
    badge_paused: "Пауза",
    action_pause: "Пауза",
    action_resume: "Возобновить",
    action_delete: "Удалить",
    toast_paused: "На паузе",
    toast_resumed: "Снова активна",
    toast_deleted: "Ссылка удалена",
    toast_update_fail: "Не удалось обновить",
    toast_open_tg: "Открой Mini App из Telegram",
    toast_paste_link: "Вставь ссылку",
    toast_need_kufar: "Нужна ссылка kufar.by",
    toast_need_start: "Сначала нажми /start в боте",
    toast_exists: "Такая ссылка уже есть",
    toast_added: "Ссылка добавлена",
    toast_add_fail: "Не удалось добавить",
    note_kufar_only: "Пока парсим только поиск аренды на Kufar.",
    note_need_start: "Нужен /start в боте, потом снова добавь ссылку.",
    note_duplicate: "Дубликат не сохраняем.",
    note_added: "Добавлено. Управление — в списке ниже.",
    note_api_error: "Ошибка. Проверь, что backend запущен и api= доступен.",
    confirm_delete: "Удалить эту ссылку?",
  },
  by: {
    hello_guest: "Прывітанне — гэта твой кабінет",
    hello_user: "{name}, гэта твой кабінет",
    eyebrow: "Разумны пошук арэнды",
    hero_title: "Новая кватэра — раней за іншых",
    hero_text:
      "FlatStalker сочыць за аб'явамі кругласутачна і адразу прысылае падыходныя варыянты ў Telegram.",
    status_title: "Маніторынг актыўны",
    status_sub: "Абнаўляем аб'явы",
    filters_title: "Фільтры",
    filters_hint:
      "Наладзіш адзін раз — далей бот будзе прысылаць падыходныя аб'явы па РБ.",
    city_label: "Горад",
    city_placeholder: "Напрыклад, Мінск",
    price_from: "Цана ад",
    price_to: "Цана да",
    rooms_label: "Пакоі",
    rooms_aria: "Колькасць пакояў",
    rooms_any: "Любое",
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
    how_2: "Задаеш горад, бюджэт і пакоі",
    how_3: "Атрымліваеш апавяшчэнне са спасылкай на новую аб'яву",
    caps_title: "Магчымасці",
    cap1_title: "Сочым за аб'явамі",
    cap1_text: "Новыя варыянты не трэба шукаць уручную",
    cap2_title: "Хуткія апавяшчэнні",
    cap2_text: "Спасылка прыходзіць адразу, без доўгіх зводак",
    cap3_title: "Твае ўмовы",
    cap3_text: "Горад, цана, пакоі і спасылка на пошук",
    cap4_title: "Па ўсёй РБ",
    cap4_text: "Пад пошук арэнды ў Беларусі",
    pricing_title: "Тарыфы",
    plan_free: "Базавая хуткасць праверкі аб'яў",
    plan_plus: "Хутчэй знаходзіць новыя аб'явы",
    plan_pro: "Максімальная хуткасць праверкі",
    plan_badge_current: "Твой",
    plan_current_loading: "Загружаем твой тарыф…",
    plan_current: "Твой тарыф: {plan}. Праверка {interval}.",
    plan_current_guest: "Адкрый кабінет з Telegram, каб убачыць тарыф.",
    plan_current_need_start: "Спачатку націсні /start у боце.",
    plan_current_error: "Не ўдалося загрузіць тарыф.",
    plan_interval_every: "Праверка {interval}",
    interval_seconds: "кожныя {n} сек",
    interval_minutes: "кожныя {n} хв",
    pricing_note: "Цэны і аплату дададзім пазней.",
    faq_q: "Ёсць пытанне па FlatStalker?",
    faq_a: "Напішы",
    footer_meta: "РБ · АРЭНДА",
    badge_active: "Актыўная",
    badge_paused: "Паўза",
    action_pause: "Паўза",
    action_resume: "Аднавіць",
    action_delete: "Выдаліць",
    toast_paused: "На паўзе",
    toast_resumed: "Зноў актыўная",
    toast_deleted: "Спасылка выдалена",
    toast_update_fail: "Не ўдалося абнавіць",
    toast_open_tg: "Адкрый Mini App з Telegram",
    toast_paste_link: "Устаў спасылку",
    toast_need_kufar: "Патрэбна спасылка kufar.by",
    toast_need_start: "Спачатку націсні /start у боце",
    toast_exists: "Такая спасылка ўжо ёсць",
    toast_added: "Спасылка дададзена",
    toast_add_fail: "Не ўдалося дадаць",
    note_kufar_only: "Пакуль парсім толькі пошук арэнды на Kufar.",
    note_need_start: "Патрэбны /start у боце, потым зноў дадай спасылку.",
    note_duplicate: "Дублікат не захоўваем.",
    note_added: "Дададзена. Кіраванне — у спісе ніжэй.",
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
const hello = document.getElementById("hello");
const linkForm = document.getElementById("link-form");
const linkNote = document.getElementById("link-note");
const linkSubmit = document.getElementById("link-submit");
const toast = document.getElementById("toast");
const roomsInput = document.getElementById("rooms-input");
const roomButtons = document.querySelectorAll(".seg-btn");
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
      const raw = me?.intervals?.[name];
      intervalEl.textContent = raw
        ? t("plan_interval_every", { interval: formatInterval(raw) })
        : "";
    }
  });

  if (!currentEl) return;
  if (!chatID()) {
    currentEl.textContent = t("plan_current_guest");
    return;
  }
  if (!me) {
    currentEl.textContent = t(meStatusKey);
    return;
  }
  currentEl.textContent = t("plan_current", {
    plan: me.plan_label || String(me.plan || "").toUpperCase(),
    interval: formatInterval(me.interval),
  });
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

  document.querySelectorAll("[data-i18n-placeholder]").forEach((el) => {
    const key = el.dataset.i18nPlaceholder;
    if (key) el.setAttribute("placeholder", t(key));
  });

  document.querySelectorAll("[data-i18n-aria]").forEach((el) => {
    const key = el.dataset.i18nAria;
    if (key) el.setAttribute("aria-label", t(key));
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

roomButtons.forEach((button) => {
  button.addEventListener("click", () => {
    roomButtons.forEach((item) => item.classList.remove("is-active"));
    button.classList.add("is-active");
    if (roomsInput) roomsInput.value = button.dataset.rooms || "";
  });
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

function chatID() {
  return user?.id ?? null;
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
    linksCount.textContent = links.length ? String(links.length) : "";
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

async function loadMe() {
  const chatId = chatID();
  const currentEl = document.getElementById("plan-current");
  if (!chatId) {
    me = null;
    meStatusKey = "plan_current_guest";
    renderPlans();
    return;
  }

  meStatusKey = "plan_current_loading";
  if (currentEl) currentEl.textContent = t(meStatusKey);

  try {
    const res = await fetch(`${API_BASE}/api/me?chat_id=${encodeURIComponent(chatId)}`);
    const body = await res.json().catch(() => ({}));
    if (res.status === 404) {
      me = null;
      meStatusKey = "plan_current_need_start";
      renderPlans();
      return;
    }
    if (!res.ok) {
      throw new Error(body.error || `HTTP ${res.status}`);
    }
    me = body;
    meStatusKey = "plan_current_loading";
    renderPlans();
  } catch (err) {
    console.error(err);
    me = null;
    meStatusKey = "plan_current_error";
    renderPlans();
  }
}

async function loadLinks() {
  const chatId = chatID();
  if (!chatId) {
    linksEmptyKey = "links_open_tg";
    if (linksEmpty) {
      linksEmpty.textContent = t(linksEmptyKey);
      linksEmpty.hidden = false;
    }
    return;
  }

  try {
    const res = await fetch(`${API_BASE}/api/links?chat_id=${encodeURIComponent(chatId)}`);
    const body = await res.json().catch(() => ({}));
    if (res.status === 404) {
      links = [];
      linksEmptyKey = "links_need_start";
      renderLinks();
      return;
    }
    if (!res.ok) {
      throw new Error(body.error || `HTTP ${res.status}`);
    }
    links = Array.isArray(body.links) ? body.links : [];
    linksEmptyKey = "links_empty";
    renderLinks();
  } catch (err) {
    console.error(err);
    linksEmptyKey = "links_load_error";
    if (linksEmpty) {
      linksEmpty.textContent = t(linksEmptyKey);
      linksEmpty.hidden = false;
    }
  }
}

async function setPaused(id, paused) {
  const chatId = chatID();
  if (!chatId) return;

  const res = await fetch(`${API_BASE}/api/links/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ chat_id: chatId, paused }),
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return body;
}

async function deleteLink(id) {
  const chatId = chatID();
  if (!chatId) return;

  const res = await fetch(`${API_BASE}/api/links/${id}`, {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ chat_id: chatId }),
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return body;
}

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
    if (action === "toggle") {
      const nextPaused = !link.paused;
      await setPaused(id, nextPaused);
      link.paused = nextPaused;
      renderLinks();
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
      showToast(t("toast_deleted"));
    }
  } catch (err) {
    console.error(err);
    showToast(t("toast_update_fail"));
  } finally {
    button.disabled = false;
  }
});

linkForm?.addEventListener("submit", async (event) => {
  event.preventDefault();

  const chatId = chatID();
  if (!chatId) {
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

  if (linkSubmit) linkSubmit.disabled = true;

  try {
    const res = await fetch(`${API_BASE}/api/links`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ chat_id: chatId, url }),
    });
    const body = await res.json().catch(() => ({}));
    if (!res.ok) {
      if (res.status === 404) {
        showToast(t("toast_need_start"));
        if (linkNote) linkNote.textContent = t("note_need_start");
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
      await loadLinks();
    }
    linkNote?.classList.add("is-flash");
    setTimeout(() => linkNote?.classList.remove("is-flash"), 700);
  } catch (err) {
    console.error(err);
    showToast(t("toast_add_fail"));
    if (linkNote) linkNote.textContent = t("note_api_error");
  } finally {
    if (linkSubmit) linkSubmit.disabled = false;
  }
});

applyLanguage();
loadMe();
loadLinks();

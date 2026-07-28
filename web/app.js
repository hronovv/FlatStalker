const tg = window.Telegram?.WebApp;

const params = new URLSearchParams(location.search);
const API_BASE =
  params.get("api") ||
  window.FLATSTALKER_API ||
  "http://127.0.0.1:8080";

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

if (user?.first_name) {
  hello.textContent = `${user.first_name}, это твой кабинет`;
}

roomButtons.forEach((button) => {
  button.addEventListener("click", () => {
    roomButtons.forEach((item) => item.classList.remove("is-active"));
    button.classList.add("is-active");
    if (roomsInput) roomsInput.value = button.dataset.rooms || "";
  });
});

let toastTimer;
let links = [];

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
            <span class="link-badge">${paused ? "Пауза" : "Активна"}</span>
            <div class="link-actions">
              <button type="button" class="link-action" data-action="toggle" data-id="${link.id}">
                ${paused ? "Возобновить" : "Пауза"}
              </button>
              <button type="button" class="link-action is-danger" data-action="delete" data-id="${link.id}">
                Удалить
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

async function loadLinks() {
  const chatId = chatID();
  if (!chatId) {
    if (linksEmpty) {
      linksEmpty.textContent = "Открой кабинет из Telegram, чтобы видеть ссылки.";
      linksEmpty.hidden = false;
    }
    return;
  }

  try {
    const res = await fetch(`${API_BASE}/api/links?chat_id=${encodeURIComponent(chatId)}`);
    const body = await res.json().catch(() => ({}));
    if (res.status === 404) {
      links = [];
      renderLinks();
      if (linksEmpty) {
        linksEmpty.textContent = "Сначала нажми /start в боте.";
        linksEmpty.hidden = false;
      }
      return;
    }
    if (!res.ok) {
      throw new Error(body.error || `HTTP ${res.status}`);
    }
    links = Array.isArray(body.links) ? body.links : [];
    renderLinks();
  } catch (err) {
    console.error(err);
    if (linksEmpty) {
      linksEmpty.textContent = "Не удалось загрузить ссылки.";
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
      showToast(nextPaused ? "На паузе" : "Снова активна");
      return;
    }

    if (action === "delete") {
      const confirmed =
        typeof tg?.showConfirm === "function"
          ? await new Promise((resolve) => {
              tg.showConfirm("Удалить эту ссылку?", resolve);
            })
          : window.confirm("Удалить эту ссылку?");
      if (!confirmed) return;

      await deleteLink(id);
      links = links.filter((item) => Number(item.id) !== id);
      renderLinks();
      showToast("Ссылка удалена");
    }
  } catch (err) {
    console.error(err);
    showToast("Не удалось обновить");
  } finally {
    button.disabled = false;
  }
});

linkForm?.addEventListener("submit", async (event) => {
  event.preventDefault();

  const chatId = chatID();
  if (!chatId) {
    showToast("Открой Mini App из Telegram");
    return;
  }

  const url = String(new FormData(linkForm).get("url") || "").trim();
  if (!url) {
    showToast("Вставь ссылку");
    return;
  }
  if (!/kufar\.by/i.test(url)) {
    showToast("Нужна ссылка kufar.by");
    if (linkNote) {
      linkNote.textContent = "Пока парсим только поиск аренды на Kufar.";
    }
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
        showToast("Сначала нажми /start в боте");
        if (linkNote) linkNote.textContent = "Нужен /start в боте, потом снова добавь ссылку.";
        return;
      }
      if (res.status === 400) {
        const msg = body.error || "Некорректная ссылка Kufar";
        showToast(msg);
        if (linkNote) linkNote.textContent = msg;
        return;
      }
      throw new Error(body.error || `HTTP ${res.status}`);
    }

    if (body.created === false) {
      showToast("Такая ссылка уже есть");
      if (linkNote) linkNote.textContent = "Дубликат не сохраняем.";
    } else {
      showToast("Ссылка добавлена");
      if (linkNote) linkNote.textContent = "Добавлено. Управление — в списке ниже.";
      linkForm.reset();
      await loadLinks();
    }
    linkNote?.classList.add("is-flash");
    setTimeout(() => linkNote?.classList.remove("is-flash"), 700);
  } catch (err) {
    console.error(err);
    showToast("Не удалось добавить");
    if (linkNote) {
      linkNote.textContent =
        "Ошибка. Проверь, что backend запущен и api= доступен.";
    }
  } finally {
    if (linkSubmit) linkSubmit.disabled = false;
  }
});

loadLinks();

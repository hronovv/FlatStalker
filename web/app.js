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
      throw new Error(body.error || `HTTP ${res.status}`);
    }

    if (body.created === false) {
      showToast("Такая ссылка уже есть");
      if (linkNote) linkNote.textContent = "Дубликат не сохраняем. Список: /links";
    } else {
      showToast("Ссылка добавлена");
      if (linkNote) linkNote.textContent = "Добавлено. В боте: /links";
      linkForm.reset();
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

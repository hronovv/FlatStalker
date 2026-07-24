const tg = window.Telegram?.WebApp;

if (tg) {
  tg.ready();
  tg.expand();
  if (tg.setHeaderColor) tg.setHeaderColor("#070707");
  if (tg.setBackgroundColor) tg.setBackgroundColor("#070707");
}

const user = tg?.initDataUnsafe?.user;
const hello = document.getElementById("hello");
const form = document.getElementById("filters-form");
const note = document.getElementById("form-note");
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

form?.addEventListener("submit", (event) => {
  event.preventDefault();
  note?.classList.add("is-flash");
  showToast("Сохранение скоро будет доступно");
  setTimeout(() => note?.classList.remove("is-flash"), 700);
});

const tg = window.Telegram?.WebApp;

if (tg) {
  tg.ready();
  tg.expand();
  if (tg.setHeaderColor) tg.setHeaderColor("#080808");
  if (tg.setBackgroundColor) tg.setBackgroundColor("#080808");
}

const user = tg?.initDataUnsafe?.user;
const greeting = document.getElementById("greeting");
const lead = document.getElementById("lead");
const cta = document.getElementById("cta-main");
const status = document.getElementById("status");

const botUsername = "flatstalker_bot";
cta.href = `https://t.me/${botUsername}`;

if (user?.first_name) {
  greeting.textContent = `${user.first_name}, ты в системе`;
  lead.textContent = "Лоты аренды по Беларуси — сразу пушем со ссылкой.";
  if (status) status.textContent = "USER LINKED";
}

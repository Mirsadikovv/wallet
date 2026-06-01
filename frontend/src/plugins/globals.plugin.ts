import type { App } from "vue";

const translations: Record<string, string> = {
  save: "Сохранить",
  cancel: "Отмена",
  back: "Назад",
  loading: "Загрузка...",
  yes: "Да",
  no: "Нет",
  delete: "Удалить",
  restore: "Восстановить",
  created_at: "Дата создания",

  wallet_page_title: "Кошельки",
  wallet_create_title: "Создать кошелёк",
  wallet_view_title: "Кошелёк",
  wallet_list: "Мои кошельки",
  create_wallet: "Создать кошелёк",
  wallet_type: "Тип кошелька",
  network: "Сеть",
  address: "Адрес",
  balance: "Баланс",
  seqno: "Seqno",
  is_active: "Активен",
  transactions: "Транзакции",
  send_ton: "Отправить TON",
  recipient: "Адрес получателя",
  amount: "Сумма (TON)",
  comment: "Комментарий",
  hash: "Хэш транзакции",
  from: "Отправитель",
  to: "Получатель",
  fee: "Комиссия",
  type: "Тип",
  success: "Успешно",

  account_page_title: "Счета",
  account_create_title: "Создать счёт",
  account_edit_title: "Редактировать счёт",
  account_view_title: "Счёт",
  account_list: "Мои счета",
  create_account: "Создать счёт",
  account_name: "Название счёта",
  currency: "Валюта",
  user_id: "ID пользователя",

  page_for_create: "Создание",
  page_for_edit: "Редактирование",
  page_for_view: "Просмотр",
};

const lang = {
  languages: [{ id: 1, name: "Русский", code: "ru" }],
  _currentLang: { id: 1, name: "Русский", code: "ru" },
};

declare module "@vue/runtime-core" {
  interface ComponentCustomProperties {
    $tl: (key: string) => string;
    $canPage: (pageName: string) => boolean;
    $lang: typeof lang;
  }
}

export function installGlobals(app: App) {
  app.config.globalProperties.$tl = (key: string) => translations[key] ?? key;
  app.config.globalProperties.$canPage = (_name: string) => true;
  app.config.globalProperties.$lang = lang;
}

# new-vue-page

Создай новый Vue 3 модуль с именем `$ARGUMENTS` в `frontend/src/modules/`.

## Что создать

```
frontend/src/modules/$ARGUMENTS/
├── pages/
│   ├── Page.vue     ← список с таблицей (q-table)
│   ├── Create.vue   ← форма создания
│   ├── Edit.vue     ← форма редактирования
│   └── View.vue     ← просмотр (q-markup-table, без формы)
├── service/
│   └── index.ts     ← API-сервис (класс + синглтон-экспорт)
└── routes.ts        ← роуты модуля
```

## Требования

- Строго следуй CLAUDE_STYLE_FRONTEND.md
- Только `<script setup lang="ts">`, только Composition API
- Все компоненты Quasar 2 (q-table, q-form, q-input и т.д.)
- Типы — отдельные интерфейсы в service/index.ts рядом с методами
- Axios через `$axios` из Quasar plugin
- Роуты подключаются через `router.addRoute` или экспорт в routes.ts

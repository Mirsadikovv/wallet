# CLAUDE_STYLE — Frontend (Vue 3 + TypeScript)

> Правила для кода в `idoctor_front/`. Строго следуй при написании нового кода.

---

## Стек

- **Фреймворк**: Vue 3 — только Composition API, только `<script setup lang="ts">`
- **Язык**: TypeScript 5.7
- **UI**: Quasar 2
- **State**: Pinia
- **Router**: Vue Router 4
- **HTTP**: Axios (настроен в `plugins/axios.plugin.ts`)
- **CSS**: UnoCSS + Sass
- **Сборка**: Vite + Bun
- **Форматирование**: Prettier — tabs, 100 chars, double quotes, semicolons, trailing commas

---

## Структура

```
idoctor_front/src/
├── modules/
│   └── <ModuleName>/
│       ├── pages/
│       │   ├── Page.vue    ← список с таблицей
│       │   ├── Create.vue  ← форма создания
│       │   ├── Edit.vue    ← форма редактирования
│       │   └── View.vue    ← просмотр (q-markup-table, без Form)
│       ├── service/
│       │   └── index.ts    ← API-сервис (класс + синглтон)
│       ├── routes.ts       ← роуты модуля
│       └── components/     ← компоненты модуля (если нужны)
├── components/
│   ├── quasar/             ← обёртки над Quasar компонентами
│   └── ...                 ← глобальные компоненты
├── store/
│   ├── auth-store.ts
│   └── language-store.ts
├── composables/
│   └── use<Name>.ts
├── plugins/
│   ├── axios.plugin.ts
│   ├── install.ts
│   ├── pinia.plugin.ts
│   ├── quasar.plugin.ts
│   └── router.plugin.ts
├── common/
│   ├── try.ts              ← @Try декоратор
│   ├── Notify.ts           ← уведомления
│   ├── validator.ts        ← валидаторы форм
│   ├── guard.ts
│   ├── permission.ts
│   ├── workflow.ts
│   └── local.sorage.ts
├── service/
│   └── index.ts            ← реэкспорт всех сервисов + общие типы
└── styles/
    └── telegram-app.scss
```

---

## Именование

| Что | Стиль | Пример |
|-----|-------|--------|
| Файлы страниц | `PascalCase.vue` | `Create.vue`, `Edit.vue`, `Page.vue`, `View.vue` |
| Файлы компонентов | `PascalCase.vue` | `ResponsiveTable.vue`, `IconBtn.vue` |
| Файлы сервисов | `index.ts` в папке `service/` | `service/index.ts` |
| Файлы роутов | `routes.ts` | `routes.ts` |
| Файлы composables | `use<Name>.ts` | `useAppNavigation.ts`, `useBreakpoints.ts` |
| Файлы сторов | `kebab-case-store.ts` | `auth-store.ts`, `language-store.ts` |
| Типы данных | `PascalCase` + суффикс `Type` | `SupplierType`, `SupplierCreateType` |
| Типы страниц | `<Entity>PageData` | `SupplierPageData`, `OrderPageData` |
| Общий тип страницы | `PageDataType<T>` | `PageDataType<SupplierType>` |
| Переменные | `camelCase` | `supplierPage`, `orderModel`, `authStore` |
| Функции | `camelCase` | `save`, `page`, `getSupplierByID` |
| Имена роутов | `SCREAMING_CASE` | `SUPPLIER_PAGE`, `SUPPLIER_CREATE`, `SUPPLIER_VIEW`, `SUPPLIER_EDIT` |
| Группы роутов (meta) | `SCREAMING_CASE` + `_GROUP` | `SUPPLIER_GROUP`, `ORDER_GROUP` |
| Экземпляры сервисов | реэкспорт как `PascalCase` | `export { supplierService as SupplierService }` |
| JSON-поля в типах | `snake_case` | `created_at`, `deleted_at`, `client_id` |

---

## Паттерны

### Сервис — класс + `@Try` на каждом методе + синглтон

```typescript
class SupplierService {
    @Try({
        async onError(err) {
            (await import("@/common/Notify")).ErrorNotify(
                err?.response?.data.message || err.message,
            );
        },
    })
    async page(params: SupplierSearchParams = {}) {
        const searchParams = new URLSearchParams();
        Object.entries(params).forEach(([key, value]) => {
            if (value !== undefined) searchParams.append(key, String(value));
        });
        const { data } = await api.get<SupplierPageData>(`/supplier/page?${searchParams}`);
        return data;
    }

    @Try({
        async onError(err) {
            (await import("@/common/Notify")).ErrorNotify(
                err?.response?.data.message || err.message,
            );
        },
    })
    async create(dto: SupplierCreateType) {
        const { data } = await api.post<IdType>(`/supplier`, dto);
        return data;
    }

    // Каждый метод — отдельный @Try
}

const supplierService = new SupplierService();
export { supplierService as SupplierService };
```

`@Try` перехватывает ошибки и возвращает `undefined`. Страницы всегда проверяют результат:

```typescript
const response = await SupplierService.create(model);
if (!response) return false;
```

### Типы в сервисе

```typescript
export type SupplierType = {
    id: number;
    name: string;
    created_at: string;
    deleted_at?: string;
};

export type SupplierCreateType = {
    name: string;
};

export type SupplierUpdateType = {
    name: string;
};

export type SupplierPartialType = Partial<SupplierType>;
export type SupplierPageData = PageDataType<SupplierType>;

export type SupplierSearchParams = {
    name?: string;
    include_deleted?: boolean;
    only_deleted?: boolean;
    page?: number;
    perpage?: number;
    limit?: number;
};
```

### Роуты модуля

```typescript
const supplierPageRoute: RouteRecordRaw = {
    path: "suppliers",
    name: "SUPPLIER_PAGE",
    component: () => import("@module/Supplier/pages/Page.vue"),
    meta: {
        title: "supplier_page_title",
        activeLinkGroup: "SUPPLIER_GROUP",
        sidebar: {
            label: "supplier_page_title",
            icon: "business",
            isExpandedGroup: false,
        },
    },
};

const supplierCreateRoute: RouteRecordRaw = {
    path: "suppliers/create",
    name: "SUPPLIER_CREATE",
    props: true,
    component: () => import("@module/Supplier/pages/Create.vue"),
    meta: { title: "supplier_create_title", activeLinkGroup: "SUPPLIER_GROUP" },
};

const supplierViewRoute: RouteRecordRaw = {
    path: "suppliers/:id",
    name: "SUPPLIER_VIEW",
    props: true,
    component: () => import("@module/Supplier/pages/View.vue"),
    meta: { title: "supplier_view_title", activeLinkGroup: "SUPPLIER_GROUP" },
};

const supplierEditRoute: RouteRecordRaw = {
    path: "suppliers/:id/edit",
    name: "SUPPLIER_EDIT",
    props: true,
    component: () => import("@module/Supplier/pages/Edit.vue"),
    meta: { title: "supplier_edit_title", activeLinkGroup: "SUPPLIER_GROUP" },
};

export function SupplierRoutes(sort: number): RouteRecordRaw[] {
    return [supplierPageRoute, supplierCreateRoute, supplierViewRoute, supplierEditRoute].map(
        (route) => {
            if (route?.meta?.sidebar) {
                return { ...route, meta: { ...route.meta, sort } };
            }
            return route;
        },
    );
}
```

### Page.vue — список

```vue
<script setup lang="ts">
import { ref } from "vue";
import { XxxService, type XxxPageData } from "../service";
import PageLoading from "@/components/PageLoading.vue";
import LoadingSkeleton from "@/components/LoadingSkeleton.vue";
import AppFooter from "@/components/AppFooter.vue";
import ResponsiveTable from "@/components/quasar/table/ResponsiveTable.vue";
import TablePaginate from "@/components/quasar/table/TablePaginate.vue";
import IconBtn from "@/components/quasar/btn/IconBtn.vue";
import { useAppNavigation } from "@/composables/useAppNavigation";
import { useTelegramViewport } from "@/composables/useTelegramViewport";
import { useAuthStore } from "@/store/auth-store";

const authStore = useAuthStore();
const { toggleLeftDrawer, setLang, logout } = useAppNavigation();
const { containerStyle } = useTelegramViewport();

const xxxPage = ref<XxxPageData>({
    data: [],
    totalRows: 0,
    currentPage: 0,
    pageSize: 0,
    totalPages: 0,
});

const pikers = ref({});

async function page(query: string = "") {
    const params = new URLSearchParams(query);
    const searchParams = {
        page: params.get("page") ? +params.get("page")! : undefined,
        perpage: params.get("perpage") ? +params.get("perpage")! : undefined,
        name: params.get("name") || undefined,
    };

    const response = await XxxService.page(searchParams);
    if (!response) return;
    xxxPage.value = response;
}
</script>

<template>
    <PageLoading :find="page" #="{ fetch, loading }">
        <LoadingSkeleton v-if="loading" />

        <q-layout view="hHh Lpr lff" v-else>
            <q-page-container>
                <q-page :style="containerStyle" class="bg-gray-100 text-gray-900 overflow-auto p-4">
                    <ResponsiveTable :models="xxxPage" hasOrder :loading="loading">
                        <template #field:thead>{{ $tl("field_label") }}</template>
                        <template #field="{ model }">{{ model.field }}</template>

                        <template #edit="{ model }">
                            <div class="text-center">
                                <IconBtn
                                    v-if="$canPage('XXX_EDIT')"
                                    :to="{ name: 'XXX_EDIT', params: { id: model.id } }"
                                    icon="edit"
                                />
                                <IconBtn
                                    v-if="$canPage('XXX_VIEW')"
                                    :to="{ name: 'XXX_VIEW', params: { id: model.id } }"
                                    icon="visibility"
                                />
                            </div>
                        </template>

                        <template #tfoot="{ totalPages }">
                            <TablePaginate
                                v-model:pikers="pikers"
                                :total="totalPages"
                                @page="fetch"
                            />
                        </template>
                    </ResponsiveTable>
                </q-page>
            </q-page-container>

            <AppFooter
                :username="authStore.user?.username"
                :languages="$lang.languages"
                :current-language-id="$lang._currentLang?.id"
                :show-add-button="true"
                :add-button-route="{ name: 'XXX_CREATE' }"
                add-button-icon="add"
                @toggle-drawer="toggleLeftDrawer"
                @go-to-profile="toggleLeftDrawer"
                @set-lang="setLang"
                @logout="logout"
            />
        </q-layout>
    </PageLoading>
</template>

<style scoped>
@import "@/styles/telegram-app.scss";
</style>
```

### Create.vue — форма создания

```vue
<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { XxxService, type XxxCreateType } from "../service";
import Form from "@/components/quasar/form/Form.vue";
import Button from "@/components/quasar/btn/Button.vue";
import Input from "@/components/quasar/form/Input.vue";
import Title from "@/components/Title.vue";
import AppFooter from "@/components/AppFooter.vue";
import { useAppNavigation } from "@/composables/useAppNavigation";
import { useTelegramViewport } from "@/composables/useTelegramViewport";
import { useAuthStore } from "@/store/auth-store";
import { formRequired } from "@/common/validator";

const router = useRouter();
const authStore = useAuthStore();
const { toggleLeftDrawer, setLang, logout } = useAppNavigation();
const { containerStyle } = useTelegramViewport();

const xxxModel = ref<Partial<XxxCreateType>>({});

async function save(model: XxxCreateType) {
    const response = await XxxService.create(model);
    if (!response) return false;
    router.push({ name: "XXX_PAGE" });
    return true;
}
</script>

<template>
    <q-layout view="hHh Lpr lff">
        <q-page-container>
            <q-page :style="containerStyle" class="bg-gray-100 text-gray-900 overflow-auto p-4">
                <div class="flex! gap-x-4 items-center mb-3">
                    <q-btn flat color="accent" icon="arrow_back" @click="router.back()" />
                    <q-breadcrumbs>
                        <q-breadcrumbs-el
                            :label="$tl('xxx_list')"
                            icon="icon_name"
                            :to="{ name: 'XXX_PAGE' }"
                        />
                        <q-breadcrumbs-el :label="$tl('page_for_create')" />
                    </q-breadcrumbs>
                </div>
                <Form v-model="xxxModel" :save="save">
                    <template #title>
                        <Title class="mb-5">{{ $tl("create_xxx") }}</Title>
                    </template>

                    <template #field_name="{ model }">
                        <Input
                            v-model="model.field_name"
                            label="Field Name"
                            class="col-lg-6 col-md-6 col-12"
                            :rules="[formRequired()]"
                        />
                    </template>

                    <template #actions="{ loading }">
                        <div class="row my-row justify-center">
                            <q-space />
                            <div class="col-auto max-w-200px! w-100%!">
                                <Button :loading="loading" type="submit" class="w-100%!">
                                    {{ $tl("save") }}
                                </Button>
                            </div>
                        </div>
                    </template>
                </Form>
            </q-page>
        </q-page-container>

        <AppFooter
            :username="authStore.user?.username"
            :languages="$lang.languages"
            :current-language-id="$lang._currentLang?.id"
            :show-add-button="false"
            @toggle-drawer="toggleLeftDrawer"
            @go-to-profile="toggleLeftDrawer"
            @set-lang="setLang"
            @logout="logout"
        />
    </q-layout>
</template>

<style scoped>
@import "@/styles/telegram-app.scss";
</style>
```

### Edit.vue — форма редактирования

Отличия от Create.vue: `PageLoading` + `LoadingSkeleton`, `props: { id }`, загрузка данных перед рендером.

```vue
<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { XxxService, type XxxUpdateType, type XxxType } from "../service";
// ... те же импорты что и Create.vue + PageLoading, LoadingSkeleton

export interface Props {
    id: number;
}
const { id } = defineProps<Props>();

const router = useRouter();
// ... те же composables

const xxxModel = ref<XxxType>({ id: 0, name: "", created_at: "" });

async function getXxxByID() {
    const data = await XxxService.getByID(+id);
    if (data) xxxModel.value = data;
}

async function save(model: XxxUpdateType) {
    const response = await XxxService.update(+id, model);
    if (!response) return false;
    router.push({ name: "XXX_PAGE" });
    return true;
}
</script>

<template>
    <PageLoading :find="getXxxByID" #="{ loading }">
        <LoadingSkeleton v-if="loading" />
        <q-layout view="hHh Lpr lff" v-else>
            <!-- та же структура что Create.vue -->
        </q-layout>
    </PageLoading>
</template>
```

### View.vue — просмотр (q-markup-table, без Form)

```vue
<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { XxxService, type XxxType } from "../service";
import PageLoading from "@/components/PageLoading.vue";
import LoadingSkeleton from "@/components/LoadingSkeleton.vue";
import AppFooter from "@/components/AppFooter.vue";
import { useAppNavigation } from "@/composables/useAppNavigation";
import { useTelegramViewport } from "@/composables/useTelegramViewport";
import { useAuthStore } from "@/store/auth-store";

export interface Props {
    id: number | string;
}
const { id } = defineProps<Props>();

const router = useRouter();
const authStore = useAuthStore();
const { toggleLeftDrawer, setLang, logout } = useAppNavigation();
const { containerStyle } = useTelegramViewport();

const xxxModel = ref<XxxType>({} as XxxType);

const loadXxx = async () => {
    const data = await XxxService.getByID(+id);
    xxxModel.value = data;
};
</script>

<template>
    <PageLoading :find="loadXxx" #="{ loading }">
        <LoadingSkeleton v-if="loading" />

        <q-layout view="hHh Lpr lff" v-else>
            <q-page-container>
                <q-page :style="containerStyle" class="bg-gray-100 text-gray-900 overflow-auto p-4">
                    <div class="flex! gap-x-4 items-center mb-3">
                        <q-btn flat color="accent" icon="arrow_back" @click="router.back()" />
                        <q-breadcrumbs>
                            <q-breadcrumbs-el
                                :label="$tl('xxx_list')"
                                icon="icon_name"
                                :to="{ name: 'XXX_PAGE' }"
                            />
                        </q-breadcrumbs>
                        <q-space />
                    </div>

                    <q-markup-table separator="cell" flat bordered>
                        <tbody>
                            <tr>
                                <td class="font-bold text-left">{{ $tl("field_label") }}</td>
                                <td>{{ xxxModel?.field }}</td>
                            </tr>
                        </tbody>
                    </q-markup-table>
                </q-page>
            </q-page-container>

            <AppFooter
                :username="authStore.user?.username"
                :languages="$lang.languages"
                :current-language-id="$lang._currentLang?.id"
                :show-add-button="true"
                :add-button-route="{ name: 'XXX_CREATE' }"
                add-button-icon="add"
                @toggle-drawer="toggleLeftDrawer"
                @go-to-profile="toggleLeftDrawer"
                @set-lang="setLang"
                @logout="logout"
            />
        </q-layout>
    </PageLoading>
</template>

<style scoped>
@import "@/styles/telegram-app.scss";
</style>
```

### Pinia store

```typescript
type State = {
    xxx: XxxType | undefined;
};

export const useXxxStore = defineStore("xxx", {
    state(): State {
        return {
            xxx: undefined,
        };
    },
    actions: {
        setXxx(data: XxxType) {
            this.xxx = data;
        },
        removeXxx() {
            this.xxx = undefined;
        },
    },
    getters: {
        hasXxx(state): Readonly<boolean> {
            return !!state.xxx;
        },
        getXxx(state): Readonly<XxxType> {
            return state.xxx as Readonly<XxxType>;
        },
    },
});
```

### Composable

```typescript
export function useXxx() {
    const router = useRouter();
    const authStore = useAuthStore();

    async function doSomething() {
        // ...
    }

    return {
        doSomething,
    };
}
```

---

## Deprecated — запрещённые паттерны

### Vue 3 — удалённое из Vue 2

```typescript
// ❌ $listeners — слит в $attrs в Vue 3
this.$listeners.click()
// ✅
$attrs.onClick

// ❌ $scopedSlots — слит в $slots в Vue 3
$scopedSlots.default()
// ✅
$slots.default?.()

// ❌ $on / $off / $once — удалены в Vue 3
this.$on('event', handler)
// ✅ — использовать mitt или reactive ref
import mitt from 'mitt'
const emitter = mitt()

// ❌ Vue.set / Vue.delete — не нужны с ref/reactive
Vue.set(obj, 'key', value)
// ✅
obj.key = value  // реактивность работает автоматически

// ❌ .native модификатор — удалён
@click.native="handler"
// ✅
@click="handler"

// ❌ v-model с .sync — заменён на v-model:prop
:value.sync="val"
// ✅
v-model:value="val"

// ❌ $children — удалён в Vue 3
this.$children[0].method()
// ✅ — использовать ref + defineExpose
const childRef = ref()
childRef.value?.method()
```

### TypeScript — устаревшие паттерны

```typescript
// ❌ namespace — устарел, не использовать
namespace MyModule { ... }

// ❌ Enum — предпочесть const object (лучше tree-shaking)
enum Status { Active = 'ACTIVE' }
// ✅
const Status = { Active: 'ACTIVE' } as const
type Status = typeof Status[keyof typeof Status]

// ❌ as any — скрывает ошибки типов
const val = something as any
// ✅
const val = something as unknown as SpecificType
```

### Axios — устаревшие паттерны

```typescript
// ❌ axios.create внутри сервиса — инстанс создаётся один раз в плагине
const api = axios.create({ baseURL: '...' })

// ✅ — импортировать готовый инстанс
import { api } from '@/plugins/axios.plugin'

// ❌ прямой импорт axios
import axios from 'axios'
axios.get('/url')

// ✅ — только через настроенный инстанс
import { api } from '@/plugins/axios.plugin'
api.get('/url')
```

### Pinia — устаревшие паттерны

```typescript
// ❌ mapState / mapActions — Vue 2 Vuex-стиль
...mapState(useStore, ['value'])

// ✅ — деструктурировать с storeToRefs
const store = useXxxStore()
const { value } = storeToRefs(store)
const { action } = store
```

---

## Запрещено

- ❌ Options API — только `<script setup lang="ts">`
- ❌ `defineComponent` — только SFC
- ❌ `this` в компонентах
- ❌ `.then()` цепочки — только `async/await`
- ❌ `try/catch` в сервисах вручную — только `@Try` декоратор
- ❌ Vuex — только Pinia
- ❌ `console.log` в коде
- ❌ `any` в типах сервисов и компонентов
- ❌ `interface` для описания форм/сервисов — только `type`
- ❌ Комментарии в коде (исключение: `// TODO:`)
- ❌ Дублировать логику между страницами — выносить в composables или сервисы
- ❌ `$listeners`, `$scopedSlots`, `$on/$off/$once`, `$children` — удалены в Vue 3
- ❌ `.native` модификатор и `.sync` — удалены в Vue 3
- ❌ `Vue.set` / `Vue.delete` — не нужны с Composition API
- ❌ `axios` импортировать напрямую — только через инстанс из плагина
- ❌ `enum` — использовать `const` объект с `as const`

---

## Комментарии

Комментарии **не пишутся**. Единственное исключение — `// TODO:`.

Если комментарий всё-таки нужен — писать **по-русски**.

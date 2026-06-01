export function formRequired(message = "Поле обязательно") {
  return (val: string | null | undefined) => !!val || message;
}

export function formMinLength(min: number, message?: string) {
  return (val: string | null | undefined) =>
    (!!val && val.length >= min) || (message ?? `Минимум ${min} символов`);
}

export function formMaxLength(max: number, message?: string) {
  return (val: string | null | undefined) =>
    !val || val.length <= max || (message ?? `Максимум ${max} символов`);
}

export function formPositiveNumber(message = "Введите положительное число") {
  return (val: string | null | undefined) => {
    if (!val) return message;
    const num = parseFloat(val);
    return (!isNaN(num) && num > 0) || message;
  };
}

export function formTONAddress(message = "Некорректный TON-адрес") {
  return (val: string | null | undefined) => {
    if (!val) return message;
    return /^[A-Za-z0-9_-]{48}$/.test(val) || message;
  };
}

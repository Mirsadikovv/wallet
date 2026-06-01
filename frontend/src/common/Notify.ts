import { Notify } from "quasar";

export function SuccessNotify(message: string) {
  Notify.create({
    type: "positive",
    message,
    position: "top",
    timeout: 3000,
  });
}

export function ErrorNotify(message: string) {
  Notify.create({
    type: "negative",
    message,
    position: "top",
    timeout: 5000,
  });
}

export function WarnNotify(message: string) {
  Notify.create({
    type: "warning",
    message,
    position: "top",
    timeout: 4000,
  });
}

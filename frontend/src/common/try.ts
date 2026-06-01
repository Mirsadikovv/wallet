type TryOptions = {
  onError?: (err: unknown) => void | Promise<void>;
};

export function Try(options: TryOptions = {}) {
  return function (_target: object, _propertyKey: string, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value;

    descriptor.value = async function (...args: unknown[]) {
      try {
        return await originalMethod.apply(this, args);
      } catch (err) {
        if (options.onError) {
          await options.onError(err);
        }
        return undefined;
      }
    };

    return descriptor;
  };
}

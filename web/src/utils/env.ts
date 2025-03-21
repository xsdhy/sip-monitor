/**
 * 获取环境变量，支持类型安全
 * @param key 环境变量key
 * @param defaultValue 默认值
 */
export function getEnv<T extends string | boolean | number>(
  key: string,
  defaultValue?: T
): T {
  const value = import.meta.env[key] as unknown;

  if (value === undefined) {
    if (defaultValue !== undefined) {
      return defaultValue;
    }
    throw new Error(`Environment variable ${key} is not defined`);
  }

  // 根据默认值类型或已知类型进行转换
  if (defaultValue !== undefined) {
    const type = typeof defaultValue;
    if (type === 'boolean') {
      return (value === 'true') as unknown as T;
    }
    if (type === 'number') {
      return Number(value) as unknown as T;
    }
  }

  return value as T;
}

/**
 * 是否为生产环境
 */
export const isProd = getEnv('VITE_MODE', 'production') === 'production';

/**
 * 是否为开发环境
 */
export const isDev = getEnv('VITE_MODE', 'production') === 'development';


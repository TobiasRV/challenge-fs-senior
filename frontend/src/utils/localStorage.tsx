export const setLsItem = (key: string, value: any) => {
    localStorage.setItem(key, JSON.stringify(value))
}

export const getLsItem = (key: string): any => {
    const value = localStorage.getItem(key);
    return value ? JSON.parse(value) : null;
}

export const removeLsItem = (key: string) => localStorage.removeItem(key);

export const clearLs = () => localStorage.clear()
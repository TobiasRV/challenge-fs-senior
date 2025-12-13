const hasLocalStorage = () => typeof window !== "undefined" && typeof window.localStorage !== "undefined" && typeof window.localStorage.getItem === "function";

export const setLsItem = (key: string, value: any) => {
    if (!hasLocalStorage()) {
        return;
    }
    try {
        window.localStorage.setItem(key, JSON.stringify(value));
    } catch (err) {
    }
}

export const getLsItem = (key: string): any => {
    if (!hasLocalStorage()) {
        return null;
    }
    try {
        const value = window.localStorage.getItem(key);
        const parsed = value ? JSON.parse(value) : null;
        return parsed;
    } catch (err) {
        return null;
    }
}

export const removeLsItem = (key: string) => {
    if (!hasLocalStorage()) {
        return;
    }
    window.localStorage.removeItem(key);
};

export const clearLs = () => {
    if (!hasLocalStorage()) {
        return;
    }
    window.localStorage.clear();
}
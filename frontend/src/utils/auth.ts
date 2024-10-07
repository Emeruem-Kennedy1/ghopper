const storeToken = (token: string) => {
    if (navigator.cookieEnabled) {
        document.cookie = `token=${token}; path=/; max-age=86400; SameSite=Strict`;
    } else {
        localStorage.setItem('token', token);
    }
};

const getToken = () => {
    if (navigator.cookieEnabled) {
        const cookie = document.cookie.split(';').find((cookie) => cookie.includes('token'));
        return cookie?.split('=')[1];
    }
    return localStorage.getItem('token');
}

const removeToken = () => {
    if (navigator.cookieEnabled) {
        document.cookie = `token=; path=/; max-age=0; SameSite=Strict`;
    } else {
        localStorage.removeItem('token');
    }
}

export { storeToken, getToken, removeToken };
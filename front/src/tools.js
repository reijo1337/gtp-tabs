import jwtDecode from "jwt-decode"

export function parse_json(data) {
    return data.text().then(function (text) {
        return text ? JSON.parse(text) : {}
    }).catch((err) => {
        console.log(err);
        return {};
    });
}

export function parse_date(date) {
    const year = date.slice(0, 4);
    const mounth = date.slice(4, 6);
    const day = date.slice(6, 8);
    return day + "." + mounth + "." + year;
}

export function updater() {
    const token = localStorage.getItem("refreshToken");
    console.log("refresh");
    const url = "http://localhost:5000/auth?refresh_token="+token;
    fetch(url)
        .then( res => {
            if (res.status === 200) {
                return parse_json(res);
            } else {
                return res.json();
            }
        })
        .then(json => {
            if (json.error) {
                throw new Error(json.error);
            }
            localStorage.setItem("accessToken", json.accessToken);
            localStorage.setItem("refreshToken", json.refreshToken);
            localStorage.setItem("login", this.state.login);
            clearInterval(this._tokenUpdater);
            const token = json.accessToken;
            let tokenData = jwtDecode(token);
            let interval = (tokenData.exp - (Date.now().valueOf() / 1000))-10;

            this._tokenUpdater = setInterval(updater.bind(this),interval*1000);
        })
        .catch((error) => {
            alert("Cant refresh token: " + error.message);
            localStorage.setItem("accessToken", "");
            localStorage.setItem("refreshToken", "");
            localStorage.setItem("login", "");
        });
}
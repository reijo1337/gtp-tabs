import React, {Component} from 'react';
import {parse_json, updater} from "../../tools";
import jwtDecode from "jwt-decode";

class CheckVkCode extends Component {
    constructor(props) {
        super(props);
        const urlParams = new URLSearchParams(window.location.search);
        const myParam = urlParams.get('code');
        debugger;
        if (myParam !== null) {
            this.url = "http://127.0.0.1:9090/auth/vk?code=" + myParam;
            this.getAccess();
        } else {
            window.location = "http://127.0.0.1:3000";
        }
    }

    getAccess = () => {
        debugger;
        let response = fetch(this.url)
            .then(res => {
                debugger;
                if (res.status === 200) {
                    return parse_json(res);
                } else {
                    return res.json();
                }
            })
            .then(json => {
                debugger;
                if (json.error) {
                    throw new Error(json.error);
                }
                localStorage.setItem("accessToken", json.tokens.accessToken);
                localStorage.setItem("refreshToken", json.tokens.refreshToken);
                localStorage.setItem("profileID", json.profile_id);
                localStorage.setItem("login", "VK"+json.profile_id.toString());
                clearInterval(this._tokenUpdater);
                const token = json.tokens.accessToken;
                let tokenData = jwtDecode(token);
                let interval = (tokenData.exp - (Date.now().valueOf() / 1000))-10;

                this._tokenUpdater = setInterval(updater.bind(this),interval*1000);
                window.location = "http://localhost:3000";

            })
            .catch((error) => {
                alert(error.message);
                window.location = "http://localhost:3000";

            });
        return response;
    };

    render() {
        return (<div></div>);
    }
}

export default CheckVkCode;
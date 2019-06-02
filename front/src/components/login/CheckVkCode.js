import React, {Component} from 'react';
import VkAuth from "./VkAuth";
import {parse_json, updater} from "../../tools";
import jwtDecode from "jwt-decode";

class CheckVkCode extends Component {
    constructor(props) {
        super(props);
        const urlParams = new URLSearchParams(window.location.search);
        const myParam = urlParams.get('code');
        debugger;
        if (myParam !== null) {
            const data = {
                "client_id": "6978271",
                "client_secret": "DKdWDto5gJU4ViGrJW4d",
                "redirect_uri": "http://localhost:3000/vk",
                "code": myParam,
            };
            this.url = "https://oauth.vk.com/access_token?" + VkAuth.encodeQueryData(data);
            this.getAccess();
        } else {
            window.location = "http://127.0.0.1:3000";
        }
    }

    getAccess = async () => {

        let response = await fetch(this.url, {
            method: "get",
            mode: "cors",
        })
            .then(res => {
                debugger;
                if (res.status === 0) {
                    return parse_json(res);
                } else {
                    return res.json();
                }
            })
            .then(json => {
                if (json.error) {
                    throw new Error(json.error);
                }
                this.access_token = json.access_token;
                this.expires_in = json.expires_in;
                this.user_id = json.user_id;
                alert("vk ok");
                this.loginOrRegister();
            })
            .catch((error) => {
                alert(error.message)
            });
        return response;
    };

    loginOrRegister = async () => {
        const data = JSON.stringify({
            user_id: this.user_id,
        });
        const url = "http://localhost:9090/auth/vk";
        fetch(url, {
            method: "post",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: data
        })
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
                localStorage.setItem("accessToken", json.tokens.accessToken);
                localStorage.setItem("refreshToken", json.tokens.refreshToken);
                localStorage.setItem("profileID", json.profile_id);
                localStorage.setItem("login", this.state.login);
                clearInterval(this._tokenUpdater);
                const token = json.tokens.accessToken;
                let tokenData = jwtDecode(token);
                let interval = (tokenData.exp - (Date.now().valueOf() / 1000))-10;

                this._tokenUpdater = setInterval(updater.bind(this),interval*1000);
                window.location = "http://127.0.0.1:3000";
            })
            .catch((error) => {
                alert("Проблемы с доступом в джойказино: " + error.message);
            });
    };

    render() {
        return (<div></div>);
    }
}

export default CheckVkCode;
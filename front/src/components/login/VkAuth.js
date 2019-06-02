import React, {Component} from 'react';
import {Button} from "react-bootstrap";

class VkAuth extends Component {
    constructor(props){
        super(props);
        const data = {
            "client_id": "6978271",
            "redirect_uri": "http://localhost:3000/vk",
            "response_type": "code",
            "v": "5.95",
        };
        this.url = "https://oauth.vk.com/authorize?" + VkAuth.encodeQueryData(data);
    }

    static encodeQueryData(data) {
        const ret = [];
        for (let d in data)
            ret.push(encodeURIComponent(d) + '=' + encodeURIComponent(data[d]));
        return ret.join('&');
    }

    render() {
        return (
            <Button
                block
                bsSize="large"
                href={this.url}
            >
                Авторизация через VK
            </Button>
        )
    }
}

export default VkAuth;
import React, {Component} from 'react';
import {Button, Form, FormControl, FormGroup, Modal} from "react-bootstrap";
import {parse_json, updater} from "../../tools";
import jwtDecode from "jwt-decode";

class Login extends Component {
    constructor(props) {
        super(props);
        this.url = "http://localhost:9090/auth";
        this.state = {
            login: "",
            password: "",
        }
    }
    render() {
        return (
            <Modal
                {...this.props}
                size="lg"
                aria-labelledby="contained-modal-title-vcenter"
                centered
            >
                <Modal.Header closeButton>
                    <Modal.Title id="contained-modal-title-vcenter">
                        Авторизация
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <form onSubmit={this.handleSubmit}>
                        <FormGroup controlId="login" >
                            <Form.Label>Логин</Form.Label>
                            <FormControl
                                autoFocus
                                type="login"
                                value={this.state.login}
                                onChange={this.handleChange}
                            />
                        </FormGroup>
                        <FormGroup controlId="password" >
                            <Form.Label>Пароль</Form.Label>
                            <FormControl
                                value={this.state.password}
                                onChange={this.handleChange}
                                type="password"
                            />
                        </FormGroup>
                        <Button
                            block
                            bsSize="large"
                            disabled={!this.validateForm()}
                            type="submit"
                        >
                            Авторизоваться
                        </Button>
                    </form>
                </Modal.Body>
            </Modal>
        );
    }

    handleChange = event => {
        this.setState({
            [event.target.id]: event.target.value
        });
    };

    validateForm() {
        return this.state.login.length > 0 && this.state.password.length > 0;
    }

    handleSubmit = event => {
        event.preventDefault();
        const login = this.state.login;
        const password = this.state.password;
        const data = JSON.stringify({
            login: login,
            password: password,
        });
        debugger;
        fetch(this.url, {
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
                this.props.setAuth();
                this.props.setClose();
            })
            .catch((error) => {
                alert("Проблемы с доступом в джойказино: " + error.message);
            });
    };

}

export default Login;
import React, {Component} from 'react';
import {Button, Form, Nav, Navbar} from "react-bootstrap";
import Register from "../login/Register";
import Login from "../login/Login";
import jwtDecode from "jwt-decode"

class TopHeader extends Component {
    constructor(...args) {
        super(...args);
        const token = localStorage.getItem("accessToken");
        if (token === null || token === "") {
            this.state = {
                registerShow: false,
                loginShow: false,
                login: '',
                password: '',
                authorized: false,
            };
            return
        }
        let tokenData = jwtDecode(token);
        let interval = (tokenData.exp - (Date.now().valueOf() / 1000))-10;
        if (interval < 0) {
            localStorage.setItem("accessToken", "");
            localStorage.setItem("refreshToken", "");
            localStorage.setItem("login", "");
            this.state = {
                registerShow: false,
                loginShow: false,
                login: '',
                password: '',
                authorized: false,
            }
        } else {
            const login = localStorage.getItem("login");
            this.state = {
                registerShow: false,
                loginShow: false,
                login: login,
                password: '',
                authorized: true,
            };
        }
    }

    render() {
        let topButtons;
        if (this.state.authorized) {
            topButtons = <div>
                <Button
                    variant="outline-info"
                >{this.state.login}</Button>
                <Button
                    variant="outline-info"
                    onClick={this.logout}
                >
                    Выйти
                </Button>
            </div>
        } else {
            topButtons = <div>
                <Button
                    variant="outline-info"
                    onClick={() => this.setState({loginShow: true})}
                >Войти</Button>
                <Button
                    variant="outline-info"
                    onClick={() => this.setState({registerShow: true})}
                >Регистрация</Button>
            </div>;
        }
        return (
            <div>
                <Navbar bg="dark" variant="dark">
                    <Navbar.Brand href="/">Огромное хранилище табулатур</Navbar.Brand>
                    <Nav className="mr-auto">
                        <Nav.Link href="/feedback">Обратная связь</Nav.Link>
                    </Nav>
                    <Form inline>
                        {topButtons}
                    </Form>
                </Navbar>
                <Register
                    show={this.state.registerShow}
                    onHide={this.registerClose.bind(this)}
                    setAuth={this.authorized.bind(this)}
                    setClose={this.registerClose.bind(this)}
                />
                <Login
                    show={this.state.loginShow}
                    setAuth={this.authorized.bind(this)}
                    setClose={this.loginClose.bind(this)}
                    onHide={this.loginClose.bind(this)}/>
            </div>
        );
    }
    authorized = () => {
        this.setState({
            authorized: true,
            login: localStorage.getItem("login"),
        });
    };
    registerClose = () => this.setState({ registerShow: false });
    loginClose = () => this.setState({ loginShow: false });
    logout = () => {
        this.setState({
            authorized: false,
            login: "",
        });
        localStorage.setItem("accessToken", "");
        localStorage.setItem("refreshToken", "");
        localStorage.setItem("login", "");
    }
}

export default TopHeader;
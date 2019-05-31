import React, {Component} from 'react';
import {Button, Form, Nav, Navbar} from "react-bootstrap";
import Register from "../login/Register";
import Login from "../login/Login";

class TopHeader extends Component {
    constructor(...args) {
        super(...args);

        const login = localStorage.getItem("login");
        console.log(login);
        this.state = {
            registerShow: false,
            loginShow: false,
            login: login,
            authorized: (login !== null && login !== ""),
        };
    }
    render() {
        let topButtons;
        if (this.state.authorized) {
            topButtons = <div>
                <Button
                    variant="outline-info"
                >{this.state.login}</Button>
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
        let registerClose = () => this.setState({ registerShow: false });
        let loginClose = () => this.setState({ loginShow: false });
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
                    onHide={registerClose}
                />
                <Login
                    show={this.state.loginShow}
                    onHide={loginClose}/>
            </div>
        );
    }
}

export default TopHeader;
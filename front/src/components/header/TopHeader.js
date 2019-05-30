import React, {Component} from 'react';
import {Button, Form, Nav, Navbar} from "react-bootstrap";

class TopHeader extends Component {
    render() {
        let topButtons;
        topButtons = <div>
            <Button variant="outline-info">Войти</Button>
            <Button variant="outline-info">Регистрация</Button>
        </div>;
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
            </div>
        );
    }
}

export default TopHeader;
import React, {Component} from 'react';
import {Button, Form, Modal} from "react-bootstrap";

class Login extends Component {
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
                    <Form>
                        <Form.Group controlId="formBasicLogin">
                            <Form.Label>Имя пользователя</Form.Label>
                            <Form.Control type="text" placeholder="Имя пользователя" />
                        </Form.Group>

                        <Form.Group controlId="formBasicPassword">
                            <Form.Label>Пароль</Form.Label>
                            <Form.Control type="password" placeholder="Пароль" />
                        </Form.Group>
                    </Form>
                </Modal.Body>
                <Modal.Footer>
                    <Button onClick={this.props.onHide}>Отмена</Button>
                    <Button>Войти</Button>
                </Modal.Footer>
            </Modal>
        );
    }
}

export default Login;
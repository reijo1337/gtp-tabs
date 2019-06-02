import React, {Component} from 'react';
import {Container, Jumbotron} from "react-bootstrap";

class Glagna extends Component{
    render() {
        return (
            <Jumbotron fluid>
            <Container>
                <h1>Добро пожаловать</h1>
                <p>
                    Gpt Tabs это портал для обмена табулатурами.
                </p>
                <p>
                    Здесь вы можете искать и скачивать табулатуры к интересующим вас песням, а так же загружать их самостоятельно.
                </p>
                <p>
                    Живи, люби, твори.
                </p>
            </Container>
        </Jumbotron>
        );
    }
}

export default Glagna;
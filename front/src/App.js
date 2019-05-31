import React, {Component} from 'react';
import Header from "./components/header/Header";
import {Container, Jumbotron} from "react-bootstrap";

class App extends Component {
  render(){
    let body;
    body = <Jumbotron fluid>
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
  </Jumbotron>;
  return (
    <div>
      <Header/>
      {body}
    </div>
  );
  }

}

export default App;

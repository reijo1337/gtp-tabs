import React, {Component} from 'react';
import Header from "./components/header/Header";
import Glagna from "./components/glagna/Glagna";
import { Route } from 'react-router-dom';
import AddFile from "./components/AddFile";

class App extends Component {
  render(){
  return (
    <div>
      <Header/>
      <Route exact path="/" component={Glagna}/>
      <Route path="/upload" component={AddFile}/>
    </div>
  );
  }

}

export default App;
